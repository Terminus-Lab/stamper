package executor

import (
	"context"

	"github.com/Terminus-Lab/stamper/internal/config"
	"github.com/Terminus-Lab/stamper/internal/llm"
	"github.com/rs/zerolog"
)

type Executor struct {
	llmClient llm.LLMClient
	logger    *zerolog.Logger
}

func New(llmClient llm.LLMClient, logger *zerolog.Logger) *Executor {
	return &Executor{
		llmClient: llmClient,
		logger:    logger,
	}
}

func (e *Executor) Run(ctx context.Context, cfg config.StamperConfig) (*llm.LLMResponse, error) {
	prompt := e.buildPrompt()

	resp, err := e.llmClient.InvokeModel(ctx, llm.LLMRequest{
		Prompt:      prompt,
		MaxTokens:   cfg.ModelConfig.MaxToken,
		Temperature: cfg.ModelConfig.Temperature,
	})

	if err != nil {
		e.logger.Error().
			Err(err).
			Msg("LLM call failed")
	}

	return &llm.LLMResponse{
		Content:    resp.Content,
		StopReason: resp.StopReason,
	}, nil

}

func (j *Executor) buildPrompt() string {
	//TODO: Add prompt
	return "sdasda"
}
