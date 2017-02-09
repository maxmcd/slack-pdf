package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("POST")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	_ = mediaType
	if err != nil {
		fmt.Println(err)
	}
	mr := multipart.NewReader(req.Body, params["boundary"])

	// cmd := exec.Command("unoconv", "-f", "pdf", "--stdin", "--stdout")
	cmd := exec.Command("cat")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	var filename string

	go func() {
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				stdin.Close()
				break
			}
			if err != nil {
				// panic(err)
				fmt.Println(err)
				break
			}
			if part.FormName() == "file" {
				filename = part.FileName() + ".pdf"
				w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
				io.Copy(stdin, part)
			}
		}
	}()

	io.Copy(w, stdout)

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
	fmt.Println(filename)

	// _, _ = mr, mediaType
}
