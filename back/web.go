package groupietracker

import (
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
		http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusBadRequest)
		return
	}

	// Récupérer les valeurs du formulaire
	name := r.FormValue("id")           // Pour chercher par nom ou ID
	searchTerm := r.FormValue("search") // Terme de recherche
	creation, _ := strconv.Atoi(r.FormValue("creationDate"))
	member, _ := strconv.Atoi(r.FormValue("memberCount"))
	firstAlbum, _ := strconv.Atoi(r.FormValue("firstAlbum"))
	locations := r.FormValue("locations")

	filterOptions := FilterOptions{
		CreationDate: creation,
		FirstAlbum:   firstAlbum,
		MemberCount:  member,
		Locations:    locations,
	}

	var data interface{}

	// Gestion d'un artiste spécifique (par ID ou nom)
	if name != "" {
		var id int
		if g.isInt(name) { // Si c'est un ID
			id, _ = strconv.Atoi(name)
		} else { // Si c'est un nom
			id = g.GetArtistIDByName(name)
		}

		if id > 0 {
			artistData, err := g.LoadArtistDetails(id)
			if err != nil {
				http.Error(w, "Erreur lors du chargement de l'artiste", http.StatusInternalServerError)
				return
			}
			data = artistData
		}
	} else if searchTerm != "" { // Gestion de la recherche
		results := g.SearchArtists(searchTerm)
		if len(results) == 0 {
			log.Println("Aucun artiste trouvé pour la recherche :", searchTerm)
		}
		data = results
	} else if creation > 0 || member > 0 || firstAlbum > 0 || locations != "" { // Gestion des filtres
		artists, err := g.FilterArtists(filterOptions)
		if err != nil {
			http.Error(w, "Erreur lors de l'application des filtres", http.StatusInternalServerError)
			return
		}
		data = artists
	} else { // Afficher tous les artistes si aucune recherche ni filtre
		artists, err := g.GetAllArtists()
		if err != nil {
			http.Error(w, "Erreur lors du chargement des artistes", http.StatusInternalServerError)
			return
		}
		data = artists
	}

	// Parsing et rendu du template
	tmpl, err := template.ParseFiles(html)
	if err != nil {
		http.Error(w, "Erreur de parsing du template : "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
	}
	log.Printf("Fin de la requête : %s, données envoyées : %+v", r.URL.Path, data)
}
