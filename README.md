# goploader
Simple client-server application to upload files/text

## Principle

The [server](https://github.com/Depado/goploader/tree/master/server) part is a simple HTTP server. It will wait for data directly on the port you defined in your configuration file (see [usage](#usage) for more information). When receiving data, it will automatically stream the data to a file on disk (not loading it all in RAM) and insert an entry in the sqlite database. Every minute, a goroutine queries the database and check for entries that are older than 30 minutes. If it finds such entries, the file and entry will be deleted. Note that it doesn't currently set a limit on the size it can receive. I use nginx in front of this Go program with a data limit. If you don't have to have a reverse proxy in front of it, you'll have to handle the size limit yourself.

The [client](https://github.com/Depado/goploader/tree/master/client) is a simple command line program that will read on Stdin or an argument you pass to the binary. When installing the binary, I named it `goploader`. It allows to do things like :

```sh
$ cat file.txt | goploader
$ goploader < file.txt
$ goploader file.txt
```

## Usage

Compile and push the server binary to your own server. Modify the configuration file (see [conf.yml.example](https://github.com/Depado/goploader/blob/master/conf.yml.example)) and name it `conf.yml` OR specify the path to your configuration file with the `-c` argument when executing the binary.
