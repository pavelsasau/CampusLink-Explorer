package main

import (
	"encoding/json" // Décoder JSON | Читать JSON
	"fmt"           // Formater le texte | Формировать текст
	"net/http"      // Requêtes HTTP | HTTP-запросы
	"regexp"        // Valider le code | Проверить код
	"strings"       // Traiter le texte | Обрабатывать текст
	"time"          // Gérer le délai | Управлять ожиданием
)

// Informations du pays | Информация о стране
type CountryInfo struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`

	Capital []string `json:"capital"`

	CapitalInfo struct {
		LatLng []float64 `json:"latlng"`
	} `json:"capitalInfo"`
}

// Réponse pour la carte | Ответ для карты
type CountryPreview struct {
	Name        string  `json:"name"`
	Capital     string  `json:"capital"`
	Temperature float64 `json:"temperature"`
	Wind        float64 `json:"wind"`
}

// Afficher les informations du pays | Показать информацию о стране
func countryPreviewPage(w http.ResponseWriter, r *http.Request) {

	// Lire le code du pays | Получить код страны
	code := strings.ToUpper(r.URL.Query().Get("code"))

	validCode := regexp.MustCompile(`^[A-Z]{2}$`)

	if !validCode.MatchString(code) {
		http.Error(w, "Code invalide", http.StatusBadRequest)
		return
	}

	// Appeler l'API des pays | Запросить API стран
	client := http.Client{Timeout: 6 * time.Second}

	url := "https://restcountries.com/v3.1/alpha/" + code

	response, err := client.Get(url)

	if err != nil {
		http.Error(w, "Service indisponible", http.StatusBadGateway)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "Pays indisponible", http.StatusBadGateway)
		return
	}

	// Décoder les pays | Прочитать данные стран
	var countries []CountryInfo

	if err := json.NewDecoder(response.Body).Decode(&countries); err != nil {
		http.Error(w, "JSON invalide", http.StatusBadGateway)
		return
	}

	if len(countries) == 0 ||
		len(countries[0].Capital) == 0 ||
		len(countries[0].CapitalInfo.LatLng) < 2 {

		http.Error(w, "Capitale indisponible", http.StatusNotFound)
		return
	}

	country := countries[0]

	// Coordonnées de la capitale | Координаты столицы
	lat := country.CapitalInfo.LatLng[0]
	lon := country.CapitalInfo.LatLng[1]

	// Réutiliser la météo existante | Использовать готовую функцию погоды
	weather, err := GetWeatherData(lat, lon)

	if err != nil {
		http.Error(w, "Météo indisponible", http.StatusBadGateway)
		return
	}

	// Préparer la réponse | Подготовить ответ
	result := CountryPreview{
		Name:        country.Name.Common,
		Capital:     country.Capital[0],
		Temperature: weather.Current.Temperature2m,
		Wind:        weather.Current.WindSpeed10m,
	}

	// Envoyer le JSON | Отправить JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		fmt.Println("Erreur JSON :", err)
	}
}
