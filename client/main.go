package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	var err error
	var client *http.Client
	var scheme string
	var datasource io.Reader

	args := os.Args[1:]
	datasource = io.TeeReader(os.Stdin, os.Stdout)
	if len(args) > 0 {
		datasource, err = os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}
	}

	client = &http.Client{}
	req, err := http.NewRequest("POST", scheme+"https://up.depado.eu", datasource)
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
