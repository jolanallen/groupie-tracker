package main

import (
	groupietracker "groupietracker/back"
)

func main() {
	var g groupietracker.Groupie
	g.Init()
	g.GetAllArtists()
	g.Web()

}
