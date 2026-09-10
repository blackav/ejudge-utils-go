package ejtypes

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

func (v RunReviewRecommendedStatus) String() string {
	switch v {
	case RunReviewRecommendedStatusOk:
		return "ok"
	case RunReviewRecommendedStatusIgnore:
		return "ignore"
	case RunReviewRecommendedStatusDisqualify:
		return "disqualify"
	case RunReviewRecommendedStatusReject:
		return "reject"
	case RunReviewRecommendedStatusSummon:
		return "summon"
	default:
		return ""
	}
}

func (v RunReviewRecommendedStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(v.String())
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

var ErrInvalidRunReviewRecommendedStatus = errors.New("invalid recommended status")

func (v *RunReviewRecommendedStatus) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	switch {
	case strings.EqualFold(j, "ok"):
		*v = RunReviewRecommendedStatusOk
	case strings.EqualFold(j, "ignore"):
		*v = RunReviewRecommendedStatusIgnore
	case strings.EqualFold(j, "disqualify"):
		*v = RunReviewRecommendedStatusDisqualify
	case strings.EqualFold(j, "reject"):
		*v = RunReviewRecommendedStatusReject
	case strings.EqualFold(j, "summon"):
		*v = RunReviewRecommendedStatusSummon
	default:
		return ErrInvalidRunReviewRecommendedStatus
	}
	return nil
}

var runReviewStatusValues = []string{
	"",
	"requested_review",
	"waiting_review",
	"reviewing",
	"waiting_approval",
	"complete",
	"canceled",
	"failed",
}

func (v RunReviewStatus) String() string {
	if int(v) <= 0 || int(v) >= len(runReviewStatusValues) {
		return ""
	}
	return runReviewStatusValues[int(v)]
}

func (v RunReviewStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(v.String())
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

var ErrInvalidRunReviewStatus = errors.New("invalid review status")

func (v *RunReviewStatus) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	switch {
	case strings.EqualFold(j, "requested_review"):
		*v = RunReviewStatusRequestedReview
	case strings.EqualFold(j, "waiting_review"):
		*v = RunReviewStatusWaitingReview
	case strings.EqualFold(j, "reviewing"):
		*v = RunReviewStatusReviewing
	case strings.EqualFold(j, "waiting_approval"):
		*v = RunReviewStatusWaitingApproval
	case strings.EqualFold(j, "complete"):
		*v = RunReviewStatusComplete
	case strings.EqualFold(j, "canceled"):
		*v = RunReviewStatusCanceled
	case strings.EqualFold(j, "failed"):
		*v = RunReviewStatusFailed
	default:
		return ErrInvalidRunReviewStatus
	}
	return nil
}

var runReviewPurposeValues = []string{
	"",
	"review",
	"help",
	"judge_help",
}

func (v RunReviewPurpose) String() string {
	if int(v) <= 0 || int(v) >= len(runReviewPurposeValues) {
		return ""
	}
	return runReviewPurposeValues[int(v)]
}

func (v RunReviewPurpose) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(v.String())
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

var ErrInvalidRunReviewPurpose = errors.New("invalid review purpose")

func (v *RunReviewPurpose) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	switch {
	case strings.EqualFold(j, "review"):
		*v = RunReviewPurposeReview
	case strings.EqualFold(j, "help"):
		*v = RunReviewPurposeHelp
	case strings.EqualFold(j, "judge_help"):
		*v = RunReviewPurposeJudgeHelp
	default:
		return ErrInvalidRunReviewPurpose
	}
	return nil
}
