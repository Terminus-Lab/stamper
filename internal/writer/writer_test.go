package writer

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/Terminus-Lab/stamper/internal/domain"
)

func TestAppend(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "out*.jsonl")
	path := f.Name()
	f.Close()

	w, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	conv := domain.Conversation{
		ConversationID: "conv-1",
		Turns: []domain.Turn{
			{UserQuery: "hello", Answer: "world"},
		},
	}

	if err := w.Append(conv, "pass", ""); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var got domain.Conversation
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	if got.ConversationID != "conv-1" {
		t.Errorf("conversation_id: got %q, want %q", got.ConversationID, "conv-1")
	}
	if got.Annotation != "pass" {
		t.Errorf("human_annotation: got %q, want %q", got.Annotation, "pass")
	}
	if len(got.Turns) != 1 || got.Turns[0].UserQuery != "hello" {
		t.Errorf("turns: unexpected value %+v", got.Turns)
	}
}

func TestAppendPreservesExtraFields(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "out*.jsonl")
	path := f.Name()
	f.Close()

	w, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	raw := `{"conversation_id":"conv-2","turns":[],"source":"internal","score":42}`
	var conv domain.Conversation
	if err := json.Unmarshal([]byte(raw), &conv); err != nil {
		t.Fatalf("unmarshal input: %v", err)
	}

	if err := w.Append(conv, "fail", ""); err != nil {
		t.Fatalf("Append: %v", err)
	}
	w.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var out map[string]json.RawMessage
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	for _, key := range []string{"source", "score", "human_annotation"} {
		if _, ok := out[key]; !ok {
			t.Errorf("missing field %q in output", key)
		}
	}

	var annotation string
	json.Unmarshal(out["human_annotation"], &annotation)
	if annotation != "fail" {
		t.Errorf("human_annotation: got %q, want %q", annotation, "fail")
	}
}

func TestAppendMultiple(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "out*.jsonl")
	path := f.Name()
	f.Close()

	w, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	convs := []domain.Conversation{
		{ConversationID: "a", Turns: []domain.Turn{{UserQuery: "q1", Answer: "a1"}}},
		{ConversationID: "b", Turns: []domain.Turn{{UserQuery: "q2", Answer: "a2"}}},
	}
	annotations := []string{"pass", "review"}

	for i, c := range convs {
		if err := w.Append(c, annotations[i], ""); err != nil {
			t.Fatalf("Append %d: %v", i, err)
		}
	}
	w.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	lines := splitLines(data)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	for i, line := range lines {
		var got domain.Conversation
		if err := json.Unmarshal([]byte(line), &got); err != nil {
			t.Fatalf("line %d unmarshal: %v", i, err)
		}
		if got.ConversationID != convs[i].ConversationID {
			t.Errorf("line %d: id got %q want %q", i, got.ConversationID, convs[i].ConversationID)
		}
		if got.Annotation != annotations[i] {
			t.Errorf("line %d: annotation got %q want %q", i, got.Annotation, annotations[i])
		}
	}
}

func splitLines(data []byte) []string {
	var lines []string
	start := 0
	for i, b := range data {
		if b == '\n' {
			if i > start {
				lines = append(lines, string(data[start:i]))
			}
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, string(data[start:]))
	}
	return lines
}
