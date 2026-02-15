FROM golang:1.25 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o todo-app ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/my-scheduler .
COPY --from=builder /build/web ./web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=

CMD ["./my-scheduler"]