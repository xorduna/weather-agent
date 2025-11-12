package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

const API_URL = "http://api.weatherapi.com/v1"

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

// Struct for the "current" field
type CurrentWeather struct {
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
	Location Location       `json:"location"`
	Current  CurrentWeather `json:"current"`
}

// Struct for the "day" field inside "forecastday"
type Day struct {
	Condition Condition `json:"condition"`
	// Ommitted fields
}

// Struct for each "forecastday"
type ForecastDay struct {
	Date string `json:"date"`
	Day  Day    `json:"day"`
	// ommitted fields
}

// Struct for the "forecast" field
type Forecast struct {
	Forecastday []ForecastDay `json:"forecastday"`
}

// Struct for the full forecast response
type WeatherForecastResponse struct {
	Forecast Forecast `json:"forecast"`
	// ...other fields omitted...
}

func GetCurrentWeather(ctx context.Context, location string) (CurrentWeather, error) {
	weatherKey := os.Getenv("WEATHER_API_KEY")

	slog.InfoContext(ctx, "Fetching weather forecast", "location", location)
	url := fmt.Sprintf("%s/v1/current.json?key=%s&q=%s&aqi=no", API_URL, weatherKey, location)
	resp, err := http.Get(url)
	if err != nil {
		return CurrentWeather{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Returns error if status is not 200
		body, _ := io.ReadAll(resp.Body)
		return CurrentWeather{}, errors.New(string(body))
	}

	var weather WeatherResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&weather); err != nil {
		return CurrentWeather{}, err
	}
	return weather.Current, nil

}

func GetWeatherForecast(ctx context.Context, location string, days int, alerts bool, air_quality bool) (Forecast, error) {
	weatherKey := os.Getenv("WEATHER_API_KEY")

	slog.InfoContext(ctx, "Fetching weather forecast", "location", location, "days", days)
	url := fmt.Sprintf("%s/v1/forecast.json?key=%s&q=%s&days=%d&aqi=%s&alerts=%s",
		API_URL,
		weatherKey, location, days,
		boolToYesNo(air_quality), boolToYesNo(alerts),
	)
	resp, err := http.Get(url)
	if err != nil {
		// Error can occur here if the HTTP request fails (network, DNS, etc.)
		return Forecast{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Error can occur here if the API returns a non-200 status (bad request, unauthorized, etc.)
		body, _ := io.ReadAll(resp.Body)
		return Forecast{}, errors.New(string(body))
	}

	var forecastResp WeatherForecastResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&forecastResp); err != nil {
		// Error can occur here if the response body is not valid JSON or doesn't match the struct
		return Forecast{}, err
	}

	return forecastResp.Forecast, nil
}

// Helper function to convert bool to "yes"/"no" string for API params
func boolToYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
