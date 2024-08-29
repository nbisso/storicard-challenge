FROM golang-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./src .

RUN go build -o main .

FROM alpine

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]

