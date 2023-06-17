# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# RUN stage

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
EXPOSE 8888
CMD ["/app/main"]
