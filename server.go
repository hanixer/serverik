package serverik

import (
	"fmt"
	"log"
	"net"
	"os"
)

type RequestHandler interface {
	HandleRequest(*HttpRequest, net.Conn)
}

func handleRequestFunc(conn net.Conn, handler func(*HttpRequest, net.Conn)) {
	request, _ := ReadRequest(conn)
	handler(request, conn)
}

func handleRequest(conn net.Conn, handler RequestHandler) {
	request, err := ReadRequest(conn)

	if err != nil {
		log.Println("ReadRequest error.", err)
		conn.Close()
		return
	}

	handler.HandleRequest(request, conn)
}

func Serve(port int, handler RequestHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println("net.Listen", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept.", err)
			return
		}

		handleRequest(conn, handler)
	}
}

func ServeFunc(dir string, port int, handler func(*HttpRequest, net.Conn)) {
	err := os.Chdir(dir)
	if err != nil {
		log.Println("chdir failed. ", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println("handle error", listener)
	}

	for {
		conn, err := listener.Accept()
		log.Println("Got connection!", conn.RemoteAddr())
		if err != nil {
			log.Println("error.", err)
			return
		}

		handleRequestFunc(conn, handler)
	}
}
