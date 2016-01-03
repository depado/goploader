package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/cheggaaa/pb"
	flag "github.com/ogier/pflag"
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
	var noprogress bool
	var clip bool

	flag.BoolVarP(&tee, "tee", "t", false, "Displays stdin to stdout")
	flag.BoolVarP(&noprogress, "noprogress", "n", false, "Never display a progress bar")
	flag.BoolVarP(&clip, "clipboard", "c", false, "Copy the returned URL directly to the clipboard (needs xclip or xsel)")
	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		var f *os.File
		var fi os.FileInfo

		if f, err = os.Open(args[0]); err != nil {
			fmt.Println("Could not open", args[0])
			os.Exit(1)
		}
		if noprogress {
			datasource = f
		} else {
			if fi, err = f.Stat(); err != nil {
				fmt.Println("Could not stat", args[0])
				os.Exit(1)
			}
			bar := pb.New64(fi.Size()).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
			bar.ShowPercent = true
			bar.ShowSpeed = true
			bar.ShowTimeLeft = true
			bar.Start()
			datasource = io.TeeReader(f, bar)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		datasource = os.Stdin
	}
	if tee {
		datasource = io.TeeReader(datasource, os.Stdout)
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
	if clip {
		clipboard.WriteAll(string(ret))
		fmt.Print("Copied URL to clipboard\n")
	} else {
		fmt.Print(string(ret))
	}
}
