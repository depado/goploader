# Build Step
FROM golang:1.24-alpine@sha256:c8c5f95d64aa79b6547f3b626eb84b16a7ce18a139e3e9ca19a8c078b85ba80d AS builder

# Dependencies
RUN apk update && apk add --no-cache git

# Source
WORKDIR $GOPATH/src/github.com/Depado/goploader
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /tmp/gpldr github.com/Depado/goploader/server

# Final Step
FROM gcr.io/distroless/static@sha256:2e114d20aa6371fd271f854aa3d6b2b7d2e70e797bb3ea44fb677afec60db22c
COPY --from=builder /tmp/gpldr /go/bin/gpldr

VOLUME [ "/data" ]
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["/go/bin/gpldr"]
