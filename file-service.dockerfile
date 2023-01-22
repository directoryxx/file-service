# base go image
FROM golang:1.19-alpine as builder

RUN mkdir /app

RUN apk add build-base librdkafka-dev pkgconf

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=1 go build -tags musl -o ./build/fileApp ./internal/app/api

RUN CGO_ENABLED=1 go build -tags musl -o ./build/fileAppCron ./internal/app/cron

RUN CGO_ENABLED=1 go build -tags musl -o ./build/fileAppWorker ./internal/app/worker

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN chmod +x /app/build/fileApp

RUN chmod +x /app/build/fileAppCron

RUN chmod +x /app/build/fileAppWorker

# build a tiny docker image
FROM alpine:latest

RUN apk add --no-cache supervisor tzdata

ENV TZ=Asia/Jakarta

RUN mkdir /app

RUN mkdir /migration

COPY ./docker/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

COPY --from=builder /app/build/fileApp /app

COPY --from=builder /app/build/fileAppCron /app

COPY --from=builder /app/build/fileAppWorker /app

COPY --from=builder /go/bin/migrate /bin/migrate

COPY ./db/ /migration

# COPY ./.env /.env

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]

