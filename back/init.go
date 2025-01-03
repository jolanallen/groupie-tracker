package groupietracker

func (a *Artists) init() {
    a.Id   = "https://groupietrackers.herokuapp.com/api"
    a.Image = "https://groupietrackers.herokuapp.com/api/artists"
	a.Name =
	a.Members =
	a.CreationDate =
    a.FirstAlbum =
    a.Locations =
    a.ConcertDates =
    a.Relations =
}


func LoadData(g *) []ArtistsData {
	Artists := make([]ArtistsData, 0)

API := [4]string{
    "https://groupietrackers.herokuapp.com/api/artists",
    "https://groupietrackers.herokuapp.com/api/locations",
    "https://groupietrackers.herokuapp.com/api/dates",
    "https://groupietrackers.herokuapp.com/api/relation",
}

var ListButStr = regexp.MustCompile(`{([^{}]*)}`)

for index, link := range API {
    response, err := http.Get(link)
    if err != nil {
        fmt.Printf("Rest failed %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        result := ListButStr.FindAllString(string(data), -1)
        switch index {
        case 0:
            Artists = Get_ArtitsData(result, Artists) //get Artist Data => name, images, members, creation date ,first album
        case 1:
            Artists = Get_LocationsData(result, Artists) // get concert locations
        case 2:
            Artists = Get_DatesData(result, Artists) //get concert date
        case 3:
            Artists = Get_RelationData(result, Artists) //get relations =>date + location
        }
    }
}
return Artists
}
