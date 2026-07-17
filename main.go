package main

import (
	"log"
	"net/http"
	"portfolio/db"
	"portfolio/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db.Init()
	defer db.DB.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Home)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
