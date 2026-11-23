package main

import (
	"encoding/json"
	"fmt"
	"net/http" // Requêtes HTTP | HTTP-запросы
	"net/url"  // Encoder URL | Кодировать URL
	"strings"  // Traiter texte | Обработка текста
	"time"     // Gestion du temps | Работа со временем
)

type City struct { // Données ville | Данные города
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GeocodingResponse struct { // Réponse API | Ответ API
	Results []City `json:"results"` // Liste des villes | Список городов
}

func SearchCity(name string) ([]City, error) { // Chercher ville | Найти город

	name = strings.TrimSpace(name) // Supprimer espaces | Убрать пробелы

	if len([]rune(name)) < 2 { // Vérifier le nom | Проверить название
		return nil, fmt.Errorf("nom de ville trop court")
	}

	client := http.Client{Timeout: 5 * time.Second} // Délai 5s | Ожидание 5 сек.

	apiURL := "https://geocoding-api.open-meteo.com/v1/search?name=" +
		url.QueryEscape(name) + "&count=5&language=fr&format=json"

	response, err := client.Get(apiURL) // Appeler API | Запросить API
	if err != nil {                     // Vérifier erreur | Проверить ошибку
		return nil, err
	}
	defer response.Body.Close() // Fermer réponse | Закрыть ответ

	if response.StatusCode != http.StatusOK { // Vérifier statut | Проверить статус
		return nil, fmt.Errorf("erreur API : %d", response.StatusCode)
	}

	var data GeocodingResponse // Stocker données | Хранить данные

	err = json.NewDecoder(response.Body).Decode(&data) // Décoder JSON | Декодировать JSON
	if err != nil {                                    // Vérifier JSON | Проверить JSON
		return nil, err
	}

	return data.Results, nil // Retourner villes | Вернуть города
}
