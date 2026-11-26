package main

import (
	"fmt"
	"strings"
)

// Génère un graphique SVG simple pour les températures horaires
func GenerateTemperatureChart(times []string, temps []float64) string {
	if len(temps) == 0 {
		return "<p>Aucune donnée disponible pour le graphique.</p>"
	}

	// Dimensions du SVG
	width := 600
	height := 200
	padding := 40

	// Trouver la température min et max pour adapter l'échelle
	minTemp := temps[0]
	maxTemp := temps[0]
	for _, t := range temps {
		if t < minTemp {
			minTemp = t
		}
		if t > maxTemp {
			maxTemp = t
		}
	}

	// Éviter une division par zéro si toutes les températures sont identiques
	if minTemp == maxTemp {
		maxTemp += 1.0
	}

	// Dimensions utiles pour le dessin
	plotWidth := float64(width - 2*padding)
	plotHeight := float64(height - 2*padding)

	var points []string
	stepX := plotWidth / float64(len(temps)-1)

	// Construction des points (X, Y)
	for i, t := range temps {
		x := float64(padding) + float64(i)*stepX
		// Inverser l'axe Y car en SVG le 0 est en haut
		y := float64(height-padding) - ((t - minTemp) / (maxTemp - minTemp)) * plotHeight
		points = append(points, fmt.Sprintf("%.1f,%.1f", x, y))
	}

	polylinePoints := strings.Join(points, " ")

	// Construction du code SVG final
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg viewBox="0 0 %d %d" width="100%%" height="100%%" style="background: #f8fafc; border-radius: 8px;">`, width, height))
	
	// Titre du graphique
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="25" font-size="14" font-family="sans-serif" font-weight="bold" fill="#1e293b" text-anchor="middle">Évolution des températures (prochaines heures)</text>`, width/2))

	// Ligne de la courbe
	sb.WriteString(fmt.Sprintf(`<polyline fill="none" stroke="#2563eb" stroke-width="3" points="%s" />`, polylinePoints))

	// Affichage des points et valeurs min/max sur les côtés
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="10" fill="#64748b">Min: %.1f°C</text>`, padding, height-10, minTemp))
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="10" fill="#64748b" text-anchor="end">Max: %.1f°C</text>`, width-padding, height-10, maxTemp))

	sb.WriteString(`</svg>`)

	return sb.String()
}