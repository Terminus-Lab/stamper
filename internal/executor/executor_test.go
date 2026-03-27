package executor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Terminus-Lab/stamper/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// repoPromptPath resolves conf/summarize_prompt.tmpl from the executor package test cwd.
func repoPromptPath(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "..", "conf", "summarize_prompt.tmpl")
}

func render(t *testing.T, conv domain.Conversation) string {
	t.Helper()
	tmpl, err := loadTemplate(repoPromptPath(t))
	require.NoError(t, err)
	result, err := renderPrompt(tmpl, conv)
	require.NoError(t, err)
	return result
}

func TestPrompt_ContainsInstructions(t *testing.T) {
	conv := domain.Conversation{
		ConversationID: "c1",
		Turns:          []domain.Turn{{Query: "What is Python?", Answer: "A high-level programming language."}},
	}

	prompt := render(t, conv)

	assert.Contains(t, prompt, "human annotator")
	assert.Contains(t, prompt, "2-3 sentences")
	assert.Contains(t, prompt, "no preamble")
}

func TestPrompt_ContainsTurnContent(t *testing.T) {
	conv := domain.Conversation{
		Turns: []domain.Turn{
			{Query: "What is Python?", Answer: "A high-level programming language."},
			{Query: "Is it hard?", Answer: "Not at all."},
		},
	}

	prompt := render(t, conv)

	assert.Contains(t, prompt, "What is Python?")
	assert.Contains(t, prompt, "A high-level programming language.")
	assert.Contains(t, prompt, "Is it hard?")
	assert.Contains(t, prompt, "Not at all.")
}

func TestPrompt_TurnNumbering(t *testing.T) {
	conv := domain.Conversation{
		Turns: []domain.Turn{
			{Query: "Q1", Answer: "A1"},
			{Query: "Q2", Answer: "A2"},
			{Query: "Q3", Answer: "A3"},
		},
	}

	prompt := render(t, conv)

	assert.Contains(t, prompt, "Turn 1")
	assert.Contains(t, prompt, "Turn 2")
	assert.Contains(t, prompt, "Turn 3")
}

func TestPrompt_EmptyTurns(t *testing.T) {
	conv := domain.Conversation{ConversationID: "empty", Turns: []domain.Turn{}}

	prompt := render(t, conv)

	assert.NotEmpty(t, prompt)
	assert.False(t, strings.Contains(prompt, "Turn 1"))
}

func TestLoadTemplate_CustomFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom.tmpl")
	err := os.WriteFile(path, []byte("Custom: {{range .Turns}}{{.Query}}{{end}}"), 0o644)
	require.NoError(t, err)

	tmpl, err := loadTemplate(path)
	require.NoError(t, err)

	conv := domain.Conversation{Turns: []domain.Turn{{Query: "hello", Answer: "world"}}}
	result, err := renderPrompt(tmpl, conv)
	require.NoError(t, err)
	assert.Equal(t, "Custom: hello", result)
}

func TestLoadTemplate_MissingFile(t *testing.T) {
	_, err := loadTemplate("/nonexistent/path/prompt.tmpl")
	assert.Error(t, err)
}

func TestLoadTemplate_EmptyPath(t *testing.T) {
	_, err := loadTemplate("")
	assert.Error(t, err)
}
