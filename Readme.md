# DNS focus

dns focus is an open source DNS server that allows you to easily block unwanted domains for your focus

## Table des matières

- [Features](#Features)
- [Prérequis](#prérequis)
- [Installation](#installation)
- [Configuration](#configuration)
- [Utilisation](#utilisation)
- [Licence](#licence)

## Features

- Blocage d'adresses IP indésirables
- Supporte IPv4
- Facile à configurer et à utiliser

## Prerequisites

- Go 1.18 ou supérieur
- Un système d'exploitation compatible (Linux, macOS)

## Installation

1. Clonez le dépôt :

   ```bash
   git clone https://github.com/RemyMach/dns-server.git
   ```


2. Compile the project

    `go build -o main`

3. Run the project

    `./main`

## Configuration

```json
{
    "DomainsBlocked": [
        "youtube.com",
        "twitter.com",
        "instagram.com",
        "twitch.com",
        "twitch.tv",
        "linkedin.com",
        "tiktok.com"
    ]
}
```

## Utilisation

- Pour démarrer le serveur DNS, exécutez simplement la commande suivante :

```bash
./NomDuProjetDNS
```
Configurer votre système ou votre appareil pour utiliser le serveur DNS en suivant les instructions spécifiques à votre système d'exploitation.


## Licence

dns-server- est distribué sous la licence MIT. Voir le fichier LICENSE pour plus d'informations.