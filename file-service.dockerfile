# base go image
FROM golang:1.19-alpine as builder

RUN mkdir /app

RUN apk add build-base librdkafka-dev pkgconf

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=1 go build -tags musl -o ./build/fileApp ./internal/app/api

RUN go build -o ./build/fileAppCron ./internal/app/cron

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN chmod +x /app/build/fileApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/build/fileApp /app

COPY --from=builder /app/build/fileAppCron /app

COPY --from=builder /go/bin/migrate /bin/migrate

COPY ./.env /.env

CMD [ "/app/fileApp" ]