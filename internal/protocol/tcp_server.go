package protocol

import (
	"fmt"
	"net"
	"runtime"
	"strings"

	"internal/app"
)

type TCPHandler interface {
	Handle(net.Conn)
}

//Define the TCP Server
func TCPServer(listener net.Listener, handler TCPHandler, l app.Logger) {
	l.Output(2, fmt.Sprintf("TCP: listening on %s", listener.Addr()))

	for {
		// Accept waits for and returns the next connection to the listener.
		clientConn, err := listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				l.Output(2, fmt.Sprintf("NOTICE: temporary Accept() failure - %s", err))
				runtime.Gosched()									//Gosched yields the processor, allowing other goroutines to run.
				continue
			}
			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
				l.Output(2, fmt.Sprintf("ERROR: listener.Accept() - %s", err))
			}
			break
		}
		go handler.Handle(clientConn)
	}

	l.Output(2, fmt.Sprintf("TCP: closing %s", listener.Addr()))
}
