# Build Step
FROM golang:1.18-alpine AS builder

# Dependencies
RUN apk update && apk add --no-cache upx git
RUN go get github.com/GeertJohan/go.rice/rice

# Source
WORKDIR $GOPATH/src/github.com/Depado/goploader
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Embed
RUN rice embed-go -i=github.com/Depado/goploader/server

# Build
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /tmp/gpldr github.com/Depado/goploader/server
RUN upx /tmp/gpldr


# Final Step
FROM gcr.io/distroless/static
COPY --from=builder /tmp/gpldr /go/bin/gpldr

VOLUME [ "/data" ]
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["/go/bin/gpldr"]
