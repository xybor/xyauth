# Introduction

Xyauth acts as an authentication provider which is compatible with oauth2 and openid.

# Technique

- Language: Golang, HTML, CSS, Javascript.
- Database: PostgreSQL, MongoDB.
- Deployment: Docker, Docker Compose.

// - Cloud: AWS.

# Get started

If you don't want to setup everything, let start with [Docker](#docker).

## Prerequisites

- Ubuntu 18.04
- Golang 1.18
- Postgres
- MongoDB
- (Optional) OpenSSL

## Setup database

This [article](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-postgresql-on-ubuntu-18-04) will guide you to setup the postgres database.

Then setup your new user and password [here](https://ubiq.co/database-blog/create-user-postgresql/).

After all, you need to determine the user, password, and database name to access to your database.

## Generate certificates

_If you already have your own certificates, ignore this section._

Following this [guide](https://devopscube.com/create-self-signed-certificates-openssl/) to create a self-signed certificate. The output is the private key `server.key` and the public key `server.crt`.

The following command will help you to generate a temporary certificate:

```bash
make cert-gen
```

## Setup Environment Variables

Please refer [.env.example](./.env.example) to setup the environment variables.

```bash
export key=value

# The following command help to setup the certificate as environment variable.
export key=`cat file.key`
```

Specially, if the value of `general.environment` in [config file](./configs/default.ini) is `dev`, you can create `.env` file which is similar to [.env.example](./.env.example) instead of using `export` commands.

The `.env` file has the higher priority than `shell` environment variables.

## Run the server

```bash
make run
```

# Docker

## Prerequisites

- Docker
- (Optional) OpenSSL

You must [generate certificate](#generate-certificates) and [setup environment variables](#setup-environment-variables) before the following steps.

## Generate docker compose file

```bash
make docker-gen
```

## Build the image

```bash
make docker-build
```

## Start docker compose

```bash
make docker-start
```

## Stop docker compose

```bash
make docker-stop
```

## Clean docker containers and images

```bash
make docker-clean
```
