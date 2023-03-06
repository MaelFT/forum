package main

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
	"net/http"
	"html/template"
	"fmt"
	controllers "forum/controllers"
	models "forum/models"
)

type Data struct {
	Users models.Users
	Connected int
	Error string
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

	if err := forumRepository.TableLike(); err != nil {
        log.Fatal(err)
    }

	if err := forumRepository.TableCategories(); err != nil {
        log.Fatal(err)
    }

	fs := http.FileServer(http.Dir("./views/assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/index", Handler)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/signup", controllers.Register)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/post", controllers.Post)
	http.HandleFunc("/addPost", controllers.AddPost)
	http.HandleFunc("/category", controllers.Category)
	http.HandleFunc("/addCategory", controllers.AddCategory)
	http.HandleFunc("/", controllers.PageNotFound)
	fmt.Println("Localhost:8080 open")
	http.ListenAndServe(":8080", nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	tmpl := template.Must(template.ParseFiles("./views/index.html")) // Affiche la page
	data := Data {
		Connected: 0,
	}
    
	if err != nil || c.Value == "" {
		fmt.Println(c, err)
		err = tmpl.Execute(w, data)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintln(w, err)
		}
		return
	}

	// Affiche dans le terminal l'activité sur le site
	switch r.Method {
	case "GET":
		fmt.Println("GET")
	case "POST": // Gestion d'erreur
		if err := r.ParseForm(); err != nil {
			return
		}
	}

	data = Data {
		Connected: 1,
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, err)
	}
}