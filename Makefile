version = 1.0.0

.PHONY: all server-release clean prod

all:
	go build -o client/client github.com/Depado/goploader/client
	go build -o server/server github.com/Depado/goploader/server

amd64-release:
	-mkdir -p goploader-server/
	cp -r server/assets/ goploader-server/
	cp -r server/templates/ goploader-server/
	go build -o goploader-server/server github.com/Depado/goploader/server
	tar czf server-without-clients_linux_amd64.tar.gz goploader-server/
	mkdir -p goploader-server/assets/clients/
	gox -output="goploader-server/assets/clients/client_{{.OS}}_{{.Arch}}" github.com/Depado/goploader/client
	tar czf server-with-clients_linux_amd64.tar.gz goploader-server/

arm-release:
	-mkdir -p goploader-server/
	cp -r server/assets/ goploader-server/
	cp -r server/templates/ goploader-server/
	go build -o goploader-server/server github.com/Depado/goploader/server
	tar czf server-without-clients_linux_arm.tar.gz goploader-server/
	mkdir -p goploader-server/assets/clients/
	gox -output="goploader-server/assets/clients/client_{{.OS}}_{{.Arch}}" github.com/Depado/goploader/client
	tar czf server-with-clients_linux_arm.tar.gz goploader-server/

prod:
	go build github.com/Depado/goploader/server
	cp -f server ~/goploader/server
	cp -r server/templates ~/goploader/
	cp -r server/assets ~/goploader/
	-mkdir tmp
	gox -output="tmp/client_{{.OS}}_{{.Arch}}" github.com/Depado/goploader/client
	-mkdir -p ~/goploader/assets/{clients,servers}
	-cp server-with*.tar.gz ~/goploader/assets/servers/
	-cp tmp/* ~/goploader/assets/clients/
	-rm -rf tmp

clean:
	-rm client/client
	-rm server/server
	-rm server-with*.tar.gz
	-rm -r goploader-server
