package forum

import (
	"net/http"
	"html/template"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"github.com/satori/go.uuid"
	models "forum/models"
)

type RegisterData struct {
	Error string
}

func Register(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
    
	if err != nil {
		fmt.Println(c, err)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	
	tmpl := template.Must(template.ParseFiles("./views/signup.html")) // Affiche la page

	// Affiche dans le terminal l'activité sur le site
	switch r.Method {
	case "GET":
		fmt.Println("GET")
	case "POST": // Gestion d'erreur
		if err := r.ParseForm(); err != nil {
			return
		}
	}

	data := RegisterData {}

	db, err := sql.Open("sqlite3", "forum.db")
    if err != nil {
        fmt.Println(err)
    }

	forumRepository := NewSQLiteRepository(db)

	uname := r.Form.Get("uname")
	mail := r.Form.Get("mail")
	pwd := r.Form.Get("pwd")
	cookie := uuid.NewV4().String()

	fmt.Println(uname, mail, pwd)

	users := models.Users {
		Username: uname,
		Mail: mail,
		Password: pwd,
		Cookie: cookie,
	}
	
	if len(uname) > 1 && len(mail) > 1 && len(pwd) > 1 {
		user, err := forumRepository.CreateUser(users)
		fmt.Println(user, "\n", err)
	} else {
		err = errors.New("len")
	}
	
	if err != nil {
		data = RegisterData {
			Error: "Register failed",
		}
		fmt.Println("Register failed")
    } else {
		data = RegisterData {
			Error: "Register sucess",
		}
		fmt.Println("Register sucess")
	
		// Cookies
		http.SetCookie(w, &http.Cookie{
			Name:       "session_token",
			Value:      cookie,
			Expires:    time.Now().Add(3600 * time.Second),
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, err)
	}
}