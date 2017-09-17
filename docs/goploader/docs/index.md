# Introduction

This is a page about the goploader project. The philosophy of goploader is that 
it should be easy to share files from your terminal. If you want to share a file
with a few of your friends for example.

This project is composed of a client and a server. Both are free (as in speech 
and as in beer) and you can install, tweak or do whatever you want with the 
code. If you want people to use your service with the official client, they will
have to modify the ~/.config/goploader.conf.yml file that is created on the 
first run of the client.

 
## Features

- HTTPS only using Let's Encrypt and Caddy Server
- Upload directly from stdin
- Upload a file by giving it as an argument to the client
- All files are encrypted upon reception and decrypted only when served
- The key to decrypt files is not saved on the server
- Curl compliant


!!! note 
    The HTTPS part applies to [gpldr.in](https://gpldr.in), but since it depends
    on the way the application is served and deployed, I can't guarantee it will be
    the case on other gpldr instances.

## Motivation

First of all, I created goploader as a fun exercise to improve my Golang skills.
I wanted to create something useful. Goploader was inspired by the way 0bin 
works, but I knew I couldn't encrypt files in the browser mainly because I 
didn't want to on a browser. 

Basically I merged omploader (which is a dead  service now) and 0bin, and the 
idea was to have a client that would allow painless file upload while keeping 
the fact that the server never knows what's being uploaded. I wanted to be able 
to upload files directly from my terminal, with or without a client (curl 
compatibility).

## How it works

The server and client aren't using a peculiar API, instead they just use 
multipart forms. This way, a curl command is simple to understand, and simple to
use. This ensures that anyone needing to upload files to the server can do it 
even without a graphic interface or without advanced knowledge in networking.

### Process

When the server receives an upload request, it does multiple steps.

- Create a unique identifier for the resource (1)*
- Create a unique key to encrypt the received data (2)*
- Create an AES cipher that uses the key generated in the second step*
- Pipe the file reception to that AES cipher*
- Pipe the output of the AES cipher to a file on disk*
- Process the whole file through that pipe
- Store the unique identifier (1) along with the name of the file and delete 
date in database
- Return the URL of the resource to the client in the form 
`[scheme][name server]/v/(1)/(2)`

All the steps with a * at the end are executed before receiving the file.

!!! warning "About the encryption key"
    The server never stores the encryption/decryption key. It is generated once 
    to process the file and is immediatly forgotten once given back to the user.

!!! info
    The file is never read fully in memory. It is buffered, and immediatly piped
    to the AES cipher. This ensures that the file is never fully stored in RAM 
    unencrypted. The data is streamed through the AES cipher which streams its
    result to the hard drive.

!!! tip
    The length of encryption/decryption key can be customised in the 
    configuration file of the server.

## Technologies

The server and client try to keep things minimal.
            
### Client

- Network part is pure Go
- [github.com/atotto/clipboard](https://github.com/atotto/clipboard) - Clipboard
 capabilities
- [github.com/cheggaaa/pb](https://github.com/cheggaaa/pb) - Progress bar 
capabilities
- [github.com/ogier/pflag](https://github.com/ogier/pflag) - Posix flags
- [github.com/mitchellh/go-homedir](https://github.com/mitchellh/go-homedir) - 
Home directory detection

### Server
                    
- [github.com/GeertJohan/go.rice](https://github.com/GeertJohan/go.rice) - 
Resource embedding
- [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin) - HTTP framework
- [github.com/boltdb/bolt](https://github.com/boltdb/bolt) - Embedded key/value 
datastore
- [github.com/ogier/pflag](https://github.com/ogier/pflag) - Posix flags
- [github.com/dchest/uniuri](https://github.com/dchest/uniuri) - Unique URI 
generator


## Disclaimer and License

As specified multiple times to people interested in the project, I don't provide
any kind of warranty. People might tamper with the code on their own server, the
only thing I can tell is that the version that runs on my machine is the exact 
version that is hosted on the Github repository. As a side note, I do store the
IP of people uploading files to [gpldr.in](https://gpldr.in) for debugging 
purposes and legal requirements, although I never store the decryption key.

```
The MIT License (MIT)

Copyright (c) 2015

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```    