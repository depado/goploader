# Build Step
FROM golang:1.24-alpine@sha256:43c094ad24b6ac0546c62193baeb3e6e49ce14d3250845d166c77c25f64b0386 AS builder

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
FROM gcr.io/distroless/static@sha256:95ea148e8e9edd11cc7f639dc11825f38af86a14e5c7361753c741ceadef2167
COPY --from=builder /tmp/gpldr /go/bin/gpldr

VOLUME [ "/data" ]
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["/go/bin/gpldr"]
