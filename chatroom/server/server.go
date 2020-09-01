package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

var (
	enteringChannel = make(chan *User)
	leavingChannel  = make(chan *User)
	messageChannel  = make(chan string, 8)
)

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		return
	}
	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringChannel:
			users[user] = struct{}{}
		case user := <-leavingChannel:
			delete(users, user)
			close(user.MessageChannel)
		case msg := <-messageChannel:
			for user := range users {
				user.MessageChannel <- msg
			}
		}
	}
}

type User struct {
	Id             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func (user *User) String() string {
	return strconv.Itoa(user.Id)
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	user := &User{
		Id:             GenUserId(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	go sendMessage(conn, user.MessageChannel)

	user.MessageChannel <- "Welcome, " + user.String()
	messageChannel <- "User: " + user.String() + " has enter"

	enteringChannel <- user

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messageChannel <- user.String() + ": " + input.Text()
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误：", err)
	}

	leavingChannel <- user
	messageChannel <- "user: " + user.String() + " has left."
}

func GenUserId() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10000)
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
