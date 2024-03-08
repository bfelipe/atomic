package mock

import (
	"net"
	"time"
)

type MockConn struct {
	input string
}

func NewMockConn(input string) net.Conn {
	return &MockConn{input: input}
}

func (m *MockConn) Read(b []byte) (int, error) {
	n := copy(b, []byte(m.input))
	return n, nil
}

func (m *MockConn) Write(b []byte) (int, error) {
	return len(b), nil
}

func (m *MockConn) Close() error {
	return nil
}

func (m *MockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 12345}
}

func (m *MockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("10.0.0.1"), Port: 54321}
}

func (m *MockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}
