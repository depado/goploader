# Build Step
FROM golang:1.25-alpine@sha256:ecb803821dbac0983b03a516adb8407b7dbadfb8771434dc5162a7c471dd27a2 AS builder

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
FROM gcr.io/distroless/static@sha256:87bce11be0af225e4ca761c40babb06d6d559f5767fbf7dc3c47f0f1a466b92c
COPY --from=builder /tmp/gpldr /go/bin/gpldr

VOLUME [ "/data" ]
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["/go/bin/gpldr"]
