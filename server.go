package main

import (
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

type PostStruct struct {
	Buffer string
}

func serve(port int) {
	virtual_fs := &assetfs.AssetFS{
		Asset: Asset,
		AssetDir: AssetDir,
		AssetInfo: AssetInfo}
	http.Handle("/static/", http.FileServer(virtual_fs))
	http.HandleFunc("/timeline/compose", guiHandler)
	http.HandleFunc("/timeline", handler)
	fmt.Printf("Listening on port %d\n" +
		"POST JSON sources to http://localhost:%d/timeline\n" +
		"Compose timelines at http://localhost:%d/timeline/compose\n", port, port, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func guiHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := Asset("static/index.html")
	fmt.Fprintf(w, "%s\n", string(bytes))
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlePost(&w, r)
	case "GET":
		handleGet(&w, r)
	}
}

func handleGet(w *http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(*w, "GET request\nRequest struct = %s\n", r)	
}

func handlePost(w *http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(*w, "Can't read POST request body: %s\n", err)
		return
	}

	//write to tmpfile
	tmpfile, err := ioutil.TempFile("", "timeline") //use const filePrefix?
	if (err != nil) {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(body); err != nil {
		log.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	//process the file
	var mu sync.Mutex
	ch := make(chan Result)
	go processFile(tmpfile.Name(), ch)
	mu.Lock()
	result := <-ch
	mu.Unlock()
	fmt.Println(result.Message)
}