# nsi_survey_server

env GOOS=linux GOARCH=amd64 go build

### Dev environment

The development env uses two bind mounts to connect the current local workspace and local microauth package to the container. To deploy development environment,
setup auth public key in .devcontainer/pk.pem and env variables in .devcontainer/devcontainer.env and run:

    docker-compose -f deploy/dev/docker-compose.yaml up -d
    docker exec -it NSISERVER_DEV bash

Env variables:

    PORT=
    DBUSER=
    DBPASS=
    DBNAME=
    DBHOST=host.docker.internal
    DBSTORE=pgx
    DBDRIVER=postgres
    DBSSLMODE=
    DBPORT=
    IPPK=.devcontainer/pk.pem
