package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nu7hatch/gouuid"

	"github.com/Depado/goploader/server/conf"
	"github.com/Depado/goploader/server/utils"
)

var db gorm.DB
var timeLimit = 2 * time.Hour

// ResourceEntry represents the data stored in the database
type ResourceEntry struct {
	gorm.Model
	Key string
}

func create(w http.ResponseWriter, r *http.Request) {
	var err error
	remote := r.Header.Get("x-forwarded-for")
	if r.Method == "GET" {
		log.Printf("[INFO][%s]\tIssued a GET request\n", remote)
		http.ServeFile(w, r, "client_linux_x86-64")
		return
	}
	log.Printf("[INFO][%s]\tReceiving data\n", remote)
	u, err := uuid.NewV4()
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring creation of uuid : %s\n", remote, err)
		http.Error(w, http.StatusText(503), 503)
		return
	}
	br := bufio.NewReaderSize(r.Body, 512)
	path := path.Join(conf.C.UploadDir, u.String())
	file, err := os.Create(path)
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring file creation : %s\n", remote, err)
		http.Error(w, http.StatusText(503), 503)
		return
	}
	defer file.Close()
	wr, err := io.Copy(file, br)
	if err != nil {
		log.Printf("[ERROR][%s]\tDuring writing file : %s\n", remote, err)
		http.Error(w, http.StatusText(503), 503)
		return
	}
	e := ResourceEntry{}
	e.Key = u.String()
	db.Create(&e)
	log.Printf("[INFO][%s]\tCreated %s file and entry (%v bytes written)\n", remote, u.String(), wr)
	fmt.Fprint(w, "http://"+conf.C.NameServer+"/view/"+u.String()+"\n")
	return
}

func view(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/view/"):]
	re := ResourceEntry{}
	db.Where(&ResourceEntry{Key: id}).First(&re)
	if re.Key == "" {
		log.Printf("[INFO][%s]\tNot found : %s", r.Header.Get("x-forwarded-for"), id)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	http.ServeFile(w, r, conf.C.UploadDir+re.Key)
}

func monit() {
	tc := time.NewTicker(1 * time.Minute)
	for {
		res := []ResourceEntry{}
		db.Find(&res, "created_at < ?", time.Now().Add(-timeLimit))
		db.Unscoped().Where("created_at < ?", time.Now().Add(-timeLimit)).Delete(&ResourceEntry{})
		if len(res) > 0 {
			log.Printf("[INFO][System]\tFlushing %d DB entries and files.\n", len(res))
		}
		for _, re := range res {
			err := os.Remove(path.Join(conf.C.UploadDir, re.Key))
			if err != nil {
				log.Printf("[ERROR][System]\tWhile deleting : %v", err)
			}
		}
		<-tc.C
	}
}

func main() {
	var err error

	confPath := flag.String("c", "conf.yml", "Local path to configuration file.")
	flag.Parse()

	if err = conf.Load(*confPath); err != nil {
		log.Fatal(err)
	}

	if err = utils.EnsureDir(conf.C.UploadDir); err != nil {
		log.Fatal(err)
	}

	db, err = gorm.Open("sqlite3", conf.C.DB)
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&ResourceEntry{})
	log.Printf("[INFO][System]\tStarted goploader server on port %d\n", conf.C.Port)
	go monit()
	log.Println("[INFO][System]\tStarted monitoring of files and db entries")

	http.HandleFunc("/", create)
	http.HandleFunc("/view/", view)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.C.Port), nil)
}
