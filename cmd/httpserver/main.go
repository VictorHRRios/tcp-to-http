package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/VictorHRRios/http/internal/headers"
	"github.com/VictorHRRios/http/internal/request"
	"github.com/VictorHRRios/http/internal/response"
	"github.com/VictorHRRios/http/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w response.Writer, req *request.Request) {
		h := headers.NewHeaders()
		h.Set("content-type", "text/html")
		var buff bytes.Buffer
		after, found := strings.CutPrefix(req.RequestLine.RequestTarget, "/httpbin")
		if found {

			t := headers.NewHeaders()
			t.Set("X-Content-SHA256", "0")
			t.Set("X-Content-Length", "0")
			h.Del("content-length")
			h.Set("transfer-encoding", "chunked")
			w.WriteStatusLine(response.StatusOk)
			w.WriteHeaders(h)
			w.AddTrailers(t)
			resp, err := http.Get("https://httpbin.org" + after)
			if err != nil {
				fmt.Print(err)
			}
			defer resp.Body.Close()
			buff := make([]byte, 1024)
			buffRead := bytes.Buffer{}
			w.WriteBody([]byte{})
			for {
				n, err := resp.Body.Read(buff)
				w.ChunkSize = n
				if n > 0 {
					w.WriteChunkedBody(buff[:n])
				}
				buffRead.Write(buff[:n])
				if err != nil {
					if err == io.EOF {
						read := sha256.Sum256(buffRead.Bytes())
						w.Trailers.Set("X-Content-SHA256", hex.EncodeToString(read[:]))
						w.Trailers.Set("X-Content-Length", strconv.Itoa(buffRead.Len()))
						w.WriteChunkedBodyDone()
						break
					}
					fmt.Print(err)
				}
			}
			return
		}
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			buff.WriteString(htmlBodyBR)
			w.WriteStatusLine(response.StatusBadRequest)
		case "/myproblem":
			buff.WriteString(htmlBodyISE)
			w.WriteStatusLine(response.StatusInternalServerError)
		default:
			buff.WriteString(htmlBodyOk)
			w.WriteStatusLine(response.StatusOk)
		}
		w.WriteHeaders(h)
		_, err := w.WriteBody(buff.Bytes())
		if err != nil {
			fmt.Print(err)
		}
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

const htmlBodyISE = `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`

const htmlBodyBR = `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`
const htmlBodyOk = `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`
