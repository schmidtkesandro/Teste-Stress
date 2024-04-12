FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o teste-stress

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/teste-stress .

ENTRYPOINT ["./teste-stress"]