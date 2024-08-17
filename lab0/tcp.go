package lab0

import (
	"net"
)

type listenerInterface func(string, int, handlerInterface)

type handlerInterface func(conn net.Conn)

func TCPListener(host string, port int, handler handlerInterface) {
}

func TCPHandler(conn net.Conn) {
}
