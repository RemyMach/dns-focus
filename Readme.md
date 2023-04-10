# DNS focus

dns focus is an open source DNS server that allows you to easily block unwanted domains for your focus

## Table of contents

- [Features](#Features)
- [Requirements](#Requirements)
- [Install](#Install)
- [Configuration](#configuration)
- [Usage](#Usage)
    - [Start with your go environment](#Start-with-your-go-environment)
    - [Start with Docker](#Start-with-docker)
    - [Env variables](#env-variables)
- [Licence](#licence)

## Features

- Blocking unwanted domain
- IPv4 support
- udp support
- Easy to set up and use
- dns proxy

## Requirements

- Go 1.16 or higher
- A compatible operating system (Linux, macOS)

## Install

1. Clone the repo :

   ```bash
   git clone https://github.com/RemyMach/dns-focus.git
   ```


2. Compile the project (to use in your go environment)

    `go build -o main`

3. Configure your system or device to use the DNS server by following the instructions specific to your operating system.

    - for mac you can simply use with your wifi network
    ```bash
    networksetup -setdnsservers Wi-Fi 127.0.0.1
    ```

    - for linux you can simply modify the file /etc/resolv.conf
        - comment the actual dns server
        - add your dns server
            ```
            nameserver 127.0.0.1
            ```

    Once you have finished using the server, remember to reset the basic dns server
    - for example to use google dns server on mac
        ```bash
        networksetup -setdnsservers Wi-Fi 8.8.8.8
        ```



## Configuration

you can block domains by adding domain names to those already present or make your own configuration file with your blocked domains

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

## Usage

### mode
- focus

### Start-with-your-go-environment

- To start the DNS server in any mode through the google dns server :

```bash
./main {mode} --proxy
```

- To start the DNS server in focus mode with google proxy (prefered mode to work in focus mode):
```bash
./main focus --proxy
```

- To start the DNS server in focus mode with google proxy and a specify config json file for example config/config.json:
```bash
./main focus --file "config/config.json" --proxy
```

- To start the DNS server in focus mode with only your dns server to resolve Dns request (for example if you want no cache on the domain name resolution)
```bash
./main focus
```

### Start-with-docker
- To start the DNS server with docker :

1. make sur to have only 127.0.0.1 on your dns server config, otherwise it might not work even if 127.0.0.1 is specified before the other dns server

2. copy the config file with the docker start command `cp .env.example .env`

3. modify DOCKER_APP_COMMAND if you want to use an other command to start your dns-focus

4. start the app in docker
```sh
docker compose up --build
```

#### env-variables
**Description :** DOCKER_APP_COMMAND is useful only to start the program in the docker env, this command will be used to start your container

## Licence

DNS focus is distributed under the MIT license.