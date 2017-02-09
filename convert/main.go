package main

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kr/pty"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("POST")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}

func writeCmdOutput(res http.ResponseWriter, pipeReader *io.PipeReader) {
	buffer := make([]byte, 1024)
	for {
		n, err := pipeReader.Read(buffer)
		if err != nil {
			// pipeReader.Close()
			fmt.Println(err)
			break
		}

		data := buffer[0:n]
		res.Write(data)
		// if f, ok := res.(http.Flusher); ok {
		// 	f.Flush()
		// }
		//reset buffer
		for i := 0; i < n; i++ {
			buffer[i] = 0
		}
	}
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
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

	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}

	go func() {
		// if err := cmd.Start(); err != nil {
		// 	panic(err)
		// }

		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			io.Copy(stdin, p)
			fmt.Printf("%#v\n", p.Header)
			// w.Write(p)
			// slurp, err := ioutil.ReadAll(p)
			// if err != nil {
			//  log.Fatal(err)
			// }
		}

		// io.Copy(w, stdout)
		if err := cmd.Wait(); err != nil {
			panic(err)
		}
	}()
	io.Copy(w, f)
	time.Sleep(time.Second * 3)

	_, _ = mr, mediaType
}
