package config

import (
	env "github.com/Terminus-Lab/stamper/internal/utils"
)

type StamperConfig struct {
	OpenAIKey           string
	AzureOpenAIEndpoint string
	OllamaBaseURL       string
	LLMFamily           string
	ModelId             string
	ModelConfig         ModelConfig
}

type ModelConfig struct {
	MaxToken    int
	Temperature float64
}

func LoadConfig() *StamperConfig {
	return &StamperConfig{
		OpenAIKey:           env.GetString("OPEN_AI_KEY", ""),
		AzureOpenAIEndpoint: env.GetString("AZURE_OPENAI_ENDPOINT", ""),
		OllamaBaseURL:       env.GetString("OLLAMA_BASE_URL", "http://localhost:11434/v1"),
		LLMFamily:           env.GetString("LLM_FAMILY", "openai_platform"),
		ModelId:             env.GetString("MODEL_ID", ""),
		ModelConfig: ModelConfig{
			MaxToken:    1000,
			Temperature: 0.0,
		},
	}
}
