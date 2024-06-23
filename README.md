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

We implement helper methods to handle the boilerplate code required to use Tools with OpenAI's Chat API.

This avoids needing to manually specify the function schema because we can infer this directly from the function signature at runtime.

```go
package main

type SearchParams struct {
  Query     string   `json:"query" jsonschema:"required,description=The search string"`
  Languages []string `json:"languages" jsonschema:"description=Restrict to a specific language"`
}

func SearchGitHub(params SearchParams) (string, error) {
  // Do the search
  return "Found 'openai-ts'" nil
}

func main() {
  search := function.From(SearchGitHub)

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
      search.AsTool(),
    },
  })
}
```

### `function.From`

`function.From` is a helper function that creates a `function.FunctionCall` from a standard Go `func`. It must follow the signature `func (params T) (string, error)` where `T` is a struct that represents the input parameters, and returns a string and an error.

The string output will be sent as a reply to the assistant to continue the conversation.

```go
search := function.From(SearchGitHub)
```

### `FunctionCall.AsTool`

`FunctionCall.AsTool` is a method that converts a `FunctionCall` into a `Tool` that can be used in the `ChatCompletionRequest`.

This will include the function schema in the the request allowing the LLM to request a call to this function if it is the best tool to handle the user's request.

```go
openai.ChatCompletionRequest{
  Tools: []openai.Tool{
    search.AsTool(),
  },
}
```

### `FunctionCall.Call`

A helper that can be used to call the function from an LLM tool call response, and return a new conversation completion message in order to continue the conversation.

```go
func SearchGitHub(params SearchParams) (string, error) {
  // Do the search
  return "Found 'openai-ts'" nil
}

search := function.From(SearchGitHub)

func main() {
  tools := function.AnyFunctions{
    search,
  }
  dialogue := []openai.ChatCompletionMessage{
    {
      Role:    openai.ChatMessageRoleUser,
      Content: "Search GitHub for the query 'openai' in Typescript.",
    },
  }
  resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
    Model: "gpt-4o",
    Messages: dialogue,
    Tools: tools.AsTools(),
  })

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

  // Continue the conversation
}
```
