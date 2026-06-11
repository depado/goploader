# Build Step
FROM golang:1.26-alpine@sha256:a6a091eac01ceac4b97496fe2957a49b6cdd83365337d5f46f6f73710424e805 AS builder

# Dependencies
RUN apk update && apk add --no-cache git

# Source
WORKDIR $GOPATH/src/github.com/Depado/goploader
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' -o /tmp/gpldr github.com/Depado/goploader/server

# Final Step
FROM gcr.io/distroless/static@sha256:3592aa8171c77482f62bbc4164e6a2d141c6122554ace66e5cc910cadb961ff0
COPY --from=builder /tmp/gpldr /go/bin/gpldr

VOLUME [ "/data" ]
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["/go/bin/gpldr"]
