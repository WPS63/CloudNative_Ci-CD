// Demonstration of channels with a chat application
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Chat is a server that lets clients chat with each other.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client struct { 
	//Converting type client
	clientChan chan<- string // an outgoing message channel
	clientName string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.clientChan <- msg
			}

		case cli := <-entering:
			clients[cli] = true

			//Displays current list of chatters to new client
			cli.clientChan <- "Current chatters: "
			for c := range clients {
				cli.clientChan <- c.clientName
			}

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.clientChan)
		}
	}
}

func handleConn(conn net.Conn) {
	//Instantiating client
	var cli client
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	//Setting username
	inputName := bufio.NewReader(conn)
	fmt.Fprintln(conn, "Enter name: ")
	name, _ := inputName.ReadString('\n')
	name = strings.TrimSuffix(name, "\n")
	

	//Assigning outgoing message channel variable
	cli.clientChan = ch
	
	//Assigning client name variable 
	cli.clientName = name

	//Sending name to channel
	ch <- "You are " + name
	//Sending name to messages case to be broadcast within the broadcast function
	messages <- name + " has arrived"
	//Sending client struct to entering case to be broadcast within the broadcast function
	entering <- cli

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- name + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- cli
	messages <- name + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
