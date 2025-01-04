package groupietracker



var a Artists


func (g *Groupie) Init() {
    a.Image =""
	a.Name =""
	a.Members =[]string{""}
	a.CreationDate = 0

	g.TemplateHome = "front/templates/index.html"
    g.TemplateArtist = "front/templates/artists.html"
    g.TemplateApropos = "front/templates/Apropos.html"
	
}

