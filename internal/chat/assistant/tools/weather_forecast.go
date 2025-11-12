package tools

import (
	"context"
	"encoding/json"
	"github.com/acai-travel/tech-challenge/internal/weather"
	"github.com/openai/openai-go/v2"
	"log/slog"
)

type WeatherForecastTool struct{}

func (w *WeatherForecastTool) Name() string {
	return "get_weather_forecast"
}

func (w *WeatherForecastTool) Description() string {
	return "Get weather forecast at the given location for the given days"
}

func (w *WeatherForecastTool) Parameters() openai.FunctionParameters {
	return openai.FunctionParameters{
		"type": "object",
		"properties": map[string]any{
			"location": map[string]string{
				"type":        "string",
				"description": "Given location",
			},
			"days": map[string]string{
				"type":        "integer",
				"description": "Number of days to forecast (1-7)",
			},
		},
		"required": []string{"location"},
	}
}

func (w *WeatherForecastTool) Execute(ctx context.Context, args ...string) (string, error) {

	var parameters struct {
		Location string `json:"location"`
		Days     int    `json:"days"`
	}

	if err := json.Unmarshal([]byte(args[0]), &parameters); err != nil {
		return "failed to parse tool call arguments: " + err.Error(), nil
	}

	slog.InfoContext(ctx, "Executing WeatherForecastTool", "location", parameters.Location, "days", parameters.Days)

	forecast, err := weather.GetWeatherForecast(ctx, parameters.Location, parameters.Days, false, false)
	if err != nil {
		// Return error as a user-friendly message
		return "Weather forecast service error: " + err.Error(), nil
	}

	response := "Weather forecast for " + parameters.Location + ":\n"

	for _, day := range forecast.Forecastday {
		slog.InfoContext(ctx, "Forecast", "date", day.Date, "condition", day.Day.Condition.Text)
		response += day.Date + ": " + day.Day.Condition.Text + "\n"
	}

	return response, nil
}
