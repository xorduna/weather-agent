package tools

import (
	"context"
	"encoding/json"
	"github.com/acai-travel/tech-challenge/internal/weather"
	"github.com/openai/openai-go/v2"
	"log/slog"
)

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

func (w *WeatherTool) Execute(ctx context.Context, args ...string) (string, error) {

	var parameters struct {
		Location string `json:"location"`
	}

	if err := json.Unmarshal([]byte(args[0]), &parameters); err != nil {
		return "failed to parse tool call arguments: " + err.Error(), nil
	}

	slog.InfoContext(ctx, "Executing WeatherTool", "location", parameters.Location)

	currentWeather, err := weather.GetCurrentWeather(ctx, parameters.Location)
	if err != nil {
		return "", err
	}

	return currentWeather.Condition.Text, nil
}
