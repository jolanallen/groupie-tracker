package groupietracker


// import (
//     "encoding/json"
//     "fmt"
//     "io/ioutil"
// 	"net/http"
// 	"net/url"
// )

// type Location struct {
// 	Lat string `json:"lat"`
// 	Lon string `json:"lon"`
// }

// func (g *Groupie) LocationApi(relation *Relations) error {

// 	if relation == nil {
// 		return fmt.Errorf("relation cannot be nil")
// 	}

// 	var urls []string // Pour stocker toutes les URLs

// 	for city := range relation.DatesLocations {
// 		fmt.Printf("Ville trouvée : %s\n", city)
// 		params.Add("q", city)



// 		// Décoder la réponse JSON
// 		var locations []Location
// 		if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
// 			fmt.Errorf("erreur lors du décodage JSON : %v", err)
// 		}

// 		// Vérifier si des résultats ont été trouvés
// 		if len(locations) == 0 {
// 			fmt.Errorf("aucune coordonnée trouvée pour la ville : %s", city)
// 		}

// 		// Retourner les coordonnées

// 		// Construire une requête vers l'API de géolocalisation

// 		u := "https://www.openstreetmap.org/search?query=" + city + "#map=10/" + locations[0].Lat + "/" + locations[0].Lon
// 		urls = append(urls, url)
// 	}

// 	// Concaténer toutes les URLs dans une seule chaîne (si nécessaire)
// 	g.apimaps = bbox
// 	return nil
// }

