package main

import "github.com/ivanvanderbyl/smith/pkg/command"

func main() {
	b := command.NewBuilder()
	b.AddCommand("search_github", "Search GitHub for the given query", []command.Arg{{Name: "query", ValueDescriptor: "The search string"}, {Name: "lang", ValueDescriptor: "Restrict to a specific language"}})
	b.AddCommand("search_google", "Search Google for the given query", []command.Arg{{Name: "query", ValueDescriptor: "Search term"}})

	println(b.Prompt())
}
