package groupietracker

import (
	"fmt"
	"html/template"
	"log"
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
	log.Printf("Début Request : %s %s", r.Method, r.URL.Path)

	defer log.Println("Fin Request")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("id")
	fmt.Println("caca", name)
	var id int
	if !g.isInt(name) {
		id = g.GetArtistIDByName(name)

	}
	url := r.URL.Path
	var data interface{}
	fmt.Println("popo", data)

	// Extraction des options de filtre
	creation, err := strconv.Atoi(r.FormValue("creationDate"))
	if err != nil && r.FormValue("creationDate") != "" {
		http.Error(w, "Invalida 'creationDate' parameter", http.StatusBadRequest)
		return
	}

	member, err := strconv.Atoi(r.FormValue("memberCount"))
	if err != nil && r.FormValue("memberCount") != "" {
		http.Error(w, "Invalidb 'memberCount' parameter", http.StatusBadRequest)
		return
	}
	firstAlbum, err := strconv.Atoi(r.FormValue("firstAlbum"))
	if err != nil && r.FormValue("firstAlbum") != "" {
		http.Error(w, "Invalidc 'firstAlbum' parameter", http.StatusBadRequest)
		return
	}

	locations := r.FormValue("locations")

	fmt.Printf(" falbum: %d", firstAlbum)
	fmt.Printf(" member: %d", member)
	fmt.Println("location :  ", locations)

	filterOptions := FilterOptions{
		CreationDate: creation,
		FirstAlbum:   firstAlbum,
		MemberCount:  member,
		Locations:    locations,
	}

	// Gestion d'un artiste spécifique
	if id < 53 && id > 0 {
		fmt.Println("id?")
		artistID := id

		artistData, err := g.LoadArtistDetails(artistID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("voui", artistData)

		data = artistData
		fmt.Println("nana", data)
	} else if url == "/" {
		fmt.Println("i2?")
		// Si aucun filtre n'est appliqué, afficher tous les artistes
		if creation == 0 && firstAlbum == 0 && member == 0 && len(locations) == 0 {
			// Pas de filtre actif, on charge tous les artistes
			artists, err := g.GetAllArtists()
			if err != nil {
				http.Error(w, "Erreur de chargement des artistes : "+err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println(artists)
			data = artists

		} else {
			fmt.Println("i4?")
			// Appliquer les filtres

			searchTerm := r.FormValue("search")

			fmt.Printf("searchterm : %s\n", searchTerm)
			artists, err := g.FilterArtists(filterOptions)
			g.SearchArtists(searchTerm)
			//g.SearchArtists(artists, filterOptions)
			fmt.Println("i5?")
			fmt.Println(" filterartistes ce que ça renvoie : ", artists)

			if err != nil {
				http.Error(w, "Erreur de recherche : "+err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("i6?")
			data = artists
		}
	}
	log.Printf("Requête reçue : %s, ID : %d", r.URL.Path, id)
	log.Printf("Options de filtre : %+v", filterOptions)

	fmt.Println("voici toute la data ", data)
	tmpl, err := template.ParseFiles(html)
	fmt.Println("voici le template ", html)
	if err != nil {
		http.Error(w, "Erreur de parsing du template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		fmt.Printf("Erreur d'exécution du template: %v\n", err)
		http.Error(w, "Erreur d'exécution du template", http.StatusInternalServerError)
	}
	log.Printf("Method: %s, Path: %s", r.Method, r.URL.Path)
}
