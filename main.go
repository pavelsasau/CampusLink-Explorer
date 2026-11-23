package main

import (
	"fmt"           // Print text | Вывод текста
	"html/template" // HTML templates | HTML-шаблоны
	"log"           // Log errors | Вывод ошибок
	"net/http"      // HTTP server | HTTP-сервер
	"strings"       // Work with text | Работа с текстом
)

func main() { // Main function | Главная функция

	http.HandleFunc("/", homePage)         // Home page | Главная страница
	http.HandleFunc("/search", searchPage) // Search page | Страница поиска

	fmt.Println("Serveur démarré !")            // Server started | Сервер запущен
	fmt.Println("Ouvrez http://localhost:8080") // Website address | Адрес сайта

	log.Fatal(http.ListenAndServe(":8080", nil)) // Start server | Запустить сервер
}

func homePage(w http.ResponseWriter, r *http.Request) { // Home page | Главная страница

	if r.URL.Path != "/" { // Check URL | Проверить адрес
		http.NotFound(w, r) // Return 404 | Вернуть ошибку 404
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8") // HTML response | HTML-ответ

	fmt.Fprintln(w, "<h1>Weather Explorer</h1>")                                  // Page title | Заголовок страницы
	fmt.Fprintln(w, "<form action='/search' method='GET'>")                       // Search form | Форма поиска
	fmt.Fprintln(w, "<input name='city' placeholder='Nom de la ville' required>") // City input | Ввод города
	fmt.Fprintln(w, "<button type='submit'>Rechercher</button>")                  // Search button | Кнопка поиска
	fmt.Fprintln(w, "</form>")                                                    // Close form | Закрыть форму
}

func searchPage(w http.ResponseWriter, r *http.Request) { // Search results | Результаты поиска

	cityName := strings.TrimSpace(r.URL.Query().Get("city")) // Get city name | Получить название города

	if len([]rune(cityName)) < 2 { // Validate input | Проверить ввод
		http.Error(w, "Nom de ville invalide", http.StatusBadRequest) // Return 400 | Вернуть ошибку 400
		return
	}

	cities, err := SearchCity(cityName) // Search cities | Найти города

	if err != nil { // Check API error | Проверить ошибку API
		log.Println(err)                                                       // Log error | Записать ошибку
		http.Error(w, "Erreur du service de recherche", http.StatusBadGateway) // Return 502 | Вернуть ошибку 502
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8") // HTML response | HTML-ответ

	page := `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Résultats - Weather Explorer</title>
</head>
<body>
    <h1>Résultats de recherche</h1>

    <p>Recherche : {{.Query}}</p>

    {{if .Cities}}
        {{range .Cities}}
            <p>
                {{.Name}} — {{.Country}} :
                {{.Latitude}}, {{.Longitude}}
            </p>
        {{end}}
    {{else}}
        <p>Aucune ville trouvée.</p>
    {{end}}

    <a href="/">Retour à l'accueil</a>
</body>
</html>`

	tmpl, err := template.New("results").Parse(page) // Parse HTML template | Подготовить HTML-шаблон

	if err != nil { // Check template error | Проверить ошибку шаблона
		http.Error(w, "Erreur du serveur", http.StatusInternalServerError) // Return 500 | Вернуть ошибку 500
		return
	}

	data := struct { // Page data | Данные страницы
		Query  string
		Cities []City
	}{
		Query:  cityName,
		Cities: cities,
	}

	if err := tmpl.Execute(w, data); err != nil { // Render results | Показать результаты
		log.Println(err) // Log error | Записать ошибку
	}
}
