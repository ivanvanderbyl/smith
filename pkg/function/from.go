package function

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/sashabaranov/go-openai"
)

type FunctionCall[Args any] struct {
	Name        string
	Description string
	Arguments   reflect.Type
	Schema      *jsonschema.Schema
}

type FuncErr[Args any] func(Args) error

func From[I any](input FuncErr[I]) (*FunctionCall[I], error) {
	fn := reflect.ValueOf(input)

	args := fn.Type().In(0)

	r := &jsonschema.Reflector{
		Anonymous:      true,
		DoNotReference: true,
	}
	schema := r.ReflectFromType(args)

	return &FunctionCall[I]{
		Name:      getFunctionName(input),
		Schema:    schema,
		Arguments: args,
	}, nil
}

// AsTool converts the function call into an OpenAI tool definition
func (fc *FunctionCall[I]) AsTool() openai.Tool {
	return openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: fc.AsFunctionDefinition(),
	}
}

func (fc *FunctionCall[I]) AsFunctionDefinition() *openai.FunctionDefinition {
	return &openai.FunctionDefinition{
		Name:        fc.Name,
		Description: fc.Description,
		Parameters:  fc.Schema,
	}
}

func getFunctionName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}
