FROM golang:1.24-alpine

RUN apk add --no-cache curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xz \
    && mv migrate /usr/local/bin/migrate \
    && go install github.com/air-verse/air@v1.61.7 \
    && go install golang.org/x/vuln/cmd/govulncheck@latest \
    && go install github.com/securego/gosec/v2/cmd/gosec@v2.22.4

WORKDIR /app

CMD ["./entrypoint.sh"]
