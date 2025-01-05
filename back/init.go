package groupietracker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var a Artists

func (g *Groupie) Init() {
	a.Image = ""
	a.Name = ""
	a.Members = []string{""}
	a.CreationDate = 0

	g.TemplateHome = "front/templates/index.html"
	g.TemplateArtist = "front/templates/artists.html"
	g.TemplateApropos = "front/templates/Apropos.html"

}

func (g *Groupie) GetAllArtists() ([]Artists, error) {
	// Appel à l'API ou à la base de données pour récupérer tous les artistes
	url := "http://groupietrackers.herokuapp.com/api/artists"
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
	//fmt.Printf("Type of data: %T\n", artistsList)
	//fmt.Printf("Content: %+v\n", artistsList)
	return artistsList, nil
}
