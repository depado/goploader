version = 1.0.0

.PHONY: all server-release clean prod

all:
	go build -o client/client github.com/Depado/goploader/client
	go build -o server/server github.com/Depado/goploader/server

clients:
	-mkdir -p clients
	gox -output="clients/client_{{.OS}}_{{.Arch}}" github.com/Depado/goploader/client
	tar czf clients.tar.gz clients

amd64-release:
	-mkdir goploader-server
	rice embed-go -i=github.com/Depado/goploader/server
	go build -o goploader-server/server-standalone github.com/Depado/goploader/server
	tar czf server-standalone_amd64.tar.gz goploader-server
	rm -r goploader-server
	rice clean -i=github.com/Depado/goploader/server
	-mkdir -p goploader-server/
	cp -r server/assets/ goploader-server/
	cp -r server/templates/ goploader-server/
	go build -o goploader-server/server github.com/Depado/goploader/server
	tar czf server_amd64.tar.gz goploader-server/

arm-release:
	-mkdir goploader-server
	rice embed-go -i=github.com/Depado/goploader/server
	go build -o goploader-server/server-standalone github.com/Depado/goploader/server
	tar czf server-standalone_arm.tar.gz goploader-server
	rm -r goploader-server
	rice clean -i=github.com/Depado/goploader/server
	-mkdir -p goploader-server/
	cp -r server/assets/ goploader-server/
	cp -r server/templates/ goploader-server/
	go build -o goploader-server/server github.com/Depado/goploader/server
	tar czf server_arm.tar.gz goploader-server/

prod:
	rice embed-go -i=github.com/Depado/goploader/server
	go build -o server/server github.com/Depado/goploader/server
	cp -f server/server ~/goploader/server
	-mkdir -p ~/goploader/releases/servers
	-cp server*.tar.gz ~/goploader/releases/servers/
	rice clean -i=github.com/Depado/goploader/server

clean:
	-rm client/client
	-rm server/server
	-rm server*.tar.gz
	-rm -r goploader-server
	rice clean -i=github.com/Depado/goploader/server
