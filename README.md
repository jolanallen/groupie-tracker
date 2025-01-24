# GROUPIE-TRACKER

Groupie-Tracker est une application web qui permet de suivre les artistes, leurs concerts, et leurs informations grâce à une interface conviviale et réactive.

## Structure du projet

L'arborescence du projet est organisée comme suit :

```
GROUPIE-TRACKER/
├── back/               # Contient les fichiers côté serveur en Go
│   ├── init.go         # Initialisation de l'application
│   ├── logic.go        # Gestion de la logique métier
│   ├── server.go       # Serveur HTTP
│   ├── structure.go    # Définition des structures de données
│   └── web.go          # Gestion des routes web
├── front/              # Contient les fichiers front-end
│   ├── .vscode/        # Configuration spécifique à l'éditeur
│   │   └── settings.json
│   ├── css/            # Feuilles de style CSS
│   │   ├── images/     # Dossier pour les images liées au style
│   │   ├── artist.css  # Styles pour la page Artiste
│   │   └── Home.css    # Styles pour la page d'accueil
│   ├── js/             # Fichiers JavaScript
│   │   └── location.js # Gestion des localisations
│   ├── templates/      # Templates HTML
│   │   ├── Artist.html # Page de détails d'un artiste
│   │   └── Home.html   # Page d'accueil
│   └── utils/          # Ressources utilitaires
│       ├── favicon.ico
│       ├── header.png
│       ├── image.png
│       ├── logo.png
│       └── nav.gif
├── go.mod              # Fichier des dépendances Go
├── main.go             # Point d'entrée de l'application
└── README.md           # Documentation du projet
```

---

## Fonctionnalités

- Affichage des informations des artistes (biographie, images, etc.).
- Suivi des concerts et des événements.
- Navigation simple et intuitive entre les pages.
- Gestion des données via un backend en **Go**.
- Interface utilisateur dynamique utilisant **HTML**, **CSS**, et **JavaScript**.

---

## Installation

1. Clonez le dépôt :
   ```bash
   git clone https://github.com/votre-utilisateur/GROUPIE-TRACKER.git
   ```
2. Naviguez dans le dossier du projet :
   ```bash
   cd GROUPIE-TRACKER
   ```
3. Installez les dépendances Go :
   ```bash
   go mod tidy
   ```
4. Lancez l'application :
   ```bash
   go run main.go
   ```

L'application sera disponible sur [http://localhost:8080](http://localhost:8080).

---

## Technologies utilisées

- **Backend** : Go (Golang)
- **Frontend** : HTML, CSS, JavaScript
- **Templates** : Fichiers HTML pour le rendu côté serveur
- **Styles** : CSS personnalisé

---

## Contribuer

Les contributions sont les bienvenues ! Suivez ces étapes :

1. Forkez le projet.
2. Créez une branche pour vos modifications :
   ```bash
   git checkout -b feature/ma-fonctionnalite
   ```
3. Effectuez vos changements et validez-les :
   ```bash
   git commit -m "Ajout d'une nouvelle fonctionnalité"
   ```
4. Poussez vos modifications :
   ```bash
   git push origin feature/ma-fonctionnalite
   ```
5. Créez une Pull Request sur GitHub.

---

## Auteurs

- **Allen jolan, Chereau Marino** –

---