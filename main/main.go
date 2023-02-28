package main

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
	"net/http"
	"html/template"
	"fmt"
	controllers "forum/controllers"
	// models "forum/models"
)

type Data struct {

}

const fileName = "forum.db"

func main() {
    db, err := sql.Open("sqlite3", fileName)
    if err != nil {
        log.Fatal(err)
    }

	forumRepository := controllers.NewSQLiteRepository(db)

	if err := forumRepository.TableUsers(); err != nil {
        log.Fatal(err)
    }

	if err := forumRepository.TablePosts(); err != nil {
        log.Fatal(err)
    }

	if err := forumRepository.TableComments(); err != nil {
        log.Fatal(err)
    }

	fs := http.FileServer(http.Dir("../views/assets/"))
	http.Handle("../assets/", http.StripPrefix("../assets/", fs))
	http.HandleFunc("/", Handler)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/register", controllers.Register)
	fmt.Println("Localhost:8080 open")
	http.ListenAndServe(":8080", nil)
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