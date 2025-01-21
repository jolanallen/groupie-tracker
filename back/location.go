package groupietracker

/*
import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Structure pour représenter une localisation retournée par Nominatim
type Location struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

// Fonction pour appeler l'API de géolocalisation et générer les URLs OpenStreetMap
func (g *Groupie) LocationApi(relation *Relations) error {

	// Vérification que la relation n'est pas nulle
	if relation == nil {
		return fmt.Errorf("relation cannot be nil")
	}

	var urls []string // Pour stocker les URLs OpenStreetMap

	// Parcourir les villes des relations
	for city := range relation.DatesLocations {
		fmt.Printf("Ville trouvée : %s\n", city)

		// Construire les paramètres de la requête
		apiURL := "https://nominatim.openstreetmap.org/search"
		params := url.Values{}
		params.Add("q", city)
		params.Add("format", "json")

		// Effectuer la requête GET
		resp, err := http.Get(apiURL + "?" + params.Encode())
		if err != nil {
			return fmt.Errorf("erreur lors de la requête vers l'API : %v", err)
		}
		defer resp.Body.Close()

		// Décoder la réponse JSON
		var locations []Location
		if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
			return fmt.Errorf("erreur lors du décodage JSON : %v", err)
		}

		// Vérifier si des résultats ont été trouvés
		if len(locations) == 0 {
			fmt.Printf("Aucune coordonnée trouvée pour la ville : %s\n", city)
			continue
		}

		// Construire l'URL OpenStreetMap
		url := fmt.Sprintf("https://www.openstreetmap.org/#map=10/%s/%s", locations[0].Lat, locations[0].Lon)
		urls = append(urls, url)
	}

	// Stocker les URLs dans le champ apimaps de Groupie
	g.apimaps = urls
	return nil
}
*/
