package groupietracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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

// GetArtists récupère et retourne un artiste spécifique ou tous les artistes
func (g *Groupie) GetArtist(id ...int) ([]Artists, error) {
	// Récupère tous les artistes
	artists, err := g.GetAllArtists()
	if err != nil {
		return nil, err
	}

	// Si un ID est spécifié, retourne uniquement cet artiste
	if len(id) > 0 {
		//fmt.Println(artists)
		for _, artist := range artists {
			//fmt.Println(artist)
			if artist.Id == id[0] {
				//fmt.Println(artist.Id)
				return []Artists{artist}, nil
			}
		}
		return nil, nil
	}
	return artists, nil
}

// LoadArtistDetails charge les détails complets d'un artiste (incluant les relations)
func (g *Groupie) LoadArtistDetails(id int) (*Artists, error) {
	// Récupère l'artiste de base
	artists, err := g.GetArtist(id)
	if err != nil || len(artists) == 0 {
		return nil, err
	}
	artist := artists[0]

	// Récupère les relations
	relations, err := g.GetRelations(id)
	if err != nil {
		return nil, err
	}

	// Ajoute les relations à l'artiste
	artist.DatesLocations = relations.DatesLocations

	return &artist, nil
}

// GetRelations récupère les relations (dates et lieux) pour un artiste
func (g *Groupie) GetRelations(id int) (*Relations, error) {
	url := fmt.Sprintf("http://groupietrackers.herokuapp.com/api/relation/%d", id)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erreur avec le code de statut %d", response.StatusCode)
	}

	var relations Relations
	err = json.NewDecoder(response.Body).Decode(&relations)
	if err != nil {
		return nil, err
	}

	return &relations, nil
}

// SortArtists trie la liste des artistes selon les options spécifiées
func (g *Groupie) FilterArtists(filterOptions FilterOptions, searchTerm string) ([]Artists, error) {
	// Récupérer tous les artistes depuis l'API
	artists, err := g.GetAllArtists()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des artistes : %v", err)
	}

	// Étape 1 : Appliquer les filtres
	filteredArtists := make([]Artists, 0)

	for _, artist := range artists {
		// Vérification du terme de recherche (sur nom ou membres)
		if searchTerm != "" {
			termMatch := strings.Contains(strings.ToLower(artist.Name), strings.ToLower(searchTerm))
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), strings.ToLower(searchTerm)) {
					termMatch = true
					break
				}
			}
			if !termMatch {
				continue
			}
		}

		// Filtrer par date de création
		if artist.CreationDate < filterOptions.CreationDate || artist.CreationDate > filterOptions.CreationDate {
			continue
		}

		// Filtrer par date du premier album
		firstAlbumYear := 1963
		fmt.Sscanf(artist.FirstAlbum, "%d", &firstAlbumYear)
		if firstAlbumYear < filterOptions.FirstAlbum || firstAlbumYear > filterOptions.FirstAlbum {
			continue
		}

		// Filtrer par nombre de membres
		memberCount := len(artist.Members)
		if memberCount < filterOptions.MemberCount || memberCount > filterOptions.MemberCount {
			continue
		}

		// Filtrer par lieux de concerts
		if len(filterOptions.Locations) > 0 {
			locationMatch := false
			for _, location := range filterOptions.Locations {
				if _, exists := artist.DatesLocations[location]; exists {
					locationMatch = true
					break
				}
			}
			if !locationMatch {
				continue
			}
		}

		// Ajouter l'artiste s'il satisfait à tous les filtres
		filteredArtists = append(filteredArtists, artist)
	}

	return filteredArtists, nil
}

// FilterArtists filtre les artistes selon différents critères
func (g *Groupie) SearchArtists(artists []Artists, filters map[string]string) []Artists {
	var results []Artists

	for _, artist := range artists {
		match := true

		// Recherche par nom de l'artiste
		if name, ok := filters["name"]; ok && name != "" {
			if !strings.Contains(strings.ToLower(artist.Name), strings.ToLower(name)) {
				match = false
			}
		}

		// Recherche par membres
		if member, ok := filters["member"]; ok && member != "" {
			memberMatch := false
			for _, m := range artist.Members {
				if strings.Contains(strings.ToLower(m), strings.ToLower(member)) {
					memberMatch = true
					break
				}
			}
			if !memberMatch {
				match = false
			}
		}

		// Recherche par emplacements
		if location, ok := filters["location"]; ok && location != "" {
			locationMatch := false
			for loc := range artist.DatesLocations {
				if strings.Contains(strings.ToLower(loc), strings.ToLower(location)) {
					locationMatch = true
					break
				}
			}
			if !locationMatch {
				match = false
			}
		}

		// Recherche par date de création
		if creationDate, ok := filters["creationDate"]; ok && creationDate != "" {
			year := 0
			fmt.Sscanf(creationDate, "%d", &year)
			if year != 0 && artist.CreationDate != year {
				match = false
			}
		}

		// Recherche par date du premier album
		if firstAlbum, ok := filters["firstAlbum"]; ok && firstAlbum != "" {
			if !strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(firstAlbum)) {
				match = false
			}
		}

		if match {
			results = append(results, artist)
		}
	}

	return results
}

func (g *Groupie) GetArtistIDByName(groupName string) int {
	// Parcourt la liste des artistes
	artists, _ := g.GetAllArtists()

	id, _ := strconv.Atoi(groupName)
	if id == 0 {
		for _, artist := range artists {
			// Si le nom de l'artiste correspond au nom recherché
			if artist.Name == groupName {
				// Retourne l'ID de l'artiste trouvé
				return artist.Id
			}
		}
	}
	return id
}

func (g *Groupie) isInt(value interface{}) bool {
	switch value.(type) {
	case int:
		return true
	default:
		return false
	}
}

func (g *Groupie) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
