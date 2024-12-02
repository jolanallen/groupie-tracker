fetch(url)
  .then(response => response.json())
  .then(data => {
    console.log(data);  // Débogage pour vérifier les données reçues
    
    // Vérification que la clé 'artists' existe et est un tableau
    if (Array.isArray(data)) {
      console.log("Il y a des artistes dans la réponse.");
      data.forEach(artist => {
        // Affichage des informations de chaque artiste
        console.log(`Nom du groupe : ${artist.name}`);
        console.log(`Membres : ${artist.members.join(', ')}`);
        console.log(`Date de création : ${artist.creationDate}`);
        console.log(`Premier album : ${artist.firstAlbum}`);
        console.log(`Image : ${artist.image}`);
        console.log(`Lien des concerts : ${artist.concertDates}`);
        console.log('-------------------');
      });
    } else {
      console.error("La réponse ne contient pas un tableau d'artistes.");
    }
  })
  .catch(error => console.error("Erreur lors de la récupération des données :", error));
