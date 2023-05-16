package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromptGenerationForCommands(t *testing.T) {
	a := assert.New(t)

	b := NewBuilder()
	b.AddCommand("cmd1", "Command 1", []Arg{{Name: "start_date", ValueDescriptor: "Start date of search query"}, {Name: "arg2", ValueDescriptor: "arg2"}})
	b.AddCommand("cmd2", "Command 2", []Arg{{Name: "arg1", ValueDescriptor: "arg1"}, {Name: "arg2", ValueDescriptor: "arg2"}})
	b.AddCommand("cmd3", "Empty args cmd", []Arg{})

	result := b.Prompt()
	a.Equal(`Command 1: "cmd1", args: "start_date": "Start date of search query", "arg2": "arg2"
Command 2: "cmd2", args: "arg1": "arg1", "arg2": "arg2"
Empty args cmd: "cmd3", args: none
`, result)

	result = b.NumberedList()
	a.Equal(`1. Command 1: "cmd1", args: "start_date": "Start date of search query", "arg2": "arg2"
2. Command 2: "cmd2", args: "arg1": "arg1", "arg2": "arg2"
3. Empty args cmd: "cmd3", args: none
`, result)
}

func TestCommandRegistration(t *testing.T) {
	a := assert.New(t)

	b := NewBuilder()
	a.NoError(b.AddCommand("cmd1", "Command 1", nil), "command should be registered")
	a.Error(b.AddCommand("cmd1", "Command 1", nil), "command should be registered")
}
