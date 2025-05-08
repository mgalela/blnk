FROM golang:1.23-alpine as build-env

WORKDIR /go/src/app

COPY . .

RUN go build -o qmonitor ./main.go

FROM alpine:latest

COPY --from=build-env /go/src/app/qmonitor /usr/local/bin/qmonitor

RUN chmod +x /usr/local/bin/qmonitor

CMD ["qmonitor"]

EXPOSE 8081
