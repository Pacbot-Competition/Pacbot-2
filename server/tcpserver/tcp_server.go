package tcpserver

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
)

/*
The design for this TCP server is heavily inspired by
https://www.youtube.com/watch?v=qJQrrscB1-4 (check out that
video if you're interested in writing a more general one)
*/

// Keep track of the number of open TCP connections
var openTcpClients int = 0

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
}

// Create a new TCP server, buffering up to 10 messages
func NewTcpServer(listenAddr string) *TcpServer {
	return &TcpServer{
		listenAddr: listenAddr,
		quitCh:     make(chan struct{}),
		readCh:     make(chan Message, 10),
	}
}

// Initialize the TCP server and handle connections andd messages
func (s *TcpServer) TcpStart(wgQuit *sync.WaitGroup) {

	// Start the TCP connection
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		fmt.Println("\033[35m\033[1mERR:  Failed to initialize the TCP server. Quitting...\033[0m")
		return
	}

	// Increment the wait group counter, and defer closing it upon returning
	wgQuit.Add(1)
	defer wgQuit.Done()

	// Close the listener upon exiting the function or (less ideally) crashing
	defer listener.Close()
	s.listener = listener

	// Run the accept loop for TCP connections
	go s.tcpAcceptLoop()

	// Block on the quit channel as long as we haven't quit yet
	<-s.quitCh

	// Close the read channel once we have quit
	close(s.readCh)

	// Log that the TCP server successfully quit
	fmt.Println("\033[35mLOG:  TCP server successfully quit\033[0m")
}

// Accept incoming TCP connections
func (s *TcpServer) tcpAcceptLoop() {
	for {
		// Accept an incoming connection request
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		// Increment the open clients, and print out debug info
		muTcp.Lock()
		{
			openTcpClients++
			fmt.Printf("\033[32m[%d -> %d] robot connected at %s\033[0m\n", openTcpClients-1, openTcpClients, conn.RemoteAddr().String())
			muTcp.Unlock()
		}
		go s.tcpReadLoop(conn)
	}
}

// Continue reading messages from a connection
func (s *TcpServer) tcpReadLoop(conn net.Conn) {

	// Close the connection when necessary
	defer conn.Close()

	for {

		// Read new messages
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {

			// If the connection ends, log and return
			if err == io.EOF {
				muTcp.Lock()
				{
					openTcpClients--
					fmt.Printf("\033[31m[%d -> %d] robot disconnected at %s\033[0m\n", openTcpClients+1, openTcpClients, conn.RemoteAddr().String())
				}
				muTcp.Unlock()
				return
			}

			// If the connection forcefully ends, log and return
			if _, ok := err.(*net.OpError); ok {
				muTcp.Lock()
				{
					openTcpClients--
					fmt.Printf("\033[31m[%d -> %d] robot vanished at %s\033[0m\n", openTcpClients+1, openTcpClients, conn.RemoteAddr().String())
				}
				muTcp.Unlock()
				return
			}

			// Log read errors
			fmt.Println("\tread error: ", err)
			continue
		}

		// Send a message to the channel for logging
		s.readCh <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}

		// For testing purposes (if a message 'q' is sent, kick the connection)
		if bytes.Equal(buf[:n], []byte("q")) {
			muTcp.Lock()
			{
				openTcpClients--
				fmt.Printf("\033[31m[%d -> %d] robot quit at %s\033[0m\n", openTcpClients+1, openTcpClients, conn.RemoteAddr().String())
			}
			muTcp.Unlock()
			return
		}

		// For testing purposes (if a message is received, send '[ACK]' to the client)
		conn.Write([]byte("[ACK]\n"))
	}
}

// Quit the TCP server
func (s *TcpServer) Quit() {

	// Add an object to its quit channel to allow it to quit
	s.quitCh <- struct{}{}
}

// Print out messages that are received
func (s *TcpServer) Printer() {
	for msg := range s.readCh {
		fmt.Printf("\033[2m\033[32m| TCP from %s: %s`\033[0m\n", msg.from, string(msg.payload))
	}
}
