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

	assert.Equal(t, d["c1"], true)
}
