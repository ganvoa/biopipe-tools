FROM golang:1.17.1-bullseye AS builder

WORKDIR /go/src/app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build -o /usr/bin/biopipe ./cmd/cli/biopipe.go

FROM debian:bullseye-slim

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates
COPY --from=builder /usr/bin/biopipe /usr/bin/biopipe

CMD ["/usr/bin/biopipe"]