FROM golang:1.18-alpine

WORKDIR /app

# Copy the go.mod and go.sum files first
COPY go.mod ./
COPY go.sum ./

# Download necessary Go modules
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o sensors main.go

CMD ["./sensors"]
