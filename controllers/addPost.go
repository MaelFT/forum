package forum

import (
	"net/http"
	"html/template"
	"fmt"
	models "forum/models"
)

type AddPostData struct {
	Users models.Users
	Connected int
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

		title := r.Form.Get("title")
		category := r.Form.Get("category")
		content := r.Form.Get("content")

		fmt.Println(title, category, content)
	
		data = AddPostData {
			Connected: 1,
		}

		err = tmpl.Execute(w, data)

		if err != nil {
			fmt.Println(err)
			fmt.Fprintln(w, err)
		}
	}
}