FROM  alpine:3.19.1

#install goose
WORKDIR /opt
RUN apk --no-cache add curl
RUN curl -fsSL \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    GOOSE_INSTALL=/opt sh -s v3.20.0

#add migrations to docker
RUN mkdir -p /opt/go/app/db
ADD ./migrations /opt/go/app/db/migrations

WORKDIR /opt/go/app

