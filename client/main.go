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

	"github.com/Depado/goploader/client/screenshot"
	"github.com/atotto/clipboard"
	"github.com/cheggaaa/pb"
	flag "github.com/ogier/pflag"

	"github.com/Depado/goploader/client/conf"
)

var (
	bar     *pb.ProgressBar
	name    string
	verbose bool
	service = "https://up.depado.eu/"
)

func debugf(a ...interface{}) {
	if verbose {
		fmt.Println(a...)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initBar(f *os.File) {
	fi, err := f.Stat()
	if err != nil {
		fmt.Println("Could not stat", f.Name())
		os.Exit(1)
	}
	bar = pb.New64(fi.Size()).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowPercent = true
	bar.ShowSpeed = true
	bar.ShowTimeLeft = true
	bar.Start()
}

func initUnknownBar() {
	bar = pb.New64(0).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.ShowCounters = true
	bar.ShowBar = false
	bar.Start()
}

func main() {
	var err error
	var datasource io.Reader
	var file string

	var tee bool
	var progress bool
	var clip bool
	var argname string
	var screen bool
	var delay time.Duration
	var window bool
	var lifetime string

	flag.BoolVarP(&tee, "tee", "t", false, "Displays stdin to stdout")
	flag.BoolVarP(&progress, "progress", "p", false, "Displays a progress bar")
	flag.BoolVarP(&clip, "clipboard", "c", false, "Copy the returned URL directly to the clipboard (needs xclip or xsel)")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Activates the debug mode")
	flag.StringVarP(&argname, "name", "n", "", "Specify the filename you want")
	flag.StringVarP(&lifetime, "lifetime", "l", "1d", "Specify the lifetime of your file (30m, 1h, 6h, 1d, 1w)")
	flag.BoolVarP(&screen, "screenshot", "s", false, "Screenshot and uploads your current screen (Need the `import` command)")
	flag.DurationVarP(&delay, "delay", "d", 0, "Define a delay before the program executes (including taking the screenshot)")
	flag.BoolVarP(&window, "window", "w", false, "Click on the window you want to screenshot (only works with -s/--screenshot option)")

	flag.Parse()
	args := flag.Args()
	debugf("Debug mode is activated")

	if err = conf.Load(); err != nil {
		log.Fatal(err)
	}

	if delay != 0 {
		debugf("Waiting", delay)
		time.Sleep(delay)
	}
	if len(args) > 0 {
		file = args[0]
	}
	if screen {
		debugf("Executing screenshot")
		file = "/tmp/tmp-goploader-screen.png"
		if err = screenshot.Do(file, window); err != nil {
			log.Fatal(err)
		}
		defer func() {
			debugf("Removing temporary screenshot", file)
			err = os.Remove(file)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	if file != "" {
		debugf("Main datasource is", file)
		f, err := os.Open(file)
		check(err)
		defer f.Close()
		name = f.Name()
		datasource = f
		if progress {
			debugf("Initialization of known progress bar")
			initBar(f)
		}
	} else {
		debugf("Main datasource is stdin")
		name = "stdin"
		datasource = os.Stdin
		if progress {
			debugf("Initalization of unknown progress bar")
			initUnknownBar()
		}
	}
	if tee {
		debugf("Main datasource is now a TeeReader")
		datasource = io.TeeReader(datasource, os.Stdout)
	}
	if argname != "" {
		debugf("Name was given as argument")
		name = argname
	}

	debugf("Initialization pipe")
	r, w := io.Pipe()
	multipartWriter := multipart.NewWriter(w)
	go func() {
		debugf("Started the goroutine to pipe data")
		var part io.Writer
		defer w.Close()
		defer multipartWriter.Close()
		multipartWriter.WriteField("duration", lifetime)
		if part, err = multipartWriter.CreateFormFile("file", name); err != nil {
			log.Fatal(err)
		}
		if progress {
			part = io.MultiWriter(part, bar)
		}
		if _, err = io.Copy(part, datasource); err != nil {
			log.Fatal(err)
		}
	}()

	debugf("Executing multipart post")
	resp, err := http.Post(conf.C.Service, multipartWriter.FormDataContentType(), r)
	check(err)
	defer resp.Body.Close()
	debugf("Multipart post is done, reading data")
	ret, err := ioutil.ReadAll(resp.Body)
	check(err)
	if clip {
		debugf("Copying to clipboard")
		clipboard.WriteAll(string(ret))
		fmt.Print("Copied URL to clipboard\n")
	} else {
		fmt.Print(string(ret))
	}
}
