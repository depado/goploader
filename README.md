# goploader

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/contains-technical-debt.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)

![Go Version](https://img.shields.io/badge/go-1.18-brightgreen.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/depado/goploader)](https://goreportcard.com/report/github.com/depado/goploader)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/depado/goploader/blob/master/LICENSE)
[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://saythanks.io/to/depado)

> [!WARNING]  
> This repository is maintained as is but most of the tech used for this project
> is now outdated. Use at your own risk. 

## Introduction

Goploader's ultimate goal is to make file sharing easy and painless. This project is composed of a server and a client, both written in Go. The main things to remember about the project are :
 - Sharing stuff from your terminal should be easy
 - Sharing stuff without a terminal should be easy
 - Privacy matters

## Build from source

Make sure you have go installed on your machine.

### Client

```shell
$ git clone https://github.com/depado/goploader.git
$ cd goploader
$ go build -o gpldr ./client/
```

### Server

```shell
$ git clone https://github.com/depado/goploader.git
$ cd goploader
$ go build -o goploader-server ./server/
$ ./goploader-server
```

## Downloads

All the downloads are available in the [releases tab](https://github.com/depado/goploader/releases)
of this repository. 

## Documentation

All the documentation is available at [depado.github.io/goploader/](https://depado.github.io/goploader/).

## License

All the software in this repository is released under the MIT License. See 
[LICENSE](https://github.com/depado/goploader/blob/master/LICENSE) for details.
