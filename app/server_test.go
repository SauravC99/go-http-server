package app

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		request  string
		expected RequestHeaders
	}{
		{
			request: "GET / HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: test-agent\r\n\r\n",
			expected: RequestHeaders{
				Method: "GET",
				Path:   "/",
				Host:   "localhost:4221",
				Agent:  "test-agent",
			},
		},
		{
			request: "GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: test-agent\r\nContent-Type: text/plain\r\nAccept-Encoding: gzip\r\n\r\nbody content",
			expected: RequestHeaders{
				Method:         "GET",
				Path:           "/index.html",
				Host:           "localhost:4221",
				Agent:          "test-agent",
				ContentType:    "text/plain",
				AcceptEncoding: "gzip",
				Body:           "body content",
			},
		},
		{
			request: "POST /files/hello.txt HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: test-agent\r\nAccept-Encoding: gzip\r\nContent-Length: 11\r\nContent-Type: application/octet-stream\r\n\r\n",
			expected: RequestHeaders{
				Method:         "POST",
				Path:           "/files/hello.txt",
				Host:           "localhost:4221",
				Agent:          "test-agent",
				ContentLength:  "11",
				ContentType:    "application/octet-stream",
				AcceptEncoding: "gzip",
			},
		},
	}

	for _, tt := range tests {
		req := []byte(tt.request)
		got, err := parseHeaders(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if fmt.Sprintf("%T", *got) != fmt.Sprintf("%T", RequestHeaders{}) {
			t.Fatalf("return value is not type RequestHeaders{}. got=%T", got)
		}

		if got.Inspect() != tt.expected.Inspect() {
			t.Errorf("expected response\n %+v, got\n %+v", got.Inspect(), tt.expected.Inspect())
		}
	}
}

func TestParseFilePath(t *testing.T) {
	tests := []struct {
		directory string
		path      string
		expected  string
	}{
		{
			directory: "/store",
			path:      "/files/hello.txt",
			expected:  "/store/hello.txt",
		},
		{
			directory: "/store/",
			path:      "/files/hello.txt",
			expected:  "/store/hello.txt",
		},
		{
			directory: "/",
			path:      "/files/hello.txt",
			expected:  "/hello.txt",
		},
	}

	for _, tt := range tests {
		got := parseFilePath(&tt.directory, tt.path)
		if got != tt.expected {
			t.Errorf("path has wrong value. got=%s, want=%s", got, tt.expected)
		}
	}
}

type mockConnection struct {
	data []byte
}

// LocalAddr implements net.Conn.
func (m mockConnection) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr implements net.Conn.
func (m mockConnection) RemoteAddr() net.Addr {
	return nil
}

// SetDeadline implements net.Conn.
func (m mockConnection) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline implements net.Conn.
func (m mockConnection) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline implements net.Conn.
func (m mockConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

func (m mockConnection) Read(b []byte) (int, error) {
	return 0, nil
}

func (m *mockConnection) Write(b []byte) (int, error) {
	m.data = append(m.data, b...)
	return 0, nil
}

func (m mockConnection) Close() error {
	return nil
}

func (m mockConnection) Inspect() string {
	return string(m.data)
}

func TestRespond200Plain(t *testing.T) {
	conn := &mockConnection{}
	want := STATUS_200_OK + CONTENT_PLAIN + CONTENT_LENGTH + "4" + END_HEADER_BLOCK + "test"

	respond200Plain(conn, "test")

	if conn.Inspect() != want {
		t.Errorf("respond200Plain() = %v, want %v", conn.Inspect(), want)
	}
}

func TestRespond200App(t *testing.T) {
	conn := &mockConnection{}
	want := STATUS_200_OK + CONTENT_APP + CONTENT_LENGTH + "8" + END_HEADER_BLOCK + "test app"

	respond200App(conn, "test app")

	if conn.Inspect() != want {
		t.Errorf("respond200App() = %v, want %v", conn.Inspect(), want)
	}
}

func TestRespond200Encode(t *testing.T) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	writer.Write([]byte("test encode"))
	writer.Close()
	encoded := buffer.String()
	encodedLen := strconv.Itoa(len(encoded))

	conn := &mockConnection{}
	want := STATUS_200_OK + CONTENT_ENCODING + "gzip" + END_HEADER_LINE + CONTENT_PLAIN + CONTENT_LENGTH + encodedLen + END_HEADER_BLOCK + encoded

	respond200Encode(conn, "test encode", "gzip")

	if conn.Inspect() != want {
		t.Errorf("respond200Encode() = %v, want %v", conn.Inspect(), want)
	}
}

func TestRespond201(t *testing.T) {
	conn := &mockConnection{}
	want := STATUS_201_CREATED + CONTENT_TYPE + "app/octet-stream" + END_HEADER_LINE + LOCATION_HEADER + "/files/test.txt" + END_HEADER_BLOCK + "this is test"

	respond201(conn, "/files/test.txt", "app/octet-stream", "this is test")

	if conn.Inspect() != want {
		t.Errorf("respond201() = %v, want %v", conn.Inspect(), want)
	}
}

func TestRespond404(t *testing.T) {
	conn := &mockConnection{}
	want := STATUS_404_ERR + END_HEADER_LINE

	respond404(conn)

	if conn.Inspect() != want {
		t.Errorf("respond404() = %v, want %v", conn.Inspect(), want)
	}
}

func TestRespond405(t *testing.T) {
	conn := &mockConnection{}
	allowed := "Allow: GET, POST" + END_HEADER_LINE
	want := STATUS_405_NOTALLOW + allowed + END_HEADER_LINE

	respond405(conn)

	if conn.Inspect() != want {
		t.Errorf("respond405() = %v, want %v", conn.Inspect(), want)
	}
}

func TestRespond500(t *testing.T) {
	conn := &mockConnection{}
	want := STATUS_500_ERR + END_HEADER_LINE

	respond500(conn)

	if conn.Inspect() != want {
		t.Errorf("respond500() = %v, want %v", conn.Inspect(), want)
	}
}
