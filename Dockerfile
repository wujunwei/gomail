FROM golang:1.18 as build

WORKDIR /app
COPY . /app

RUN go env -w GOPROXY=goproxy.io && go mod tidy \
    && CGO_ENABLED=0 go build -o ./bin/mailserver main.go

FROM alpine:latest
COPY --from=build /app/bin/mailserver  /usr/bin/mailserver

CMD ["/app/bin/mailserver"]
