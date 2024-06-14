package lab0

import (
	"net"
)

type listenerInterface func(string, int, handlerInterface)

type handlerInterface func(conn net.Conn)

func TcpListener(host string, port int, handler handlerInterface) {

}

func TcpHandler(conn net.Conn) {

}
