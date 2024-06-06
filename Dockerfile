FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY main.go ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ayaka .

FROM debian:slim

COPY --from=builder /app/ayaka /bin/ayaka

RUN chmod +x /bin/ayaka

RUN set -e \
    && apt-get update \
    && apt-get install -y \
        ca-certificates \
    && apt-get clean \
    && apt-get autoremove \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /tmp/* \
    && rm -rf /var/tmp/* \
    && rm -rf /var/cache/apt/archives/* \
    && groupadd -g 1100 app \
    && useradd -u 1100 -g 1100 -s /bin/bash -m app

USER app

RUN ayaka -h
