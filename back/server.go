package groupietracker

import (
	"fmt"
	"net/http"
)
func (g *Groupie) Web() {


fmt.Println("Serveur démarré sur http://localhost:3666")

css := http.FileServer(http.Dir("front/css/"))
http.Handle("/css/", http.StripPrefix("/css/", css))

js := http.FileServer(http.Dir("front/js/"))
http.Handle("/js/", http.StripPrefix("/js/", js))

images := http.FileServer(http.Dir("front/utiles/"))
http.Handle("/utiles/", http.StripPrefix("/utiles/", images))


	// routes et fonctions associées
	http.HandleFunc("/", g.Home)
	http.HandleFunc("/Artists", g.Artist)
	http.HandleFunc("/Apropos", g.Apropos)
	
	//démarrage du server web
	err := http.ListenAndServe(":3666", nil)
	if err != nil {
		fmt.Println("Erreur lors du démarrage du serveur:", err)
		return
	}
	
}