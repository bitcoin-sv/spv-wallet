# Get Golang for builder
FROM golang:1.24.1 as builder

# Set the working directory
WORKDIR /go/src/github.com/bitcoin-sv/spv-wallet

COPY . ./

# Build binary
RUN GOOS=linux go build -o spvwallet cmd/main.go

# Get runtime image
FROM registry.access.redhat.com/ubi9-minimal

# Version
LABEL version="1.0" name="SPVWallet"

# Set working directory
WORKDIR /

# Copy binary to runner
COPY --from=builder /go/src/github.com/bitcoin-sv/spv-wallet/engine .

# Set entrypoint
ENTRYPOINT ["/spvwallet"]
