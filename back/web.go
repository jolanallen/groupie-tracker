package groupietracker




import (
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

	if Id!= "" {
        g.RequestById(Id)
    }

	data := Artists{
        Name:        a.Name,
        Image:       a.Image,
        Members:     a.Members,
        CreationDate: a.CreationDate,
		FirstAlbum:  a.FirstAlbum,
		Relations:   a.Relations,
        DatesLocations:   a.DatesLocations,
	}
	
	// Exécution du template sans données supplémentaires (nil)
	tmpl.Execute(w, data)
	
}





