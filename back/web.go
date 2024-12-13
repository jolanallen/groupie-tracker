package groupietracker

import (
	"fmt"
	"html/template"
	"net/http"
)
func (a *Artists) Web() {
	fmt.Println("Serveur démarré sur http://localhost:3031")

	css := http.FileServer(http.Dir("./web/css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))

	js := http.FileServer(http.Dir("./web/js"))
	http.Handle("/js/", http.StripPrefix("/js/", js))

	images := http.FileServer(http.Dir("./web/utiles"))
	http.Handle("/utiles/", http.StripPrefix("/utiles/", images))

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		fmt.Println("Erreur de chargement du template:", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	err = http.ListenAndServe(":3031", nil)
	if err != nil {
		fmt.Println("Erreur lors du démarrage du serveur:", err)
		return
	}
}