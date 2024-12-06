const apiURL = "https://groupietrackers.herokuapp.com/api";
const ProxyURL = "https://cors-anywhere.herokuapp.com/"; // Ou utilisez un autre proxy comme AllOrigins

function Fetch() {
    fetch(ProxyURL + apiURL)
        .then(response => {
            if (!response.ok) {
                throw new Error(`Erreur HTTP ! Statut : ${response.status}`);
            }
            return response.json(); // La réponse principale est en JSON
        })
        .then(data => {
            console.log("Données principales récupérées :", data);

            // Accéder à la liste des artistes via l'URL fournie
            const artistsURL = data.artists;
            console.log("URL des artistes :", artistsURL);

            // Faire un nouvel appel pour récupérer les artistes
            return fetch(ProxyURL + artistsURL);
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Erreur HTTP ! Statut : ${response.status}`);
            }
            return response.json(); // La réponse des artistes est en JSON
        })
        .then(artists => {
            console.log("Liste des artistes :", artists);

            // Vérifier si 'artists' est bien un tableau
            if (Array.isArray(artists)) {
                console.log("Artistes récupérés avec succès !");
                artists.forEach(artist => {
                    console.log(`Nom du groupe : ${artist.name}`);
                    console.log(`Membres : ${artist.members.join(', ')}`);
                    console.log(`Date de création : ${artist.creationDate}`);
                    console.log(`Premier album : ${artist.firstAlbum}`);
                    console.log(`Image : ${artist.image}`);
                    console.log(`Lien des concerts : ${artist.concertDates}`);
                    console.log('-------------------');
                });
            } else {
                console.error("La réponse des artistes n'est pas un tableau.");
            }
        })
        .catch(error => console.error("Erreur lors de la récupération des données :", error));
}
