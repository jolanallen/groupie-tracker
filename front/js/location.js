
function getCityInfo() {
    const city = document.getElementById("city").value;

    // Effectuer une requête GET au backend Go pour obtenir les infos de la ville
    fetch(`http://localhost:8080/city?city=${city}`)
        .then(response => response.json())
        .then(data => {
            if (data.city) {
                document.getElementById("result").innerHTML = 
                    `<h2>${data.city}</h2>
                      <p>Population: ${data.population}</p>
                      <p>Pays: ${data.country}</p>`;
            } else {
                document.getElementById("result").innerHTML = "Aucune donnée trouvée.";
            }
        })
        .catch(error => {
            console.error('Erreur:', error);
            document.getElementById("result").innerHTML = "Erreur lors de la récupération des données.";
        });
}
