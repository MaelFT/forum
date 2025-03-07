package forum

import (
	"net/http"
	"html/template"
	"fmt"
	models "forum/models"
)

type PageNotFoundData struct {
	Users models.Users
	Connected int
	Error string
}

func PageNotFound(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	tmpl := template.Must(template.ParseFiles("./views/404.html")) // Affiche la page
	data := PageNotFoundData {
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

	data = PageNotFoundData {
		Connected: 1,
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, err)
	}
}