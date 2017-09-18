If you don't want to download a pre-compiled version of goploader client you can
build it from source. Otherwise you can directly download a pre-compiled binary
[here](/client/install.md).

## Prerequisites

- A recent version Go installed on your machine
- `$GOPATH` should be set to the appropriate directory
- The `dep` tool must be installed : `go get -u github.com/golang/dep/cmd/dep`

## Clone the repo

```shell
$ mkdir -p $GOPATH/src/github.com/Depado/
$ cd $GOPATH/src/github.com/Depado/
$ git clone https://github.com/Depado/goploader.git
$ cd goploader/
$ dep ensure
```

## Build 

Now all you have to do is build the client :

```shell
$ cd 
$ go build -i -o gpldr github.com/Depado/goploader/client
```

## Next steps

Ensure you have `$GOPATH/bin` added to your standard `$PATH` so you can have
access to the binary you just built. Then head over to 
[the client's documentation](/client/documentation.md) to learn how to use it.