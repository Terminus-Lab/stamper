package executor

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"os"
	"text/template"

	"github.com/Terminus-Lab/stamper/internal/config"
	"github.com/Terminus-Lab/stamper/internal/domain"
	"github.com/Terminus-Lab/stamper/internal/llm"
	"github.com/rs/zerolog"
)

//go:embed default_prompt.tmpl
var defaultPromptFS embed.FS

type Executor struct {
	llmClient llm.LLMClient
	cfg       *config.StamperConfig
	logger    *zerolog.Logger
	tmpl      *template.Template
}

func New(llmClient llm.LLMClient, cfg *config.StamperConfig, logger *zerolog.Logger) (*Executor, error) {
	tmpl, err := loadTemplate(cfg.PromptFile)
	if err != nil {
		return nil, err
	}
	return &Executor{
		llmClient: llmClient,
		cfg:       cfg,
		logger:    logger,
		tmpl:      tmpl,
	}, nil
}

func (e *Executor) Run(ctx context.Context, conv domain.Conversation) (string, error) {
	prompt, err := renderPrompt(e.tmpl, conv)
	if err != nil {
		return "", fmt.Errorf("render prompt: %w", err)
	}

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

func loadTemplate(path string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}

	if path != "" {
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read prompt file %q: %w", path, err)
		}
		return template.New("prompt").Funcs(funcMap).Parse(string(raw))
	}

	raw, err := defaultPromptFS.ReadFile("default_prompt.tmpl")
	if err != nil {
		return nil, fmt.Errorf("read embedded prompt: %w", err)
	}
	return template.New("prompt").Funcs(funcMap).Parse(string(raw))
}

func renderPrompt(tmpl *template.Template, conv domain.Conversation) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, conv); err != nil {
		return "", err
	}
	return buf.String(), nil
}
