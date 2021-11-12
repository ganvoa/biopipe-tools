FROM golang:1.17.1-bullseye AS builder

WORKDIR /go/src/app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build -o /usr/bin/biopipe ./cmd/cli/biopipe.go

CMD ["/usr/bin/biopipe"]