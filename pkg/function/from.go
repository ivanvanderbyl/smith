package function

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/sashabaranov/go-openai"
)

type AnyFunction interface {
	CanCall(tc openai.ToolCall) bool
	Call(tc openai.ToolCall) (openai.ChatCompletionMessage, error)
	AsTool() openai.Tool
}

type FunctionCall[Args any, Result ~string] struct {
	Name        string
	Description string
	Arguments   reflect.Type
	Schema      *jsonschema.Schema
	calls       map[string]*Result
	impl        FuncErr[Args, Result]
}

type FuncErr[Args any, Result ~string] func(Args) (Result, error)

func From[I any, O ~string](input FuncErr[I, O]) (*FunctionCall[I, O], error) {
	fn := reflect.ValueOf(input)

	args := fn.Type().In(0)

	if args.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected function input to be a struct, got %v", args)
	}

	r := &jsonschema.Reflector{
		Anonymous:      true,
		DoNotReference: true,
	}
	schema := r.ReflectFromType(args)

	return &FunctionCall[I, O]{
		Name:      getFunctionName(input),
		Schema:    schema,
		Arguments: args,
		impl:      input,
	}, nil
}

// AsTool converts the function call into an OpenAI tool definition
func (fc *FunctionCall[I, O]) AsTool() openai.Tool {
	return openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: fc.AsFunctionDefinition(),
	}
}

func (fc *FunctionCall[I, O]) AsFunctionDefinition() *openai.FunctionDefinition {
	return &openai.FunctionDefinition{
		Name:        fc.Name,
		Description: fc.Description,
		Parameters:  fc.Schema,
	}
}

func (fc *FunctionCall[I, O]) CanCall(tc openai.ToolCall) bool {
	return tc.Function.Name == fc.Name
}

func (fc *FunctionCall[I, O]) Call(tc openai.ToolCall) (openai.ChatCompletionMessage, error) {
	nextMessage := openai.ChatCompletionMessage{}

	if !fc.CanCall(tc) {
		return nextMessage, fmt.Errorf("unexpected function name %q, expected %q", tc.Function.Name, fc.Name)
	}

	args := new(I)
	err := json.Unmarshal([]byte(tc.Function.Arguments), args)
	if err != nil {
		return nextMessage, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	result, err := fc.impl(*args)
	if err != nil {
		return nextMessage, err
	}

	if fc.calls == nil {
		fc.calls = make(map[string]*O)
	}
	fc.calls[tc.ID] = &result

	nextMessage = openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    string(result),
		Name:       tc.Function.Name,
		ToolCallID: tc.ID,
	}

	return nextMessage, nil
}

func getFunctionName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}
