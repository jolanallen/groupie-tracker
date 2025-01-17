package groupietracker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func (g *Groupie) Init() {
	// Initialise les chemins des templates
	g.TemplateHome = "front/templates/index.html"
	g.TemplateArtist = "front/templates/artists.html"
	g.TemplateApropos = "front/templates/Apropos.html"

	// Valider les chemins des fichiers de template
	templates := []string{g.TemplateHome, g.TemplateArtist, g.TemplateApropos}
	for _, template := range templates {
		if !fileExists(template) {
			log.Printf("Template manquant : %s", template)
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (g *Groupie) GetAllArtists() ([]Artists, error) {
	const apiBaseURL = "http://groupietrackers.herokuapp.com/api/"
	url := fmt.Sprintf("%sartists", apiBaseURL)

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des artistes : %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur avec le code de statut %d", response.StatusCode)
	}

	var artistsList []Artists
	err = json.NewDecoder(response.Body).Decode(&artistsList)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de la réponse : %v", err)
	}

	// Validation des données récupérées
	for _, artist := range artistsList {
		if artist.Id == 0 || artist.Name == "" {
			return nil, fmt.Errorf("artiste invalide dans la réponse API : %+v", artist)
		}
	}

	return artistsList, nil
}
