package streaming_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ivanvanderbyl/smith/pkg/streaming"
	"github.com/stretchr/testify/assert"
)

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
			expected: "0:Hello, World!\n",
		},
		{
			name:     "FunctionCallPart",
			part:     streaming.FunctionCallPart[string]{Name: "GET_WEATHER", Args: "Sydney, Australia"},
			expected: `1:{"function_call":{"name":"GET_WEATHER","args":"Sydney, Australia"}}`,
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
