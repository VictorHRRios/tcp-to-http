package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
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
