package main

import (
	"encoding/json" // Decode JSON | Декодировать JSON
	"fmt"           // Format errors | Формировать ошибки
	"net/http"      // HTTP requests | HTTP-запросы
	"strings"       // Work with text | Работа с текстом
	"time"          // Request timeout | Время ожидания
)

type GeoRegion struct { // Region data | Данные региона
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Population  int    `json:"population"`
	Settlements int    `json:"settlements"`
}

type RegionsResponse struct { // API response | Ответ API
	Country string   `json:"country"`
	Regions []Region `json:"regions"`
}

func GetRegions(countryCode string) ([]Region, error) {

	countryCode = strings.ToUpper(strings.TrimSpace(countryCode)) // Normalize code | Нормализовать код

	if len(countryCode) != 2 { // Check country code | Проверить код страны
		return nil, fmt.Errorf("code pays invalide")
	}

	apiURL := "https://api.placedb.org/v1/regions/" + countryCode + ".json" // API address | Адрес API

	client := http.Client{
		Timeout: 10 * time.Second, // Set timeout | Установить время ожидания
	}

	response, err := client.Get(apiURL) // Send request | Отправить запрос
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() // Close response | Закрыть ответ

	if response.StatusCode != http.StatusOK { // Check status | Проверить статус
		return nil, fmt.Errorf("erreur API : %d", response.StatusCode)
	}

	var data RegionsResponse // Store response | Хранить ответ

	err = json.NewDecoder(response.Body).Decode(&data) // Read JSON | Прочитать JSON
	if err != nil {
		return nil, err
	}

	return data.Regions, nil // Return regions | Вернуть регионы
}
