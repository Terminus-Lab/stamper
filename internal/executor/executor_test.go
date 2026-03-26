package executor

import (
	"strings"
	"testing"

	"github.com/Terminus-Lab/stamper/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestBuildPrompt_ContainsInstructions(t *testing.T) {
	conv := domain.Conversation{
		ConversationID: "c1",
		Turns: []domain.Turn{
			{Query: "What is Python?", Answer: "A high-level programming language."},
		},
	}

	prompt := buildPrompt(conv)

	assert.Contains(t, prompt, "human annotator")
	assert.Contains(t, prompt, "2-3 sentences")
	assert.Contains(t, prompt, "no preamble")
}

func TestBuildPrompt_ContainsTurnContent(t *testing.T) {
	conv := domain.Conversation{
		ConversationID: "c1",
		Turns: []domain.Turn{
			{Query: "What is Python?", Answer: "A high-level programming language."},
			{Query: "Is it hard?", Answer: "Not at all."},
		},
	}

	prompt := buildPrompt(conv)

	assert.Contains(t, prompt, "What is Python?")
	assert.Contains(t, prompt, "A high-level programming language.")
	assert.Contains(t, prompt, "Is it hard?")
	assert.Contains(t, prompt, "Not at all.")
}

func TestBuildPrompt_TurnNumbering(t *testing.T) {
	conv := domain.Conversation{
		Turns: []domain.Turn{
			{Query: "Q1", Answer: "A1"},
			{Query: "Q2", Answer: "A2"},
			{Query: "Q3", Answer: "A3"},
		},
	}

	prompt := buildPrompt(conv)

	assert.Contains(t, prompt, "Turn 1")
	assert.Contains(t, prompt, "Turn 2")
	assert.Contains(t, prompt, "Turn 3")
}

func TestBuildPrompt_EmptyTurns(t *testing.T) {
	conv := domain.Conversation{
		ConversationID: "empty",
		Turns:          []domain.Turn{},
	}

	prompt := buildPrompt(conv)

	assert.NotEmpty(t, prompt)
	assert.False(t, strings.Contains(prompt, "Turn 1"))
}
