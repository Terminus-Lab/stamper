package llm

type LLMFamily string

const (
	FamilyOpenAI         LLMFamily = "openai"
	FamilyOpenAIPlatform LLMFamily = "openai_platform"
	FamilyOllama         LLMFamily = "ollama"
)

type LLMRequest struct {
	Prompt      string
	MaxTokens   int
	Temperature float64
}

type LLMResponse struct {
	Content    string
	StopReason string
}
