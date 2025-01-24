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
			if artist.CreationDate != filterOptions.CreationDate {
				continue // Ignorer les artistes avec une date de création inférieure
			}
		}

		// Filtre sur FirstAlbum
		if filterOptions.FirstAlbum != 0 {
			// Extraire les 4 derniers caractères de FirstAlbum
			albumYearStr := artist.FirstAlbum[len(artist.FirstAlbum)-4:] // Prend les 4 derniers caractères
			var albumYear int
			_, err := fmt.Sscanf(albumYearStr, "%d", &albumYear)
			fmt.Println("c'est filtre", filterOptions.FirstAlbum)
			if err != nil {
				return nil, fmt.Errorf("Erreur lors de la lecture de FirstAlbum pour l'artiste %s: %v", artist.Name, err)
			}
			if albumYear != filterOptions.FirstAlbum {
				continue // Ignorer les artistes avec une année de premier album différente
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
			towns := g.GetSingleTownFilter(artist.DatesLocations)

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

/*   Get Single Town
1) affiche sous forme de slice de string le noms des villes

*/

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
//(request)
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
		towns := g.GetSingleTown(location)    // Appel à GetSingleTown pour chaque clé
		allTowns = append(allTowns, towns...) // Ajouter les villes trouvées à la liste
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
