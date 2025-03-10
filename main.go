package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/sauravc99/go-http-server/app"
)

func main() {
	directoryPtr := flag.String("directory", "", "Directory for file hosting (download and upload)")
	portPtr := flag.String("port", "4221", "Port to bind to")
	flag.Parse()

	fmt.Println("Server started")

	listener, err := net.Listen("tcp", "0.0.0.0:"+*portPtr)
	if err != nil {
		fmt.Println("Failed to bind to port "+*portPtr, err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on port " + *portPtr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Connection accept")
		go app.ConnectAndRespond(conn, directoryPtr)
	}
}
