FROM golang:1.26.10-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates curl

WORKDIR /app

COPY --from=builder /app/server .

ENV PORT=8080
ENV INTERNAL_PORT=9090
EXPOSE 8080 9090

CMD ["./server"]
