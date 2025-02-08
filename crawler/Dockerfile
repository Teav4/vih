FROM golang:1.18-alpine

WORKDIR /app

# Copy go.mod first
COPY go.mod ./

# Run go mod tidy to ensure dependencies
RUN go mod tidy

# Copy the rest of the code
COPY . .

# Build
RUN go build -o crawler

CMD ["./crawler"]
