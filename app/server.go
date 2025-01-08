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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn) //"go" so while one connection is being processed, the server can accept and process additional connections without waiting for the first one to complete.
	}
}

// conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
func handleRequest(conn net.Conn) {
	defer conn.Close() //conn.Close() will be called automatically when the function below finishes (defer function) => connection will be closed when response is sent, prevent loading on the browser (if not, client might not know if when response is completed, so keep loading)
	req := make([]byte, 1024)
	n, err := conn.Read(req) //n is the number of bytes need
	if err != nil {
		fmt.Println("Error reading connection: ", err.Error())
		os.Exit(1)
	}
	message := string(req[:n])
	request := strings.Split(message, " ")
	path := request[1] // split by space and extract second elenment (index 1)
	response := ""

	if path == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	} else if strings.HasPrefix(path, "/echo") {
		text := path[6:]                                                                                                       // from req, "/echo/abc" => text = abc (from index 6)
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(text), text) // khớp value sau vào những ký hiệu %, &d for interger, %s for string and %f for float. Trả về string
	} else if path == "/user-agent" {
		lines := strings.Split(message, "\r\n") // khong the dung index cu the vi trong 1 request co the co nhieu thanh phan khac
		var userAgent string
		for _, line := range lines {
			if strings.HasPrefix(line, "User-Agent: ") {
				userAgent = strings.TrimPrefix(line, "User-Agent: ")
				break
			}
		}
		response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
	} else if strings.HasPrefix(path, "/files/") {
		fileName := strings.TrimPrefix(path, "/files/")
		dir := os.Args[2]
		if request[0] == "POST" {
			lines := strings.Split(message, "\r\n")
			file, _ := os.Create(dir + fileName)
			file.WriteString(lines[len(lines)-1])
			response = "HTTP/1.1 201 Created\r\n\r\n"
		} else {
			data, err := os.ReadFile(dir + fileName)
			if err != nil {
				response = "HTTP/1.1 404 Not Found\r\n\r\n"
			} else {
				response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(data), data)
			}
		}
	} else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error sending response to client", err.Error())
		os.Exit(1)
	}
}
