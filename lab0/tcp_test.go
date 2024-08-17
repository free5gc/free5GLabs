package lab0

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTcpFunction(t *testing.T) {
	// Type Assertion
	var _ listenerInterface = TCPListener
	var _ handlerInterface = TCPHandler

	go TCPListener("127.0.0.1", 8080, TCPHandler)

	time.Sleep(5 * time.Second)

	connectionSlice := make([]net.Conn, 0)

	connNum := 10

	for i := 0; i < connNum; i++ {
		func() {
			buf := make([]byte, 1024)
			conn, err := net.Dial("tcp", "127.0.0.1:8080")
			require.NoError(t, err)

			_, err = conn.Write([]byte("OK\n"))
			require.NoError(t, err)

			var n int
			n, err = conn.Read(buf)
			require.NoError(t, err)

			require.Equal(t, "OK\n", string(buf[:n]))

			connectionSlice = append(connectionSlice, conn)
		}()
	}

	for _, conn := range connectionSlice {
		err := conn.Close()
		require.NoError(t, err)
	}
}
