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

func (g *Groupie) Request(w http.ResponseWriter, r *http.Request, html string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("id")

	var id int
	if !g.isInt(name) {
		id = g.GetArtistIDByName(name)

	}
	url := r.URL.Path
	var data interface{}
	var err error

	// Extraction des options de filtre
	creation, _ := strconv.Atoi(r.FormValue("creationDate"))
	firstAlbum, _ := strconv.Atoi(r.FormValue("firstAlbum"))
	member, _ := strconv.Atoi(r.FormValue("memberCount"))

	locations := r.FormValue("locations")

	fmt.Printf(" creation: %d", creation)
	fmt.Printf(" falbum: %d", firstAlbum)
	fmt.Printf(" member: %d", member)
	fmt.Printf("location %d", locations)

	filterOptions := FilterOptions{
		CreationDate: creation,
		FirstAlbum:   firstAlbum,
		MemberCount:  member,
		Locations:    locations,
	}

	// Gestion d'un artiste spécifique
	if id < 53 && id > 0 {
		artistID := id

		artistData, err := g.LoadArtistDetails(artistID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		/*
			relations := &Relations{
				ID:             artistData.Id,
				DatesLocations: artistData.DatesLocations,
			}

			if err = g.LocationApi(relations); err != nil {
				fmt.Printf("Erreur LocationApi: %v\n", err)
			}
		*/
		data = artistData
	} else if url == "/" {

		// Si aucun filtre n'est appliqué, afficher tous les artistes
		if creation == 0 && firstAlbum == 0 && member == 0 && len(locations) == 0 {
			// Pas de filtre actif, on charge tous les artistes
			artists, err := g.GetAllArtists()
			if err != nil {
				http.Error(w, "Erreur de chargement des artistes : "+err.Error(), http.StatusInternalServerError)
				return
			}
			data = artists

		} else {

			// Appliquer les filtres

			searchTerm := r.FormValue("search")

			fmt.Printf(" searchterm : %d", searchTerm)

			artists, err := g.FilterArtists(filterOptions, searchTerm)

			//g.SearchArtists(artists, filterOptions)

			fmt.Println(" filterartistes ce que ça renvoie : ", artists)

			if err != nil {
				http.Error(w, "Erreur de recherche : "+err.Error(), http.StatusInternalServerError)
				return
			}
			data = artists
		}
	}
	fmt.Println("voici toute la data ", data)
	var vide interface{}
	if data == vide {

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
