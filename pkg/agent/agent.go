package agent

import (
	"github.com/ivanvanderbyl/smith/pkg/function"
	"github.com/sashabaranov/go-openai"
)

const MAX_TOOL_CALLS = 10

type (
	Agent struct {
		Tasks        []Task
		SystemPrompt string
		llm          *openai.Client
		Tools        function.AnyFunctions
	}

	AgentOption func(*Agent)

	Task struct {
		History       []ChatMessage
		ToolCallCount int
	}

	Step struct {
		ID           string
		PreviousStep *Step
		NextSteps    []*Step
	}

	ChatMessage struct {
		Content string
		Role    Role
	}

	Role string
)

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
	RoleMemory    Role = "memory"
)

// ShouldContinue returns true if the task should continue to run.
func (t *Task) ShouldContinue() bool {
	return t.ToolCallCount < MAX_TOOL_CALLS
}

func WithSystemPrompt(prompt string) AgentOption {
	return func(a *Agent) {
		a.SystemPrompt = prompt
	}
}

func (a *Agent) CreateTask() *Task {
	initialMessages := []ChatMessage{
		{
			Content: a.SystemPrompt,
			Role:    RoleSystem,
		},

		// Include previous history
	}

	return &Task{
		History:       initialMessages,
		ToolCallCount: 0,
	}
}

func (s *Step) IsLast() bool {
	return len(s.NextSteps) == 0
}
