FROM golang:1.18-alpine as build
LABEL maintainer="Roch D'Amour <roch.damour@gmail.com>"
MAINTAINER Roch D'Amour <roch.damour@gmail.com>

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk add --no-cache --update ca-certificates

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o technewsbot main.go

FROM scratch

COPY --from=build /app/technewsbot /technewsbot
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/technewsbot"]
