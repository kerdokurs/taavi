FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/taavi cmd/http/main.go

FROM alpine:latest AS runner

WORKDIR /app

ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip

COPY --from=builder /app/taavi /app/taavi
COPY --from=builder /app/templates /app/templates

ENTRYPOINT [ "/app/taavi" ]
