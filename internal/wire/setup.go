package wire

import (
	"context"
	"fmt"

	"github.com/Terminus-Lab/stamper/internal/config"
	"github.com/Terminus-Lab/stamper/internal/llm"
	"github.com/Terminus-Lab/stamper/internal/llm/azure"
	"github.com/Terminus-Lab/stamper/internal/llm/ollama"
	"github.com/Terminus-Lab/stamper/internal/llm/openaiplatform"
)

func GetLLMClient(ctx context.Context, cfg *config.StamperConfig) (llm.LLMClient, error) {
	family := llm.LLMFamily(cfg.LLMFamily)
	switch family {
	case llm.FamilyOpenAIPlatform:
		if cfg.OpenAIKey == "" {
			return nil, fmt.Errorf("OPEN_AI_KEY is required for LLM_FAMILY=%s (MODEL_ID=%s)", cfg.LLMFamily, cfg.ModelId)
		}
		client, err := openaiplatform.NewClient(ctx, cfg.OpenAIKey, cfg.ModelId)
		if err != nil {
			return nil, fmt.Errorf("failed to create OpenAI Platform client for model %s: %w", cfg.ModelId, err)
		}
		return client, nil
	case llm.FamilyOllama:
		client, err := ollama.NewClient(ctx, cfg.OllamaBaseURL, cfg.ModelId)
		if err != nil {
			return nil, fmt.Errorf("failed to create Ollama client for model %s: %w", cfg.ModelId, err)
		}
		return client, nil
	case llm.FamilyOpenAI:
		if cfg.OpenAIKey == "" {
			return nil, fmt.Errorf("OPEN_AI_KEY required for openai model %s", cfg.ModelId)
		}
		if cfg.AzureOpenAIEndpoint == "" {
			return nil, fmt.Errorf("AZURE_OPENAI_ENDPOINT required for openai model %s", cfg.ModelId)
		}
		client, err := azure.NewClient(cfg.OpenAIKey, cfg.ModelId, cfg.AzureOpenAIEndpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to create Azure OpenAI client for model %s: %w", cfg.ModelId, err)
		}
		return client, err
	default:
		return nil, fmt.Errorf("unsupported model family: %s", cfg.LLMFamily)
	}
}
