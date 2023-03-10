package forum

import (
	"database/sql"
	"fmt"
	models "forum/models"
	"html/template"
	"net/http"
	"strconv"
)

type PostData struct {
	Users       models.Users
	SessionUser models.Users
	Posts       models.Posts
	Connected   int
	Error       string
	// like ?
}

func Post(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
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

	data := PostData{
		Users:     *user,
		Posts:     *post,
		Connected: 0,
	}

	c, err := r.Cookie("session_token")

	if err != nil || c.Value == "" {
		fmt.Println("ici : ", c, err)
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
	case "POST":
		c, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println(c, err)
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			if c.Value == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
			}

			db, err := sql.Open("sqlite3", "forum.db")
			if err != nil {
				fmt.Println(err)
			}
			forumRepository := NewSQLiteRepository(db)

			post_id, _ := strconv.Atoi(r.FormValue("post_id"))

			user, _ := forumRepository.GetUserByCookie(c.Value)

			likeValue, _ := strconv.Atoi(r.FormValue("value"))

			like := models.Like{
				Value:   likeValue,
				User_ID: user.ID,
				Post_ID: int64(post_id),
			}
			Like(like, likeValue)
			http.Redirect(w, r, "/post?id="+strconv.Itoa(post_id), 302)
		}
		return
	}
	session_user, err := forumRepository.GetUserByCookie(c.Value)
	if err != nil {
		fmt.Println(err)
	}

	data = PostData{
		Users:       *user,
		SessionUser: *session_user,
		Posts:       *post,
		Connected:   1,
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, err)
	}

	like := r.Form.Get("like")

	fmt.Println(like)
}
