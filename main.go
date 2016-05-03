package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

		fmt.Fprintln(w, "Ok")
	})

	http.ListenAndServe(":http-alt", nil)
}
