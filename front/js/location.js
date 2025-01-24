async function getCoordinates(city) {
  try {
      const response = await fetch(`https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(city)}`);
      const data = await response.json();

      if (data.length === 0) {
          return null;
      }
      return {
          lat: parseFloat(data[0].lat),
          lon: parseFloat(data[0].lon),
      };
  } catch (error) {
      return null;
  }
}

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

      const delta = 0.05; 
      const bbox = `${coordinates.lon - delta},${coordinates.lat - delta},${coordinates.lon + delta},${coordinates.lat + delta}`;

     
      const iframe = container.querySelector('iframe');
      if (iframe) {
          iframe.src = `https://www.openstreetmap.org/export/embed.html?bbox=${bbox}&layer=mapnik`;
      } else {
          console.warn("Aucun iframe trouvé pour ce conteneur.");
      }
  }
}

document.addEventListener("DOMContentLoaded", updateMaps);


