package annotator

import (
	"io"
	"strings"
	"testing"
)

func TestReadKeyFrom_ValidKeys(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"p", "pass"},
		{"r", "review"},
		{"f", "fail"},
		{"x", "skip"},
	}
	for _, tc := range cases {
		got, err := ReadKeyFrom(strings.NewReader(tc.input))
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.input, err)
		}
		if got != Outcome(tc.expected) {
			t.Errorf("input %q: got %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestReadKeyFrom_SkipsInvalidThenReads(t *testing.T) {
	got, err := ReadKeyFrom(strings.NewReader("zzp"))
	if err != nil {
		t.Fatal(err)
	}
	if got != OutcomePass {
		t.Errorf("got %q, want %q", got, OutcomePass)
	}
}

func TestReadKeyFrom_EOFReturnsError(t *testing.T) {
	_, err := ReadKeyFrom(strings.NewReader(""))
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}
