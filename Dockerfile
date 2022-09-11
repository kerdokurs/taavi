FROM bufbuild/buf AS buf-generate

WORKDIR /app

COPY . .
RUN go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

RUN buf generate

FROM golang:1.18-alpine AS go-builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY *.go ./
COPY --from=buf-generate ./pkg/ ./

RUN go build -o taavi

CMD [ "./taavi" ]
