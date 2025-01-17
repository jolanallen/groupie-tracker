package groupietracker


type Groupie struct {
	Name      string
	Id        int
	Relations string
	TemplateHome string
	TemplateArtist string
	TemplateApropos string
}

type Artists struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Image       string   `json:"image"`
	Members     []string `json:"members"`
	CreationDate int     `json:"creationDate"`
	FirstAlbum  string   `json:"firstAlbum"` 
	Relations   string   `json:"relations"`    
	DatesLocations map[string][]string      `json:"datesLocations"`  
}

type Relations struct {
	ID             int           `json:"id"`
	DatesLocations map[string][]string      `json:"datesLocations"`
}

type SortOptions struct {
	// Champ sur lequel effectuer le tri (name, creation, members, etc.)
	Field string

	// Direction du tri (asc ou desc)
	Direction string
}

type ArtistList struct {
	Artists []Artists `json:"artists"`
}