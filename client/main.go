package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Depado/goploader/client/conf"
)

func main() {
	var err error
	var client *http.Client
	var scheme string
	var datasource io.Reader

	if err = conf.Load("conf.yml"); err != nil {
		conf.C = conf.Configuration{Host: "up.depado.eu", TLS: true}
	}
	if conf.C.TLS {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
		scheme = "https"
	} else {
		client = &http.Client{}
		scheme = "http"
	}

	args := os.Args[1:]
	if len(args) > 0 {
		datasource, err = os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		datasource = os.Stdin
	}

	req, err := http.NewRequest("POST", scheme+"://"+conf.C.Host, datasource)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(ret))
}
