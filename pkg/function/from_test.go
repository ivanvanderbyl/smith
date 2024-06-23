package function_test

import (
	"fmt"
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/ivanvanderbyl/smith/pkg/function"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

type GetWeatherParams struct {
	City string `json:"city" jsonschema:"required,description=The city to get the weather for"`
}

func GetWeather(input GetWeatherParams) (string, error) {
	return fmt.Sprintf("The weather in %q is sunny and 25ºC", input.City), nil
}

func TestFrom(t *testing.T) {
	a := assert.New(t)
	fn, err := function.From(GetWeather)
	a.Nil(err)
	a.NotNil(fn)
	a.Equal("GetWeather", fn.Name)
	a.NotNil(fn.Schema)
	tool := fn.AsTool()
	a.NotNil(tool)
	a.Equal("GetWeather", tool.Function.Name)
	a.NotNil(tool.Function.Parameters)
	params, ok := tool.Function.Parameters.(*jsonschema.Schema)
	a.True(ok)
	city, ok := params.Properties.Get("city")
	a.True(ok)
	a.Equal("The city to get the weather for", city.Description)
}

func TestCanCall(t *testing.T) {
	a := assert.New(t)
	fn, err := function.From(GetWeather)
	a.Nil(err)
	a.NotNil(fn)
	a.True(fn.IsCallable(openai.ToolCall{ID: "call_1", Type: openai.ToolTypeFunction, Function: openai.FunctionCall{Name: "GetWeather", Arguments: "{\"City\":\"Sydney\"}"}}))
	a.False(fn.IsCallable(openai.ToolCall{ID: "call_1", Type: openai.ToolTypeFunction, Function: openai.FunctionCall{Name: "Search", Arguments: "{\"City\":\"Sydney\"}"}}))
}

func TestCall(t *testing.T) {
	a := assert.New(t)
	fn, err := function.From(GetWeather)
	a.Nil(err)
	a.NotNil(fn)
	result, err := fn.Call(openai.ToolCall{ID: "call_1", Type: openai.ToolTypeFunction, Function: openai.FunctionCall{Name: "GetWeather", Arguments: "{\"City\":\"Sydney\"}"}})
	a.Nil(err)
	a.Equal(openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    "The weather in \"Sydney\" is sunny and 25ºC",
		Name:       "GetWeather",
		ToolCallID: "call_1",
	}, result)
}
