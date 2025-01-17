package groupietracker

import (
	"fmt"
)

func (g *Groupie) LocationApi(relation *Relations) error {
	if relation == nil {
		return fmt.Errorf("relation cannot be nil")
	}

	var urls []string // Pour stocker toutes les URLs
	for city := range relation.DatesLocations {
		fmt.Printf("Ville trouvée : %s\n", city)

		// Construire une requête vers l'API de géolocalisation
		url := "https://nominatim.openstreetmap.org/ui/search.html?q=" + city
		urls = append(urls, url)
	}

	// Concaténer toutes les URLs dans une seule chaîne (si nécessaire)
	g.apimaps = urls
	return nil
}
