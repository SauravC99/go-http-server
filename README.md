# go-http-server

> http server written in Go using concurrency

go-http-server is a HTTP server written in GoLang. By default the server will run on localhost with port 4221, and an option to specify which port to bind to.

It is able to handle multiple requests concurrently using Goroutines. It handles downloading and uploading data to the server by the way of GET and POST requests.


## Installation
Clone the repository to your computer:
```
git clone https://github.com/SauravC99/go-http-server.git
```

Compile the program:
```
go build main.go
```

### Alternate Installation

You can clone the repository using the instructions above and run the `server.sh` script. The script will compile the program and run it.


## Usage (compiled)
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


## Usage (script)
Run the script:
```console
$ ./server.sh
Server started
Listening on port 4221
```
Run with `-h` to see avaliable commands:
```console
$ ./server.sh -h
Usage of /tmp/tmp.5V70Zti4Q5:
  -directory string
        Directory for file hosting (download and upload)
  -port string
        Port to bind to (default "4221")
```