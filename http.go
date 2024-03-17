package atomic

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
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

type StatusCode string

const (
	CONTINUE                        StatusCode = "100 Continue"
	SWITCHING_PROTOCOLS             StatusCode = "101 Switching Protocols"
	PROCESSING                      StatusCode = "102 Processing"
	EARLY_HINTS                     StatusCode = "103 Early Hints"
	OK                              StatusCode = "200 OK"
	CREATED                         StatusCode = "201 Created"
	ACCEPTED                        StatusCode = "202 Accepted"
	NON_AUTHORITATIVE_INFORMATION   StatusCode = "203 Non-Authoritative Information"
	NO_CONTENT                      StatusCode = "204 No Content"
	RESET_CONTENT                   StatusCode = "205 Reset Content"
	PARTIAL_CONTENT                 StatusCode = "206 Partial Content"
	MULTI_STATUS                    StatusCode = "207 Multi-Status"
	ALREADY_REPORTED                StatusCode = "208 Already Reported"
	IM_USED                         StatusCode = "226 IM Used"
	MULTIPLE_CHOICES                StatusCode = "300 Multiple Choices"
	MOVED_PERMANENTLY               StatusCode = "301 Moved Permanently"
	FOUND                           StatusCode = "302 Found"
	SEE_OTHER                       StatusCode = "303 See Other"
	NOT_MODIFIED                    StatusCode = "304 Not Modified"
	TEMPORARY_REDIRECT              StatusCode = "307 Temporary Redirect"
	PERMANENT_REDIRECT              StatusCode = "308 Permanent Redirect"
	BAD_REQUEST                     StatusCode = "400 Bad Request"
	UNAUTHORIZED                    StatusCode = "401 Unauthorized"
	PAYMENT_REQUIRED                StatusCode = "402 Payment Required"
	FORBIDDEN                       StatusCode = "403 Forbidden"
	NOT_FOUND                       StatusCode = "404 Not Found"
	METHOD_NOT_ALLOWED              StatusCode = "405 Method Not Allowed"
	NOT_ACCEPTABLE                  StatusCode = "406 Not Acceptable"
	PROXY_AUTHENTICATION_REQUIRED   StatusCode = "407 Proxy Authentication Required"
	REQUEST_TIMEOUT                 StatusCode = "408 Request Timeout"
	CONFLICT                        StatusCode = "409 Conflict"
	GONE                            StatusCode = "410 Gone"
	LENGTH_REQUIRED                 StatusCode = "411 Length Required"
	PRECONDITION_FAILED             StatusCode = "412 Precondition Failed"
	CONTENT_TOO_LARGE               StatusCode = "413 Content Too Large"
	URI_TOO_LONG                    StatusCode = "414 URI Too Long"
	UNSUPPORTED_MEDIA_TYPE          StatusCode = "415 Unsupported Media Type"
	RANGE_NOT_SATISFIABLE           StatusCode = "416 Range Not Satisfiable"
	EXPECTATION_FAILED              StatusCode = "417 Expectation Failed"
	IM_A_TEAPOT                     StatusCode = "418 I'm a teapot"
	MISDIRECTED_REQUEST             StatusCode = "421 Misdirected Request"
	UNPROCESSABLE_CONTENT           StatusCode = "422 Unprocessable Content"
	LOCKED                          StatusCode = "423 Locked"
	FAILED_DEPENDENCY               StatusCode = "424 Failed Dependency"
	TOO_EARLY                       StatusCode = "425 Too Early"
	UPGRADE_REQUIRED                StatusCode = "426 Upgrade Required"
	PRECONDITION_REQUIRED           StatusCode = "428 Precondition Required"
	TOO_MANY_REQUESTS               StatusCode = "429 Too Many Requests"
	REQUEST_HEADER_FIELDS_TOO_LARGE StatusCode = "431 Request Header Fields Too Large"
	UNAVAILABLE_FOR_LEGAL_REASONS   StatusCode = "451 Unavailable For Legal Reasons"
	INTERNAL_SERVER_ERROR           StatusCode = "500 Internal Server Error"
	NOT_IMPLEMENTED                 StatusCode = "501 Not Implemented"
	BAD_GATEWAY                     StatusCode = "502 Bad Gateway"
	SERVICE_UNAVAILABLE             StatusCode = "503 Service Unavailable"
	GATEWAY_TIMEOUT                 StatusCode = "504 Gateway Timeout"
	HTTP_VERSION_NOT_SUPPORTED      StatusCode = "505 HTTP Version Not Supported"
	VARIANT_ALSO_NEGOTIATES         StatusCode = "506 Variant Also Negotiates"
	INSUFFICIENT_STORAGE            StatusCode = "507 Insufficient Storage"
	LOOP_DETECTED                   StatusCode = "508 Loop Detected"
	NOT_EXTENDED                    StatusCode = "510 Not Extended"
	NETWORK_AUTHENTICATION_REQUIRED StatusCode = "511 Network Authentication Required"
)

type Response struct {
	headers    map[string]string
	body       string
	statusCode StatusCode
}

func (r *Response) initHeaders() {
	if r.headers == nil {
		r.headers = map[string]string{}
	}
}

func (r *Response) SetHeader(key string, value string) *Response {
	r.initHeaders()
	r.headers[key] = value
	return r
}

func (r *Response) SetBody(content string, contentType string) *Response {
	r.body = content
	r.initHeaders()
	r.headers["Content-Type"] = contentType
	r.headers["Content-Length"] = strconv.Itoa(len(content))
	return r
}

func (r *Response) SetStatusCode(statusCode StatusCode) *Response {
	r.statusCode = statusCode
	return r
}

func (r Response) String() string {
	var str strings.Builder
	str.WriteString("HTTP/1.1")
	str.WriteString(" ")
	str.WriteString(string(r.statusCode))
	str.WriteString("\\r\\n")
	var headers string
	for k, v := range r.headers {
		headers += k + ": " + v + "\\r\\n"
	}
	if len(headers) != 0 {
		headers += "\\r\\n"
	} else {
		headers += "\\r\\n\\r\\n"
	}
	str.WriteString(headers)
	str.WriteString(r.body)
	return str.String()
}

func (r Response) Enconde() []byte {
	return []byte(r.String())
}
