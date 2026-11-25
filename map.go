package main

import (
	"log"      // Journal des erreurs | Журнал ошибок
	"net/http" // Serveur HTTP | HTTP-сервер
	"os"       // Lire les fichiers | Читать файлы
	"regexp"   // Rechercher les pays | Находить страны
	"strings"  // Modifier le texte | Изменять текст
)

// Expression régulière pour trouver les pays | Поиск стран в SVG
var countryPath = regexp.MustCompile(`<path id="([a-z]{2})"[^>]*/>`)

// Afficher la carte interactive | Показать интерактивную карту
func mapPage(w http.ResponseWriter, _ *http.Request) {

	// Charger le fichier SVG | Загрузить SVG-файл
	svg, err := os.ReadFile("static/world.svg")

	if err != nil {
		log.Println(err)
		http.Error(w, "Carte indisponible", http.StatusInternalServerError)
		return
	}

	// Convertir le fichier en texte | Преобразовать файл в текст
	content := string(svg)

	// Parcourir les pays et ajouter des liens | Перебрать страны и добавить ссылки
	content = countryPath.ReplaceAllStringFunc(content, func(path string) string {

		// Trouver le code du pays | Найти код страны
		match := countryPath.FindStringSubmatch(path)

		// Convertir fr en FR | Преобразовать fr в FR
		code := strings.ToUpper(match[1])

		// Créer un lien vers le pays | Создать ссылку на страницу страны
		link := `<a href="/country?code=` + code + `" target="_top">`

		// Ajouter le lien autour du pays | Добавить ссылку вокруг страны
		return link + path + `</a>`
	})

	// Styles de la carte | Стили карты
	style := `<style>
        path {
            fill: #2563a6;
            stroke: #ffffff;
            stroke-width: 0.6;
            cursor: pointer;
        }

        a:hover path {
            fill: #38bdf8;
        }
    </style>`

	// Ajouter les styles au SVG | Добавить стили в SVG
	content = strings.Replace(content, ">", ">"+style, 1)

	// Définir le type du fichier | Указать тип файла
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")

	// Envoyer la carte au navigateur | Отправить карту браузеру
	_, _ = w.Write([]byte(content))
}
