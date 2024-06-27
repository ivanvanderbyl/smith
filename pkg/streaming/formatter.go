// streaming implements a stream encoder that is compatible with the Vercel @ai-sdk.
package streaming

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type (
	StreamPart interface {
		isStreamPart()
		Code() Code
		fmt.Stringer
	}

	Code int8

	TextPart struct {
		Value string
	}

	FunctionCallPart[T any] struct {
		Name string
		Args T
	}
)

func (TextPart) isStreamPart() {}

func (TextPart) Code() Code {
	return TextPartCode
}

func (t TextPart) String() string {
	return t.Value
}

func (FunctionCallPart[T]) isStreamPart() {}

func (FunctionCallPart[T]) Code() Code {
	return FunctionCallPartCode
}

type functionCall[T any] struct {
	Name string `json:"name"`
	Args T      `json:"args"`
}

type functionCallWrapper[T any] struct {
	Fn functionCall[T] `json:"function_call"`
}

func (f FunctionCallPart[T]) String() string {
	v := functionCallWrapper[T]{
		Fn: functionCall[T]{
			Name: f.Name,
			Args: f.Args,
		},
	}
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(v)
	return buf.String()
}

const (
	TextPartCode Code = iota
	FunctionCallPartCode
	DataPartCode
	ErrorPartCode
	AssistantMessagePartCode
	AssistantControlDataPartCode
	DataMessagePartCode
	ToolCallsPartCode
	MessageAnnotationsPartCode
	ToolCallPartCode
	ToolResultPartCode
)

// Marshal formats a stream part into a string
func Marshal(part StreamPart) ([]byte, error) {
	out := bytes.NewBuffer(nil)

	value := part.String()
	if strings.HasSuffix(value, "\n") {
		value = strings.TrimSuffix(value, "\n")
	}

	_, err := fmt.Fprintf(out, "%d:%s\n", part.Code(), value)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
