package main

import(
	"net/http"
	"html/template"
	"fmt"
)

func main() {
	fs := http.FileServer(http.Dir("./static/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}

type data struct {
	Test []int
}

var frontData data 

func Handler(w http.ResponseWriter, r *http.Request) {

	frontData.Test = []int{0,1,2,3}

	tmpl := template.Must(template.ParseFiles("./static/login.html")) // Affiche la page

	// Affiche dans le terminal l'activité sur le site
	switch r.Method {
	case "GET":
		fmt.Println("GET")
	case "POST": // Gestion d'erreur
		if err := r.ParseForm(); err != nil {
			return
		}
	}

	uname := r.Form.Get("uname")
	pwd := r.Form.Get("pwd")

	fmt.Println(uname, pwd)

	tmpl.Execute(w, frontData)

}