package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ivanvanderbyl/smith/pkg/function"
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

type Result struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
}

func (r Result) String() string {
	return fmt.Sprintf("%s (%s): %s", r.Name, r.Language, r.Description)
}

func SearchGitHub(params SearchParams) (string, error) {
	log.Println("> Searching GitHub for", params.Query)
	return Result{
		Name:        "openai",
		Description: "OpenAI API client for Go",
		Language:    "ts",
	}.String(), nil
}

type GetWeatherParams struct {
	City string `json:"city" jsonschema:"required,description=The city to get the weather for"`
	Unit string `json:"unit" jsonschema:"description=The unit of temperature to return,enum=celsius,fahrenheit"`
}

func GetWeather(input GetWeatherParams) (string, error) {
	if input.Unit == "" {
		input.Unit = "celsius"
	}

	return fmt.Sprintf("The weather in %q is sunny and 25º %s", input.City, input.Unit), nil
}

func realMain() error {
	apiKey := os.Getenv("OPENAI_API_API")

	debugClient := &http.Client{
		// Transport: &httpdebugutil.Transport{},
	}

	// Setup tools
	searchGithub, err := function.From(SearchGitHub)
	if err != nil {
		return err
	}

	getWeather, err := function.From(GetWeather)
	if err != nil {
		return err
	}

	tools := function.AnyFunctions{
		searchGithub,
		getWeather,
	}

	ctx := context.Background()
	cfg := openai.DefaultConfig(apiKey)
	cfg.HTTPClient = debugClient
	client := openai.NewClientWithConfig(cfg)

	dialogue := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleUser, Content: "What is the weather in Sydney today?"},
	}

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: dialogue,
		Tools:    tools.AsTools(),
	})
	if err != nil {
		return err
	}

	displayMessages(resp.Choices)

	// Append OpenAI's response to the dialogue so that our tools can reply to it
	msg := resp.Choices[0].Message
	dialogue = append(dialogue, msg)
	for _, toolCall := range msg.ToolCalls {
		for _, tool := range tools {
			if !tool.CanCall(toolCall) {
				continue
			}

			nextMessage, err := tool.Call(toolCall)
			if err != nil {
				return err
			}

			dialogue = append(dialogue, nextMessage)
		}
	}

	resp, err = client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    "gpt-4o",
		Messages: dialogue,
		Tools:    tools.AsTools(),
	})
	if err != nil {
		return err
	}
	displayMessages(resp.Choices)

	return nil
}

func displayDialog(dialogue []openai.ChatCompletionMessage) {
	for _, message := range dialogue {
		switch message.Role {
		case openai.ChatMessageRoleUser:
			fmt.Printf("> User: %q\n", message.Content)
		case openai.ChatMessageRoleAssistant:
			fmt.Printf("> Assistant: %q\n", message.Content)
		case openai.ChatMessageRoleSystem:
			fmt.Printf("> System: %q\n", message.Content)
		case openai.ChatMessageRoleTool:
			fmt.Printf("> Tool (%s): %q\n", message.Name, message.Content)
		}
	}
}

func displayMessages(choices []openai.ChatCompletionChoice) {
	for _, choice := range choices {
		switch choice.Message.Role {
		case openai.ChatMessageRoleUser:
			fmt.Printf("> User: %q\n", choice.Message.Content)
		case openai.ChatMessageRoleAssistant:
			if choice.Message.Content != "" {
				fmt.Printf("> Assistant: %+v\n", choice.Message.Content)
			}
			displayToolCall(choice)
		case openai.ChatMessageRoleSystem:
			fmt.Printf("> System: %q\n", choice.Message.Content)
		}

	}
}

func displayToolCall(choice openai.ChatCompletionChoice) {
	for _, toolCall := range choice.Message.ToolCalls {
		fmt.Printf("> Function Call (%s): %s(%+v)\n", toolCall.ID, toolCall.Function.Name, toolCall.Function.Arguments)
	}
}

func expectToolCall(choices []openai.ChatCompletionChoice, functionName string) error {
	didCall := false
	for _, choice := range choices {
		for _, toolCall := range choice.Message.ToolCalls {
			if toolCall.Function.Name == functionName {
				didCall = true
				break
			}
		}
	}

	if !didCall {
		return fmt.Errorf("expected function call %q, but did not receive it", functionName)
	}

	return nil
}
