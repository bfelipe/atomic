package atomic

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
)

type Request struct {
	frames    []string
	startLine string
	headers   map[string]string
	body      []byte
	method    string
	path      string
}

func (r *Request) parseFrames(conn net.Conn) error {
	buf := make([]byte, 1024)
	numOfBytes, err := conn.Read(buf)
	if err != nil {
		if numOfBytes == 0 {
			return errors.New("connection closed prematurely")
		}
		return fmt.Errorf("error while reading connection %s", err)
	}
	r.frames = strings.Split(string(buf), "\r\n")
	return nil
}

func (r *Request) parseStartLine() error {
	if len(r.frames) == 0 {
		return errors.New("request has no frames")
	}
	r.startLine = r.frames[0]
	parts := strings.SplitAfter(r.startLine, " ")
	if len(parts) != 3 {
		return errors.New("invalid start line format")
	}
	r.method = strings.TrimSpace(parts[0])
	r.path = strings.TrimSpace(parts[1])
	return nil
}

func (r *Request) parseHeaders() {
	r.headers = map[string]string{}
	for _, header := range r.frames[1:] {
		pair := strings.Split(header, " ")
		if len(pair) < 2 {
			break
		}
		key := pair[0][:len(pair[0])-1]
		r.headers[key] = pair[1]
	}
}

func (r *Request) parseBody() {
	if r.frames[len(r.frames)-1] != "" {
		r.body = []byte(r.frames[len(r.frames)-1])
		r.body = bytes.Trim(r.body, "\x00")
	}
}

func (r *Request) Decode(conn net.Conn) error {
	if err := r.parseFrames(conn); err != nil {
		return err
	}
	if err := r.parseStartLine(); err != nil {
		return err
	}
	r.parseHeaders()
	r.parseBody()
	return nil
}

func (r Request) Body() []byte {
	return r.body
}

func (r Request) Headers() map[string]string {
	return r.headers
}

func (r Request) Header(key string) string {
	return r.headers[key]
}

func (r Request) Method() string {
	return r.method
}

func (r Request) Path() string {
	return r.path
}

func (r Request) StartLine() string {
	return r.startLine
}
