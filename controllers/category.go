package forum

import (
	"net/http"
	"html/template"
	"fmt"
	models "forum/models"
)

type CategoryData struct {
	Users models.Users
	Connected int
	Error string
}

func Category(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	tmpl := template.Must(template.ParseFiles("./views/category.html")) // Affiche la page
	data := CategoryData {
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

	data = CategoryData {
		Connected: 1,
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, err)
	}
}