FROM golang:1.23-alpine

RUN apk add --no-cache curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xz \
    && mv migrate /usr/local/bin/migrate \
    && go install github.com/air-verse/air@v1.61.7

WORKDIR /app

CMD ["./entrypoint.sh"]
