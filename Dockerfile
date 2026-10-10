# Build Step
FROM golang:1.27-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS builder

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
FROM gcr.io/distroless/static@sha256:58133991db06659feaabe0f4e97a35cebf15ef4ea08f8a4c6d2ee5f75e4aa6a0
COPY --from=builder /tmp/gpldr /go/bin/gpldr

VOLUME [ "/data" ]
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["/go/bin/gpldr"]
