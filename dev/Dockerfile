
FROM golang:1.23

RUN apt-get update && apt-get install -y inotify-tools

ENV ROOT_PATH /go/src/github.com/jdstrand/language-checker
WORKDIR $ROOT_PATH
COPY go.mod ./
COPY go.sum ./

RUN go mod download

ENTRYPOINT ["./dev/autoreload.sh"]
