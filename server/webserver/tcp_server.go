package webserver

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

// Credit for this specific TCP server implementation goes to
// https://www.youtube.com/watch?v=qJQrrscB1-4

// Keep track of the number of open TCP connections
var NumOpenTCPClients int = 0

/*
Protect important TCP variables from race conditions with a mutex
- it will lock the resource until it is written to
*/
var muTcp sync.Mutex

// Keep track of who sent the message, and the content
type Message struct {
	from    string
	payload []byte
}

// TCP server, with a message channel and quit channel
type TcpServer struct {
	listenAddr string
	listener   net.Listener
	quitCh     chan struct{}
	readCh     chan Message
	tcpSendCh  <-chan []byte
	conns      map[net.Conn]struct{}
}

// Create a new TCP server, buffering up to 10 messages
func NewTcpServer(listenAddr string, _tcpSendCh <-chan []byte) *TcpServer {
	return &TcpServer{
		listenAddr: listenAddr,
		quitCh:     make(chan struct{}),
		readCh:     make(chan Message, 200),
		tcpSendCh:  _tcpSendCh,
		conns:      make(map[net.Conn](struct{})),
	}
}

// Initialize the TCP server and handle connections and messages
func (s *TcpServer) TcpStart() error {

	// Start the TCP connection
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	// Close the listener upon exiting the function or (less ideally) crashing
	defer listener.Close()
	s.listener = listener

	// Run the accept loop for TCP connections
	go s.tcpAcceptLoop()

	// Run the send loop for TCP connections
	go s.tcpSendLoop()

	// Block on the quit channel as long as we haven't quit yet
	<-s.quitCh

	// Close the read channel once we have quit
	close(s.readCh)

	// No errors
	return nil
}

// Accept incoming TCP connections
func (s *TcpServer) tcpAcceptLoop() {
	for {
		// Accept an incoming connection request
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		// Increment the open clients, and print out debug info
		muTcp.Lock()
		NumOpenTCPClients++
		s.conns[conn] = struct{}{}
		log.Printf("\033[32m[%d -> %d] robot connected at %s\033[0m\n", NumOpenTCPClients-1, NumOpenTCPClients, conn.RemoteAddr().String())
		muTcp.Unlock()
		go s.tcpReadLoop(conn)
	}
}

// Continue reading messages from a connection
func (s *TcpServer) tcpReadLoop(conn net.Conn) {

	// Close the connection when necessary
	defer func() {
		conn.Close()
		muTcp.Lock()
		NumOpenTCPClients--
		delete(s.conns, conn)
		log.Printf("\033[31m[%d -> %d] robot quit at %s\033[0m\n", NumOpenTCPClients+1, NumOpenTCPClients, conn.RemoteAddr().String())
		muTcp.Unlock()
	}()

	for {
		// Read new messages
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {

			// Handle EOF (connection closure)
			if err == io.EOF {
				muTcp.Lock()
				NumOpenTCPClients--
				delete(s.conns, conn)
				log.Printf("\033[31m[%d -> %d] robot disconnected at %s\033[0m\n", NumOpenTCPClients+1, NumOpenTCPClients, conn.RemoteAddr().String())
				muTcp.Unlock()
				return
			}

			// Handle network operational errors, such as timeouts or connection failures
			if opErr, ok := err.(*net.OpError); ok {
				// If it's a timeout, retry a few times (backoff strategy or a simple retry)
				if opErr.Op == "read" && opErr.Err.Error() == "i/o timeout" {
					// Timeout error - retry reading a few more times
					log.Printf("\033[33m[%d] Timeout error with robot at %s. Retrying...\033[0m\n", NumOpenTCPClients, conn.RemoteAddr().String())
					continue // Retry the read operation
				}

				// If the error indicates a temporary failure (e.g., network unreachable), we could log it and try again
				if opErr.Op == "read" && strings.Contains(opErr.Err.Error(), "forcibly closed") {
					// You might want to add a retry limit or backoff strategy
					return // Retry reading from the connection
				}

				// For other types of net.OpErrors (like connection reset), log the error and keep the connection open
				log.Printf("\033[31m[%d] Network operation error with robot at %s: %v. Continuing...\033[0m\n", NumOpenTCPClients, conn.RemoteAddr().String(), opErr)
				continue // Keep trying to read
			}

			// Log any other read errors that aren't EOF or network operation errors
			log.Println("\tRead error: ", err)
			continue
		}

		// Send a message to the channel for logging
		s.readCh <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}

		// For testing purposes (if a message 'q' is sent, kick the connection)
		if bytes.Equal(buf[:n], []byte("q")) {
			return
		}
	}
}

// Send out messages to the TCP client
func (s *TcpServer) tcpSendLoop() {
	// For testing purposes (if a message is received, send '[ACK]' to the client)
	for msg := range s.tcpSendCh {
		for conn := range s.conns {
			conn.Write(msg)
		}
	}
}

// Print out messages that are received
func (s *TcpServer) Printer() {
	for msg := range s.readCh {
		log.Printf("\033[2m\033[32m| TCP from %s: %s`\033[0m\n", msg.from, string(msg.payload))
	}
}
