package groupietracker

type Artists struct {
    Id    int               `json:"id"`
    Image string            `json:"image"`
	Name string             `json:"name"`
	Members []string        `json:"members"`
	CreationDate int        `json:"creationDate"`
    FirstAlbum string       `json:"firstAlbum"`
    Locations string        
    ConcertDates string     `json:"concertDates"`
    Relations string        `json:"relations"`
}

type Identifiant struct {
    Id int  `json:"id"`
}

type RelationData struct {
	Locations []string      `json:"locations"`
	Dates     []string      `json:"dates"`
}