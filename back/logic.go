package groupietracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (g *Groupie) RequestById(Id string) {
	

	url := fmt.Sprintf("http://groupietrackers.herokuapp.com/api/artists/%s", Id)

	Request, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data: ", err)
		return
	}
	defer Request.Body.Close()

	if Request.StatusCode != http.StatusOK {
		fmt.Println("Error status code: ", Request.StatusCode)
		return
	}

	data, err := ioutil.ReadAll(Request.Body)
	if err != nil {
		fmt.Println("error reading body: ", err)
		return
	}

	var responseData Artists
	err = json.Unmarshal(data, &responseData)
	if err != nil {
		fmt.Println("Error Unmarshal: ", err)
		return
	}

	g.Relations = responseData.Relations

	// Affichage des informations de l'artiste
	fmt.Println("ID: ", responseData.Id)
	fmt.Println("Name: ", responseData.Name)
	fmt.Println("Image: ", responseData.Image)
	fmt.Println("Members: ", responseData.Members)
	fmt.Println("Creation Date: ", responseData.CreationDate)
	fmt.Println("First Album: ", responseData.FirstAlbum)
	fmt.Println("Relations: ", g.Relations)
	
    a.Image = responseData.Image
    a.Members = responseData.Members
    a.CreationDate = responseData.CreationDate


	g.RequestRelation()
}

func (g *Groupie) RequestRelation() {
	if g.Relations == "" {
		fmt.Println("Relations is not set")
		return
	}

	Requestbis, err := http.Get(g.Relations)
	if err != nil {
		fmt.Println("Error fetching relations: ", err)
		return
	}
	defer Requestbis.Body.Close()

	if Requestbis.StatusCode != http.StatusOK {
		fmt.Println("Error status code: ", Requestbis.StatusCode)
		return
	}

	data, err := ioutil.ReadAll(Requestbis.Body)
	if err != nil {
		fmt.Println("error reading body: ", err)
		return
	}

	var responseDataBis Relations
	err = json.Unmarshal(data, &responseDataBis)
	if err != nil {
		fmt.Println("Error Unmarshal: ", err)
		return
	}
	jsonData, err := json.MarshalIndent(responseDataBis.DatesLocations, "", "")
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println("ID: ", responseDataBis.ID)
	fmt.Println("name:", )
	fmt.Println("Relations JSON: ", string(jsonData))
}
