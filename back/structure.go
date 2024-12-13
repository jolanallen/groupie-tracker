package groupietracker

type Artists struct {
    Id    int    
    Image string 
	Name string
	Members []string
	CreationDate int
    FirstAlbum string
    Locations string
    ConcertDates string
    Relations string
}

func (a *Artists)Run(){
    
}