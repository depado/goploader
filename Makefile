version = 1.0.0

.PHONY: all server-release clean prod

all:
	go build -o client/client github.com/Depado/goploader/client
	go build -o server/server github.com/Depado/goploader/server

amd64-release:
	rice embed-go -i=github.com/Depado/goploader/server
	go build -o server-standalone_amd64 github.com/Depado/goploader/server
	-rm server/*.rice-box.go
	-mkdir -p goploader-server/
	cp -r server/assets/ goploader-server/
	cp -r server/templates/ goploader-server/
	go build -o goploader-server/server github.com/Depado/goploader/server
	tar czf server-without-clients_linux_amd64.tar.gz goploader-server/
	mkdir -p goploader-server/assets/clients/
	gox -output="goploader-server/releases/clients/client_{{.OS}}_{{.Arch}}" github.com/Depado/goploader/client
	tar czf server-with-clients_linux_amd64.tar.gz goploader-server/

arm-release:
	rice embed-go -i=github.com/Depado/goploader/server
	go build -o server-standalone_arm github.com/Depado/goploader/server
	-rm server/*.rice-box.go
	-mkdir -p goploader-server/
	cp -r server/assets/ goploader-server/
	cp -r server/templates/ goploader-server/
	go build -o goploader-server/server github.com/Depado/goploader/server
	tar czf server-without-clients_linux_arm.tar.gz goploader-server/
	mkdir -p goploader-server/assets/clients/
	gox -output="goploader-server/releases/clients/client_{{.OS}}_{{.Arch}}" github.com/Depado/goploader/client
	tar czf server-with-clients_linux_arm.tar.gz goploader-server/

prod:
	rice embed-go -i=github.com/Depado/goploader/server
	go build -o server/server github.com/Depado/goploader
	cp -f server/server ~/goploader/server
	-mkdir tmp
	gox -output="tmp/client_{{.OS}}_{{.Arch}}" github.com/Depado/goploader/client
	-mkdir -p ~/goploader/releases/clients/
	-mkdir -p ~/goploader/releases/servers/
	-cp server-with*.tar.gz ~/goploader/releases/servers/
	-cp server-standalone_arm ~/goploader/releases/servers/
	-cp tmp/* ~/goploader/releases/clients/
	-rm -rf tmp
	-rm server/*.rice-box.go

clean:
	-rm client/client
	-rm server/server
	-rm server-with*.tar.gz
	-rm -r goploader-server
	-rm server/*.rice-box.go
