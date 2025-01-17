package groupietracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// GetArtists récupère et retourne un artiste spécifique ou tous les artistes
func (g *Groupie) GetArtists(id ...int) ([]Artists, error) {
    // Récupère tous les artistes
    artists, err := g.GetAllArtists()
    if err != nil {
        return nil, err
    }

    // Si un ID est spécifié, retourne uniquement cet artiste
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

// SortArtists trie la liste des artistes selon les options spécifiées
func (g *Groupie) SortArtists(artists []Artists, options SortOptions) []Artists {
    // Crée une copie du slice pour ne pas modifier l'original
    sortedArtists := make([]Artists, len(artists))
    copy(sortedArtists, artists)

    // Applique le tri en fonction du champ spécifié
    switch options.Field {
    case "name":
        // Tri par nom
        sort.Slice(sortedArtists, func(i, j int) bool {
            if options.Direction == "desc" {
                return strings.ToLower(sortedArtists[i].Name) > strings.ToLower(sortedArtists[j].Name)
            }
            return strings.ToLower(sortedArtists[i].Name) < strings.ToLower(sortedArtists[j].Name)
        })

    case "creation":
        // Tri par date de création
        sort.Slice(sortedArtists, func(i, j int) bool {
            if options.Direction == "desc" {
                return sortedArtists[i].CreationDate > sortedArtists[j].CreationDate
            }
            return sortedArtists[i].CreationDate < sortedArtists[j].CreationDate
        })

    case "members":
        // Tri par nombre de membres
        sort.Slice(sortedArtists, func(i, j int) bool {
            if options.Direction == "desc" {
                return len(sortedArtists[i].Members) > len(sortedArtists[j].Members)
            }
            return len(sortedArtists[i].Members) < len(sortedArtists[j].Members)
        })
    }

    return sortedArtists
}

// FilterArtists filtre les artistes selon différents critères
func (g *Groupie) FilterArtists(artists []Artists, filters map[string]string) []Artists {
    var filtered []Artists

    for _, artist := range artists {
        include := true

        // Filtre par nom
        if name, ok := filters["name"]; ok && name != "" {
            if !strings.Contains(strings.ToLower(artist.Name), strings.ToLower(name)) {
                include = false
            }
        }

        // Filtre par année de création
        if year, ok := filters["year"]; ok && year != "" {
            yearInt := 0
            fmt.Sscanf(year, "%d", &yearInt)
            if yearInt != 0 && artist.CreationDate != yearInt {
                include = false
            }
        }

        if include {
            filtered = append(filtered, artist)
        }
    }

    return filtered
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
    // Récupère l'artiste de base
    artists, err := g.GetArtists(id)
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
