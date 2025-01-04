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
	// Utilisation de template.Must pour charger et exécuter le template
	tmpl := template.Must(template.ParseFiles(html))
	
	// Exécution du template sans données supplémentaires (nil)
	tmpl.Execute(w, nil)
	
}





