package groupietracker

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// home , index.html
func (g *Groupie) Home(w http.ResponseWriter, r *http.Request) {
	// Fonction pour afficher les différents templates HTML (Page d'accueil)
	g.Request(w, r, g.TemplateHome)
}

// artiste en particulier
func (g *Groupie) Artist(w http.ResponseWriter, r *http.Request) {
	// Fonction pour afficher le template de l'artiste spécifique
	g.Request(w, r, g.TemplateArtist)
}

// g.FilterArtists()
func (g *Groupie) Request(w http.ResponseWriter, r *http.Request, html string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	url := r.URL.Path
	var data interface{} // Changé pour supporter à la fois un Artist et []Artists
	var err error

	sortField := r.FormValue("sortField")
	sortDir := r.FormValue("sortDir")

	if sortField == "" {
		sortField = "name"
	}
	if sortDir == "" {
		sortDir = "asc"
	}

	options := SortOptions{
		Field:     sortField,
		Direction: sortDir,
	}

	// Gestion d'un artiste spécifique
	if id != "" {
		artistID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "ID invalide", http.StatusBadRequest)
			return
		}

		artistData, err := g.LoadArtistDetails(artistID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		relations := &Relations{
			ID:             artistData.Id,
			DatesLocations: artistData.DatesLocations,
		}

		if err = g.LocationApi(relations); err != nil {
			fmt.Printf("Erreur LocationApi: %v\n", err)
		}

		data = artistData

	} else if url == "/" {
		// Page d'accueil - liste de tous les artistes
		artists, err := g.GetAllArtists()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Tri des artistes
		data = g.SortArtists(artists, options)
	}

	// Vérification finale des données
	if data == nil {
		http.Error(w, "Aucune donnée trouvée", http.StatusNotFound)
		return
	}

	// Parsing et exécution du template
	tmpl, err := template.ParseFiles(html)
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}
