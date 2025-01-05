package groupietracker

import  ( 
    "regexp"
    "net/http"
    "fmt"
    "io/ioutil"
)


func (a *Artists) LoadData() {


    API := [2]string{
        "https://groupietrackers.herokuapp.com/api/artists",
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
                            Artists = Get_ArtistsData(result, Artists) //get Artist Data => name, images, members, creation date ,first album
                    case 1:
                            Artists = Get_RelationData(result, Artists) //get relations =>date + location
            }
        }
    }
    return fmt.Println(Artists)
}

func (a *Artists) Get_ArtistsData(){


}

func (a *Artists) Get_RelationData(){

    
}