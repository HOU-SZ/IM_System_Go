package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>Please input your name:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error: ", err)
		return false
	}

	return true
}

func (client *Client) PublicMessage() {
	var chatMsg string

	fmt.Println(">>>>Please input message, or exit by input \"exit\".")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			_, err := client.conn.Write([]byte(chatMsg + "\n"))
			if err != nil {
				fmt.Println("conn.Write error: ", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>Please input message, or exit by input \"exit\".")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) QueryOnlineUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error: ", err)
		return
	}
}

func (client *Client) PrivateMessage() {
	var remoteName string
	var chatMsg string

	client.QueryOnlineUsers()
	fmt.Println(">>>>Please input user name, or exit by input \"exit\".")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>Please input message, or exit by input \"exit\".")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write error: ", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>Please input message, or exit by input \"exit\".")
			fmt.Scanln(&chatMsg)
		}

		client.QueryOnlineUsers()
		fmt.Println(">>>>Please input user name, or exit by input \"exit\".")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) DealResponse() {
	// Once the are some message in client.conn, copy it to stdout, block and listen forever
	io.Copy(os.Stdout, client.conn)

	// Equivalent to the following code
	// buffer := make([]byte, 4096)
	// client.conn.Read(buffer)
	// fmt.Println(buffer)
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
			client.PublicMessage()
			break
		case 2:
			client.PrivateMessage()
			break
		case 3:
			client.UpdateName()
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

	// Create a goroutine to process the message from server
	go client.DealResponse()

	fmt.Println("Connect to server success...")

	// todo
	client.Run()
}
