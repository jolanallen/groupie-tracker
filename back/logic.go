package groupietracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
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

	for _, artist := range artistsList {
		if artist.Id == 0 || artist.Name == "" {
			return nil, fmt.Errorf("artiste invalide dans la réponse API : %+v", artist)
		}
	}

	return artistsList, nil
}

// GetArtists récupère et retourne un artiste spécifique ou tous les artistes
func (g *Groupie) GetArtist(id ...int) ([]Artists, error) {

	artists, err := g.GetAllArtists()
	if err != nil {
		return nil, err
	}

	if len(id) > 0 {
		for _, artist := range artists {
			if artist.Id == id[0] {
				return []Artists{artist}, nil
			}
		}
		return nil, nil
	}
	return artists, nil
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


// LoadArtistDetails charge les détails complets d'un artiste (incluant les relations)
func (g *Groupie) LoadArtistDetails(id int) (*Artists, error) {
	artists, err := g.GetArtist(id)
	if err != nil || len(artists) == 0 {
		return nil, err
	}
	artist := artists[0]

	relations, err := g.GetRelations(id)
	if err != nil {
		return nil, err
	}

	artist.DatesLocations = relations.DatesLocations

	return &artist, nil
}

func (g *Groupie) FilterArtists(filterOptions FilterOptions) ([]Artists, error) {
	artists, _ := g.GetAllArtists()

	var filteredArtists []Artists

	// Appliquer les filtres sur les artistes
	for _, artist := range artists {
		// Filtre sur CreationDate
		if filterOptions.CreationDate != 0 {
			if artist.CreationDate != filterOptions.CreationDate {
				continue // Ignorer les artistes avec une date de création inférieure
			}
		}

		if filterOptions.FirstAlbum != 0 {
			albumYearStr := artist.FirstAlbum[len(artist.FirstAlbum)-4:]
			var albumYear int
			_, err := fmt.Sscanf(albumYearStr, "%d", &albumYear)
			fmt.Println("c'est filtre", filterOptions.FirstAlbum)
			if err != nil {
				return nil, fmt.Errorf("Erreur lors de la lecture de FirstAlbum pour l'artiste %s: %v", artist.Name, err)
			}
			if albumYear != filterOptions.FirstAlbum {
				continue
			}
		}

		// Filtre sur MemberCount
		if filterOptions.MemberCount > 0 {
			if len(artist.Members) != filterOptions.MemberCount {
				continue
			}
		}

		// Filtre sur Locations
		if filterOptions.Locations != "" {
			towns := g.GetSingleTownFilter(artist.DatesLocations)
			found := false
			for _, town := range towns {
				if town == filterOptions.Locations {
					found = true
					break
				}
			}

			if !found {
				continue
			}
		}

		filteredArtists = append(filteredArtists, artist)
	}

	// Tri selon les options spécifiées
	switch {
	case filterOptions.CreationDate != 0:
		sort.Slice(filteredArtists, func(i, j int) bool {
			return filteredArtists[i].CreationDate < filteredArtists[j].CreationDate
		})
	case filterOptions.FirstAlbum != 0:
		sort.Slice(filteredArtists, func(i, j int) bool {
			var yearI, yearJ int
			fmt.Sscanf(filteredArtists[i].FirstAlbum[len(filteredArtists[i].FirstAlbum)-4:], "%d", &yearI)
			fmt.Sscanf(filteredArtists[j].FirstAlbum[len(filteredArtists[j].FirstAlbum)-4:], "%d", &yearJ)
			return yearI < yearJ
		})
	case filterOptions.MemberCount != 0:
		sort.Slice(filteredArtists, func(i, j int) bool {
			return len(filteredArtists[i].Members) < len(filteredArtists[j].Members)
		})
	case filterOptions.Locations != "":
		sort.Slice(filteredArtists, func(i, j int) bool {
			locI := len(filteredArtists[i].DatesLocations)
			locJ := len(filteredArtists[j].DatesLocations)
			return locI < locJ
		})
	}

	return filteredArtists, nil
}

func (g *Groupie) GetArtistIDByName(groupName string) int {
	artists, _ := g.GetAllArtists()
	id, _ := strconv.Atoi(groupName)
	if id == 0 {
		for _, artist := range artists {
			if artist.Name == groupName {
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

func (g *Groupie) SearchArtists(searchTerm string) []Artists {
	artists, _ := g.GetAllArtists()
	var results []Artists

	// Met le searchTerm en minuscule pour assurer l'insensibilité à la casse
	searchTerm = strings.ToLower(searchTerm)

	for _, artist := range artists {
		// Vérifie si le searchTerm correspond à l'un des champs
		if strings.Contains(strings.ToLower(artist.Name), searchTerm) ||
			g.containsInSliceInsensitive(artist.Members, searchTerm) ||
			g.containsInMapKeysInsensitive(artist.DatesLocations, searchTerm) ||
			strings.Contains(strings.ToLower(g.GetLastFourChars(artist.FirstAlbum)), searchTerm) ||
			g.creationDateMatches(artist.CreationDate, searchTerm) {
			results = append(results, artist)
		}
	}

	return results
}

func (g *Groupie) GetSingleTownFilter(datesLocations map[string][]string) []string {
	var allTowns []string

	// Parcourir chaque entrée de la map et obtenir les villes correspondantes
	for location := range datesLocations {
		towns := g.GetSingleTown(location)
		allTowns = append(allTowns, towns...) 
	}

	// Éliminer les doublons en utilisant une map
	uniqueTowns := make(map[string]struct{})
	for _, town := range allTowns {
		uniqueTowns[town] = struct{}{}
	}

	// Convertir les clés de la map uniqueTowns en une slice
	var result []string
	for town := range uniqueTowns {
		result = append(result, town)
	}

	return result
}

// Fonction pour vérifier si un terme est présent dans un slice de chaînes (insensible à la casse)
func (g *Groupie) containsInSliceInsensitive(slice []string, searchTerm string) bool {
	for _, item := range slice {
		if strings.Contains(strings.ToLower(item), searchTerm) {
			return true
		}
	}
	return false
}

// Recherche dans les clés de la map datesLocations
func (g *Groupie) containsInMapKeysInsensitive(data map[string][]string, searchTerm string) bool {
	// Vérifier chaque clé dans la map de manière insensible à la casse
	for key := range data {
		if strings.Contains(strings.ToLower(key), strings.ToLower(searchTerm)) {
			return true
		}
	}
	return false
}

// Fonction qui retourne les villes extraites de la clé des datesLocations
func (g *Groupie) GetSingleTown(location string) []string {
	var towns []string
	// Diviser la chaîne au niveau du tiret "-"
	parts := strings.Split(location, "-")
	for _, part := range parts {
		// Remplacer les underscores "_" par des espaces " "
		city := strings.ReplaceAll(part, "_", " ")
		towns = append(towns, city)
	}
	return towns
}

// Fonction qui extrait les 4 derniers caractères d'un album
func (g *Groupie) GetLastFourChars(albumDate string) string {
	if len(albumDate) >= 4 {
		return albumDate[len(albumDate)-4:]
	}
	return ""
}

// Fonction pour vérifier la date de création (si l'année correspond)
func (g *Groupie) creationDateMatches(creationDate int, searchTerm string) bool {
	// Convertir la création de la date en string pour faire la comparaison
	return strings.Contains(fmt.Sprintf("%d", creationDate), searchTerm)
}
