package main

import (
	"fmt"
	"net"
	"os"
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

	recievedData := make([]byte, 1024)

	_, err = conn.Read(recievedData)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Read incoming data:\n", string(recievedData))

	// need 2 sets of \r\n for end of headers section
	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}

	conn.Close()
}
