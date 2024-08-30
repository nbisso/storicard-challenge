FROM golang:1.21.6-alpine as builder

RUN apk add alpine-sdk

WORKDIR /app

COPY ./src /app

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -o api -tags musl

FROM alpine:latest as runner

WORKDIR /root/

COPY --from=builder /app/api /api

ENTRYPOINT /api