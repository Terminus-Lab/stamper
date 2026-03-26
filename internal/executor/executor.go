package executor

import (
	"context"
	"fmt"
	"strings"

	"github.com/Terminus-Lab/stamper/internal/config"
	"github.com/Terminus-Lab/stamper/internal/domain"
	"github.com/Terminus-Lab/stamper/internal/llm"
	"github.com/rs/zerolog"
)

type Executor struct {
	llmClient llm.LLMClient
	cfg       *config.StamperConfig
	logger    *zerolog.Logger
}

func New(llmClient llm.LLMClient, cfg *config.StamperConfig, logger *zerolog.Logger) *Executor {
	return &Executor{
		llmClient: llmClient,
		cfg:       cfg,
		logger:    logger,
	}
}

func (e *Executor) Run(ctx context.Context, conv domain.Conversation) (string, error) {
	prompt := buildPrompt(conv)

	resp, err := e.llmClient.InvokeModel(ctx, llm.LLMRequest{
		Prompt:      prompt,
		MaxTokens:   e.cfg.ModelConfig.MaxToken,
		Temperature: e.cfg.ModelConfig.Temperature,
	})
	if err != nil {
		e.logger.Error().Err(err).Msg("LLM call failed")
		return "", err
	}

	return resp.Content, nil
}

func buildPrompt(conv domain.Conversation) string {
	var sb strings.Builder
	sb.WriteString("You are helping a human annotator evaluate an AI conversation.\n")
	sb.WriteString("Summarize this conversation in 2-3 sentences.\n")
	sb.WriteString("Focus on: what the user asked, whether the agent's responses were accurate and helpful, and any notable issues.\n\n")

	for i, turn := range conv.Turns {
		fmt.Fprintf(&sb, "Turn %d\n", i+1)
		fmt.Fprintf(&sb, "User: %s\n", turn.Query)
		fmt.Fprintf(&sb, "Agent: %s\n\n", turn.Answer)
	}

	sb.WriteString("Provide only the summary, no preamble.")
	return sb.String()
}
