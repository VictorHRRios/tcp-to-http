package response

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/VictorHRRios/http/internal/headers"
)

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

const (
	WriterStatusLine = iota
	WriterStatusHeaders
	WriterStatusBody
)

type Writer struct {
	WriterState int
	Conn        net.Conn
	Headers     headers.Headers
	StatusCode  StatusCode
	ChunkSize   int
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	var buff bytes.Buffer
	var err error
	_, err = buff.WriteString(fmt.Sprintf("%X\r\n", w.ChunkSize))
	if err != nil {
		return 0, err
	}
	_, err = buff.Write(p)
	if err != nil {
		return 0, err
	}
	_, err = buff.Write([]byte("\r\n"))
	if err != nil {
		return 0, err
	}

	n, err := w.WriteBody(buff.Bytes())
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	var buff bytes.Buffer
	var err error
	_, err = buff.WriteString("0\r\n")
	if err != nil {
		return 0, err
	}
	_, err = buff.WriteString("\r\n")
	if err != nil {
		return 0, err
	}
	n, err := w.WriteBody(buff.Bytes())
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.WriterState != WriterStatusLine {
		return fmt.Errorf("Writer in wrong order got=%d, expected=%d\n", w.WriterState, WriterStatusLine)
	}
	w.StatusCode = statusCode
	w.WriterState = WriterStatusHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.WriterState != WriterStatusHeaders {
		return fmt.Errorf("Writer in wrong order got=%d, expected=%d\n", w.WriterState, WriterStatusHeaders)
	}
	w.WriterState = WriterStatusBody
	w.Headers = headers
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.WriterState != WriterStatusBody {
		return 0, fmt.Errorf("Writer in wrong order got=%d, expected=%d\n", w.WriterState, WriterStatusBody)
	}
	if err := WriteStatusLine(w.Conn, w.StatusCode); err != nil {
		return 0, err
	}
	if _, ok := w.Headers["transfer-encoding"]; !ok {
		w.Headers.Set("content-length", strconv.Itoa(len(p)))
	}
	if err := WriteHeaders(w.Conn, w.Headers); err != nil {
		return 0, err
	}
	n, err := w.Conn.Write(p)
	if err != nil {
		return 0, err
	}
	return n, nil
}

const HTTPV = "HTTP/1.1"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOk:
		if _, err := w.Write([]byte(HTTPV + " 200 OK\r\n")); err != nil {
			return err
		}
		return nil
	case StatusBadRequest:
		if _, err := w.Write([]byte(HTTPV + " 400 Bad Request\r\n")); err != nil {
			return err
		}
		return nil
	case StatusInternalServerError:
		if _, err := w.Write([]byte(HTTPV + " 500 Internal Server Error\r\n")); err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = strconv.Itoa(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		if _, err := w.Write([]byte(key + ": " + value + "\r\n")); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("\r\n")); err != nil {
		return err
	}
	return nil
}
