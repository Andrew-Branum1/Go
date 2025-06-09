package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

// !+
func main() {
	userName := flag.String("u", "default", "Username")
	address := flag.String("a", "localhost:8080", "IP:PORT")
	flag.Parse()

	//log.Printf("Username %s", *userName)
	//log.Printf("Address %s", *address)

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()

	//fmt.Fprintf(conn, "%s\n", *userName)
	sendUsername(conn, *userName)
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func sendUsername(dst io.Writer, user string) {
	if _, err := io.WriteString(dst, user+"\n"); err != nil {
		log.Fatal(err)
	}
}
