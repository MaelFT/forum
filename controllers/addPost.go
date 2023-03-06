package forum

import (
	"net/http"
	"html/template"
	"database/sql"
	"fmt"
	models "forum/models"
)

type AddPostData struct {
	Users models.Users
	Error string
}

func AddPost(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
    
	if err != nil {
		fmt.Println(c, err)
		http.Redirect(w, r, "/index", http.StatusFound)
	} else {
		if c.Value == "" {
			http.Redirect(w, r, "/index", http.StatusFound)
		}
		tmpl := template.Must(template.ParseFiles("./views/addPost.html"))

		switch r.Method {
		case "GET":
			fmt.Println("GET")
		case "POST": // Gestion d'erreur
			if err := r.ParseForm(); err != nil {
				return
			}
		}

		data := AddPostData {}

		db, err := sql.Open("sqlite3", "forum.db")
		if err != nil {
			fmt.Println(err)
		}

		forumRepository := NewSQLiteRepository(db)

		user, err := forumRepository.GetUserByCookie(c.Value)
		if err != nil {
			fmt.Println(err)
		}
	
		data = AddPostData {
			Users: *user,
			Error: "",
		}

		err = tmpl.Execute(w, data)

		if err != nil {
			fmt.Println(err)
			fmt.Fprintln(w, err)
		}
	}
}