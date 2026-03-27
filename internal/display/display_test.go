package display

import (
	"bytes"
	"testing"

	"github.com/Terminus-Lab/stamper/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	conv := domain.Conversation{
		ConversationID: "c1",
		Turns: []domain.Turn{
			{UserQuery: "Hi", Answer: "Hello"},
		},
	}

	var buf bytes.Buffer
	Reader(&buf, conv, 1, 1)
	output := buf.String()
	assert.Contains(t, output, "Conversation 1 / 1  - 1 turns")
	assert.Contains(t, output, "Turn 1")
	assert.Contains(t, output, "User: Hi")
	assert.Contains(t, output, "Answer: Hello")
	assert.Contains(t, output, "[p] pass   [r] review   [f] fail   [s]summarize   [x] skip")
}
