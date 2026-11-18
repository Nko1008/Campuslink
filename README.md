# CampusLink

Projet **FR01 - HTML/CSS** réalisé dans le cadre du Bachelor 1 Informatique.

## Présentation

CampusLink est une interface web permettant de consulter différentes informations liées au campus et de faciliter la vie étudiante.

Le projet propose plusieurs fonctionnalités :

- Consulter les informations des salles.
- Découvrir les équipements disponibles.
- Consulter les incidents déclarés.
- Afficher le détail d'un incident.
- Déclarer un nouvel incident.
- Consulter les transports à proximité du campus et rechercher des arrêts en Île-de-France.

## Les différentes pages

Le projet présente plusieurs écrans :

- **Tableau de bord** : page d'accueil et accès aux différentes fonctionnalités.
- **Salles** : consultation des salles du campus.
- **Équipements** : consultation des équipements disponibles.
- **Incidents** : affichage des incidents déclarés.
- **Détail d'un incident** : consultation des informations détaillées d'un incident.
- **Déclaration d'un incident** : formulaire permettant de signaler un problème.
- **Transports** : recherche d'arrêts, consultation des lignes de transport, des prochains passages et des perturbations.

## Lancement du projet

Le projet utilise un serveur développé en Go. Il est donc nécessaire d'avoir Go installé sur son ordinateur.

Pour lancer le projet :

1. Ouvrir le dossier `CampusLink` dans Visual Studio Code.
2. Ouvrir un terminal dans ce dossier.
3. Lancer la commande :

```powershell
go run .
```

4. Ouvrir le lien suivant dans un navigateur :

```text
http://localhost:8080
```

Le serveur doit rester actif pendant toute l'utilisation du site.

**Important :** il ne faut pas ouvrir directement les fichiers HTML avec le navigateur. Les pages doivent être consultées à travers le serveur Go pour que les fonctionnalités utilisant les routes API puissent fonctionner correctement.

## Technologies utilisées

- **HTML** : structure des différentes pages.
- **CSS** : mise en forme et présentation de l'interface.
- **Flexbox** : organisation des éléments.
- **CSS Grid** : création des mises en page.
- **Media Queries** : adaptation aux différentes tailles d'écran.
- **JavaScript** : gestion des interactions et des requêtes vers les API.
- **Go** : création du serveur web et gestion des requêtes côté serveur.
- **API PRIM – Île-de-France Mobilités** : récupération des informations sur les transports en commun.

## Arborescence

```text
CampusLink/
├── index.html
├── salles.html
├── equipements.html
├── incidents.html
├── incident-detail.html
├── declarer-incident.html
├── transports.html
├── style.css
├── main.go
└── go.mod
```

## Organisation

Les différentes pages HTML correspondent aux différents écrans du projet.

Le fichier `style.css` contient les styles communs aux pages afin de conserver une présentation cohérente sur l'ensemble du site.

Le fichier `main.go` contient le serveur Go. Il permet notamment de servir les pages du site et de gérer les différentes routes API.

Le fichier `go.mod` définit le module Go et ses dépendances.

La navigation permet de passer d'un écran à l'autre grâce aux liens HTML.

## Fonctionnalité Transports

La page Transports utilise l'API PRIM d'Île-de-France Mobilités pour récupérer des informations actualisées sur le réseau de transports.

Elle permet notamment de :

- Rechercher un arrêt en Île-de-France.
- Consulter les lignes de bus et de tramway à proximité du campus.
- Consulter les lignes de RER et de train.
- Afficher les arrêts situés à proximité du campus.
- Consulter les prochains passages.
- Afficher les perturbations du réseau.

Les requêtes vers l'API sont effectuées côté serveur en Go afin de ne pas exposer la clé API dans le code HTML ou JavaScript.

## Responsive

L'interface est prévue pour être consultée sur ordinateur et sur téléphone.

Les éléments sont organisés avec CSS, Flexbox et Grid afin de permettre une adaptation de la mise en page selon la largeur de l'écran.

## Configuration de l'API

La fonctionnalité Transports nécessite une clé API PRIM d'Île-de-France Mobilités.

Cette clé doit être configurée dans une variable d'environnement nommée `PRIM_API_KEY`.

Sous Windows PowerShell :

```powershell
$env:PRIM_API_KEY = "VOTRE_CLE_API"
```

Cette commande doit être exécutée dans le même terminal que celui utilisé pour lancer le serveur Go.

Sans cette clé, les fonctionnalités liées aux transports ne pourront pas récupérer les données de l'API.
