package main

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
	"net/http"
	"html/template"
	"fmt"
	forum "forum/controllers"
)

type Data struct {

}

const fileName = "forum.db"

func main() {
    db, err := sql.Open("sqlite3", fileName)
    if err != nil {
        log.Fatal(err)
    }

	forumRepository := forum.NewSQLiteRepository(db)

	if err := forumRepository.Migrate(); err != nil {
        log.Fatal(err)
    }

	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("./assets", http.StripPrefix("./assets", fs))
	http.HandleFunc("/", Handler)
	http.HandleFunc("/controller/login", forum.Login)
	http.HandleFunc("/controller/register", forum.Register)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Localhost:8080 open")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./views/index.html")) // Affiche la page

	// Affiche dans le terminal l'activité sur le site
	switch r.Method {
	case "GET":
		fmt.Println("GET")
	case "POST": // Gestion d'erreur
		if err := r.ParseForm(); err != nil {
			return
		}
	}

	data := Data {}

	err := tmpl.Execute(w, data)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, err)
	}
}