package llm

import "context"

type LLMClient interface {
	InvokeModel(ctx context.Context, request LLMRequest) (*LLMResponse, error)
	InvokeModelWithRetry(ctx context.Context, request LLMRequest) (*LLMResponse, error)
}
