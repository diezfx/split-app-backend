VERSION 0.7
FROM golang:1.21
WORKDIR /go-workdir

build:
    COPY . .
    RUN go build -o output/example ./cmd/split-app-backend
    SAVE ARTIFACT output/example AS LOCAL local-output/go-example

docker:
    COPY +build/example .
    ENTRYPOINT ["/go-workdir/example"]
    SAVE IMAGE go-example:latest --push