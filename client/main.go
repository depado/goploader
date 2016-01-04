package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/cheggaaa/pb"
	flag "github.com/ogier/pflag"
)

const (
	service = "https://up.depado.eu/"
	method  = "POST"
)

func main() {
	var err error
	var datasource io.Reader
	var bar *pb.ProgressBar

	var tee bool
	var progress bool
	var clip bool
	var name string

	flag.BoolVarP(&tee, "tee", "t", false, "Displays stdin to stdout")
	flag.BoolVarP(&progress, "progress", "p", false, "Displays a progress bar")
	flag.BoolVarP(&clip, "clipboard", "c", false, "Copy the returned URL directly to the clipboard (needs xclip or xsel)")
	flag.StringVarP(&name, "name", "n", "", "Specify the filename you want")
	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		var f *os.File
		var fi os.FileInfo

		if f, err = os.Open(args[0]); err != nil {
			fmt.Println("Could not open", args[0])
			os.Exit(1)
		}
		defer f.Close()
		if fi, err = f.Stat(); err != nil {
			fmt.Println("Could not stat", args[0])
			os.Exit(1)
		}
		if name == "" {
			name = fi.Name()
		}
		datasource = f
		if progress {
			bar = pb.New64(fi.Size()).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
			bar.ShowPercent = true
			bar.ShowSpeed = true
			bar.ShowTimeLeft = true
			bar.Start()
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		if name == "" {
			name = "stdin"
		}
		datasource = os.Stdin
	}
	if tee {
		datasource = io.TeeReader(datasource, os.Stdout)
	}

	r, w := io.Pipe()
	multipartWriter := multipart.NewWriter(w)
	contentType := multipartWriter.FormDataContentType()
	go func() {
		var part io.Writer
		defer w.Close()
		if part, err = multipartWriter.CreateFormFile("file", name); err != nil {
			log.Fatal(err)
		}
		multiWriter := part
		if progress {
			multiWriter = io.MultiWriter(part, bar)
		}
		if _, err = io.Copy(multiWriter, datasource); err != nil {
			log.Fatal(err)
		}
		if err = multipartWriter.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	resp, err := http.Post(service, contentType, r)
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
