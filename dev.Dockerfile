FROM golang:latest

RUN apt-get update && apt-get install vim
RUN go get github.com/go-delve/delve/cmd/dlv@latest
RUN go get github.com/go-delve/delve/cmd/dlv@master

ENTRYPOINT ["/workspaces/nsi_survey_server"]
