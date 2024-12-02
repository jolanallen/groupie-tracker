package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	fmt.Println("Serveur démarré sur http://localhost:3030")


	css := http.FileServer(http.Dir("./web/css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	js := http.FileServer(http.Dir("./web/js"))
	http.Handle("/js/", http.StripPrefix("/js/", js))

	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tmpl.Execute(w, nil)
	})

	
	http.ListenAndServe(":3030", nil)
}
