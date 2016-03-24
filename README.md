# goploader

![Go Version](https://img.shields.io/badge/go-1.5-brightgreen.svg)
![Go Version](https://img.shields.io/badge/go-1.6-brightgreen.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Depado/goploader)](https://goreportcard.com/report/github.com/Depado/goploader)
[![codebeat badge](https://codebeat.co/badges/0faefc03-91a4-41e7-a955-ccd8c1b096cd)](https://codebeat.co/projects/github-com-depado-goploader)
[![Build Status](https://drone.depado.eu/api/badges/Depado/goploader/status.svg)](https://drone.depado.eu/Depado/goploader)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Depado/goploader/blob/master/LICENSE)
[![Docs](https://img.shields.io/badge/docs-up.depado.eu-blue.svg)](https://up.depado.eu/)


Goploader is a client-server application that is intended to ease the process of uploading files and sharing them.

All the documentation and information about this project are on [up.depado.eu](https://up.depado.eu).

- [Introduction](https://up.depado.eu/#introduction)
- [Client](https://up.depado.eu/#client)
- [Curl](https://up.depado.eu/#curl)
- [Server](https://up.depado.eu/#server)

## Philosophy

Goploader is intended to be easy to use and install both on the client side, and on the server side. The server itself is painless to install, configure and run, so that anyone can host a goploader server and use it for its own needs. One of the main feature of goploader is that it encrypts files while receiving them, sending back the full url of the resource to the user (including the private key allowing to decryp the file) while saving only the file ID in the database.

To get more information, please go the official documentation on [up.depado.eu](https://up.depado.eu).
