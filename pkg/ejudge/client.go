package ejudge

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/blackav/ejudge-utils-go/pkg/ejclient"
)

type DefaultClient struct {
	Config ejclient.Config
}

const (
	listPendingReviewsPath = "/ej/api/v1/master/list-pending-reviews-json"
	startReviewPath        = "/ej/api/v1/master/start-review-json"
	heartbeatReviewPath    = "/ej/api/v1/master/heartbeat-review-json"
	finishReviewPath       = "/ej/api/v1/master/finish-review-json"
	listActiveReviewsPath  = "/ej/api/v1/master/list-active-reviews-json"
	getActiveReviewPath    = "/ej/api/v1/master/get-active-review-json"
)

func (c *DefaultClient) makeAuthHeader(q *http.Request) {
	q.Header.Add("Authorization", "Bearer AQAA"+c.Config.Token)
	q.Header.Add("X-Contest-ID", strconv.Itoa(int(c.Config.ContestID)))
}

func genericGet[T any](ctx context.Context, c *DefaultClient, fullURL string) (int, *ejclient.Reply[T], error) {
	q, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return 0, nil, err
	}
	c.makeAuthHeader(q)
	client := &http.Client{}
	r, err := client.Do(q)
	if err != nil {
		return 0, nil, err
	}
	defer r.Body.Close()
	d := json.NewDecoder(r.Body)
	j := ejclient.Reply[T]{}
	err = d.Decode(&j)
	if err != nil {
		return 0, nil, err
	}
	return r.StatusCode, &j, nil
}

func genericPost[S any, T any](ctx context.Context, c *DefaultClient, path string, req *S) (int, *ejclient.Reply[T], error) {
	fullURL := c.Config.URL + path
	jr, err := json.Marshal(req)
	if err != nil {
		return 0, nil, err
	}
	q, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewBuffer(jr))
	if err != nil {
		return 0, nil, err
	}
	c.makeAuthHeader(q)
	q.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(q)
	if err != nil {
		return 0, nil, err
	}
	defer r.Body.Close()
	d := json.NewDecoder(r.Body)
	j := ejclient.Reply[T]{}
	err = d.Decode(&j)
	if err != nil {
		return 0, nil, err
	}
	return r.StatusCode, &j, nil
}

func (c *DefaultClient) ListPendingReviews(ctx context.Context, req *ejclient.ListPendingReviewsRequest) (int, *ejclient.Reply[ejclient.ListPendingReviewsResult], error) {
	params := url.Values{}
	params.Add("date_mode", req.DateMode.String())
	fullURL := c.Config.URL + listPendingReviewsPath + "?" + params.Encode()
	return genericGet[ejclient.ListPendingReviewsResult](ctx, c, fullURL)
}

func (c *DefaultClient) StartReview(ctx context.Context, req *ejclient.StartReviewRequest) (int, *ejclient.Reply[ejclient.StartReviewResult], error) {
	return genericPost[ejclient.StartReviewRequest, ejclient.StartReviewResult](ctx, c, startReviewPath, req)
}

func (c *DefaultClient) HeartbeatReview(ctx context.Context, req *ejclient.HeartbeatReviewRequest) (int, *ejclient.Reply[ejclient.HeartbeatReviewResult], error) {
	return genericPost[ejclient.HeartbeatReviewRequest, ejclient.HeartbeatReviewResult](ctx, c, heartbeatReviewPath, req)
}

func (c *DefaultClient) FinishReview(ctx context.Context, req *ejclient.FinishReviewRequest) (int, *ejclient.Reply[ejclient.FinishReviewResult], error) {
	return genericPost[ejclient.FinishReviewRequest, ejclient.FinishReviewResult](ctx, c, finishReviewPath, req)
}

func (c *DefaultClient) ListActiveReviews(ctx context.Context, req *ejclient.ListActiveReviewsRequest) (int, *ejclient.Reply[ejclient.ListActiveReviewsResult], error) {
	params := url.Values{}
	params.Add("date_mode", req.DateMode.String())
	fullURL := c.Config.URL + listActiveReviewsPath + "?" + params.Encode()
	return genericGet[ejclient.ListActiveReviewsResult](ctx, c, fullURL)
}

func (c *DefaultClient) GetActiveReview(ctx context.Context, req *ejclient.GetActiveReviewRequest) (int, *ejclient.Reply[ejclient.GetActiveReviewResult], error) {
	params := url.Values{}
	params.Add("date_mode", req.DateMode.String())
	params.Add("review_uuid", req.ReviewUUID.String())
	fullURL := c.Config.URL + getActiveReviewPath + "?" + params.Encode()
	return genericGet[ejclient.GetActiveReviewResult](ctx, c, fullURL)
}

func New(cfg ejclient.Config) (ejclient.Client, error) {
	return &DefaultClient{
		Config: cfg,
	}, nil
}
