package atomic_test

import (
	"testing"

	"gitlab.com/bfelipe/atomic"
	"gitlab.com/bfelipe/atomic/mock"
)

func TestRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                      string
		conn                      string
		expectedStartLine         string
		expectedPath              string
		expectedMethod            string
		expectedHeaders           map[string]string
		expectedBody              string
		expectErr                 bool
		expectedErrorMessage      string
		expectInvalidStartLineErr bool
		expectNoFramesErr         bool
	}{
		{
			name:              "Valid GET request",
			conn:              "GET /path HTTP/1.1\r\nHeader1: Value1\r\n\r\n",
			expectedStartLine: "GET /path HTTP/1.1",
			expectedPath:      "/path",
			expectedMethod:    "GET",
			expectedHeaders:   map[string]string{"Header1": "Value1"},
			expectErr:         false,
		},
		{
			name:              "Valid POST request",
			conn:              "POST /path HTTP/1.1\r\nHeader1: Value1\r\n\r\nbody message",
			expectedStartLine: "POST /path HTTP/1.1",
			expectedPath:      "/path",
			expectedMethod:    "POST",
			expectedHeaders:   map[string]string{"Header1": "Value1"},
			expectedBody:      "body message",
			expectErr:         false,
		},
		{
			name:              "Valid PUT request",
			conn:              "PUT /path HTTP/1.1\r\nHeader1: Value1\r\n\r\nbody message",
			expectedStartLine: "PUT /path HTTP/1.1",
			expectedPath:      "/path",
			expectedMethod:    "PUT",
			expectedHeaders:   map[string]string{"Header1": "Value1"},
			expectedBody:      "body message",
			expectErr:         false,
		},
		{
			name:              "Valid DELETE request",
			conn:              "DELETE /path HTTP/1.1\r\nHeader1: Value1\r\n\r\n",
			expectedStartLine: "DELETE /path HTTP/1.1",
			expectedPath:      "/path",
			expectedMethod:    "DELETE",
			expectedHeaders:   map[string]string{"Header1": "Value1"},
			expectErr:         false,
		},
		{
			name:              "Valid PATCH request",
			conn:              "PATCH /path HTTP/1.1\r\nHeader1: Value1\r\n\r\nbody message",
			expectedStartLine: "PATCH /path HTTP/1.1",
			expectedPath:      "/path",
			expectedMethod:    "PATCH",
			expectedHeaders:   map[string]string{"Header1": "Value1"},
			expectedBody:      "body message",
			expectErr:         false,
		},
		{
			name:              "Missing headers",
			conn:              "GET /path HTTP/1.1\r\n\r\n",
			expectedStartLine: "GET /path HTTP/1.1",
			expectedPath:      "/path",
			expectedMethod:    "GET",
			expectedHeaders:   nil,
			expectErr:         false,
		},
		{
			name:                      "Empty conn",
			conn:                      "",
			expectedErrorMessage:      "invalid start line format",
			expectErr:                 true,
			expectInvalidStartLineErr: true,
		},
		{
			name:                      "Invalid start line format",
			conn:                      " GET  /path  \r\nHeader1: Value1\r\n\r\n",
			expectedErrorMessage:      "invalid start line format",
			expectErr:                 true,
			expectInvalidStartLineErr: true,
		},
		{
			name:                      "Start line without spaces",
			conn:                      "GET/pathHTTP/1.1\r\n\r\n",
			expectedErrorMessage:      "invalid start line format",
			expectErr:                 true,
			expectInvalidStartLineErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			conn := mock.NewMockConn(tc.conn)
			var req atomic.Request
			err := req.Decode(conn)

			if tc.expectErr && err == nil {
				t.Error("expected an error, but got none")
				return
			}

			if !tc.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tc.expectErr {
				if tc.expectInvalidStartLineErr == true {
					if err.Error() != tc.expectedErrorMessage {
						t.Errorf("unexpected error got: %v expected: %v", err.Error(), tc.expectedErrorMessage)
					}
				}
				return
			}

			if !tc.expectErr {
				if req.StartLine() != tc.expectedStartLine {
					t.Errorf("unexpected start line got: %v expected: %v", req.StartLine(), tc.expectedStartLine)
					return
				}

				if req.Path() != tc.expectedPath {
					t.Errorf("unexpected path got: %v expected: %v", req.Path(), tc.expectedPath)
					return
				}

				if req.Method() != tc.expectedMethod {
					t.Errorf("unexpected method got: %v expected: %v", req.Method(), tc.expectedMethod)
					return
				}

				if len(tc.expectedHeaders) > 0 {
					for k, v := range req.Headers() {
						if req.Header(k) != v {
							t.Errorf("unexpected header %v got: %v expected: %v", k, req.Header(k), tc.expectedHeaders[k])
							return
						}
					}
				}

				if string(req.Body()) != tc.expectedBody {
					t.Errorf("unexpected body got: %v expected: %v", string(req.Body()), tc.expectedPath)
					return
				}
			}
		})
	}
}
