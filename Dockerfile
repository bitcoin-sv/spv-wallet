# Get Golang
FROM golang:1.17.6-alpine

# Version
LABEL version="1.0" name="BuxServer"

# Set the timezone to UTC
RUN apk update && \
    apk add -U tzdata build-base && \
    cp /usr/share/zoneinfo/EST5EDT /etc/localtime && \
    echo "UTC" > /etc/timezone

# Set the working directory
WORKDIR /go/bin

# Expose the port to the server
EXPOSE 3003

# Move the current files into the directory
COPY . $GOPATH/bin/

# Run the application
RUN chmod +x main
CMD ["main"]
