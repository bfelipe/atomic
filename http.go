package atomic

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

type Request struct {
	conn      net.Conn
	frames    []string
	StartLine string
	Headers   map[string]string
	Body      []byte
	Method    string
	Path      string
}

func (r *Request) parseFrames() error {
	buf := make([]byte, 1024)
	if _, err := r.conn.Read(buf); err != nil {
		return fmt.Errorf("Error while reading content %s\n", err)
	}
	r.frames = strings.Split(string(buf), "\r\n")
	return nil
}

func (r *Request) parseStartLine() error {
	if len(r.frames) == 0 {
		return fmt.Errorf("Error while parsing StartLine, request has no frames")
	}
	r.StartLine = r.frames[0]
	parts := strings.SplitAfter(r.StartLine, " ")
	if len(parts) != 3 {
		return fmt.Errorf("Invalid StarLine format %q", r.StartLine)
	}
	r.Method = strings.TrimSpace(parts[0])
	r.Path = strings.TrimSpace(parts[1])
	return nil
}

func (r *Request) parseHeaders() {
	r.Headers = map[string]string{}
	for _, header := range r.frames[1 : len(r.frames)-2] {
		pair := strings.Split(header, " ")
		if len(pair) < 2 {
			break
		}
		key := pair[0][:len(pair[0])-1]
		r.Headers[key] = pair[1]
	}
}

func (r *Request) parseBody() {
	r.Body = []byte(r.frames[len(r.frames)-1])
	r.Body = bytes.Trim(r.Body, "\x00")
}

func (r *Request) Decode(conn net.Conn) error {
	r.conn = conn
	if err := r.parseFrames(); err != nil {
		return err
	}
	if err := r.parseStartLine(); err != nil {
		return err
	}
	r.parseHeaders()
	r.parseBody()
	return nil
}
