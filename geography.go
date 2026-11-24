package main

import (
	"encoding/json" // Decode JSON | Декодировать JSON
	"fmt"           // Format errors | Формировать ошибки
	"net/http"      // HTTP requests | HTTP-запросы
	"regexp"        // Validate codes | Проверять коды
	"strings"       // Work with text | Работа с текстом
	"time"          // Request timeout | Время ожидания
)

type GeoRegion struct { // Geographic region | Географический регион
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Population  int    `json:"population"`
	Settlements int    `json:"settlements"`
}

type RegionsResponse struct { // API response | Ответ API
	Country string      `json:"country"`
	Regions []GeoRegion `json:"regions"`
}

func GetRegions(countryCode string) ([]GeoRegion, error) { // Get regions | Получить регионы

	countryCode = strings.ToUpper(strings.TrimSpace(countryCode)) // Normalize code | Нормализовать код

	validCode := regexp.MustCompile(`^[A-Z]{2}$`) // Country code format | Формат кода страны

	if !validCode.MatchString(countryCode) { // Check country code | Проверить код страны
		return nil, fmt.Errorf("code pays invalide")
	}

	apiURL := "https://api.placedb.org/v1/regions/" + countryCode + ".json" // API address | Адрес API

	client := http.Client{
		Timeout: 10 * time.Second, // Request timeout | Время ожидания запроса
	}

	response, err := client.Get(apiURL) // Send API request | Отправить запрос к API

	if err != nil { // Check connection error | Проверить ошибку соединения
		return nil, err
	}

	defer response.Body.Close() // Close response | Закрыть ответ

	if response.StatusCode != http.StatusOK { // Check HTTP status | Проверить HTTP-статус
		return nil, fmt.Errorf("erreur API : %d", response.StatusCode)
	}

	var data RegionsResponse // Store API data | Хранить данные API

	err = json.NewDecoder(response.Body).Decode(&data) // Decode JSON | Декодировать JSON

	if err != nil { // Check JSON error | Проверить ошибку JSON
		return nil, err
	}

	return data.Regions, nil // Return regions | Вернуть регионы
}
