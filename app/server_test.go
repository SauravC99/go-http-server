package app

import (
	"fmt"
	"testing"
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
