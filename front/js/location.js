//document.addEventListener("DOMContentLoaded", function () {
//  const iframe = document.getElementById("map-iframe");
//  iframe.onload = function () {
//    iframe.contentWindow.scrollTo(0, iframe.contentWindow.document.body.scrollHeight);
//  };
//});


function initMap() {
	// Créer l'objet "macarte" et l'insèrer dans l'élément HTML qui a l'ID "map"
	macarte = L.map('map').setView([lat, lon], 11);
	// Leaflet ne récupère pas les cartes (tiles) sur un serveur par défaut. Nous devons lui préciser où nous souhaitons les récupérer. Ici, openstreetmap.fr
	L.tileLayer('https://{s}.tile.openstreetmap.fr/osmfr/{z}/{x}/{y}.png', {
		// Il est toujours bien de laisser le lien vers la source des données
		attribution: 'données © OpenStreetMap/ODbL - rendu OSM France',
		minZoom: 1,
		maxZoom: 20
	}).addTo(macarte);
	// Nous parcourons la liste des villes
	for (ville in villes) {
		var marker = L.marker([villes[ville].lat, villes[ville].lon]).addTo(macarte);
	}               	
}

// Nous parcourons la liste des villes
for (ville in villes) {
	var marker = L.marker([villes[ville].lat, villes[ville].lon]).addTo(macarte);
	// Nous ajoutons la popup. A noter que son contenu (ici la variable ville) peut être du HTML
	marker.bindPopup(ville);
}               	