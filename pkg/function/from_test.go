package function_test

import (
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/ivanvanderbyl/smith/pkg/function"
	"github.com/stretchr/testify/assert"
)

type GetWeatherParams struct {
	// The city to get the weather for
	City string `json:"city" jsonschema:"required"`
}

// GetWeather gets the weather for a given city
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
	a.Equal("GetWeather", tool.Name)
	a.NotNil(tool.Parameters)
	params, ok := tool.Parameters.(*jsonschema.Schema)
	a.True(ok)
	city, ok := params.Properties.Get("city")
	a.True(ok)
	a.Equal("The city to get the weather for", city.Description)
}
