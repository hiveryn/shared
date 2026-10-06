package domain

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// Agent notifications and questions.
//
// An agent may notify the user (a phone alert through ntfy that returns at
// once) or ask a blocking question: a phone alert plus a pending question the
// user answers in the originating desktop session, by choosing one of the
// suggested answers or typing free text. The agent's call waits at most
// QuestionTimeout; the daemon enforces that expiry itself. A question is not an
// Intent: it approves nothing, has no policy and never resolves on its own to
// an answer — in particular the recommended answer is never selected for the
// user.
//
// Lengths are counted in Unicode code points (runes), here and in question.ts.

const (
	// MaxNotifyMessageLength bounds notify's shortMessage.
	MaxNotifyMessageLength = 500
	// MaxQuestionLength bounds askQuestion's question.
	MaxQuestionLength = 1000
	// MinQuestionAnswers / MaxQuestionAnswers bound the suggested answers.
	MinQuestionAnswers = 2
	MaxQuestionAnswers = 6
	// MaxQuestionAnswerLength bounds one suggested answer.
	MaxQuestionAnswerLength = 200
	// MaxQuestionResponseLength bounds the answer the user submits, chosen or
	// typed.
	MaxQuestionResponseLength = 4000

	// QuestionTimeout is how long a question stays answerable and the agent's
	// call stays open. Agent launches give the Hiveryn MCP server a longer
	// client deadline so this reply always arrives.
	QuestionTimeout = time.Hour

	// NotificationSentText is notify's exact reply after the alert was
	// accepted by the ntfy server. It is not proof the phone received it.
	NotificationSentText = "Notification sent to user"
	// QuestionTimeoutText is askQuestion's exact reply when nobody answered
	// within QuestionTimeout.
	QuestionTimeoutText = "User didn't respond within 1 hour, stop here and wait for user to get back"
)

// QuestionStatus is a question's lifecycle. A question starts pending and is
// resolved exactly once.
type QuestionStatus string

const (
	QuestionPending QuestionStatus = "pending"
	// QuestionAnswered: the user submitted an answer; the agent received it.
	QuestionAnswered QuestionStatus = "answered"
	// QuestionExpired: nobody answered within QuestionTimeout; the agent was
	// told to stop and wait for the user.
	QuestionExpired QuestionStatus = "expired"
	// QuestionCancelled: the question ended without an answer because the
	// agent stopped waiting or its session ended. Reason says which.
	QuestionCancelled QuestionStatus = "cancelled"
	// QuestionInterrupted: the daemon stopped while the question was pending.
	// The agent's call did not survive it.
	QuestionInterrupted QuestionStatus = "interrupted"
)

func (s QuestionStatus) Terminal() bool {
	switch s {
	case QuestionAnswered, QuestionExpired, QuestionCancelled, QuestionInterrupted:
		return true
	default:
		return false
	}
}

// NotifyRequest is the agent's notify call.
type NotifyRequest struct {
	Message string `json:"message"`
}

// NotifyResult is returned once the ntfy server accepted the alert.
type NotifyResult struct {
	Text string `json:"text"`
}

// AskQuestionRequest is the agent's askQuestion call. RecommendedIndex is
// 0-based into Answers and required; it is a pointer so a missing value is
// distinguishable from 0.
type AskQuestionRequest struct {
	Question         string   `json:"question"`
	Answers          []string `json:"answers"`
	RecommendedIndex *int     `json:"recommended_index"`
}

// Normalize validates the request and returns it with the question and every
// answer trimmed. Answers must be unique, ignoring case and surrounding space.
func (r AskQuestionRequest) Normalize() (AskQuestionRequest, error) {
	question := strings.TrimSpace(r.Question)
	if question == "" {
		return AskQuestionRequest{}, &ValidationError{Field: "question", Message: "is required"}
	}
	if n := utf8.RuneCountInString(question); n > MaxQuestionLength {
		return AskQuestionRequest{}, &ValidationError{Field: "question", Message: fmt.Sprintf("must be at most %d characters, got %d", MaxQuestionLength, n)}
	}
	if len(r.Answers) < MinQuestionAnswers || len(r.Answers) > MaxQuestionAnswers {
		return AskQuestionRequest{}, &ValidationError{Field: "answers", Message: fmt.Sprintf("must list %d to %d suggested answers, got %d", MinQuestionAnswers, MaxQuestionAnswers, len(r.Answers))}
	}
	answers := make([]string, len(r.Answers))
	seen := make(map[string]int, len(r.Answers))
	for i, raw := range r.Answers {
		answer := strings.TrimSpace(raw)
		if answer == "" {
			return AskQuestionRequest{}, &ValidationError{Field: fmt.Sprintf("answers[%d]", i), Message: "must not be blank"}
		}
		if n := utf8.RuneCountInString(answer); n > MaxQuestionAnswerLength {
			return AskQuestionRequest{}, &ValidationError{Field: fmt.Sprintf("answers[%d]", i), Message: fmt.Sprintf("must be at most %d characters, got %d", MaxQuestionAnswerLength, n)}
		}
		key := strings.ToLower(answer)
		if first, dup := seen[key]; dup {
			return AskQuestionRequest{}, &ValidationError{Field: fmt.Sprintf("answers[%d]", i), Message: fmt.Sprintf("duplicates answers[%d]", first)}
		}
		seen[key] = i
		answers[i] = answer
	}
	if r.RecommendedIndex == nil {
		return AskQuestionRequest{}, &ValidationError{Field: "recommended_index", Message: fmt.Sprintf("is required: the 0-based index of the recommended answer (0 to %d)", len(answers)-1)}
	}
	index := *r.RecommendedIndex
	if index < 0 || index >= len(answers) {
		return AskQuestionRequest{}, &ValidationError{Field: "recommended_index", Message: fmt.Sprintf("must be a 0-based index into answers (0 to %d), got %d", len(answers)-1, index)}
	}
	return AskQuestionRequest{Question: question, Answers: answers, RecommendedIndex: &index}, nil
}

// NormalizeNotifyMessage validates and trims notify's message.
func NormalizeNotifyMessage(message string) (string, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return "", &ValidationError{Field: "message", Message: "is required"}
	}
	if n := utf8.RuneCountInString(message); n > MaxNotifyMessageLength {
		return "", &ValidationError{Field: "message", Message: fmt.Sprintf("must be at most %d characters, got %d", MaxNotifyMessageLength, n)}
	}
	return message, nil
}

// NormalizeQuestionResponse validates and trims the user's answer.
func NormalizeQuestionResponse(answer string) (string, error) {
	answer = strings.TrimSpace(answer)
	if answer == "" {
		return "", &ValidationError{Field: "answer", Message: "is required"}
	}
	if n := utf8.RuneCountInString(answer); n > MaxQuestionResponseLength {
		return "", &ValidationError{Field: "answer", Message: fmt.Sprintf("must be at most %d characters, got %d", MaxQuestionResponseLength, n)}
	}
	return answer, nil
}

// QuestionOrigin is where a question came from, derived from the session
// record daemon-side. ArchitectKey is empty for an Action session, whose
// ExecutionID and Action name it carries instead.
type QuestionOrigin struct {
	ArchitectKey string      `json:"architect_key"`
	SessionID    string      `json:"session_id"`
	SessionType  SessionType `json:"session_type"`
	TicketID     string      `json:"ticket_id,omitempty"`
	Action       string      `json:"action,omitempty"`
	ExecutionID  string      `json:"execution_id,omitempty"`
}

// AgentQuestion is one question as the desktop renders it. It travels on the
// session event stream: type "question", status "required" carries it in raw
// (question_id, question, answers, recommended_index, origin, created_at,
// expires_at); status "resolved" carries question_id, status and, when
// answered, answer, or reason otherwise. The resolved event is published for
// every resolution, so the backlog replays a pending question only while it is
// still answerable.
type AgentQuestion struct {
	ID               string         `json:"question_id"`
	Question         string         `json:"question"`
	Answers          []string       `json:"answers"`
	RecommendedIndex int            `json:"recommended_index"`
	Origin           QuestionOrigin `json:"origin"`
	Status           QuestionStatus `json:"status"`
	CreatedAt        time.Time      `json:"created_at"`
	ExpiresAt        time.Time      `json:"expires_at"`
	Answer           string         `json:"answer,omitempty"`
	Reason           string         `json:"reason,omitempty"`
}

// AnswerQuestionRequest is the desktop's answer: one suggested answer's text
// or free text, sent verbatim to the agent.
type AnswerQuestionRequest struct {
	Answer string `json:"answer"`
}

// AskQuestionResult is the daemon's reply to the agent's blocking call, sent
// once the question resolved. Text is what the agent receives: the answer, or
// QuestionTimeoutText, or an explanation for any other end.
type AskQuestionResult struct {
	QuestionID string         `json:"question_id"`
	Status     QuestionStatus `json:"status"`
	Text       string         `json:"text"`
}
