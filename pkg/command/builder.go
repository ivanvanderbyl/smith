package command

import (
	"bytes"
	"fmt"
	"strings"
)

type (
	PrompBuilder interface {
		Prompt() string
	}

	Builder struct {
		commands []*Command
	}

	Command struct {
		Name  string
		Label string
		Args  []Arg
	}

	Arg struct {
		Name            string
		ValueDescriptor string
	}
)

func NewBuilder() *Builder {
	return &Builder{
		commands: []*Command{},
	}
}

var _ PrompBuilder = (*Builder)(nil)

// Prompt returns the prompt string for the command builder.
func (b *Builder) Prompt() string {
	return strings.Join(b.Strings(), "\n") + "\n"
}

func (b *Builder) Strings() []string {
	result := []string{}
	for _, cmd := range b.commands {
		buf := bytes.NewBuffer(nil)
		buf.WriteString(fmt.Sprintf("%s: %q, args: ", cmd.Label, cmd.Name))
		if len(cmd.Args) == 0 {
			buf.WriteString("none")
		}
		for i, arg := range cmd.Args {
			buf.WriteString(fmt.Sprintf("%q: %q", arg.Name, arg.ValueDescriptor))
			if i < len(cmd.Args)-1 {
				buf.WriteString(", ")
			}
		}
		result = append(result, buf.String())
	}
	return result
}

func (b *Builder) NumberedList() string {
	buf := bytes.NewBuffer(nil)

	for i, cmd := range b.Strings() {
		buf.WriteString(fmt.Sprintf("%d. %s", i+1, cmd))
		buf.WriteString("\n")
	}

	return buf.String()
}

// AddCommand adds a new command to the builder.
func (b *Builder) AddCommand(name, label string, args []Arg) error {
	if b.GetCommand(name) != nil {
		return fmt.Errorf("command is already registred")
	}

	b.commands = append(b.commands, &Command{
		Name:  name,
		Label: label,
		Args:  args,
	})

	return nil
}

func (b *Builder) GetCommand(name string) *Command {
	for _, cmd := range b.commands {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}
