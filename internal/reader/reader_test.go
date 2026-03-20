package reader

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestReaderLoad(t *testing.T) {
	f, err := os.CreateTemp("", "samper-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(f.Name())

	content := `{"conversation_id":"c1","turns":[{"query":"Hi","answer":"Hello"}]}
{"conversation_id":"c2","turns":[{"query":"Bye","answer":"Goodbye"}]}`
	f.WriteString(content)
	f.Close()

	logger := zerolog.Nop()
	reader := NewReader(&logger)

	convs, err := reader.Load(f.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Assert
	assert.Len(t, convs, 2)
	assert.Equal(t, "c1", convs[0].ConversationID)
	assert.Equal(t, "Hi", convs[0].Turns[0].Query)
}

func TestReaderLoad_InvalidJson(t *testing.T) {
	f, _ := os.CreateTemp("", "stamper-*.json")
	defer os.Remove(f.Name())

	f.WriteString("not a valid json")
	f.Close()

	logger := zerolog.Nop()
	reader := NewReader(&logger)

	_, err := reader.Load(f.Name())
	if err == nil {
		t.Fatalf("Expected error from invalid json. Got nil")
	}
}

func TestReaderLoad_FileNotFound(t *testing.T) {
	logger := zerolog.Nop()
	reader := NewReader(&logger)

	_, err := reader.Load("inexistentFile.txt")

	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
