package forum

import (
	"net/http"
	"html/template"
	"database/sql"
	"fmt"
	models "forum/models"
)

type PostData struct {
	Users models.Users
	Error string
}

func Post(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	tmpl := template.Must(template.ParseFiles("./views/post.html")) // Affiche la page
	data := PostData {}
    
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

	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Println(err)
	}

	forumRepository := NewSQLiteRepository(db)

	user, err := forumRepository.GetUserByCookie(c.Value)
	if err != nil {
		fmt.Println(err)
	}

	data = PostData {
		Users: *user,
		Error: "",
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, err)
	}
}