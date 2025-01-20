package groupietracker

import (
	"log"
)

func (g *Groupie) Init() {
	// Initialise les chemins des templates
	g.TemplateHome = "front/templates/Home.html"
	g.TemplateArtist = "front/templates/Artist.html"

	// Valider les chemins des fichiers de template
	templates := []string{g.TemplateHome, g.TemplateArtist}
	for _, template := range templates {
		if !g.fileExists(template) {
			log.Printf("Template manquant : %s", template)
		}
	}
}

// renvoie une liste des artistes
