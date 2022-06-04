package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

//Create a new user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// start the goroutine to monitor current channel
	go user.ListenMessage()

	return user
}

// func to monitor current user channel, once there are any message in the channel, send the message to client
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {
	// Add the user to server's OnlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// Broadcast the user online message to server's Message Channel
	this.server.Broadcast(this, "online")
}

func (this *User) Offline() {
	// Delete the user from server's OnlineMap
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// Broadcast the user offline message to server's Message Channel
	this.server.Broadcast(this, "offline")
}

func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// Query online users
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "is online\n"
			this.sendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// Change user name
		newName := strings.Split(msg, "|")[1]
		// Check if the new name exist already
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.sendMsg("The new user name has already exist.\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.Name = newName
			this.server.OnlineMap[this.Name] = this
			this.server.mapLock.Unlock()

			this.sendMsg("Your user name has updated as: " + this.Name + "\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		// Private Message
		if len(strings.Split(msg, "|")) != 3 {
			this.sendMsg("The message format is not correct, please use the foramt like: \"to|shizheng|I love you\" .\n")
			return
		}
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.sendMsg("The user name cannot be empty.\n")
			return
		}

		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.sendMsg("The user dosn't exist.\n")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.sendMsg("The message cannot be empty.\n")
			return
		}

		remoteUser.sendMsg(this.Name + " to you: " + content + "\n")

	} else {
		this.server.Broadcast(this, msg)
	}
}

func (this *User) sendMsg(msg string) {
	this.conn.Write([]byte(msg))
}
