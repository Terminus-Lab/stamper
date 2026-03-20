package reader

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestResumeLoad(t *testing.T) {
	f, err := os.CreateTemp("", "stamper-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(f.Name())
	content := `{"conversation_id":"c1","turns":[{"query":"Hi","answer":"Hello"}]}
{"conversation_id":"c2","turns":[{"query":"Bye","answer":"Goodbye"}]}`

	f.WriteString(content)
	f.Close()

	logger := zerolog.Nop()
	resume := NewResume(&logger)
	d, err := resume.Load(f.Name())

	if err != nil {
		t.Fatal("Unable to run Reader Load function")
	}

	assert.Len(t, d, 2)
	assert.Equal(t, d["c1"], true)
	assert.Equal(t, d["c2"], true)

}

func TestResumeLoad_InvalidJson(t *testing.T) {
	f, err := os.CreateTemp("", "stamper-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(f.Name())
	f.WriteString("invalid json")
	f.Close()

	logger := zerolog.Nop()
	resume := NewResume(&logger)

	_, err = resume.Load(f.Name())

	if err == nil {
		t.Fatal("Expected error for invalid json.")
	}
}

func TestReaderLoad_FileNotExist(t *testing.T) {
	logger := zerolog.Nop()
	reader := NewResume(&logger)

	d, err := reader.Load("inexistentFile.txt")

	assert.NoError(t, err)
	assert.Empty(t, d)
}
