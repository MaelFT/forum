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
	Comments    models.Comments
	Connected   int
	Error       string
	// like ?
}

func Post(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./views/post.html")) // Affiche la page

	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	data := PostData{}

	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Println(err)
	}

	forumRepository := NewSQLiteRepository(db)

	post, err := forumRepository.GetPostByID(id)
	if err != nil {
		fmt.Println(err)
	}

	comments, err := forumRepository.GetCommentByPostID(post.ID)
	if err != nil {
		fmt.Println(err)
	}

	user, err := forumRepository.GetUserByID(post.User_ID)
	if err != nil {
		fmt.Println(err)
	}

	if comments != nil {
		data = PostData{
			Users:     *user,
			Posts:     *post,
			Comments:  *comments,
			Connected: 0,
		}
	} else {
		data = PostData{
			Users:     *user,
			Posts:     *post,
			Connected: 0,
		}
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
		if err := r.ParseForm(); err != nil {
			return
		}
	}

	commentContent := r.Form.Get("comment")
	fmt.Println("commentContent", commentContent)
	comment := models.Comments{
		Content:    commentContent,
		Post_ID:    post.ID,
		User_ID:    user.ID,
	}

	if len(commentContent) > 1 {
		forumRepository.CreateComment(comment)
	}

	session_user, err := forumRepository.GetUserByCookie(c.Value)
	if err != nil {
		fmt.Println(err)
	}

	if comments != nil {
		data = PostData{
			Users:       *user,
			SessionUser: *session_user,
			Posts:       *post,
			Comments:    *comments,
			Connected:   1,
		}
	} else {
		data = PostData{
			Users:       *user,
			SessionUser: *session_user,
			Posts:       *post,
			Connected:   1,
		}
	}


	err = tmpl.Execute(w, data)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, err)
	}
}