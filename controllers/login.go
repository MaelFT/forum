package forum

import (
	"net/http"
	"html/template"
	"database/sql"
	"fmt"
	"time"
)

type LoginData struct {
	Error string
}

func Login(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
    
	if err != nil || c.Value == "" {
		fmt.Println(c, err)
	} else {
		http.Redirect(w, r, "/index", http.StatusFound)
	}

	tmpl := template.Must(template.ParseFiles("./views/login.html")) // Affiche la page

	// Affiche dans le terminal l'activité sur le site
	switch r.Method {
	case "GET":
		fmt.Println("GET")
	case "POST": // Gestion d'erreur
		if err := r.ParseForm(); err != nil {
			return
		}
	}

	data := LoginData {}

	db, err := sql.Open("sqlite3", "forum.db")
    if err != nil {
        fmt.Println(err)
    }

	forumRepository := NewSQLiteRepository(db)

	uname := r.Form.Get("uname")
	pwd := r.Form.Get("pwd")

	fmt.Println(uname, pwd)

	user, err := forumRepository.CheckUser(uname, pwd)
	
	fmt.Println(user)

	if err != nil {
		data = LoginData {
			Error: "Login failed",
		}
		fmt.Println("Login failed")
    } else {
		data = LoginData {
			Error: "Logged",
		}
		fmt.Println("Logged")

		// Cookies
		http.SetCookie(w, &http.Cookie{
			Name:       "session_token",
			Value:      user.Cookie,
			Expires:    time.Now().Add(3600 * time.Second),
		})

		http.Redirect(w, r, "/index", http.StatusFound)
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, err)
	}
}