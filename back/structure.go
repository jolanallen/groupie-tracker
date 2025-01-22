package groupietracker

type Groupie struct {
	Name      string
	Id        int
	Relations string
	//apimaps        []string
	TemplateHome   string
	TemplateArtist string
	City           []string
}

type Artists struct {
	Id             int                 `json:"id"`
	Name           string              `json:"name"`
	Image          string              `json:"image"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Relations      string              `json:"relations"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistList struct {
	Artists []Artists `json:"artists"`
}

type FilterOptions struct {
	CreationDate int
	FirstAlbum   int
	MemberCount  int
	Locations    string
	SearchQuery  string
}
