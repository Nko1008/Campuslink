# CampusLink

Projet **FR01 - HTML/CSS** réalisé dans le cadre du Bachelor 1 Informatique.

## Présentation

CampusLink est une interface statique permettant de consulter les informations du campus et de déclarer des incidents.

Le projet présente plusieurs écrans :

* Tableau de bord
* Salles
* Équipements
* Incidents
* Détail d’un incident
* Déclaration d’un incident

## Lancement

Aucune installation n'est nécessaire.

Pour lancer le projet :

1. Ouvrir le dossier du projet.
2. Ouvrir le fichier `index.html` dans un navigateur.

Le projet fonctionne comme une interface statique et ne nécessite pas de JavaScript pour son affichage.

## Technologies utilisées

* HTML
* CSS
* Flexbox
* CSS Grid
* Media Queries

## Arborescence

```text
CampusLink/
├── index.html
├── salles.html
├── equipements.html
├── incidents.html
├── incident-detail.html
├── declarer-incident.html
└── style.css
```

## Organisation

Les différentes pages HTML correspondent aux différents écrans du projet.

Le fichier `style.css` contient les styles communs aux pages afin de conserver une présentation cohérente.

La navigation permet de passer d’un écran à l’autre grâce aux liens HTML.

## Responsive

L’interface est prévue pour être consultée sur ordinateur et sur téléphone.

Les éléments sont organisés avec CSS et Flexbox afin de permettre une adaptation de la mise en page selon la largeur de l’écran.
