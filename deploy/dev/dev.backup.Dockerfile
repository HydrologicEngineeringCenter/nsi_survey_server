FROM golang:latest

RUN apt-get update && apt-get -y install vim
RUN go get github.com/go-delve/delve/cmd/dlv@latest
RUN go get github.com/go-delve/delve/cmd/dlv@master
