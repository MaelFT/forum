package forum

import (
	"database/sql"
	"fmt"
	models "forum/models"
	"html/template"
	"net/http"
)

type AddCategoryData struct {
	Users     models.Users
	Connected int
	Error     string
}

func AddCategory(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")

	if err != nil {
		fmt.Println(c, err)
		http.Redirect(w, r, "/index", http.StatusFound)
	} else {
		if c.Value == "" {
			http.Redirect(w, r, "/index", http.StatusFound)
		}
		tmpl := template.Must(template.ParseFiles("./views/addCategory.html"))

		switch r.Method {
		case "GET":
			fmt.Println("GET")
		case "POST": // Gestion d'erreur
			if err := r.ParseForm(); err != nil {
				return
			}
		}
		data := AddCategoryData{}

		db, err := sql.Open("sqlite3", "forum.db")
		if err != nil {
			fmt.Println(err)
		}

		forumRepository := NewSQLiteRepository(db)

		title := r.Form.Get("title")
		description := r.Form.Get("description")

		fmt.Println(title, description)

		user, err := forumRepository.GetUserByCookie(c.Value)

		category := models.Categories{
			Title:       title,
			Description: description,
			User_ID:     user.ID,
		}

		if len(title) > 1 && len(description) > 1 {
			forumRepository.CreateCategorie(category)
		}

		data = AddCategoryData{
			Connected: 1,
		}

		err = tmpl.Execute(w, data)

		if err != nil {
			fmt.Println(err)
			fmt.Fprintln(w, err)
		}
	}
}
