FROM golang:1.23.7

WORKDIR /app

COPY go.mod .
COPY go.sum .

# Download all dependencies. Dependencies will be cached if the go.{mod,sum} files are unchanged
RUN go mod download

COPY . .

# This container exposes port 8080 to the outside world
EXPOSE 9000

# Run the binary program produced by `go install`
CMD go run main.go