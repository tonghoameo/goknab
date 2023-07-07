# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN apk add curl
RUN go build -o main main.go

#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
# RUN stage

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
#COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migrations ./db/migrations
EXPOSE 8888
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
