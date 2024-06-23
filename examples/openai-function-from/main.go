package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ivanvanderbyl/smith/pkg/function"
	"github.com/ivanvanderbyl/smith/pkg/utils/httpdebugutil"
	"github.com/sashabaranov/go-openai"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatalln(err)
	}
}

type SearchParams struct {
	Query     string   `json:"query" jsonschema:"required,description=The search string"`
	Languages []string `json:"languages" jsonschema:"description=Restrict to a specific language"`
}

func SearchGitHub(params SearchParams) error {
	return nil
}

func realMain() error {
	fn, err := function.From(SearchGitHub)
	if err != nil {
		return err
	}

	apiKey := os.Getenv("OPENAI_API_API")

	debugClient := &http.Client{
		Transport: &httpdebugutil.Transport{},
	}

	ctx := context.Background()
	cfg := openai.DefaultConfig(apiKey)
	cfg.HTTPClient = debugClient
	client := openai.NewClientWithConfig(cfg)

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Search GitHub for the query 'openai' in Typescript.",
			},
		},
		Tools: []openai.Tool{
			fn.AsTool(),
		},
	})
	if err != nil {
		return err
	}

	log.Println(resp.Choices[0].Message.Content)
	for _, tool := range resp.Choices[0].Message.ToolCalls {
		log.Println(tool.Function.Name, tool.Function.Arguments)
	}

	return nil
}
