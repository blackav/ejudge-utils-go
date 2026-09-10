package ejtypes

import (
	"time"

	"github.com/google/uuid"
)

type RunStatus uint8

const (
	RunStatusOk               RunStatus = 0
	RunStatusCompileErr       RunStatus = 1
	RunStatusRunTimeErr       RunStatus = 2
	RunStatusTimeLimitErr     RunStatus = 3
	RunStatusPresentationErr  RunStatus = 4
	RunStatusWrongAnswerErr   RunStatus = 5
	RunStatusCheckFailed      RunStatus = 6
	RunStatusPartial          RunStatus = 7
	RunStatusAccepted         RunStatus = 8
	RunStatusIgnored          RunStatus = 9
	RunStatusDisqualified     RunStatus = 10
	RunStatusPending          RunStatus = 11
	RunStatusMemLimitErr      RunStatus = 12
	RunStatusSecurityErr      RunStatus = 13
	RunStatusStyleErr         RunStatus = 14
	RunStatusWallTimeLimitErr RunStatus = 15
	RunStatusPendingReview    RunStatus = 16
	RunStatusRejected         RunStatus = 17
	RunStatusSkipped          RunStatus = 18
	RunStatusSyncErr          RunStatus = 19
	RunStatusVirtualStart     RunStatus = 20
	RunStatusVirtualStop      RunStatus = 21
	RunStatusVirtualEmpty     RunStatus = 22
	RunStatusSummoned         RunStatus = 23
	RunStatusFullRejudge      RunStatus = 95
	RunStatusRunning          RunStatus = 96
	RunStatusCompiled         RunStatus = 97
	RunStatusCompiling        RunStatus = 98
	RunStatusAvailable        RunStatus = 99
)

type RunReviewStatus uint8

const (
	RunReviewStatusRequestedReview RunReviewStatus = 1
	RunReviewStatusWaitingReview   RunReviewStatus = 2
	RunReviewStatusReviewing       RunReviewStatus = 3
	RunReviewStatusWaitingApproval RunReviewStatus = 4
	RunReviewStatusComplete        RunReviewStatus = 5
	RunReviewStatusCanceled        RunReviewStatus = 6
	RunReviewStatusFailed          RunReviewStatus = 7
)

type RunReviewPurpose uint8

const (
	RunReviewPurposeReview    RunReviewPurpose = 1
	RunReviewPurposeHelp      RunReviewPurpose = 2
	RunReviewPurposeJudgeHelp RunReviewPurpose = 3
)

type RunReviewRecommendedStatus uint8

const (
	RunReviewRecommendedStatusOk         RunReviewRecommendedStatus = 0
	RunReviewRecommendedStatusIgnore     RunReviewRecommendedStatus = 9
	RunReviewRecommendedStatusDisqualify RunReviewRecommendedStatus = 10
	RunReviewRecommendedStatusReject     RunReviewRecommendedStatus = 17
	RunReviewRecommendedStatusSummon     RunReviewRecommendedStatus = 23
)

type RunReview struct {
	ReviewUUID              uuid.UUID                   `json:"review_uuid"`
	ContestID               int32                       `json:"contest_id"`
	RunID                   int32                       `json:"run_id"`
	SerialID                *int64                      `json:"serial_id,omitempty"`
	RunSerialID             *int64                      `json:"run_serial_id,omitempty"`
	Generation              int32                       `json:"generation"`
	Status                  RunReviewStatus             `json:"status"`
	Purpose                 RunReviewPurpose            `json:"purpose"`
	RequestUserID           *int32                      `json:"request_user_id,omitempty"`
	ModeratorUserID         *int32                      `json:"moderator_user_id,omitempty"`
	ReviewerUserID          *int32                      `json:"reviewer_user_id,omitempty"`
	ApproverUserID          *int32                      `json:"approver_user_id,omitempty"`
	CreationTime            *time.Time                  `json:"creation_time_iso,omitempty"`
	LastUpdateTime          *time.Time                  `json:"last_update_time_iso,omitempty"`
	ModerationTime          *time.Time                  `json:"moderation_time_iso,omitempty"`
	ReviewStartTime         *time.Time                  `json:"review_start_time_iso,omitempty"`
	ReviewHeartbeatTime     *time.Time                  `json:"review_heartbeat_time_iso,omitempty"`
	ReviewFinishTime        *time.Time                  `json:"review_finish_time_iso,omitempty"`
	ApprovalTime            *time.Time                  `json:"approval_time_iso,omitempty"`
	UserOpenedTime          *time.Time                  `json:"user_opened_time_iso,omitempty"`
	ModerationText          *string                     `json:"moderation_text,omitempty"`
	CustomReview            *string                     `json:"custom_review,omitempty"`
	ReviewSource            *string                     `json:"review_source,omitempty"`
	ReviewAgent             *string                     `json:"review_agent,omitempty"`
	ReviewHeartbeatStatus   *string                     `json:"review_heartbeat_status,omitempty"`
	ReviewResult            *string                     `json:"review_result,omitempty"`
	ReviewJudgeResult       *string                     `json:"review_judge_result,omitempty"`
	ReviewStatistics        *string                     `json:"review_statistics,omitempty"`
	ReviewLog               *string                     `json:"review_log,omitempty"`
	ApprovedText            *string                     `json:"approved_text,omitempty"`
	JudgeApprovedText       *string                     `json:"judge_approved_text,omitempty"`
	Model                   *string                     `json:"model,omitempty"`
	ApproverFeedback        *string                     `json:"approver_feedback,omitempty"`
	UserFeedback            *string                     `json:"user_feedback,omitempty"`
	ReviewSourceSha256      *string                     `json:"review_source_sha256,omitempty"`
	InputTokens             *int64                      `json:"input_tokens,omitempty"`
	CachedInputTokens       *int64                      `json:"cached_input_tokens,omitempty"`
	OutputTokens            *int64                      `json:"output_tokens,omitempty"`
	ReasoningTokens         *int64                      `json:"reasoning_tokens,omitempty"`
	TotalTokens             *int64                      `json:"total_tokens,omitempty"`
	ReviewRecommendedStatus *RunReviewRecommendedStatus `json:"review_recommended_status,omitempty"`
	ApproverReviewMark      *int32                      `json:"approver_review_mark,omitempty"`
	UserOpenedCount         *int32                      `json:"user_opened_count,omitempty"`
	UserReviewMark          *int32                      `json:"user_review_mark,omitempty"`
	ReviewApprovedAsIs      *bool                       `json:"review_approved_as_is,omitempty"`
	StatusApprovedAsIs      *bool                       `json:"status_approved_as_is,omitempty"`
	AIGenerationScore       *int32                      `json:"ai_generation_score,omitempty"`
}

type Run struct {
	RunID       int32      `json:"run_id"`
	ContestID   int32      `json:"contest_id"`
	RunUUID     *uuid.UUID `json:"run_uuid,omitempty"`
	Status      RunStatus  `json:"status"`
	StatusStr   string     `json:"status_str,omitempty"`
	StatusDesc  string     `json:"status_desc,omitempty"`
	UserID      *int32     `json:"user_id,omitempty"`
	UserLogin   *string    `json:"user_login,omitempty"`
	UserName    *string    `json:"user_name,omitempty"`
	ProbID      int32      `json:"prob_id"`
	ProbName    string     `json:"prob_name,omitempty"`
	ProbUUID    *uuid.UUID `json:"prob_uuid,omitempty"`
	LangID      int32      `json:"lang_id,omitempty"`
	LangName    *string    `json:"lang_name,omitempty"`
	Size        int64      `json:"size,omitempty"`
	LocaleID    *int32     `json:"locale_id,omitempty"`
	RawScore    *int32     `json:"raw_score,omitempty"`
	RawTest     *int32     `json:"raw_test,omitempty"`
	VerdictBits *uint64    `json:"verdict_bits,omitempty"`
}
