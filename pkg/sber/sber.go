package sber

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/blackav/ejudge-utils-go/pkg/aitypes"
	"github.com/google/uuid"
)

type Config struct {
	AuthURL  string
	GenURL   string
	AuthKey  string
	CertFile string
	Model    string
}

type Token struct {
	AccessToken string
	ExpiresAt   time.Time
}

type SberState struct {
	cfg        *Config
	token      atomic.Pointer[Token]
	httpClient *http.Client
}

func loadCertificates(file string) (*tls.Config, error) {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	certs, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read certs file '%s': %w", file, err)
	}
	ok := rootCAs.AppendCertsFromPEM(certs)
	_ = ok
	return &tls.Config{
		RootCAs: rootCAs,
	}, nil
}

func New(_ context.Context, cfg *Config) (aitypes.Generator, error) {
	var tr *http.Transport
	if cfg.CertFile != "" {
		tls, err := loadCertificates(cfg.CertFile)
		if err != nil {
			return nil, err
		}
		tr = &http.Transport{
			TLSClientConfig: tls,
		}
	}
	return &SberState{
		cfg: cfg,
		httpClient: &http.Client{
			Transport: tr,
		},
	}, nil
}

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func (s *SberState) RefreshToken(ctx context.Context) error {
	token := s.token.Load()
	currentTime := time.Now()
	if token != nil && currentTime.Add(time.Minute).Before(token.ExpiresAt) {
		return nil
	}

	form := url.Values{}
	form.Set("scope", "GIGACHAT_API_PERS")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.AuthURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("RqUID", uuid.NewString())
	req.Header.Add("Authorization", "Basic "+s.cfg.AuthKey)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Refresh failed: status code: %d", resp.StatusCode)
	}

	d := json.NewDecoder(resp.Body)
	rtr := refreshTokenResponse{}
	err = d.Decode(&rtr)
	if err != nil {
		return err
	}
	expiresAt := time.UnixMilli(rtr.ExpiresAt)
	s.token.Store(&Token{
		AccessToken: rtr.AccessToken,
		ExpiresAt:   expiresAt,
	})
	//fmt.Println("Token: " + rtr.AccessToken)
	//fmt.Println("Expires: " + expiresAt.Format(time.RFC3339))
	return nil
}

type listModelResponseItem struct {
	ID      string `json:"id"`
	Object  string `json:"object,omitempty"`
	Type    string `json:"type,omitempty"`
	OwnedBy string `json:"owned_by,omitempty"`
}

type listModelResponse struct {
	Data []listModelResponseItem
}

func (s *SberState) ListModels(ctx context.Context) ([]aitypes.ModelInfo, error) {
	err := s.RefreshToken(ctx)
	if err != nil {
		return nil, err
	}
	fullURL := s.cfg.GenURL + "/v1/models"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.token.Load().AccessToken)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("List failed: status code: %d", resp.StatusCode)
	}
	d := json.NewDecoder(resp.Body)
	jr := listModelResponse{}
	err = d.Decode(&jr)
	if err != nil {
		return nil, err
	}
	res := []aitypes.ModelInfo{}
	for i := range jr.Data {
		src := &jr.Data[i]
		res = append(res, aitypes.ModelInfo{
			ID:      src.ID,
			Object:  src.Object,
			Type:    src.Type,
			OwnedBy: src.OwnedBy,
		})
	}
	return res, nil
}

type completionRequestMessageContent struct {
	Text string `json:"text,omitempty"`
}

type completionRequestMessage struct {
	Role    string                            `json:"role,omitempty"`
	Content []completionRequestMessageContent `json:"content,omitempty"`
}

type completionRequest struct {
	Model    string                     `json:"model"`
	Messages []completionRequestMessage `json:"messages,omitempty"`
}

type completionResultMessageContent struct {
	Text string `json:"text,omitempty"`
}

type completionResultMessage struct {
	Role    string                           `json:"role,omitempty"`
	Content []completionResultMessageContent `json:"content,omitempty"`
}

type completionResultUsage struct {
	InputTokens  *int64 `json:"input_tokens,omitempty"`
	OutputTokens *int64 `json:"output_tokens,omitempty"`
	TotalTokens  *int64 `json:"total_tokens,omitempty"`
}

type completionResult struct {
	Model        string                    `json:"model,omitempty"`
	CreatedAt    int64                     `json:"created_at,omitempty"`
	Messages     []completionResultMessage `json:"messages,omitempty"`
	FinishReason string                    `json:"finish_reason,omitempty"`
	Usage        *completionResultUsage    `json:"usage,omitempty"`
}

func copyUsage(src *completionResultUsage) *aitypes.CompletionResultUsage {
	if src == nil {
		return nil
	}
	return &aitypes.CompletionResultUsage{
		InputTokens:  src.InputTokens,
		OutputTokens: src.OutputTokens,
		TotalTokens:  src.TotalTokens,
	}
}

func (s *SberState) SimpleCompletion(ctx context.Context, model string, text string) (*aitypes.CompletionResult, error) {
	err := s.RefreshToken(ctx)
	if err != nil {
		return nil, err
	}
	if model == "" {
		model = s.cfg.Model
	}
	fullURL := s.cfg.GenURL + "/v2/chat/completions"
	ro := completionRequest{
		Model: model,
		Messages: []completionRequestMessage{
			{
				Role: "user",
				Content: []completionRequestMessageContent{
					{
						Text: text,
					},
				},
			},
		},
	}
	jr, err := json.Marshal(ro)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewBuffer(jr))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.token.Load().AccessToken)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Completion failed: status code: %d", resp.StatusCode)
	}
	d := json.NewDecoder(resp.Body)
	scr := completionResult{}
	err = d.Decode(&scr)
	if err != nil {
		return nil, err
	}

	resultText := ""
	if len(scr.Messages) > 0 {
		if len(scr.Messages[0].Content) > 0 {
			resultText = scr.Messages[0].Content[0].Text
		}
	}

	return &aitypes.CompletionResult{
		Model:     scr.Model,
		CreatedAt: time.Unix(scr.CreatedAt, 0),
		Text:      resultText,
		Usage:     copyUsage(scr.Usage),
	}, nil
}
