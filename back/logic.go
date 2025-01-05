package groupietracker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (g *Groupie) GetArtistById(id string) (Artists, error) {

	var artist Artists

	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		return artist, fmt.Errorf("erreur lors de la récupération de l'artiste: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return artist, fmt.Errorf("échec de la récupération de l'artiste, statut HTTP: %v", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&artist)
	if err != nil {
		return artist, fmt.Errorf("erreur lors du décodage de la réponse: %v", err)
	}

	fmt.Println("ID: ", artist.Id)
	fmt.Println("Name: ", artist.Name)
	fmt.Println("Image: ", artist.Image)
	fmt.Println("Members: ", artist.Members)
	fmt.Println("Creation Date: ", artist.CreationDate)
	fmt.Println("First Album: ", artist.FirstAlbum)
	fmt.Println("Relations: ", artist.Relations)
	g.RequestRelation(a.Relations)

	return artist, nil
}

func (g *Groupie) RequestRelation(relationpath string) (map[string][]string, error) {
	var relations Relations

	if relationpath == "" {
		return nil, fmt.Errorf("le chemin des relations est vide")
	}

	resp, err := http.Get(relationpath)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des relations: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("échec de la récupération des relations, statut HTTP: %v", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&relations)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de la réponse des relations: %v", err)
	}

	if len(relations.DatesLocations) == 0 {
		fmt.Println("Aucune relation trouvée")
	}

	fmt.Println("Dates et localisations de concerts:", relations.DatesLocations)

	return relations.DatesLocations, nil
}
