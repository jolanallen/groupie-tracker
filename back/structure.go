package groupietracker

type Artists struct {
    Id    int               `json:"id"`
    Image string            `json:"image"`
	Name string             `json:"name"`
	Members []string        `json:"members"`
	CreationDate int        `json:"creationDate"`
    FirstAlbum string       `json:"firstAlbum"`
    
}

type Identifiant struct {
    Id int  `json:"id"`
}

type RelationData struct {
    DatesLocations map[string][]string    `json:"datesLocations"`
	
}

