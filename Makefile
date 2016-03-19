version = 1.0.2

.PHONY: all clients servers release clean

all:
	go build -o client/client github.com/Depado/goploader/client
	go build -o server/server github.com/Depado/goploader/server

clients:
	-mkdir -p releases/clients
	-mkdir -p releases/servers
	-rm releases/clients/*
	gox -output="releases/clients/client_{{.OS}}_{{.Arch}}" github.com/Depado/goploader/client
	tar czf releases/servers/clients.tar.gz releases/clients

servers:
	-mkdir -p releases/servers
	-mkdir goploader-server
	rice embed-go -i=github.com/Depado/goploader/server
	go build -o goploader-server/server-standalone github.com/Depado/goploader/server
	tar czf releases/servers/server-standalone_amd64.tar.gz goploader-server
	rm -r goploader-server/*
	rice clean -i=github.com/Depado/goploader/server
	cp -r server/assets/ goploader-server/
	cp -r server/templates/ goploader-server/
	go build -o goploader-server/server github.com/Depado/goploader/server
	tar czf releases/servers/server_amd64.tar.gz goploader-server/
	rm -r goploader-server/*
	rice embed-go -i=github.com/Depado/goploader/server
	GOARCH=arm go build -o goploader-server/server-standalone github.com/Depado/goploader/server
	tar czf releases/servers/server-standalone_arm.tar.gz goploader-server
	rm -r goploader-server/*
	rice clean -i=github.com/Depado/goploader/server
	cp -r server/assets/ goploader-server/
	cp -r server/templates/ goploader-server/
	GOARCH=arm go build -o goploader-server/server github.com/Depado/goploader/server
	tar czf releases/servers/server_arm.tar.gz goploader-server/
	-rm -r goploader-server

release: clients servers
	tar czf servers.tar.gz releases/servers/
	mv servers.tar.gz releases/servers/

clean:
	-rm -r releases/
	-rm server/rice-box.go
	-rm -r goploader-server
