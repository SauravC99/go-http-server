package main

import (
	"flag"
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
const CONTENT_APP string = "Content-Type: application/octet-stream\r\n"
const END_HEADER_LINE string = "\r\n"
const END_HEADER_BLOCK string = "\r\n\r\n"

func main() {
	fmt.Println("Server started")

	directoryPtr := flag.String("directory", "", "Directory for file")
	flag.Parse()

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port")
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Connection accept")
		go connectAndRespond(conn, directoryPtr)
	}
}

func connectAndRespond(connection net.Conn, directoryPtr *string) {
	headers, err := parseHeaders(connection)
	if err != nil {
		fmt.Println("Failed to read headers: ", err.Error())
	}

	// need 2 sets of \r\n for end of headers section
	if headers.Path == "/" {
		response := STATUS_200_OK + CONTENT_PLAIN + END_HEADER_LINE
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(headers.Path, "/echo/") {
		//echo the request in body
		response := STATUS_200_OK + CONTENT_PLAIN
		content := strings.TrimPrefix(headers.Path, "/echo/")
		length := "Content-Length: " + strconv.Itoa(len(content)) + END_HEADER_BLOCK

		response = response + length + content
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(headers.Path, "/user-agent") {
		response := STATUS_200_OK + CONTENT_PLAIN
		length := "Content-Length: " + strconv.Itoa(len(headers.Agent)) + END_HEADER_BLOCK
		body := headers.Agent

		response = response + length + body
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(headers.Path, "/files") {
		fPath := parseFilePath(directoryPtr, headers.Path)
		file, err := os.ReadFile(fPath)
		if err != nil {
			fmt.Println("Error reading file: ", err.Error())
			response := STATUS_404_ERR + END_HEADER_LINE
			_, err = connection.Write([]byte(response))
			if err != nil {
				fmt.Println("Error writing to connection: ", err.Error())
				os.Exit(1)
			}
			connection.Close()
			return
		}

		response := STATUS_200_OK + CONTENT_APP
		length := "Content-Length: " + strconv.Itoa(len(string(file))) + END_HEADER_BLOCK

		response = response + length + string(file)
		_, err = connection.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(headers.Path, "/test") {
		fmt.Println(headers)
		r := STATUS_200_OK + END_HEADER_LINE
		connection.Write([]byte(r))
	} else {
		response := STATUS_404_ERR + END_HEADER_LINE
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

	request := strings.Split(string(recievedData), END_HEADER_LINE)

	start_line := strings.Split(request[0], " ")
	http_method := start_line[0]
	path := start_line[1]

	host := ""
	agent := ""

	for _, line := range request {
		if strings.HasPrefix(line, "Host:") {
			temp := strings.Split(line, " ")
			host = temp[1]
		} else if strings.HasPrefix(line, "User-Agent:") {
			temp := strings.Split(line, " ")
			agent = temp[1]
		}
	}

	return &RequestHeaders{
		Method: http_method,
		Path:   path,
		Host:   host,
		Agent:  agent,
	}, nil
}

func parseFilePath(directoryPtr *string, path string) string {
	fileName := strings.TrimPrefix(path, "/files/")
	fullPath := ""

	// In case directory provided does not have '/' at end
	if strings.HasSuffix(*directoryPtr, "/") {
		fullPath = *directoryPtr + fileName
	} else {
		fullPath = *directoryPtr + "/" + fileName
	}

	return fullPath
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
