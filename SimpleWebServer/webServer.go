package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Benvingut a la meva pàgina web</h1>")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Servidor en execució a http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
