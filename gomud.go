package main

import (
	"net"
	"log"
	"bufio"
	"flag"
	"fmt"
)

var flagPath, flagAddr, accountsPath, messagesPath, areasPath string

func init() {
	// Parse command line flags
	flag.StringVar(&flagPath, "path", "./", "Path to the root directory for game files.")
	flag.StringVar(&flagAddr, "address", ":4000", "Address for the game server.")
	flag.Parse()

	// Set commonly used filesystem paths
	accountsPath = fmt.Sprintf("%s/accounts/", flagPath)
	messagesPath = fmt.Sprintf("%s/messages/", flagPath)
	areasPath = fmt.Sprintf("%s/areas/", flagPath)
}

type Client struct {
	conn 	net.Conn
	ch 		chan<- string
}

func main() {
	loadMessages()

	ln, err := net.Listen("tcp", flagAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("gomud started and running on address %s", flagAddr)

	msgchan := make(chan string)
	addchan := make(chan Client)
	rmchan := make(chan net.Conn)

	go handleMessages(msgchan, addchan, rmchan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Connection from %v", conn.RemoteAddr())
		go handleConnection(conn, msgchan, addchan, rmchan)
	}
}

/*
 * Holds a client list and handles message broadcasts.
 */
func handleMessages(msgchan <-chan string, addchan <-chan Client, rmchan <-chan net.Conn) {
	clients := make(map[net.Conn]chan<- string)
	for {
		select {
		case msg := <-msgchan:
			log.Printf("New message: %s\n", msg)
			for _, ch := range clients {
				go func(mch chan<- string) { mch <- "\033[1;33;40m" + msg + "\033[m\n\r" }(ch)
			}
		case client := <-addchan:
			log.Printf("New client %v\n", client)
			clients[client.conn] = client.ch
		case conn := <-rmchan:
			log.Printf("Client disconnected %v", conn)
			delete(clients, conn)
		}
	}
}


/*
 * Handles an individual's connection to the mud.
 */
func handleConnection(c net.Conn, msgchan chan<- string, addchan chan<- Client, rmchan chan<- net.Conn) {
	ch := make(chan string)
	msgs := make(chan string)
	color := true
	addchan <- Client{c, ch}

	go func() {
		defer close(msgs)
		bufc := bufio.NewReader(c)

		//ColorWrite(c, getMessage("login"))
		ColorStripWrite(c, getMessage("login"))
		//c.Write( getMessage("login") )

		nick, _, err := bufc.ReadLine()
		if err != nil {
			return
		}
		nickname := string(nick)

		c.Write([]byte("Welcome, " + nickname + "!\n\r\n\r"))
		msgs <- "New user, " + nickname + ", joined the chat room."

		for {
			line, _, err := bufc.ReadLine()
			if err != nil {
				break
			}
			msgs <- nickname + ": " + string(line)
		}

		msgs <- nickname + " left the chatroom."
	}()


LOOP:
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				break LOOP
			}
			msgchan <- msg
		case msg := <-ch:
			_, err := c.Write([]byte(msg))
			if err != nil {
				break LOOP
			}
		}
	}


	log.Printf("Connection from %v closed.", c.RemoteAddr())
}





