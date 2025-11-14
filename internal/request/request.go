package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/VictorHRRios/http/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	State       int
	Headers     headers.Headers
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func NewRequest() *Request {
	return &Request{
		RequestLine: RequestLine{},
		State:       requestStateInitialized,
		Headers:     headers.NewHeaders(),
	}
}

const (
	_ = iota
	requestStateInitialized
	requestStateParsingHeaders
	requestStateDone
)

const (
	HTTP_VERSION = "1.1"
	bufferSize   = 8
)

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}
func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = requestStateDone
		}
		return n, nil
	case requestStateInitialized:
		parsed, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = parsed.RequestLine
		r.State = requestStateParsingHeaders
		return n, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in done state")
	default:
		return 0, fmt.Errorf("state unknown")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	request := NewRequest()
	for request.State != requestStateDone {
		if readToIndex == len(buf) {
			buf = grow(buf)
		}
		read, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.State = requestStateDone
				break
			}
			return nil, err
		}
		readToIndex += read
		n, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if n > 0 {
			copy(buf, buf[n:readToIndex])
			readToIndex = readToIndex - n
		}
	}
	return request, nil
}

func grow(buf []byte) []byte {
	bufCopy := make([]byte, len(buf)*2)
	copy(bufCopy, buf)
	return bufCopy
}

func parseRequestLine(request []byte) (*Request, int, error) {
	firstNewLine := bytes.Index(request, []byte("\r\n"))
	if firstNewLine == -1 {
		return nil, 0, nil
	}
	requestLine := request[:firstNewLine]
	slicedRequestLine := bytes.Split(requestLine, []byte(" "))
	if len(slicedRequestLine) != 3 {
		return nil, 0, fmt.Errorf("error in parts or request line, expected=%v, actual=%v, extra=%v", 3, len(slicedRequestLine), string(request))
	}

	if !hasAllCaps(slicedRequestLine[0]) {
		return nil, 0, fmt.Errorf("error in method, needs to be in all caps")
	}

	httpVersion := bytes.TrimPrefix(slicedRequestLine[2], []byte("HTTP/"))
	if !bytes.Equal(httpVersion, []byte("1.1")) {
		return nil, 0, fmt.Errorf("error in http version, needs to be 1.1, got=%T%s", httpVersion, httpVersion)
	}

	return &Request{RequestLine: RequestLine{
		Method:        string(slicedRequestLine[0]),
		RequestTarget: string(slicedRequestLine[1]),
		HttpVersion:   string(httpVersion),
	}, State: requestStateInitialized}, firstNewLine + 2, nil
}

func hasAllCaps(request []byte) bool {
	for _, value := range request {
		if value < 'A' || value > 'Z' {
			return false
		}
	}
	return true
}
