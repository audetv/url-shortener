# syntax = docker/dockerfile:1-experimental
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./url-shortener ./urlshortener/cmd/urlshortener

# 2

FROM scratch

WORKDIR /app

COPY --from=builder /app/urlshortener/static /app/static
COPY --from=builder /app/url-shortener /app/url-shortener
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Moscow

EXPOSE 8000

ENV REGUSER_STORE=mem

CMD ["./url-shortener"]

