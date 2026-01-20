package main

import (
	art2 "ascii-art-web2/functions"
	"html/template"
	"log"
	"net/http"
)

// Route par défault, qui renvoie une erreur 404
func notFound(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/404.html"))
	w.WriteHeader(http.StatusNotFound)
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error 500 : unavailable server", http.StatusInternalServerError)
	}
}

// Route vers la page principale
func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error 500 : unavailable server", http.StatusInternalServerError)
	}
}

// Gère la transformation
func transformHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unauthorized method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Parsing Error", http.StatusBadRequest)
		return
	}

	// Récupère la phrase et le choix de bannière
	phrase := r.FormValue("phrase")
	banner := r.FormValue("banner")

	// Réalise la transformation en ASCII art
	result, err := art2.TransformPhrase(phrase, banner)
	if err != nil {
		http.Error(w, "Error 400 : "+result, http.StatusBadRequest)
		log.Println("Error 400 :", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Renvoie le résultat
	w.Write([]byte(result))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/transform", transformHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// Si l'URL n'est pas celui de la racine, on renvoie l'erreur 404
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			home(w, r)
			return
		}
		notFound(w, r)
	})

	log.Print("starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
