FROM golang:1.24.3-alpine as build-env

ENV CGO_ENABLED=1

RUN set -ex && \
    apk add --no-cache gcc musl-dev

WORKDIR /go/src/blnk

COPY . .
RUN apk add --no-cache git

RUN go build -ldflags='-s -w -extldflags "-static"' -o /blnk ./cmd/*.go

FROM debian:bullseye-slim

# Install pg_dump version 16
RUN apt-get update && apt-get install -y wget gnupg2 lsb-release gcc && \
    echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    apt-get update && \
    apt-get install -y postgresql-client-16 && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build-env /blnk /usr/local/bin/blnk

RUN chmod +x /usr/local/bin/blnk

CMD ["blnk", "start"]

EXPOSE 8080
