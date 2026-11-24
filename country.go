package main

import "net/http"

// Region data | Данные региона
type Region struct {
	Name   string
	Cities []string
}

// Country page data | Данные страницы страны
type CountryPage struct {
	Name    string
	Regions []Region
}

func countryPage(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code") // Country code | Код страны

	if code != "FR" { // Check country | Проверить страну
		http.NotFound(w, r)
		return
	}

	france := CountryPage{
		Name: "France",

		Regions: []Region{
			{
				Name: "Île-de-France",
				Cities: []string{
					"Paris",
					"Versailles",
					"Nanterre",
				},
			},
			{
				Name: "Auvergne-Rhône-Alpes",
				Cities: []string{
					"Lyon",
					"Grenoble",
					"Clermont-Ferrand",
				},
			},
			{
				Name: "Provence-Alpes-Côte d'Azur",
				Cities: []string{
					"Marseille",
					"Nice",
					"Toulon",
				},
			},
			{
				Name: "Nouvelle-Aquitaine",
				Cities: []string{
					"Bordeaux",
					"Limoges",
					"Poitiers",
				},
			},
			{
				Name: "Occitanie",
				Cities: []string{
					"Toulouse",
					"Montpellier",
					"Nîmes",
				},
			},
			{
				Name: "Bretagne",
				Cities: []string{
					"Rennes",
					"Brest",
					"Quimper",
				},
			},
		},
	}

	renderPage(w, "country.html", france) // Display country | Показать страну
}
