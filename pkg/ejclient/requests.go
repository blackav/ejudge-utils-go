package ejclient

import (
	"bytes"
	"strconv"

	"github.com/blackav/ejudge-utils-go/pkg/ejtypes"
	"github.com/google/uuid"
)

type ListPendingReviewsRequestIDs struct {
	Any        bool
	ContestIDs []int32
}

type ListPendingReviewsRequest struct {
	DateMode   RequestDateMode               `json:"date_mode,omitempty"`
	ContestIDs *ListPendingReviewsRequestIDs `json:"contests_ids,omitempty"`
}

type StartReviewRequest struct {
	ReviewUUID      uuid.UUID       `json:"review_uuid"`
	Agent           *string         `json:"agent,omitempty"`
	HeartbeatStatus *string         `json:"heartbeat_status,omitempty"`
	DateMode        RequestDateMode `json:"date_mode,omitempty"`
}

type FinishReviewRequest struct {
	ReviewUUID        uuid.UUID                           `json:"review_uuid"`
	DateMode          RequestDateMode                     `json:"date_mode,omitempty"`
	Status            *ejtypes.RunReviewStatus            `json:"status"`
	RecommendedStatus *ejtypes.RunReviewRecommendedStatus `json:"recommended_status,omitempty"`
	Result            *string                             `json:"result,omitempty"`
	JudgeResult       *string                             `json:"judge_result,omitempty"`
	Agent             *string                             `json:"agent,omitempty"`
	Statistics        *string                             `json:"statistics,omitempty"`
	InputTokens       *int64                              `json:"input_tokens,omitempty"`
	CachedInputTokens *int64                              `json:"cached_input_tokens,omitempty"`
	OutputTokens      *int64                              `json:"output_tokens,omitempty"`
	ReasoningTokens   *int64                              `json:"reasoning_tokens,omitempty"`
	TotalTokens       *int64                              `json:"total_tokens,omitempty"`
	Model             *string                             `json:"model,omitempty"`
	AIGenerationScore *int32                              `json:"ai_generation_score,omitempty"`
	Log               *string                             `json:"log,omitempty"`
}

type HeartbeatReviewRequest struct {
	ReviewUUID      uuid.UUID       `json:"review_uuid"`
	DateMode        RequestDateMode `json:"date_mode,omitempty"`
	HeartbeatStatus *string         `json:"heartbeat_status,omitempty"`
}

type ListActiveReviewsRequest struct {
	DateMode RequestDateMode `json:"date_mode,omitempty"`
}

type GetActiveReviewRequest struct {
	DateMode   RequestDateMode `json:"date_mode,omitempty"`
	ReviewUUID uuid.UUID       `json:"review_uuid"`
}

func (v ListPendingReviewsRequestIDs) String() string {
	buffer := bytes.Buffer{}
	if v.Any {
		_ = buffer.WriteByte('*')
	} else {
		sep := ""
		for _, id := range v.ContestIDs {
			buffer.WriteString(sep)
			buffer.WriteString(strconv.Itoa(int(id)))
			sep = ","
		}
	}
	return buffer.String()
}

func String(s string) *string {
	return &s
}

func Int32(x int32) *int32 {
	return &x
}

func Int64(x int64) *int64 {
	return &x
}

func Status(x ejtypes.RunReviewStatus) *ejtypes.RunReviewStatus {
	return &x
}

func RunReviewRecommendedStatus(x ejtypes.RunReviewRecommendedStatus) *ejtypes.RunReviewRecommendedStatus {
	return &x
}
