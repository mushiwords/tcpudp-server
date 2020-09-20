package server

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
)

/************************************************************
 *  TCP Server
************************************************************/
// Client holds info about connection
type TCPClient struct {
	conn   net.Conn
	Server *tcpserver
}

// TCP server
type tcpserver struct {
	address                  string // Address to open connection: localhost:9999
	config                   *tls.Config
	onNewClientCallback      func(c *TCPClient)
	onClientConnectionClosed func(c *TCPClient, err error)
	onNewMessage             func(c *TCPClient, message string)
}

// Read client data from channel
func (c *TCPClient) listen() {
	c.Server.onNewClientCallback(c)
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}
		c.Server.onNewMessage(c, message)
	}
}

// Send text message to client
func (c *TCPClient) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// Send bytes to client
func (c *TCPClient) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *TCPClient) Conn() net.Conn {
	return c.conn
}

func (c *TCPClient) Close() error {
	return c.conn.Close()
}

// Called right after server starts listening new client
func (s *tcpserver) OnNewClient(callback func(c *TCPClient)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *tcpserver) OnClientConnectionClosed(callback func(c *TCPClient, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *tcpserver) OnNewMessage(callback func(c *TCPClient, message string)) {
	s.onNewMessage = callback
}

// Listen starts network server
func (s *tcpserver) Listen() {
	var listener net.Listener
	var err error
	if s.config == nil {
		listener, err = net.Listen("tcp", s.address)
	} else {
		listener, err = tls.Listen("tcp", s.address, s.config)
	}
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client := &TCPClient{
			conn:   conn,
			Server: s,
		}
		go client.listen()
	}
}

// Creates new tcp server instance
func NewTCP(address string) *tcpserver {
	server := &tcpserver{
		address: address,
		config:  nil,
	}

	server.OnNewClient(func(c *TCPClient) {})
	server.OnNewMessage(func(c *TCPClient, message string) {})
	server.OnClientConnectionClosed(func(c *TCPClient, err error) {})

	return server
}

func NewTCPWithTLS(address string, certFile string, keyFile string) *tcpserver {
	cert, _ := tls.LoadX509KeyPair(certFile, keyFile)
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	server := &tcpserver{
		address: address,
		config:  &config,
	}

	server.OnNewClient(func(c *TCPClient) {})
	server.OnNewMessage(func(c *TCPClient, message string) {})
	server.OnClientConnectionClosed(func(c *TCPClient, err error) {})

	return server
}
