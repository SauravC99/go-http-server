package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type RequestHeaders struct {
	Method         string
	Path           string
	Host           string
	Agent          string
	ContentType    string
	ContentLength  string
	AcceptEncoding string
	Body           string
}

const STATUS_200_OK string = "HTTP/1.1 200 OK\r\n"
const STATUS_201_CREATED string = "HTTP/1.1 201 Created\r\n"
const STATUS_404_ERR string = "HTTP/1.1 404 Not Found\r\n"
const STATUS_405_NOTALLOW string = "HTTP/1.1 405 Method Not Allowed\r\n"
const STATUS_500_ERR string = "HTTP/1.1 500 Internal Server Error\r\n"
const CONTENT_PLAIN string = "Content-Type: text/plain\r\n"
const CONTENT_APP string = "Content-Type: application/octet-stream\r\n"
const END_HEADER_LINE string = "\r\n"
const END_HEADER_BLOCK string = "\r\n\r\n"
const CONTENT_LENGTH string = "Content-Length: "
const CONTENT_TYPE string = "Content-Type: "
const LOCATION_HEADER string = "Location: "
const CONTENT_ENCODING string = "Content-Encoding: "

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
		go connectAndRespond(conn, directoryPtr)
	}
}

func connectAndRespond(connection net.Conn, directoryPtr *string) {
	headers, err := parseHeaders(connection)
	if err != nil {
		fmt.Println("Failed to read headers: ", err.Error())
	}

	if headers.Path == "/" {
		respond200Plain(connection, "")
		return
	} else if strings.HasPrefix(headers.Path, "/echo/") {
		//echo the request in body
		if strings.ToLower(headers.AcceptEncoding) == "invalid-encoding" || headers.AcceptEncoding == "" {
			content := strings.TrimPrefix(headers.Path, "/echo/")
			respond200Plain(connection, content)
			return
		} else {
			content := strings.TrimPrefix(headers.Path, "/echo/")
			respond200Encode(connection, content, headers.AcceptEncoding)
			return
		}
	} else if strings.HasPrefix(headers.Path, "/user-agent") {
		content := headers.Agent
		respond200Plain(connection, content)
		return
	} else if strings.HasPrefix(headers.Path, "/files") {
		fPath := parseFilePath(directoryPtr, headers.Path)

		if headers.Method == "GET" {
			file, err := os.ReadFile(fPath)
			if err != nil {
				fmt.Println("Error reading file: ", err.Error())
				respond404(connection)
				return
			}

			content := string(file)
			respond200App(connection, content)
			return
		} else if headers.Method == "POST" {
			file, err := os.Create(fPath)
			if err != nil {
				fmt.Println("Error creating file: ", err.Error())
				respond500(connection)
				return
			}
			// Replaces null characters "\x00" with empty string; needed bc this is binary data
			parsedBody := strings.ReplaceAll(headers.Body, "\x00", "")
			_, err = file.WriteString(parsedBody)
			if err != nil {
				fmt.Println("Error writing file: ", err.Error())
				respond500(connection)
				return
			}
			file.Close()

			respond201(connection, fPath, headers.ContentType, parsedBody)
			return
		} else {
			respond405(connection)
			return
		}
	} else {
		respond404(connection)
		return
	}
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
	contentType := ""
	contentLen := ""
	acceptEncoding := ""
	body := ""

	for _, line := range request {
		switch {
		case strings.HasPrefix(strings.ToLower(line), "host:"):
			temp := strings.Split(line, " ")
			host = temp[1]
		case strings.HasPrefix(strings.ToLower(line), "user-agent:"):
			temp := strings.Split(line, " ")
			agent = temp[1]
		case strings.HasPrefix(strings.ToLower(line), "content-type:"):
			temp := strings.Split(line, " ")
			contentType = temp[1]
		case strings.HasPrefix(strings.ToLower(line), "content-length:"):
			temp := strings.Split(line, " ")
			contentLen = temp[1]
		case strings.HasPrefix(strings.ToLower(line), "accept-encoding:"):
			temp := strings.Split(line, " ")
			// Handle multiple compression scheme values
			if len(temp) > 2 {
				encodings := temp[1:]
				for _, line := range encodings {
					line = strings.ReplaceAll(line, ",", "")
					// For the time being accept gzip compression
					if strings.ToLower(line) == "gzip" {
						acceptEncoding = "gzip"
					}
				}
			} else {
				acceptEncoding = temp[1]
			}
		}
	}
	// headers section \r\n\r\n body section
	if len(strings.Split(string(recievedData), END_HEADER_BLOCK)) > 1 {
		temp := strings.Split(string(recievedData), END_HEADER_BLOCK)
		body = temp[1]
	}

	return &RequestHeaders{
		Method:         http_method,
		Path:           path,
		Host:           host,
		Agent:          agent,
		ContentType:    contentType,
		ContentLength:  contentLen,
		AcceptEncoding: acceptEncoding,
		Body:           body,
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

func respond200Plain(connection net.Conn, content string) {
	response := STATUS_200_OK + CONTENT_PLAIN
	length := CONTENT_LENGTH + strconv.Itoa(len(content)) + END_HEADER_BLOCK
	body := content

	response = response + length + body
	_, err := connection.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
	connection.Close()
}

func respond200App(connection net.Conn, content string) {
	response := STATUS_200_OK + CONTENT_APP
	length := CONTENT_LENGTH + strconv.Itoa(len(content)) + END_HEADER_BLOCK
	body := content

	response = response + length + body
	_, err := connection.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
	connection.Close()
}

func respond200Encode(connection net.Conn, content string, encode string) {
	//gzip encode
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	writer.Write([]byte(content))
	writer.Close()
	encodedContent := buffer.String()

	response := STATUS_200_OK
	encoding := CONTENT_ENCODING + encode + END_HEADER_LINE
	cType := CONTENT_PLAIN
	length := CONTENT_LENGTH + strconv.Itoa(len(encodedContent)) + END_HEADER_BLOCK
	body := encodedContent

	response = response + encoding + cType + length + body
	_, err := connection.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
	connection.Close()
}

func respond201(connection net.Conn, filepath string, cType string, content string) {
	response := STATUS_201_CREATED
	contentType := CONTENT_TYPE + cType + END_HEADER_LINE
	location := LOCATION_HEADER + filepath + END_HEADER_BLOCK
	body := content

	response = response + contentType + location + body
	_, err := connection.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
	connection.Close()
}

func respond404(connection net.Conn) {
	response := STATUS_404_ERR + END_HEADER_LINE
	_, err := connection.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
	connection.Close()
}

func respond405(connection net.Conn) {
	// Code 405 must include allow header field in response
	allowed := "Allow: GET, POST" + END_HEADER_LINE
	response := STATUS_405_NOTALLOW + allowed + END_HEADER_LINE
	_, err := connection.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
	connection.Close()
}

func respond500(connection net.Conn) {
	response := STATUS_500_ERR + END_HEADER_LINE
	_, err := connection.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
	connection.Close()
}
