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

/*   Get ALL Artists
1) récupére tout les artists
2) les renvoies
*/

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

/*   Get Artists
1) appel Get All Artists pour récupe tout les artists
2) si il y a un id dans les parametre renvoie uniquement le bon artists

*/

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

/*   Get Relations
1) récupére les relations en fonction d'un id
2) renvoie les relations de l'id

*/

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

/*   Load Artist Details
1) récupére un artiste avec Get Artist
2) récupére les relations avec Get Relations
3) renvoie artist avec toute les info des deux API

*/

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

/*   Filter Artists
1) avons pour arguments des options de filtres
2) filtre suivant mes prérogative
3) renvoie uniquement la liste des artists filtrés

*/

func (g *Groupie) FilterArtists(filterOptions FilterOptions) ([]Artists, error) {
	artists, _ := g.GetAllArtists()

	var filteredArtists []Artists

	// Appliquer les filtres sur les artistes
	for _, artist := range artists {
		// Filtre sur CreationDate
		if filterOptions.CreationDate != 0 {
			if artist.CreationDate < filterOptions.CreationDate {
				continue // Ignorer les artistes avec une date de création inférieure
			}
		}

		// Filtre sur FirstAlbum
		if filterOptions.FirstAlbum != 0 {
			var albumYear int
			_, err := fmt.Sscanf(artist.FirstAlbum, "%d", &albumYear)
			if err != nil {
				return nil, fmt.Errorf("Erreur lors de la lecture de FirstAlbum pour l'artiste %s: %v", artist.Name, err)
			}
			if albumYear < filterOptions.FirstAlbum {
				continue // Ignorer les artistes avec une année de premier album inférieure
			}
		}

		// Filtre sur MemberCount
		if filterOptions.MemberCount > 0 {
			if len(artist.Members) != filterOptions.MemberCount {
				continue // Ignorer les artistes qui n'ont pas le nombre exact de membres
			}
		}

		// Filtre sur Locations
		if filterOptions.Locations != "" {
			// Obtenir les villes visitées par l'artiste
			towns := g.GetSingleTown(artist.DatesLocations)

			// Vérifier si l'artiste a visité la ville spécifiée
			found := false
			for _, town := range towns {
				if town == filterOptions.Locations {
					found = true
					break
				}
			}

			if !found {
				continue // Ignorer l'artiste si la ville n'est pas trouvée
			}
		}

		// Ajouter l'artiste filtré à la liste finale
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
			fmt.Sscanf(filteredArtists[i].FirstAlbum, "%d", &yearI)
			fmt.Sscanf(filteredArtists[j].FirstAlbum, "%d", &yearJ)
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

/*   Get Single Town
1) affiche sous forme de slice de string le noms des villes

*/

func (g *Groupie) GetSingleTown(datesLocations map[string][]string) []string {
	var towns []string
	for town := range datesLocations {
		towns = append(towns, g.GetCityName(town))
	}
	return towns
}

/*   Search Artists
1) avons pour arguments des options de filtres
2) filtre suivant
3)
4)
5)
6)
7)
8)
9)


*/

func (g *Groupie) SearchArtists(searchTerm string) []Artists {
	// Récupère tous les artistes
	artists, _ := g.GetAllArtists()
	var results []Artists

	// Met le searchTerm en minuscule pour assurer l'insensibilité à la casse
	searchTerm = strings.ToLower(searchTerm)

	for _, artist := range artists {
		// Vérifie si le searchTerm correspond à l'un des champs
		if strings.Contains(strings.ToLower(artist.Name), searchTerm) ||
			g.containsInSliceInsensitive(artist.Members, searchTerm) ||
			g.containsInMapKeysInsensitive(artist.DatesLocations, searchTerm) ||
			strings.Contains(strings.ToLower(g.GetFirstFourCharsDatesLocations(artist.FirstAlbum)), searchTerm) ||
			g.creationDateMatches(artist.CreationDate, searchTerm) {
			results = append(results, artist)
		}
	}

	return results
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

// Fonction pour vérifier si un terme est présent dans les clés d'une map (insensible à la casse)
func (g *Groupie) containsInMapKeysInsensitive(data map[string][]string, searchTerm string) bool {
	for key := range data {
		if strings.Contains(strings.ToLower(key), searchTerm) {
			return true
		}
	}
	return false
}

// Fonction pour vérifier si une date de création correspond au searchTerm
func (g *Groupie) creationDateMatches(creationDate int, searchTerm string) bool {
	// Convertir la date de création en chaîne pour comparaison
	creationDateStr := strconv.Itoa(creationDate)
	return strings.Contains(creationDateStr, searchTerm)
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

func (g *Groupie) GetFirstFourCharsDatesLocations(word string) string {
	return word[:4]
}

func (g *Groupie) GetCityName(location string) string {
	// Remplace les tirets par des espaces
	location = strings.ReplaceAll(location, "-", " ")

	// Sépare la chaîne pour récupérer la ville (partie avant le premier espace)
	parts := strings.Split(location, " ")
	if len(parts) > 0 {
		return parts[0] // Retourne la première partie (nom de la ville)
	}
	return "" // Si aucun espace trouvé, retourne une chaîne vide
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (g *Groupie) GetCreationYears(artists []Artists) []int {
	yearSet := make(map[int]struct{}) // Utilisation d'une map pour éliminer les doublons

	// Parcourir les artistes pour collecter les années de création
	for _, artist := range artists {
		yearSet[artist.CreationDate] = struct{}{}
	}

	// Convertir la map en slice
	var creationYears []int
	for year := range yearSet {
		creationYears = append(creationYears, year)
	}

	// Trier les années dans l'ordre croissant
	sort.Ints(creationYears)

	return creationYears
} ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

/*for _, artist := range artists {
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

	filteredArtists = append(filteredArtists, artist)
}
*/
