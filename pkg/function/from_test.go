package function_test

import (
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/ivanvanderbyl/smith/pkg/function"
	"github.com/stretchr/testify/assert"
)

type GetWeatherParams struct {
	City string `json:"city" jsonschema:"required,description=The city to get the weather for"`
}

func GetWeather(input GetWeatherParams) error {
	return nil
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
