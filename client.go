package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

var serverIp string
var serverPort int

func NewClient(serverIp string, serverPort int) *Client {
	// Creat a new client
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// Connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return nil
	}

	client.conn = conn
	return client
}

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Server ip address")
	flag.IntVar(&serverPort, "port", 8888, "Server port")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("Connect to server failed...")
		return
	}

	fmt.Println("Connect to server success...")

	// todo
	select {}
}
