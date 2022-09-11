FROM golang:1.18-alpine

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o taavi

CMD [ "./taavi" ]
