package main // Main package | Главный пакет программы

import (
	"log"      // Display server errors | Выводить ошибки сервера
	"net/http" // Handle HTTP requests | Обрабатывать HTTP-запросы
	"regexp"   // Validate country codes | Проверять коды стран
	"strings"  // Process text | Обрабатывать текст
)

type CountryPage struct { // Data for country HTML page | Данные для HTML-страницы страны
	Code    string
	Regions []GeoRegion
}

func countryPage(w http.ResponseWriter, r *http.Request) { // Handle country page | Обработать страницу страны

	// Get country code from URL | Получить код страны из адреса
	code := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("code")))

	// Accept two uppercase letters | Разрешить две заглавные буквы
	validCode := regexp.MustCompile(`^[A-Z]{2}$`)

	if !validCode.MatchString(code) { // Check country code | Проверить код страны
		http.Error(w, "Code pays invalide", http.StatusBadRequest) // Return HTTP 400 | Вернуть ошибку 400
		return                                                     // Stop function | Остановить функцию
	}

	regions, err := GetRegions(code) // Request regions from API | Запросить регионы из API

	if err != nil { // Check API error | Проверить ошибку API
		log.Println(err)                                                          // Log error in terminal | Записать ошибку в терминал
		http.Error(w, "Service géographique indisponible", http.StatusBadGateway) // Return HTTP 502 | Вернуть ошибку 502
		return                                                                    // Stop function | Остановить функцию
	}

	data := CountryPage{ // Prepare page data | Подготовить данные страницы
		Code:    code,    // Country code | Код страны
		Regions: regions, // Regions from API | Регионы из API
	}

	renderPage(w, "country.html", data) // Send data to HTML template | Передать данные в HTML-шаблон
}
