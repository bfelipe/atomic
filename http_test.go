package atomic_test

import (
	"testing"

	"gitlab.com/bfelipe/atomic"
	"gitlab.com/bfelipe/atomic/mock"
)

/**

    GET:
		Empty request
		StartLine without spaces
		StartLine with missing parts
        Invalid start line format.
        Empty request.

    POST:
        Valid POST request with JSON payload and headers.
        Valid POST request with form data and headers.
        Missing content type header for non-empty body.
        Invalid JSON payload.
        Empty request.

    PUT:
        Valid PUT request with JSON payload and headers.
        Valid PUT request with form data and headers.
        Missing content type header for non-empty body.
        Invalid JSON payload.
        Empty request.

    DELETE:
        Valid DELETE request with headers and body (including empty body).
        Missing headers but valid body.
        Invalid start line format.
        Empty request.

    PATCH:
        Valid PATCH request with JSON payload and headers.
        Valid PATCH request with form data and headers.
        Missing content type header for non-empty body.
        Invalid JSON payload.
        Empty request.


**/

func TestGetRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		conn         string
		expStartLine string
		expHeaders   map[string]string
		expBody      string
		expReqMethod string
		expReqPath   string
	}{
		{
			name:         "Valid GET request",
			conn:         "GET /path HTTP/1.1\r\nHeader1: Value1\r\n\r\n",
			expStartLine: "GET /path HTTP/1.1",
			expHeaders:   map[string]string{"Header1": "Value1"},
			expBody:      "",
			expReqMethod: "GET",
			expReqPath:   "/path",
		},
		{
			name:         "Missing Headers",
			conn:         "GET /path HTTP/1.1\r\n\r\n\r\n",
			expStartLine: "GET /path HTTP/1.1",
			expHeaders:   nil,
			expBody:      "",
			expReqMethod: "GET",
			expReqPath:   "/path",
		},
		{
			name:         "Invalid StartLine",
			conn:         "Invalid",
			expReqMethod: "",
			expReqPath:   "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := mock.NewMockConn(tc.conn)
			req := atomic.Request{}
			req.Decode(conn)

			if req.StartLine != tc.expStartLine {
				t.Errorf("Unexpected StartLine got %q want %q", req.StartLine, tc.expStartLine)
				return
			}

			if tc.expHeaders != nil {
				if req.Headers == nil {
					t.Errorf("Unexpected Header got %q want %q", req.Headers, tc.expHeaders)
					return
				}
				for k, _ := range tc.expHeaders {
					if tc.expHeaders[k] != req.Headers[k] {
						t.Errorf("Unexpected Header %q got %q want %q", k, req.Headers[k], tc.expHeaders[k])
						return
					}
				}
			}

			reqBody := string(req.Body)
			if reqBody != tc.expBody {
				t.Errorf("Unexpected Body got %q want %q", reqBody, tc.expBody)
				return
			}

			if req.Method != tc.expReqMethod {
				t.Errorf("Unexpected Request Method got %q want %q", req.Method, tc.expReqMethod)
				return
			}

			if req.Path != tc.expReqPath {
				t.Errorf("Unexpected Request Path got %q want %q", req.Path, tc.expReqPath)
				return
			}

		})
	}
}
