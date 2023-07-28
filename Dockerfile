# Get Golang for builder
FROM golang:1.20.6 as builder

# Set the working directory
WORKDIR /go/src/github.com/BuxOrg/bux-server

COPY . ./

# Build binary
RUN GOOS=linux go build -o bux cmd/server/main.go

# Get runtime image
FROM registry.access.redhat.com/ubi9-minimal

# Version
LABEL version="1.0" name="Bux"

# Set working directory
WORKDIR /

# Copy binary to runner
COPY --from=builder /go/src/github.com/BuxOrg/bux-server/bux .

# Set entrypoint
ENTRYPOINT ["/bux"]
