# Smith

Smith is a Go package for building agents that are capable of understanding natural language and executing commands.

This is not a framework, rather it is a set of primitives that can be used to build a custom agent that understands a specific set of commands that can be used within a larger project.

## Motivation

After building several AI projects using LangChain and LlamaIndex, I found existing frameworks to be overly opinionated and abstracted. They often obscured important details, making debugging and understanding the underlying processes challenging.

Most modern AI APIs support similar feature sets with comparable semantics. This project aims to provide a set of primitives that simplify the implementation of interactive agents. These primitives are designed to be easily interchangeable, allowing for seamless integration of alternative implementations at any stage of development.

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
