package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// headers section \r\n\r\n body section

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

	recievedData := make([]byte, 1024) //may need to be bigger

	_, err = conn.Read(recievedData)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Read incoming data:\n", string(recievedData))

	request := strings.Split(string(recievedData), "\r\n")

	start_line := strings.Split(request[0], " ")

	http_method := start_line[0]
	path := start_line[1]

	fmt.Println(request)
	fmt.Println(start_line)
	fmt.Println(http_method)
	fmt.Println(path)

	// need 2 sets of \r\n for end of headers section
	if path == "/" {
		response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.HasPrefix(path, "/echo/") {
		//echo the request in body
		response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n"
		content := strings.TrimPrefix(path, "/echo/")
		length := "Content-Length: " + strconv.Itoa(len(content)) + "\r\n\r\n"

		response = response + length + content
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	} else {
		response := "HTTP/1.1 404 Not Found\r\n\r\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	}

	conn.Close()
}
