
package main // Paquet principal | Главный пакет

import (
    "log"      // Erreurs | Ошибки
    "net/http" // HTTP | Запросы
    "regexp"   // Validation | Проверка
    "sort"     // Tri | Сортировка
    "strings"  // Texte | Текст
)

type CountryPage struct { // Données du pays | Данные страны
    Code    string      // Code pays | Код страны
    Regions []GeoRegion // Régions | Регионы
}

func countryPage(w http.ResponseWriter, r *http.Request) { // Page du pays | Страница страны

    code := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("code"))) // Lire le code | Получить код

    validCode := regexp.MustCompile(`^[A-Z]{2}$`) // Format du code | Формат кода

    if !validCode.MatchString(code) { // Code invalide ? | Код неверный?
        http.Error(w, "Code pays invalide", http.StatusBadRequest) // Erreur 400 | Ошибка 400
        return // Arrêter | Остановить
    }

    regions, err := GetRegions(code) // Charger les régions | Загрузить регионы

    if err != nil { // Erreur API ? | Ошибка API?
        log.Println(err) // Journal | Запись ошибки
        http.Error(w, "Service géographique indisponible", http.StatusBadGateway) // Erreur 502 | Ошибка 502
        return // Arrêter | Остановить
    }

    sort.Slice(regions, func(i, j int) bool { // Trier par population | Сортировать по населению
        return regions[i].Population > regions[j].Population // Plus grand d'abord | Сначала большее
    })

    data := CountryPage{ // Préparer les données | Подготовить данные
        Code:    code,    // Code pays | Код страны
        Regions: regions, // Régions triées | Отсортированные регионы
    }

    renderPage(w, "country.html", data) // Afficher HTML | Показать страницу
}
