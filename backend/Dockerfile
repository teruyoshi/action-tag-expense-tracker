FROM golang:1.22-alpine

RUN apk add --no-cache curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xz \
    && mv migrate /usr/local/bin/migrate

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN go build -o server cmd/server/main.go

CMD ["./entrypoint.sh"]
