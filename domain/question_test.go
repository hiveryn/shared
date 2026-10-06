package domain

import (
	"errors"
	"strings"
	"testing"
)

func intPtr(v int) *int { return &v }

func TestAskQuestionRequestNormalize(t *testing.T) {
	ok, err := AskQuestionRequest{
		Question:         "  Ship it?  ",
		Answers:          []string{" Yes ", "No"},
		RecommendedIndex: intPtr(1),
	}.Normalize()
	if err != nil {
		t.Fatalf("Normalize: %v", err)
	}
	if ok.Question != "Ship it?" || ok.Answers[0] != "Yes" || *ok.RecommendedIndex != 1 {
		t.Fatalf("not normalized: %+v", ok)
	}

	cases := []struct {
		name  string
		req   AskQuestionRequest
		field string
	}{
		{"blank question", AskQuestionRequest{Question: " ", Answers: []string{"a", "b"}, RecommendedIndex: intPtr(0)}, "question"},
		{"long question", AskQuestionRequest{Question: strings.Repeat("é", MaxQuestionLength+1), Answers: []string{"a", "b"}, RecommendedIndex: intPtr(0)}, "question"},
		{"one answer", AskQuestionRequest{Question: "q", Answers: []string{"a"}, RecommendedIndex: intPtr(0)}, "answers"},
		{"too many answers", AskQuestionRequest{Question: "q", Answers: []string{"a", "b", "c", "d", "e", "f", "g"}, RecommendedIndex: intPtr(0)}, "answers"},
		{"missing answers", AskQuestionRequest{Question: "q", RecommendedIndex: intPtr(0)}, "answers"},
		{"blank answer", AskQuestionRequest{Question: "q", Answers: []string{"a", " "}, RecommendedIndex: intPtr(0)}, "answers[1]"},
		{"long answer", AskQuestionRequest{Question: "q", Answers: []string{"a", strings.Repeat("x", MaxQuestionAnswerLength+1)}, RecommendedIndex: intPtr(0)}, "answers[1]"},
		{"duplicate answer", AskQuestionRequest{Question: "q", Answers: []string{"Yes", " yes"}, RecommendedIndex: intPtr(0)}, "answers[1]"},
		{"missing recommendation", AskQuestionRequest{Question: "q", Answers: []string{"a", "b"}}, "recommended_index"},
		{"negative recommendation", AskQuestionRequest{Question: "q", Answers: []string{"a", "b"}, RecommendedIndex: intPtr(-1)}, "recommended_index"},
		{"1-based recommendation", AskQuestionRequest{Question: "q", Answers: []string{"a", "b"}, RecommendedIndex: intPtr(2)}, "recommended_index"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.req.Normalize()
			var verr *ValidationError
			if !errors.As(err, &verr) || verr.Field != tc.field {
				t.Fatalf("want validation error on %s, got %v", tc.field, err)
			}
		})
	}
}

func TestNotifyAndResponseBounds(t *testing.T) {
	if got, err := NormalizeNotifyMessage(" done "); err != nil || got != "done" {
		t.Fatalf("NormalizeNotifyMessage = %q, %v", got, err)
	}
	if _, err := NormalizeNotifyMessage(strings.Repeat("x", MaxNotifyMessageLength+1)); err == nil {
		t.Fatal("over-long notify message accepted")
	}
	if _, err := NormalizeNotifyMessage("  "); err == nil {
		t.Fatal("blank notify message accepted")
	}
	if _, err := NormalizeQuestionResponse(strings.Repeat("x", MaxQuestionResponseLength)); err != nil {
		t.Fatalf("max-length answer rejected: %v", err)
	}
	if _, err := NormalizeQuestionResponse(strings.Repeat("x", MaxQuestionResponseLength+1)); err == nil {
		t.Fatal("over-long answer accepted")
	}
}
