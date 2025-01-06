package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	req := make([]byte, 1024)
	n, err := conn.Read(req) //n is the number of bytes need
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		os.Exit(1)
	}
	message := string(req[:n])
	path := strings.Split(message, " ")[1] // split by space and extract second elenment (index 1)
	response := ""

	if path == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else if strings.HasPrefix(path, "/echo") {
		text := path[6:]                                                                                                       // from req, "/echo/abc" => text = abc (from index 6)
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(text), text) // khớp value sau vào những ký hiệu %, &d for interger, %s for string and %f for float. Trả về string
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error sending response to client", err.Error())
		os.Exit(1)
	}

}
