package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("POST")
	fmt.Println("I'm up")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	// mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	// _ = mediaType
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// mr := multipart.NewReader(req.Body, params["boundary"])
	var filename string
	file, header, err := req.FormFile("file")
	// fmt.Println(header)
	_ = header

	// cmd := exec.Command("unoconv", "-f", "pdf", "--stdin", "--stdout")
	cmd := exec.Command("cat")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	fmt.Println("io.copy")
	go func() {
		_, err := io.Copy(stdin, file)
		if err != nil {
			panic(err)
		}
	}()

	fmt.Println("initate copy to sdtout")
	_, err = io.Copy(w, stdout)
	if err != nil {
		panic(err)
	}
	fmt.Println("yo")
	fmt.Printf("%#v\n", cmd.ProcessState)

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
	fmt.Println(filename)

	// _, _ = mr, mediaType
}
