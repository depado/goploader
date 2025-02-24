# Building from source

If you don't want to download a pre-compiled version of goploader server you can
build it from source.

## Prerequisites

- A recent version Go installed on your machine
- `$GOPATH` should be set to the appropriate directory

## Clone the repo

```shell
$ mkdir -p $GOPATH/src/github.com/Depado/
$ cd $GOPATH/src/github.com/Depado/
$ git clone https://github.com/Depado/goploader.git
$ cd goploader/
```

## Build 

Now all you have to do is build the project :

```shell
$ go build -o server github.com/Depado/goploader/server
```

That's it. Now head to the [setup](install.md) part to see how to
generate or create the `conf.yml` file that is necessary for goploader to run.
