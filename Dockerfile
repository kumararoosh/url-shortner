FROM golang:1.24-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o shortener

EXPOSE 8080

CMD ["./shortener"]