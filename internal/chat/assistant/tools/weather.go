package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/openai/openai-go/v2"
	"io"
	"log/slog"
	"net/http"
	"os"
)

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

// Struct for the "current" field
type Current struct {
	LastUpdatedEpoch int       `json:"last_updated_epoch"`
	LastUpdated      string    `json:"last_updated"`
	TempC            float64   `json:"temp_c"`
	TempF            float64   `json:"temp_f"`
	IsDay            int       `json:"is_day"`
	Condition        Condition `json:"condition"`
	WindMph          float64   `json:"wind_mph"`
	WindKph          float64   `json:"wind_kph"`
	WindDegree       int       `json:"wind_degree"`
	WindDir          string    `json:"wind_dir"`
	PressureMb       float64   `json:"pressure_mb"`
	PressureIn       float64   `json:"pressure_in"`
	PrecipMm         float64   `json:"precip_mm"`
	PrecipIn         float64   `json:"precip_in"`
	Humidity         int       `json:"humidity"`
	Cloud            int       `json:"cloud"`
	FeelslikeC       float64   `json:"feelslike_c"`
	FeelslikeF       float64   `json:"feelslike_f"`
	WindchillC       float64   `json:"windchill_c"`
	WindchillF       float64   `json:"windchill_f"`
	HeatindexC       float64   `json:"heatindex_c"`
	HeatindexF       float64   `json:"heatindex_f"`
	DewpointC        float64   `json:"dewpoint_c"`
	DewpointF        float64   `json:"dewpoint_f"`
	VisKm            float64   `json:"vis_km"`
	VisMiles         float64   `json:"vis_miles"`
	Uv               float64   `json:"uv"`
	GustMph          float64   `json:"gust_mph"`
	GustKph          float64   `json:"gust_kph"`
	ShortRad         int       `json:"short_rad"`
	DiffRad          int       `json:"diff_rad"`
	DNI              int       `json:"dni"`
	GTI              int       `json:"gti"`
}

// Struct for the "location" field
type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int     `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

// Root struct for the whole response
type WeatherResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

func getWeather(location string) (string, error) {
	weatherKey := os.Getenv("WEATHER_API_KEY")

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", weatherKey, location)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Returns error if status is not 200
		body, _ := io.ReadAll(resp.Body)
		return "", errors.New(string(body))
	}

	var weather WeatherResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&weather); err != nil {
		return "", err
	}
	return weather.Current.Condition.Text, nil

}

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

	// Geat Weather Key from environment. In real implementation, I would centralize the setup of all variable in a config at main level

	var parameters struct {
		Location string `json:"location"`
	}

	if err := json.Unmarshal([]byte(args[0]), &parameters); err != nil {
		return "failed to parse tool call arguments: " + err.Error(), nil
	}

	slog.InfoContext(ctx, "Executing WeatherTool", "location", parameters.Location)

	return "Weather is fine.", nil
}
