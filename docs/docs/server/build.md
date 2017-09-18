If you don't want to download a pre-compiled version of goploader server you can
build it from source. Otherwise you can directly download a pre-compiled binary
[here](/server/downloads.md).

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

Now all you have to do is build the project :

```shell
$ go build -o server github.com/Depado/goploader/server
```

That's it. Now head to the [setup](/server/install.md) part to see how to
generate or create the `conf.yml` file that is necessary for goploader to run.

!!! tip "Embedding resources in the binary"
    If you want to include all the assets used by goploader server inside the
    generated binary you'll have to install the `rice` tool.

    ```shell
    $ go get github.com/GeertJohan/go.rice
    $ go get github.com/GeertJohan/go.rice/rice
    ```

    You can then generate a `rice-box.go` file by using this command :

    ```shell
    $ rice embed-go -i=github.com/Depado/goploader/server
    $ # Or by moving directly into the server dir
    $ cd server && rice embed-go
    ```

    This step must be executed before building the binary if you want embedded
    assets.