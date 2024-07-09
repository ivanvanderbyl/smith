package llm

type (
	LLM interface {
		CreateChatCompletion(input ChatCompletionParams) (ChatResponse, error)
		// CreateCompletionStream(input CompletionRequest) (CompletionStream, error)
	}

	ChatCompletionParams struct {
	}

	ChatResponse struct{}
)
