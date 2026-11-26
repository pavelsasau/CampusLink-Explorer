package main

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/explore", explorePage) // Page d'exploration | Страница исследования       // Home page | Главная страница
	http.HandleFunc("/search", searchPage)   // Search cities | Поиск городов
	http.HandleFunc("/city", cityPage)       // Selected city | Выбранный город
	http.HandleFunc("/country", countryPage)
	http.HandleFunc("/region", regionPage)                                                     // Country regions | Регионы страны
	http.HandleFunc("/map.svg", mapPage)                                                       // Carte interactive | Интерактивная карта                                                 // Afficher les villes d'une région | Показать города региона
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // Servir les fichiers statiques | Отдавать статические файлы

	log.Println("Serveur : http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderPage(w http.ResponseWriter, file string, data any) {
	tmpl, err := template.ParseFiles("templates/" + file) // Load HTML | Загрузить HTML
	if err != nil {
		log.Println(err)
		http.Error(w, "Erreur du serveur", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	renderPage(w, "index.html", nil) // Show homepage | Показать главную
}

func searchPage(w http.ResponseWriter, r *http.Request) {
	cityName := strings.TrimSpace(r.URL.Query().Get("city")) // Get city | Получить город

	if len([]rune(cityName)) < 2 {
		http.Error(w, "Nom de ville invalide", http.StatusBadRequest)
		return
	}

	cities, err := SearchCity(cityName) // Search API | Поиск через API
	if err != nil {
		log.Println(err)
		http.Error(w, "Erreur du service de recherche", http.StatusBadGateway)
		return
	}

	data := struct {
		Query  string
		Cities []City
	}{
		Query:  cityName,
		Cities: cities,
	}

	renderPage(w, "search.html", data) // Show results | Показать результаты
}

func cityPage(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	country := strings.TrimSpace(r.URL.Query().Get("country"))

	lat, errLat := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, errLon := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)

	if name == "" || country == "" || errLat != nil || errLon != nil ||
		math.IsNaN(lat) || math.IsNaN(lon) ||
		math.IsInf(lat, 0) || math.IsInf(lon, 0) ||
		lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		http.Error(w, "Paramètres invalides", http.StatusBadRequest)
		return
	}

	// 1. Récupérer les données météo (avec les températures horaires)
	weather, err := GetWeatherData(lat, lon)
	if err != nil {
		log.Println("Erreur météo :", err)
		http.Error(w, "Erreur lors de la récupération de la météo", http.StatusBadGateway)
		return
	}

	// 2. Générer le graphique SVG (on ne prend que les 24 premières heures pour que ce soit lisible)
	maxHours := 24
	if len(weather.Hourly.Time) < maxHours {
		maxHours = len(weather.Hourly.Time)
	}

	svgString := GenerateTemperatureChart(
		weather.Hourly.Time[:maxHours],
		weather.Hourly.Temperature2m[:maxHours],
	)

	// 3. Définir la structure "data" passée au template city.html
	data := struct {
		Name      string
		Country   string
		Latitude  float64
		Longitude float64
		Weather   *WeatherResponse
		ChartSVG  template.HTML // Intégration du type template.HTML
	}{
		Name:      name,
		Country:   country,
		Latitude:  lat,
		Longitude: lon,
		Weather:   weather,
		ChartSVG:  template.HTML(svgString),
	}

	renderPage(w, "city.html", data) // Afficher la page avec toutes les données
}

// Afficher la page d'exploration | Показать страницу исследования
func explorePage(w http.ResponseWriter, r *http.Request) {

	// Vérifier l'adresse | Проверить адрес
	if r.URL.Path != "/explore" {
		http.NotFound(w, r)
		return
	}

	// Charger la page HTML | Загрузить HTML-страницу
	renderPage(w, "explore.html", nil)
}
