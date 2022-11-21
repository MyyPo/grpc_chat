FROM golang:alpine

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

RUN mkdir /app
WORKDIR /app

COPY . .
COPY config.env .

RUN go get -d -v ./...

RUN go install -v ./...

#Setup auto-recompilation for dev stage
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

# without polling live reload DOESN'T WORK with new docker engine
ENTRYPOINT CompileDaemon -polling --build="go build -a -installsuffix cgo -o main ./cmd/server" --command=./main

EXPOSE 8080