FROM golang:1.13 AS builder

RUN mkdir /src

COPY go.mod go.sum /src/

WORKDIR /src

RUN go mod download

COPY . /src/

RUN CGO_ENABLED=0 GOOS=linux go build -o billing

FROM alpine:latest as certs

RUN apk --update add ca-certificates

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /src/billing /src/

ENTRYPOINT ["/src/billing"]

EXPOSE 80
