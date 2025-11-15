package server

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/VictorHRRios/http/internal/request"
	"github.com/VictorHRRios/http/internal/response"
)

const (
	serverStateClosed = iota
	serverStateOpen
)

type Server struct {
	listener    net.Listener
	closed      *atomic.Bool
	handlerFunc Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w response.Writer, req *request.Request)

func NewServer(listener net.Listener, handler Handler) *Server {
	var closed atomic.Bool
	closed.Store(false)
	return &Server{
		listener:    listener,
		handlerFunc: handler,
		closed:      &closed,
	}
}

func (hErr *HandlerError) Write(writer io.Writer, handlerError HandlerError) {
	header := response.GetDefaultHeaders(len(handlerError.Message))
	if err := response.WriteStatusLine(writer, handlerError.StatusCode); err != nil {
		fmt.Print(err)
	}
	if err := response.WriteHeaders(writer, header); err != nil {
		fmt.Print(err)
	}
	writer.Write([]byte(handlerError.Message))
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("an error has occurred: %s\n", err)
	}
	newServer := NewServer(listener, handler)
	go newServer.listen()
	return newServer, nil
}

func (s *Server) listen() {
	fmt.Println("Listening...")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Print("an error has occurred: ", err)
			continue
		}
		fmt.Printf("established connection with %v\n", conn.LocalAddr())

		go s.handle(conn)
	}
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if err := s.listener.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	request, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Print("an error has occurred: ", err)
		return
	}
	writer := response.Writer{
		Conn:       conn,
		Headers:    response.GetDefaultHeaders(0),
		StatusCode: response.StatusOk,
	}
	s.handlerFunc(writer, request)
}
