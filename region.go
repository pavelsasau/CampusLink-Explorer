package main

import (
	"encoding/json" // Décoder JSON | Декодировать JSON
	"fmt"           // Formater les erreurs | Формировать ошибки
	"log"           // Enregistrer les erreurs | Записывать ошибки
	"net/http"      // Gérer les requêtes HTTP | Обрабатывать HTTP-запросы
	"regexp"        // Vérifier le code pays | Проверять код страны
	"sort"          // Trier les villes | Сортировать города
	"strconv"       // Convertir les nombres | Преобразовывать числа
	"strings"       // Traiter le texte | Обрабатывать текст
	"time"          // Gérer les délais | Управлять временем ожидания
)

type GeoCity struct { // Données d'une ville | Данные города
	Name       string  `json:"name"`       // Nom de la ville | Название города
	Latitude   float64 `json:"lat"`        // Latitude | Широта
	Longitude  float64 `json:"lon"`        // Longitude | Долгота
	Population int     `json:"population"` // Population | Население
}

type CitiesResponse struct { // Réponse de l'API | Ответ API
	Settlements []GeoCity `json:"settlements"` // Liste des localités | Список населённых пунктов
}

type RegionPage struct { // Données de la page HTML | Данные HTML-страницы
	Country string    // Code du pays | Код страны
	Cities  []GeoCity // Liste des villes | Список городов
	Total   int       // Nombre de résultats | Количество результатов
}

func GetCities(countryCode string, regionID int) ([]GeoCity, error) { // Obtenir les villes | Получить города

	apiURL := fmt.Sprintf(
		"https://api.placedb.org/v1/cities/%s/%d.json",
		countryCode, regionID,
	) // Construire l'URL de l'API | Создать адрес API

	client := http.Client{Timeout: 15 * time.Second} // Limiter le délai | Ограничить ожидание

	response, err := client.Get(apiURL) // Envoyer une requête | Отправить запрос
	if err != nil {                     // Vérifier l'erreur | Проверить ошибку
		return nil, err
	}

	defer response.Body.Close() // Fermer la réponse | Закрыть ответ

	if response.StatusCode != http.StatusOK { // Vérifier le statut HTTP | Проверить HTTP-статус
		return nil, fmt.Errorf("erreur API : %d", response.StatusCode)
	}

	var data CitiesResponse // Stocker les données JSON | Хранить данные JSON

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil { // Décoder JSON | Декодировать JSON
		return nil, err
	}

	return data.Settlements, nil // Retourner les villes | Вернуть города
}

func regionPage(w http.ResponseWriter, r *http.Request) { // Afficher la région | Показать регион

	country := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("country"))) // Lire le code pays | Прочитать код страны

	validCode := regexp.MustCompile(`^[A-Z]{2}$`) // Format du code pays | Формат кода страны

	if !validCode.MatchString(country) { // Vérifier le pays | Проверить страну
		http.Error(w, "Code pays invalide", http.StatusBadRequest) // Erreur 400 | Ошибка 400
		return
	}

	regionID, err := strconv.Atoi(r.URL.Query().Get("id")) // Lire l'ID de région | Прочитать ID региона

	if err != nil || regionID <= 0 { // Vérifier l'ID | Проверить ID
		http.Error(w, "ID de région invalide", http.StatusBadRequest)
		return
	}

	cities, err := GetCities(country, regionID) // Charger les villes | Загрузить города

	if err != nil { // Vérifier l'erreur API | Проверить ошибку API
		log.Println(err) // Afficher l'erreur | Вывести ошибку
		http.Error(w, "Service indisponible", http.StatusBadGateway)
		return
	}

	sort.Slice(cities, func(i, j int) bool { // Trier par population décroissante | Сортировать по убыванию населения
		return cities[i].Population > cities[j].Population
	})

	data := RegionPage{ // Préparer les données HTML | Подготовить данные HTML
		Country: country,     // Pays sélectionné | Выбранная страна
		Cities:  cities,      // Villes trouvées | Найденные города
		Total:   len(cities), // Nombre de villes | Количество городов
	}

	renderPage(w, "region.html", data) // Afficher la page HTML | Показать HTML-страницу
}
