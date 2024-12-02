
const url = 'https://api.adviceslip.com/advice';  


function FetchRequest() {
    fetch(url)
        .then(response => response.json())
        .then(data => {
          
            let adviceElement = document.createElement('p');
            adviceElement.innerHTML = `Conseil du jour : "${data.slip.advice}"`;
           
            document.getElementById('advice-container').appendChild(adviceElement);
        })
        .catch(error => console.error('Erreur lors de la récupération du conseil:', error));
}


function DireBonjour() {
    alert("Bonjour! Prêt pour un conseil ?");
    FetchRequest();  
}
