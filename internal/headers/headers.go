package headers

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Headers map[string]string

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h Headers) GetContentLength() (int, error) {
	contentLengthString := h.Get("Content-Length")
	if contentLengthString == "" {
		return 0, nil
	}
	contentLength, err := strconv.Atoi(contentLengthString)
	if err != nil {
		return 0, fmt.Errorf("content length could not be parsed to integer: %s\n", err)
	}
	return contentLength, nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	firstNewLine := bytes.Index(data, []byte("\r\n"))
	if firstNewLine == -1 {
		return 0, false, nil
	}
	if firstNewLine == 0 {
		return 2, true, nil
	}

	requestLine := data[:firstNewLine]
	trimmedRL := bytes.Trim(requestLine, " ")
	if bytes.Count(trimmedRL, []byte(" ")) != 1 {
		return 0, false, fmt.Errorf("error: extra spacing in key values")
	}
	endOfKeyValue := bytes.IndexByte(trimmedRL, ':')
	keyByte := trimmedRL[:endOfKeyValue]
	if !keyIsValid(keyByte) {
		return 0, false, fmt.Errorf("error key:%s violates constraints", keyByte)
	}
	key := strings.ToLower(string(keyByte))
	value := string(trimmedRL[endOfKeyValue+2:])
	if existing, ok := h[key]; ok {
		value = fmt.Sprintf("%s, %s", existing, value)
	}
	h[key] = value
	return firstNewLine + 2, false, nil
}

const SPECIALCHARS = "!#$%&'*-.^_`|~"

func keyIsValid(key []byte) bool {
	for _, value := range key {
		if !((value >= 'A' && value <= 'Z') ||
			(value >= 'a' && value <= 'z') ||
			(value >= '0' && value <= '9') ||
			strings.IndexByte(SPECIALCHARS, value) >= 0) {
			return false
		}
	}
	return true
}

func NewHeaders() Headers {
	return Headers{}
}
