package tools

import "github.com/openai/openai-go/v2"

type WeatherTool struct{}

func (w *WeatherTool) Name() string {
	return "get_weather"
}

func (w *WeatherTool) Description() string {
	return "Get weather at the given location"
}

func (w *WeatherTool) Parameters() openai.FunctionParameters {
	return openai.FunctionParameters{
		"type": "object",
		"properties": map[string]any{
			"location": map[string]string{
				"type": "string",
			},
		},
		"required": []string{"location"},
	}
}

func (w *WeatherTool) Execute(args ...string) (string, error) {
	return "Weather is fine.", nil
}
