# Smith

Smith is a Go package for building agents that are capable of understanding natural language and executing commands.

This is not a framework, rather it is a set of primitives that can be used to build a custom agent that understands a specific set of commands that can be used within a larger project.

## Motivation

After building a few AI projects with langchain and llamaindex, I found that the existing frameworks were too opinionated abstracted away the details that are important, and made it very difficult to debug and understand what was happening.

Almost all AI APIs these days support a similar feature set with similar semantics, so this project aims to provide a set of primitives that make implementing interactive agents trivial and easy to swap out at any stage with an alternative implementation.

## Installation

```bash
go get github.com/ivanvanderbyl/smith
```

## Usage

```go
package main

type SearchParams struct {
	Query     string   `json:"query" jsonschema:"required,description=The search string"`
	Languages []string `json:"languages" jsonschema:"description=Restrict to a specific language"`
}

func SearchGitHub(params SearchParams) error {
  // Do the search
	return nil
}

func main() {
  myFunction := function.From(SearchGitHub)

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
			myFunction.AsTool(),
		},
	})
}
```
