package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// headers section \r\n\r\n body section

type RequestHeaders struct {
	Method string
	Path   string
	Host   string
	Agent  string
}

const STATUS_200_OK string = "HTTP/1.1 200 OK\r\n"
const STATUS_404_ERR string = "HTTP/1.1 404 Not Found\r\n"
const CONTENT_PLAIN string = "Content-Type: text/plain\r\n"
const END_HEADER_BLOCK string = "\r\n"

func main() {
	fmt.Println("Server started")

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port")
		os.Exit(1)
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Connection accept")

	connectAndRespond(conn)
}

func connectAndRespond(connection net.Conn) {
	headers, err := parseHeaders(connection)
	if err != nil {
		fmt.Println("Failed to read headers: ", err.Error())
	}

	fmt.Println(headers)

	// need 2 sets of \r\n for end of headers section
	if headers.Path == "/" {
		response := STATUS_200_OK + CONTENT_PLAIN + END_HEADER_BLOCK
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(headers.Path, "/echo/") {
		//echo the request in body
		response := STATUS_200_OK + CONTENT_PLAIN
		content := strings.TrimPrefix(headers.Path, "/echo/")
		length := "Content-Length: " + strconv.Itoa(len(content)) + "\r\n\r\n"

		response = response + length + content
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(headers.Path, "/user-agent") {
		response := STATUS_200_OK + CONTENT_PLAIN
		length := "Content-Length: " + strconv.Itoa(len(headers.Agent)) + "\r\n\r\n"
		body := headers.Agent

		response = response + length + body
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else {
		response := STATUS_404_ERR + END_HEADER_BLOCK
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	}

	connection.Close()
}

func parseHeaders(connection net.Conn) (*RequestHeaders, error) {
	recievedData := make([]byte, 1024)

	_, err := connection.Read(recievedData)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		return nil, err
	}

	request := strings.Split(string(recievedData), "\r\n")

	start_line := strings.Split(request[0], " ")
	host := strings.Split(request[1], " ")
	agent := strings.Split(request[2], " ")

	http_method := start_line[0]
	path := start_line[1]

	return &RequestHeaders{
		Method: http_method,
		Path:   path,
		Host:   host[1],
		Agent:  agent[1],
	}, nil
}

/*
func main2() {
	const STATUS_200_OK string = "HTTP/1.1 200 OK\r\n"
	const STATUS_404_ERR string = "HTTP/1.1 404 Not Found\r\n"
	const CONTENT_PLAIN string = "Content-Type: text/plain\r\n"
	const END_HEADER_BLOCK string = "\r\n"

	fmt.Println("Server started")

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port")
		os.Exit(1)
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Connection accept")

	recievedData := make([]byte, 1024) //may need to be bigger

	_, err = conn.Read(recievedData)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		os.Exit(1)
	}

	request := strings.Split(string(recievedData), "\r\n")

	start_line := strings.Split(request[0], " ")
	agent := strings.Split(request[2], " ")

	//http_method := start_line[0]
	path := start_line[1]

	// need 2 sets of \r\n for end of headers section
	if path == "/" {
		response := STATUS_200_OK + CONTENT_PLAIN + END_HEADER_BLOCK
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(path, "/echo/") {
		//echo the request in body
		response := STATUS_200_OK + CONTENT_PLAIN
		content := strings.TrimPrefix(path, "/echo/")
		length := "Content-Length: " + strconv.Itoa(len(content)) + "\r\n\r\n"

		response = response + length + content
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(path, "/user-agent") {
		response := STATUS_200_OK + CONTENT_PLAIN
		length := "Content-Length: " + strconv.Itoa(len(agent[1])) + "\r\n\r\n"
		body := agent[1]

		response = response + length + body
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else {
		response := STATUS_404_ERR + END_HEADER_BLOCK
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	}

	conn.Close()
}
*/
