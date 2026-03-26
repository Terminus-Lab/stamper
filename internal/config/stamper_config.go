package config

import (
	env "github.com/Terminus-Lab/stamper/internal/utils"
	"github.com/joho/godotenv"
)

type StamperConfig struct {
	OpenAIKey           string
	AzureOpenAIEndpoint string
	OllamaBaseURL       string
	LLMFamily           string
	ModelId             string
	ModelConfig         ModelConfig
	SummarizeEnabled    bool
	PromptFile          string
}

type ModelConfig struct {
	MaxToken    int
	Temperature float64
}

func LoadConfig() *StamperConfig {
	// .env is optional — exported env vars are always respected
	_ = godotenv.Load()

	return &StamperConfig{
		OpenAIKey:           env.GetString("OPEN_AI_KEY", ""),
		AzureOpenAIEndpoint: env.GetString("AZURE_OPENAI_ENDPOINT", ""),
		OllamaBaseURL:       env.GetString("OLLAMA_BASE_URL", "http://localhost:11434/v1"),
		LLMFamily:           env.GetString("LLM_FAMILY", "openai_platform"),
		ModelId:             env.GetString("MODEL_ID", ""),
		ModelConfig: ModelConfig{
			MaxToken:    env.GetInt("MODEL_MAX_TOKENS", 1000),
			Temperature: env.GetFloat("MODEL_TEMPERATURE", 0.0),
		},
		SummarizeEnabled: env.GetBool("STAMPER_SUMMARIZE", false),
	}
}
