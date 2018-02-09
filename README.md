# goploader

![Go Version](https://img.shields.io/badge/go-1.8-brightgreen.svg)
![Go Version](https://img.shields.io/badge/go-1.9-brightgreen.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Depado/goploader)](https://goreportcard.com/report/github.com/Depado/goploader)
[![codebeat badge](https://codebeat.co/badges/0faefc03-91a4-41e7-a955-ccd8c1b096cd)](https://codebeat.co/projects/github-com-depado-goploader)
[![Build Status](https://drone.depado.eu/api/badges/Depado/goploader/status.svg)](https://drone.depado.eu/Depado/goploader)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Depado/goploader/blob/master/LICENSE)
[![Docs](https://img.shields.io/badge/docs-gpldr.in-blue.svg)](https://gpldr.in/)
[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://saythanks.io/to/Depado)

## Introduction

Goploader's ultimate goal is to make file sharing easy and painless. This project is composed of a server and a client, both written in Go. The main things to remember about the project are :
 - Sharing stuff from your terminal should be easy
 - Sharing stuff without a terminal should be easy
 - Privacy matters

## Build from source

Make sure you have Go installed on your machine.

### Client

```shell
$ go get github.com/Depado/goploader/client
$ go build -o $GOPATH/bin/goploader github.com/Depado/goploader/client
```

### Server

```shell
$ # Move to a new directory that will be used to run the server
$ go get github.com/Depado/goploader/server
$ # The following steps are optional
$ # Execute those if you wish to embed the assets and templates into the binary
$ go get github.com/GeertJohan/go.rice/rice
$ rice embed-go -i=github.com/Depado/goploader/server
$ # End of the optional steps
$ go build github.com/Depado/goploader/server
$ # If you did not embed the resources, make sure to copy the assets and templates directories
$ cp -r $GOPATH/src/github.com/Depado/goploader/server/{assets,templates} .
$ # Execute the binary a first time to trigger the setup
$ # Or write your own conf.yml file
$ ./server
```

## Downloads

All the downloads are available at [gpldr.in](https://gpldr.in) in the [clients](https://gpldr.in/#client-downloads) and [server](https://gpldr.in/#server-downloads) sections.

### Client

| Linux         | FreeBSD | Mac OS     | Windows  |
| ------------- |---------|------------|----------|
| [Linux 64bit](https://gpldr.in/releases/clients/client_linux_amd64) | [FreeBSD 64bit](https://gpldr.in/releases/clients/client_freebsd_amd64) | [Mac OS 64bit](https://gpldr.in/releases/clients/client_darwin_amd64) | [Windows 64bit](https://gpldr.in/releases/clients/client_windows_amd64.exe) |
| [Linux 32bit](https://gpldr.in/releases/clients/client_linux_386) | [FreeBSD 32bit](https://gpldr.in/releases/clients/client_freebsd_386) | [Mac OS 32bit](https://gpldr.in/releases/clients/client_darwin_386) | [Windows 32bit](https://gpldr.in/releases/clients/client_windows_386.exe) |
| [Linux ARMv7](https://gpldr.in/releases/clients/client_linux_arm) | | | | |

## Documentation

All the documentation is available at [gpldr.in](https://gpldr.in). I intend to write a proper `README.md` file, but it takes a lot of work to transpose the existing documentation to the markdown format. So, work in progress.


## License
All the software in this repository is released under the MIT License. See [LICENSE](https://github.com/Depado/goploader/blob/master/LICENSE) for details.
