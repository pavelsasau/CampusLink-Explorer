package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Structures pour la météo (Open-Meteo Forecast)
type WeatherResponse struct {
	Current struct {
		Temperature2m      float64 `json:"temperature_2m"`
		RelativeHumidity2m int     `json:"relative_humidity_2m"`
		WindSpeed10m       float64 `json:"wind_speed_10m"`
	} `json:"current"`
	Hourly struct {
		Time          []string  `json:"time"`
		Temperature2m []float64 `json:"temperature_2m"`
	} `json:"hourly"`
}

// Structures pour la qualité de l'air (Open-Meteo Air Quality)
type AirQualityResponse struct {
	Current struct {
		EuropeanAqi int `json:"european_aqi"`
	} `json:"current"`
}

// Fonction pour récupérer la météo d'une ville via ses coordonnées
func GetWeatherData(lat, lon float64) (*WeatherResponse, error) {
	client := http.Client{Timeout: 5 * time.Second}
	
	// URL de l'API Forecast avec la température actuelle et horaire
	apiURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,wind_speed_10m&hourly=temperature_2m",
		lat, lon,
	)

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur API Météo : %d", resp.StatusCode)
	}

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, err
	}

	return &weather, nil
}

// Fonction pour récupérer la qualité de l'air d'une ville
func GetAirQualityData(lat, lon float64) (*AirQualityResponse, error) {
	client := http.Client{Timeout: 5 * time.Second}

	// URL de l'API Air Quality
	apiURL := fmt.Sprintf(
		"https://air-quality-api.open-meteo.com/v1/air-quality?latitude=%f&longitude=%f&current=european_aqi",
		lat, lon,
	)

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur API Air Quality : %d", resp.StatusCode)
	}

	var aq AirQualityResponse
	if err := json.NewDecoder(resp.Body).Decode(&aq); err != nil {
		return nil, err
	}

	return &aq, nil
}