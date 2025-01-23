// Assurez-vous que ce script est bien lié à votre page HTML

// Fonction pour récupérer les coordonnées d'une ville via l'API Nominatim
async function getCoordinates(city) {
  try {
      const response = await fetch(`https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(city)}`);
      const data = await response.json();

      if (data.length === 0) {
          console.warn(`Ville non trouvée : ${city}`);
          return null;
      }

      // Récupération de la latitude et de la longitude de la première correspondance
      return {
          lat: parseFloat(data[0].lat),
          lon: parseFloat(data[0].lon),
      };
  } catch (error) {
      console.error(`Erreur lors de la récupération des coordonnées pour ${city}:`, error);
      return null;
  }
}

// Fonction principale pour mettre à jour les iframes des cartes
async function updateMaps() {
 
  const mapContainers = document.querySelectorAll('div[style*="position: relative;"]');
  for (const container of mapContainers) {
    
      const cityNameElement = container.querySelector('h3');
      const cityName = cityNameElement ? cityNameElement.textContent : null;
      const coordinates = await getCoordinates(cityName);

      if (!coordinates) {
          console.warn(`Impossible de récupérer les coordonnées pour : ${cityName}`);
          continue;
      }

      // Calculer la bounding box autour de la ville (ajuster pour zoomer)
      const delta = 0.05; // Ajustez cette valeur pour modifier le niveau de zoom
      const bbox = `${coordinates.lon - delta},${coordinates.lat - delta},${coordinates.lon + delta},${coordinates.lat + delta}`;

      // Modifier dynamiquement l'iframe pour centrer sur la ville
      const iframe = container.querySelector('iframe');
      if (iframe) {
          iframe.src = `https://www.openstreetmap.org/export/embed.html?bbox=${bbox}&layer=mapnik`;
      } else {
          console.warn("Aucun iframe trouvé pour ce conteneur.");
      }
  }
}

// Appeler la fonction après le chargement 
document.addEventListener("DOMContentLoaded", updateMaps);


