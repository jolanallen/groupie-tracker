package groupietracker

import (
	"fmt"
	"html/template"
	"net/http"
)

func (g *Groupie) Home(w http.ResponseWriter, r *http.Request) { // fonction pour afficher les differents templates html
	g.Request(w, r, g.TemplateHome)

}
func (g *Groupie) Artist(w http.ResponseWriter, r *http.Request) { // fonction pour afficher les differents templates html
	g.Request(w, r, g.TemplateArtist) // template dans les quel on injectera les donner récupérer dans l'API

}
func (g *Groupie) Apropos(w http.ResponseWriter, r *http.Request) { // fonction pour afficher les differents templates html
	g.Request(w, r, g.TemplateApropos)

}

func (g *Groupie) Request(w http.ResponseWriter, r *http.Request, html string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Debug: Afficher l'URL
	fmt.Printf("URL Path__________________________________________________: %s\n", r.URL.Path)

	// Récupérer les données avant de parser le template
	var data interface{}
	var err error

	if id := r.FormValue("id"); id != "" {
		data, err = g.GetArtistById(id)
		fmt.Printf("Artist data______!!!!!!!!!!!!!!!!!!!!!!____________________: %v\n", data)
	} else {
		data, err = g.GetAllArtists()
		// Debug: Afficher les données récupérées
		fmt.Printf("Artists data__________________________________________________: %+v\n", data)


// Pourquoi on appel cette ligne une deuxieme fois quoi qu'il arrive.?

	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parser le template après avoir les données
	tmpl, err := template.ParseFiles(html)
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Exécuter le template avec les données
	if err := tmpl.Execute(w, data); err != nil {
		fmt.Printf("Template execution error: %v\n", err)
	}
}
