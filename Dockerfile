FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o fetch .

FROM alpine:latest

COPY --from=builder /app/fetch /usr/local/bin/fetch

ENTRYPOINT ["fetch"]