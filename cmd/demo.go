package main

import (
	"fmt"
	"log"
	"net"

	"gitlab.com/bfelipe/atomic"
)

func handle(conn net.Conn) {
	var req atomic.Request
	if err := req.Decode(conn); err != nil {
		log.Fatalf("unable to decode connection %s", err)
	}

	log.Print("Request Headers")
	for k, v := range req.Headers() {
		log.Printf("%s: %s\n", k, v)
	}
	log.Printf("Request body %s", string(req.Body()))

	var resp atomic.Response
	resp.SetStatusCode(atomic.OK)
	resp.SetHeader("request-id", "123")
	resp.SetBody(string(req.Body()), "text/plain")
	log.Printf("Response %s", resp.String())

	if _, err := conn.Write(resp.Enconde()); err != nil {
		log.Fatalf("unable to write response %v", err)
	}
	if err := conn.Close(); err != nil {
		log.Fatalf("unable to close connection %v", err)
	}

}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		fmt.Println("Failed to bind to port 3000")
	}
	defer l.Close()

	fmt.Println("Server listening on port 3000")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting request ", err.Error())
		}
		go handle(conn)
	}
}
