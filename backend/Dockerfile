FROM golang:1.18-alpine

# Install required tools
RUN apk add --no-cache protobuf protobuf-dev git

# Set working directory inside GOPATH
WORKDIR /go/src/app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Install protoc plugins and dependencies
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 \
    && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1

# Add protoc plugins to PATH
ENV PATH="$PATH:$(go env GOPATH)/bin"

# Copy proto files first
COPY proto/ proto/

# Set GO111MODULE
ENV GO111MODULE=on

# Generate Go files from proto
RUN mkdir -p proto && \
    protoc \
    --proto_path=proto \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=. \
    --grpc-gateway_opt=paths=source_relative \
    proto/*.proto

# Copy the rest of the code
COPY . .

# Run go mod tidy to ensure dependencies
RUN go mod tidy

# Build the main package
RUN CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd/main.go

EXPOSE 50051 8080

CMD ["./backend"]
