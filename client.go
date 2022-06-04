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
	flag       int
}

var serverIp string
var serverPort int

func NewClient(serverIp string, serverPort int) *Client {
	// Creat a new client
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
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

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. Send public message")
	fmt.Println("2. Send private message")
	fmt.Println("3. Update user name")
	fmt.Println("0. Exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>Please input a valid number<<<<<")
		return false
	}
}

func (client *Client) Run() {
	// If the flag number != 0, go into the loop, else exit
	for client.flag != 0 {
		// Loop until recieve a valid flag number
		for client.menu() != true {
			continue
		}

		switch client.flag {
		case 1:
			fmt.Println("Choose: Send public message")
			break
		case 2:
			fmt.Println("Choose: Send private message")
			break
		case 3:
			fmt.Println("Choose: Update user name")
			break
		}
	}
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
	client.Run()
}
