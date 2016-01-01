package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cheggaaa/pb"
)

const (
	service = "https://up.depado.eu"
	method  = "POST"
)

func main() {
	var err error
	var client *http.Client
	var datasource io.Reader
	var tee bool

	flag.BoolVar(&tee, "tee", false, "Displays stdin to stdout")
	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		var (
			f  *os.File
			fi os.FileInfo
		)

		if f, err = os.Open(args[0]); err != nil {
			fmt.Println("Could not open", args[0])
			os.Exit(1)
		}
		if fi, err = f.Stat(); err != nil {
			fmt.Println("Could not stat", args[0])
			os.Exit(1)
		}
		bar := pb.New64(fi.Size()).SetUnits(pb.U_BYTES)
		bar.Start()
		datasource = io.TeeReader(f, bar)
		if err != nil {
			log.Fatal(err)
		}
	} else if tee {
		datasource = io.TeeReader(os.Stdin, os.Stdout)
	} else {
		datasource = os.Stdin
	}

	client = &http.Client{}
	req, err := http.NewRequest(method, service, datasource)
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
