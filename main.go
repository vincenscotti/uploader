package main

import (
	"github.com/tylerb/graceful"

	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const form = `<form action="/up" enctype="multipart/form-data" method="post"><input type="file" name="f" multiple><input type="submit"></form>`

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.Write([]byte(form))
	})

	http.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil || r.MultipartForm == nil {
			errorStr := ""
			if err != nil {
				errorStr = err.Error()
			}
			http.Error(w, "Invalid request "+errorStr, http.StatusBadRequest)
			return
		}

		for _, fh := range r.MultipartForm.File["f"] {
			outfile, err := os.Create(fh.Filename)
			if err != nil {
				http.Error(w, "Error saving file: "+err.Error(), http.StatusBadRequest)
				return
			}

			infile, err := fh.Open()
			if err != nil {
				http.Error(w, "Error opening file: "+err.Error(), http.StatusBadRequest)
				return
			}

			_, err = io.Copy(outfile, infile)
			if err != nil {
				http.Error(w, "Error saving file: "+err.Error(), http.StatusBadRequest)
				return
			}
		}
	})

	l, err := net.Listen("unix", "/var/run/uploader.socket")
	if err != nil {
		log.Fatalln(err)
	}

	server := graceful.Server{
		Timeout: 1 * time.Minute,
		Server:  &http.Server{},
	}

	log.Println(server.Serve(l))
	l.Close()
}
