# Smith

Smith is a Go package for building AI assisted agents, capable of understanding natural language and executing commands. The primary backend is OpenAI, however, other backends are planned.

## Installation

```bash
go get github.com/ivanvanderbyl/smith
```

## Usage

```go
package main

import "github.com/ivanvanderbyl/smith/pkg/command"

func main() {
  b := command.NewBuilder()
  b.AddCommand("search_github", "Search GitHub for the given query", []command.Arg{{Name: "query", ValueDescriptor: "The search string"}, {Name: "lang", ValueDescriptor: "Restrict to a specific language"}})
  b.AddCommand("search_google", "Search Google for the given query", []command.Arg{{Name: "query", ValueDescriptor: "Search term"}})

  println(b.Prompt())
}
```

```stdout
Commands:
1. Search GitHub for the given query: "search_github", args: "query": "The search string", "lang": "Restrict to a specific language"
2. Search Google for the given query: "search_google", args: "query": "Search term"
You should only respond in JSON format as described below
{"command":{"name":"command name","args":[{"name":"arg name","value":"arg value"}]}}
```
