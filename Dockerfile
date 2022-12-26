############################
# STEP 1 build executable binary
############################
FROM golang:1.18-alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
# Create appuser.s
RUN adduser -D -g '' appuser
WORKDIR $GOPATH/src/github.com/samlitowitz/go-qr
RUN pwd
COPY . .
# Fetch dependencies.
# Using go mod.
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/go-qr-test cmd/go-qr/*.go

# Test Debug
FROM builder AS test-debug

# Install delver
RUN go install github.com/go-delve/delve/cmd/dlv@latest \
  && cp /go/bin/dlv /usr/local/bin/dlv

EXPOSE 8000 40000

CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "test", "github.com/samlitowitz/go-qr/mode/numeric", "--", "-test.run", "^TestEncoder_Encode$"]
CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "test", "github.com/samlitowitz/go-qr/mode/alphanumeric", "--", "-test.run", "^TestEncoder_Encode$"]
