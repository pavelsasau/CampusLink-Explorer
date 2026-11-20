
package main

import (
    "fmt"    
    "log"     
    "net/http"
)

func main() { 

    http.HandleFunc("/", homePage) // Associe l'adresse "/" à la page d'accueil

    fmt.Println("Serveur démarré !")

    fmt.Println("Ouvrez http://localhost:8080")

    log.Fatal(http.ListenAndServe(":8080", nil)) // Démarre le serveur sur le port 8080
}

func homePage(w http.ResponseWriter, r *http.Request) { // Fonction qui affiche la page d'accueil

    fmt.Fprintln(w, "Weather Explorer fonctionne !")

}
