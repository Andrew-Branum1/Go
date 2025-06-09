// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

// !+broadcaster
type client chan<- string // an outgoing message channel

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
	clients  = make(map[string]chan<- string)
)

func broadcaster() {
	//log.Print("client map: ", clients)
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for _, ch := range clients {
				ch <- msg
			}

		case cli := <-entering:
			cli <- "Active Users:"
			for username := range clients {
				cli <- username
			}

		case cli := <-leaving:
			for username, ch := range clients {
				if ch == cli {
					delete(clients, username)
					break
				}
			}
			close(cli)
		}
	}
}

//!-broadcaster

// !+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	input := bufio.NewScanner(conn)
	input.Scan()
	who := strings.TrimSpace(input.Text())

	clients[who] = ch

	ch <- "You are " + who
	messages <- who + " has arrived"

	entering <- ch

	for input.Scan() {
		message := input.Text()
		if message[0] == '/' {
			split := strings.SplitN(message[1:], " ", 2)
			//log.Print(split)

			username := split[0]
			//log.Print(username)
			dm := split[1]
			//log.Print(dm)

			channel := clients[username]

			channel <- "DM (" + who + "): " + dm

		} else {
			messages <- who + ": " + message
		}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

// !+main
func main() {
	portNumber := flag.String("p", "8080", "Port")
	flag.Parse()
	//log.Printf("Port Number: %s", *portNumber)

	address := fmt.Sprintf("localhost:%s", *portNumber)
	//log.Printf("Host Location: %s", address)

	listener, err := net.Listen("tcp", address)
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

//!-main
