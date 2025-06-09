package main

import (
	"io"
	"log"
	"net"
	"os"
	"sync"
)

// run this locally to connect to EC2 instances and display time
func main() {
	var wg sync.WaitGroup

	//log.Printf("Args: %s", os.Args)
	for _, argument := range os.Args[1:] {
		wg.Add(1)
		go connectToServer(argument, &wg)
	}

	wg.Wait()
}
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

// added param and wg.Done()
func connectToServer(address string, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := net.Dial("tcp", address)

	if err != nil {
		log.Print(err)
		return
	}

	defer conn.Close()
	mustCopy(os.Stdout, conn)
}
