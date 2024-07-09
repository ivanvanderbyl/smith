package streaming_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/ivanvanderbyl/smith/pkg/streaming"
	"github.com/stretchr/testify/assert"
)

type WeatherArgs struct {
	Location string `json:"location"`
	Units    string `json:"units"`
}

type AgentEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type Anot struct {
	Test string `json:"test"`
}

func TestEncoding(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name     string
		part     streaming.StreamPart
		expected string
	}{
		{
			name:     "TextPart",
			part:     streaming.TextPart{Value: "Hello, World!"},
			expected: "0:\"Hello, World!\"\n",
		},
		{
			name:     "FunctionCallPart",
			part:     streaming.FunctionCallPart[string]{Name: "GET_WEATHER", Args: "Sydney, Australia"},
			expected: `1:"{\"function_call\":{\"name\":\"GET_WEATHER\",\"arguments\":\"{\\\"location\\\":\\\"Sydney, Australia\\\",\\\"units\\\":\\\"celsius\\\"}\"}}`,
		},
		{
			name:     "FunctionCallPart with object",
			part:     streaming.FunctionCallPart[WeatherArgs]{Name: "GET_WEATHER", Args: WeatherArgs{Location: "Sydney, Australia", Units: "celsius"}},
			expected: `1:{"function_call":{"name":"GET_WEATHER","arguments":"{\"location\":\"Sydney, Australia\",\"units\":\"celsius\"}"}}`,
		},
		{
			name:     "AnnotationPart",
			part:     streaming.AnnotationPart[Anot]{Value: []Anot{{Test: "value"}}},
			expected: `8:[{"test":"value"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := streaming.Marshal(tt.part)
			a.NoError(err)
			a.Equal(strings.TrimSpace(tt.expected), strings.TrimSpace(string(encoded)))
			a.True(bytes.HasSuffix(encoded, []byte("\n")))
		})
	}
}
