package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// run this on EC2 instances
func main() {
	portNumber := flag.String("p", "8080", "Port")
	flag.Parse()
	//log.Printf("Port Number: %s", *portNumber)

	address := fmt.Sprintf(":%s", *portNumber)
	//log.Printf("Host Location: %s", address)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

// modified time format
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05 MST\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
