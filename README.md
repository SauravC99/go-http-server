# go-http-server

> http server written in Go using concurrency

go-http-server is a lightweight HTTP server written in GoLang. By default the server will run on localhost with port 4221, and an option to specify which port to bind to.

The server is able to handle multiple requests concurrently using Goroutines. It handles downloading and uploading data to the server by the way of GET and POST requests and can compress responses using gzip.


## Installation
Clone the repository to your computer:
```
git clone https://github.com/SauravC99/go-http-server.git
```

Build the project:
```bash
go build -o server main.go
```

### Alternate Installation

You can clone the repository using the instructions above and run the `server.sh` script. The script will compile the program and run it.
```console
$ ./server.sh
```


## Usage
Run the server:
```console
$ ./server
Server started
Listening on port 4221
```

Run with `-h` to see avaliable commands:
```console
$ ./server -h
Usage of ./server:
  -directory string
        Directory for file hosting (download and upload)
  -port string
        Port to bind to (default "4221")
```


### Command Line Flags
- `-directory`: Path to the directory for file storage (required for `/files` endpoint).
- `-port`: Port to bind the server to (default: `4221`).

```console
$ ./server -directory /path/to/directory -port 4221
```


## Endpoints
The server supports the following endpoints.

### Root Endpoint (`/`)
- **Path**: `/`
- **Method**: `GET`
- **Response**: Returns `200 OK` with an empty body.

### Echo Endpoint (`/echo/<content>`)
- **Path**: `/echo/<content>`
- **Method**: `GET`
- **Response**:
  - Returns `200 OK` with the content of the path as the response body.
  - If the client includes an `Accept-Encoding` header with a valid encoding (`gzip`), the response will be compressed.
  - If `Accept-Encoding` is invalid or not included, response will be plain text.

### User-Agent Endpoint (`/user-agent`)
- **Path**: `/user-agent`
- **Method**: `GET`
- **Response**: Returns `200 OK` with the `User-Agent` header value from the request.

### Files Endpoint (`/files/<filename>`)
- **Path**: `/files/<filename>`
- **Method**: `GET`
- **Response**:
  - Serves the requested file from the directory provided by the `-directory` flag.
  - Returns `200 OK` including the `Content-Type: application/octet-stream` header with the file content.
  - If the file does not exist, returns `404 Not Found`.
- **Method**: `POST`
- **Response**:
  - Creates a new file, `<filename>`, with the request body in the directory provided by the `-directory` flag.
  - Returns `201 Created` with the `Location` header indicating the file path.
  - If the operation fails, returns `500 Internal Server Error`.

Note the `/files` endpoint requires the `-directory` flag to be set.


## Examples
Here are some example interactions with the server using `curl`:

- **Root Request**
```bash
curl -v http://localhost:4221/
```
```bash
< HTTP/1.1 200 OK
< Content-Type: text/plain
< Content-Length: 0
<
```

- **Echo Request**
```bash
curl -v http://localhost:4221/echo/hello123
```
```bash
< HTTP/1.1 200 OK
< Content-Type: text/plain
< Content-Length: 10
<
hello123
```

- **Echo Request with gzip compression**
```bash
curl -v -H "Accept-Encoding: gzip" http://localhost:4221/echo/HelloWorld
```
```bash
< HTTP/1.1 200 OK
< Content-Encoding: gzip
< Content-Type: text/plain
< Content-Length: 34
<
��H����/�I��y
             ww
```

- **User-Agent Request**
```bash
curl -v -A "MyCustomAgent" http://localhost:4221/user-agent
```
```bash
< HTTP/1.1 200 OK
< Content-Type: text/plain
< Content-Length: 13
<
MyCustomAgent
```

- **File GET Request**
```bash
curl -v http://localhost:4221/files/hello.txt
```
```bash
< HTTP/1.1 200 OK
< Content-Type: application/octet-stream
< Content-Length: 37
<
HELLO WORLD This is my file with data
```
If the file does not exist you get `404 Not Found`:
```bash
curl -v http://localhost:4221/files/randomfile.txt
```
```bash
< HTTP/1.1 404 Not Found
<
```

- **File POST Request**
```bash
curl -X POST -d "content or file to send" http://localhost:4221/files/newfile.txt
```
```bash
< HTTP/1.1 201 Created
< Content-Type: application/x-www-form-urlencoded
< Location: ./newfile.txt
<
content or file to send
```
If you try to write a file without the `-directory` flag set you get a `500 Internal Server Error`:
```bash
< HTTP/1.1 500 Internal Server Error
<
```