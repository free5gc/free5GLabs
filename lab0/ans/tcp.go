package lab0

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type listenerInterface func(string, int, handlerInterface)

type handlerInterface func(conn net.Conn)

func TcpListener(host string, port int, handler handlerInterface) {
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("Server can't listening on port %d: %v", port, err)
	}
	defer server.Close()

	log.Printf("TCP is listening on %s:%d", host, port)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalf("Accept failed: %v", err)
		}

		log.Printf("new client accepted: %s", conn.RemoteAddr().String())
		go handler(conn)
	}
}

func TcpHandler(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()
	log.Printf("Handle Request from [%s]", clientAddr)

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("Client [%s[] Error: %v", clientAddr, err)
			return
		}
		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Printf("Reply to Client [%s] failed: %v", clientAddr, err)
			return
		}
	}
}
