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
	g.RequestHome(w, r, g.TemplateHome)
}

// artiste en particulier
func (g *Groupie) Artist(w http.ResponseWriter, r *http.Request) {
	// Fonction pour afficher le template de l'artiste spécifique
	g.RequestArtist(w, r, g.TemplateArtist)
}

//g.FilterArtists()

func (g *Groupie) RequestHome(w http.ResponseWriter, r *http.Request, html string) {
	// Analyser les formulaires et récupérer les paramètres de la requête
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Récupérer la valeur de "id" du formulaire
	//options := r.FormValue("name")
	id := r.FormValue("id")
	url := r.URL.Path
	var data interface{}
	var err error
	sortField := r.FormValue("sortField") // champ de tri
	sortDir := r.FormValue("sortDir")     // direction du tri

	// Définir les options de tri par défaut si non spécifiées
	if sortField == "" {
		sortField = "name" // tri par nom par défaut
	}
	if sortDir == "" {
		sortDir = "asc" // tri ascendant par défaut
	}

	options := SortOptions{
		Field:     sortField,
		Direction: sortDir,
	}
	// Si un "id" est fourni, récupérer les informations de l'artiste correspondant
	if id != "" {
		// Convertir l'id en entier
		artistID, err := strconv.Atoi(id)
		if err != nil {
			// Si l'id n'est pas valide, renvoyer une erreur HTTP
			fmt.Printf("Erreur de conversion de l'id en entier: %v\n", err)
			http.Error(w, "ID invalide", http.StatusBadRequest)
			return
		}

		// Récupérer les données de l'artiste avec l'ID converti
		data, err = g.GetArtist(artistID)
		if err != nil {
			// Si l'appel pour récupérer les données échoue, afficher une erreur
			fmt.Printf("Erreur de récupération des données pour l'artiste: %v\n", err)
			http.Error(w, "Erreur lors de la récupération des données de l'artiste", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Artist data : %v\n", data)
	} else if url == "/" {
		data, err = g.GetAllArtists()
		if err != nil {
			fmt.Printf("Erreur de récupération des données des artistes: %v\n", err)
			http.Error(w, "Erreur lors de la récupération des données des artistes", http.StatusInternalServerError)
			return
		}

		// Maintenant on peut utiliser options correctement
		if artists, ok := data.([]Artists); ok {
			data = g.SortArtists(artists, options)
		}
	}

	// Vérification d'erreur après récupération des données (ou après l'appel à GetAllArtists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parser le template après avoir les données
	tmpl, err := template.ParseFiles(html)
	if err != nil {
		// Si l'erreur de parsing se produit, renvoyer l'erreur
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Exécuter le template avec les données récupérées
	if err := tmpl.Execute(w, data); err != nil {
		// Si une erreur survient pendant l'exécution du template, renvoyer l'erreur
		fmt.Printf("Template execution error: %v\n", err)
	}
}

func (g *Groupie) RequestArtist(w http.ResponseWriter, r *http.Request, html string) {
	// Analyser les formulaires et récupérer les paramètres de la requête
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Récupérer la valeur de "id" du formulaire
	id := r.FormValue("id")
	url := r.URL.Path
	var data interface{}
	var err error

	// Si un "id" est fourni, récupérer les informations de l'artiste correspondant
	if id != "" {
		// Convertir l'id en entier
		artistID, err := strconv.Atoi(id)
		if err != nil {
			// Si l'id n'est pas valide, renvoyer une erreur HTTP
			fmt.Printf("Erreur de conversion de l'id en entier: %v\n", err)
			http.Error(w, "ID invalide", http.StatusBadRequest)
			return
		}

		// Récupérer les données de l'artiste avec l'ID converti
		data, err = g.GetArtist(artistID)
		fmt.Println(data)
		if err != nil {
			// Si l'appel pour récupérer les données échoue, afficher une erreur
			fmt.Printf("Erreur de récupération des données pour l'artiste: %v\n", err)
			http.Error(w, "Erreur lors de la récupération des données de l'artiste", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Artist data : %v\n", data)
	} else if url == "/" {
		// Si aucun "id" n'est spécifié, récupérer tous les artistes
		data, err = g.GetAllArtists()
		if err != nil {
			// Si l'appel pour récupérer tous les artistes échoue
			fmt.Printf("Erreur de récupération des données des artistes: %v\n", err)
			http.Error(w, "Erreur lors de la récupération des données des artistes", http.StatusInternalServerError)
			return
		}
		//fmt.Printf("Artists data__________________________________________________: %+v\n", data)
	}

	// Vérification d'erreur après récupération des données (ou après l'appel à GetAllArtists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parser le template après avoir les données
	tmpl, err := template.ParseFiles(html)
	if err != nil {
		// Si l'erreur de parsing se produit, renvoyer l'erreur
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Exécuter le template avec les données récupérées
	if err := tmpl.Execute(w, data); err != nil {
		// Si une erreur survient pendant l'exécution du template, renvoyer l'erreur
		fmt.Printf("Template execution error: %v\n", err)
	}
}

/*


	artists, err := g.GetArtists()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Tri des artistes par défaut par nom
	sortedArtists := g.SortArtists(artists, SortOptions{
		Field:     "name",
		Direction: "asc",
	})
	tmpl, err := template.ParseFiles("Home.html")
	if err != nil {
		// Si l'erreur de parsing se produit, renvoyer l'erreur
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Affiche le template avec les artistes triés
	err = tmpl.ExecuteTemplate(w, "Home.html", sortedArtists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}






*/

func (g *Groupie) handleHome(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query != "" {
		// Si une requête est présente, afficher le terme recherché
		fmt.Fprintf(w, "You searched for: %s", query)
	} else {
		// Si aucune requête n'est envoyée, afficher un message
		fmt.Fprintf(w, "No search query provided.")
	}
	artists, err := g.GetArtist()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Tri des artistes par défaut par nom
	sortedArtists := g.SortArtists(artists, SortOptions{
		Field:     "name",
		Direction: "asc",
	})
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		// Si l'erreur de parsing se produit, renvoyer l'erreur
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Affiche le template avec les artistes triés
	err = tmpl.ExecuteTemplate(w, "index.html", sortedArtists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
