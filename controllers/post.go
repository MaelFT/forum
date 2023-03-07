package forum

import (
	"net/http"
	"html/template"
	"database/sql"
	"strconv"
	"fmt"
	models "forum/models"
)

type PostData struct {
	Users models.Users
	SessionUser models.Users
	Posts models.Posts
	Connected int
	Error string
}

func Post(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./views/post.html")) // Affiche la page

	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Println(err)
	}

	forumRepository := NewSQLiteRepository(db)

	post, err := forumRepository.GetPostByID(id)
	if err != nil {
		fmt.Println(err)
	}

	user, err := forumRepository.GetUserByID(post.User_ID)
	if err != nil {
		fmt.Println(err)
	}

	data := PostData {
		Users: *user,
		Posts: *post,
		Connected: 0,
	}

	c, err := r.Cookie("session_token")
    
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

	session_user, err := forumRepository.GetUserByCookie(c.Value)
	if err != nil {
		fmt.Println(err)
	}

	data = PostData {
		Users: *user,
		SessionUser: *session_user,
		Posts: *post,
		Connected: 1,
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, err)
	}
}