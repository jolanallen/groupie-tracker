package groupietracker

import (
	"fmt"
	"html/template"
	"net/http"
)

func (g *Groupie) Home(w http.ResponseWriter, r *http.Request) {   // fonction pour afficher les differents templates html
	g.Request(w, r, g.TemplateHome)
	
}
func (g *Groupie) Artist(w http.ResponseWriter, r *http.Request) {   // fonction pour afficher les differents templates html
	g.Request(w, r, g.TemplateArtist)  // template dans les quel on injectera les donner récupérer dans l'API
	
}
func (g *Groupie) Apropos(w http.ResponseWriter, r *http.Request) {   // fonction pour afficher les differents templates html
    g.Request(w, r, g.TemplateApropos)
    
}


func (g *Groupie) Request(w http.ResponseWriter, r *http.Request, html string) {
    r.ParseForm()

    tmpl := template.Must(template.ParseFiles(html))

 
    Id := r.FormValue("id")
    urlPath := r.URL.Path
    var data interface{}

   
    if Id != "" {
        artist, err := g.GetArtistById(Id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        
        fmt.Println("ID: ", artist.Id)
        fmt.Println("Name: ", artist.Name)
        fmt.Println("Image: ", artist.Image)
        fmt.Println("Members: ", artist.Members)
        fmt.Println("Creation Date: ", artist.CreationDate)
        fmt.Println("First Album: ", artist.FirstAlbum)
        fmt.Println("Relations: ", artist.Relations)

       
        data = artist
    } else if urlPath == "/" {
        // Si aucun ID n'est spécifié, on récupère tous les artistes
        artists, err := g.GetAllArtists()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        data = artists
    }

 
    err := tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
