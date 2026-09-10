package ejclient

import (
	"github.com/blackav/ejudge-utils-go/pkg/ejtypes"
)

type ErrorReply struct {
	Num     int32   `json:"num"`
	Symbol  string  `json:"symbol"`
	Message *string `json:"message,omitempty"`
	LogID   *string `json:"log_id,omitempty"`
}

type Reply[T any] struct {
	Ok         bool        `json:"ok"`
	ServerTime int64       `json:"server_time"`
	Action     *string     `json:"action,omitempty"`
	Result     *T          `json:"result,omitempty"`
	Error      *ErrorReply `json:"error,omitempty"`
}

type ListPendingReviewsResult struct {
	Reviews []ejtypes.RunReview `json:"reviews,omitempty"`
}

type ScopedHints struct {
	Hint         *string `json:"hint,omitempty"`
	LanguageHint *string `json:"language_hint,omitempty"`
	ReviewHint   *string `json:"review_hint,omitempty"`
	JudgeHint    *string `json:"judge_hint,omitempty"`
}

type StartReviewResultDetails struct {
	SourceCode        *string      `json:"source_code,omitempty"`
	SourceLanguage    *string      `json:"source_language,omitempty"`
	ProblemStatement  *string      `json:"problem_statement,omitempty"`
	InterfaceLanguage *string      `json:"interface_language,omitempty"`
	Global            *ScopedHints `json:"global,omitempty"`
	Contest           *ScopedHints `json:"contest,omitempty"`
	Problem           *ScopedHints `json:"problem,omitempty"`
	Run               *ejtypes.Run `json:"run,omitempty"`
	CompilerMessages  *string      `json:"compiler_messages,omitempty"`
	CustomReview      *string      `json:"custom_review,omitempty"`
}

type StartReviewResult struct {
	Review  *ejtypes.RunReview        `json:"review,omitempty"`
	Details *StartReviewResultDetails `json:"details,omitempty"`
}

type FinishReviewResult struct {
	Review *ejtypes.RunReview `json:"review,omitempty"`
}

type HeartbeatReviewResult struct {
	Review *ejtypes.RunReview `json:"review,omitempty"`
}

type ListActiveReviewsResult struct {
	Reviews []ejtypes.RunReview `json:"reviews,omitempty"`
}

type GetActiveReviewResult struct {
	Review  *ejtypes.RunReview        `json:"review,omitempty"`
	Details *StartReviewResultDetails `json:"details,omitempty"`
}
