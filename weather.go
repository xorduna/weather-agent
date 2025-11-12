package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Struct for the "condition" field inside "current"

// Fetches weather data from the API. Returns WeatherResponse if status 200, else returns error.
func GetCurrentWeather(apiKey, city string) (*WeatherResponse, error) {
}

func main() {
	// Example usage of GetCurrentWeather
	apiKey := "9ee70b615fe14afabc9220929251111"
	city := "Barcelona"
	weather, err := GetCurrentWeather(apiKey, city)
	if err != nil {
		// Print error if API call fails
		fmt.Println("Error:", err)
		return
	}
	// Print weather response struct
	fmt.Printf("%+v\n", weather)
}
