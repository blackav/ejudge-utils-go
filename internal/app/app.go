package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/blackav/ejudge-utils-go/pkg/ejclient"
	"github.com/blackav/ejudge-utils-go/pkg/ejtypes"
	"github.com/blackav/ejudge-utils-go/pkg/gentypes"
	"github.com/blackav/ejudge-utils-go/pkg/slogt"
	"github.com/cbroglie/mustache"
	"github.com/google/uuid"
)

type Config struct {
	Ejudge       ejclient.Client
	Generator    gentypes.Generator
	ID           string
	TemplateFile string
}

type App struct {
	Ejudge       ejclient.Client
	Generator    gentypes.Generator
	ID           string
	TemplateFile string
}

func New(cfg Config) *App {
	return &App{
		Ejudge:       cfg.Ejudge,
		Generator:    cfg.Generator,
		ID:           cfg.ID,
		TemplateFile: cfg.TemplateFile,
	}
}

func (a *App) getPendingReviews(ctx context.Context) ([]ejtypes.RunReview, error) {
	status, resp, err := a.Ejudge.ListPendingReviews(ctx, &ejclient.ListPendingReviewsRequest{
		DateMode: ejclient.RequestDateModeISO,
		ContestIDs: &ejclient.ListPendingReviewsRequestIDs{
			Any: true,
		},
	})
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("get pending reviews http status %d", status)
	}
	if resp == nil || !resp.Ok || resp.Result == nil {
		return nil, fmt.Errorf("get pending reviews failed on server")
	}

	return resp.Result.Reviews, nil
}

func convertToMap(obj any) (map[string]any, error) {
	bb, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	res := map[string]any{}
	err = json.Unmarshal(bb, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *App) doReview(ctx context.Context, rr *ejclient.StartReviewResult) error {
	tb, err := os.ReadFile(a.TemplateFile)
	if err != nil {
		return fmt.Errorf("failed to open template file %s: %w", a.TemplateFile, err)
	}
	vars, err := convertToMap(rr)
	if err != nil {
		return err
	}
	text, err := mustache.Render(string(tb), vars)
	if err != nil {
		return err
	}

	review, err := a.Generator.SimpleCompletion(ctx, "", text)
	if err != nil {
		return err
	}
	/*
			review := aitypes.CompletionResult{
			CreatedAt: time.Now(),
			Model:     "xxx",
			Usage:     nil,
			Text:      text,
		}
	*/

	req := ejclient.FinishReviewRequest{
		ReviewUUID: rr.Review.ReviewUUID,
		DateMode:   ejclient.RequestDateModeISO,
		Status:     ejclient.Status(ejtypes.RunReviewStatusComplete),
		Result:     ejclient.String(review.Text),
		Agent:      ejclient.String(a.ID),
		Model:      ejclient.String(review.Model),
	}
	if review.Usage != nil {
		req.InputTokens = review.Usage.InputTokens
		req.OutputTokens = review.Usage.OutputTokens
		req.TotalTokens = review.Usage.TotalTokens
	}
	status, fr, err := a.Ejudge.FinishReview(ctx, &req)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("finish review http status %d", status)
	}
	if fr == nil || !fr.Ok || fr.Result == nil {
		return fmt.Errorf("finish review failed on server")
	}
	return nil
}

func (a *App) reportErrorUpstream(ctx context.Context, reviewUUID uuid.UUID, message string) {
	// loop indefinitely, until success...
	// FIXME: analyze errors
loop:
	for {
		status, fr, err := a.Ejudge.FinishReview(ctx, &ejclient.FinishReviewRequest{
			ReviewUUID: reviewUUID,
			DateMode:   ejclient.RequestDateModeISO,
			Status:     ejclient.Status(ejtypes.RunReviewStatusFailed),
			Agent:      ejclient.String(a.ID),
			Log:        ejclient.String(message),
		})
		if err == nil && status == http.StatusOK && fr != nil && fr.Ok {
			// all good
			return
		}
		if err != nil {
			slog.Error("finish review failed", slogt.Error(err))
		} else if status != http.StatusOK {
			slog.Error("finish review http failed", slog.Int("status", status))
		} else if fr == nil {
			slog.Error("finish review returned nil")
		} else if !fr.Ok {
			slog.Error("finish review failed on server")
		}
		select {
		case <-ctx.Done():
			break loop
		case <-time.After(5 * time.Second):
			continue loop
		}
	}
}

func (a *App) handleReview(ctx context.Context, rr *ejtypes.RunReview) error {
	currentTime := time.Now()
	status, resp, err := a.Ejudge.StartReview(ctx, &ejclient.StartReviewRequest{
		DateMode:        ejclient.RequestDateModeISO,
		ReviewUUID:      rr.ReviewUUID,
		Agent:           ejclient.String(a.ID),
		HeartbeatStatus: ejclient.String("review started + " + currentTime.Format(time.RFC3339)),
	})
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("handle review http status %d", status)
	}
	if resp == nil || !resp.Ok || resp.Result == nil {
		return fmt.Errorf("get pending reviews failed on server")
	}
	err = a.doReview(ctx, resp.Result)
	if err != nil {
		a.reportErrorUpstream(ctx, rr.ReviewUUID, err.Error())
		return err
	}
	return nil
}

func (a *App) Run(ctx context.Context) {
loop:
	for {
		prs, err := a.getPendingReviews(ctx)
		if err != nil {
			slog.Error("get pending reviews failed", slogt.Error(err))
			// FIXME: do exponential retry...
			select {
			case <-ctx.Done():
				break loop
			case <-time.After(5 * time.Second):
				continue loop
			}
		}
		if len(prs) == 0 {
			slog.Info("no pending reviews")
			select {
			case <-ctx.Done():
				break loop
			case <-time.After(60 * time.Second):
				continue loop
			}
		}
		for i := range prs {
			err = a.handleReview(ctx, &prs[i])
			if err != nil {
				slog.Error("review failed", slog.String("reviewUUID", prs[i].ReviewUUID.String()), slogt.Error(err))
			} else {
				slog.Info("review success", slog.String("reviewUUID", prs[i].ReviewUUID.String()))
			}
		}
	}
}
