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
                    if (artist.id === 1) {
                        console.log("Artiste trouvé :", artist);

                        // Extraire les informations de l'artiste
                        const ID = artist.id;
                        const NameGroupe = artist.name;
                        const Members = artist.members.join(', ');
                        const CDate = artist.creationDate;
                        const FAlbum = artist.firstAlbum;
                        const Image = artist.image;
                        const DConcert = artist.concertDates;

                        // Ajouter les informations à la page HTML
                        document.getElementById("affichage").innerHTML = `
                            <h2>Nom du groupe : ${NameGroupe}</h2>
                            <p>Membres : ${Members}</p>
                            <p>Date de création : ${CDate}</p>
                            <p>Premier album : ${FAlbum}</p>
                            <img src="${Image}" alt="${NameGroupe}" style="width: 200px; height: auto;">
                            <p>Dates de concert : ${DConcert}</p>
                        `;
                    }
                });
            } else {
                console.error("La réponse des artistes n'est pas un tableau.");
            }
        })
        .catch(error => console.error("Erreur lors de la récupération des données :", error));
}
