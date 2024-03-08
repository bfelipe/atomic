package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"gitlab.com/bfelipe/atomic"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Record struct {
	User      User      `json:"user"`
	Products  []Product `json:"products"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

func handle(conn net.Conn) {
	defer conn.Close()
	req := atomic.Request{}
	req.Decode(conn)

	for k, v := range req.Headers {
		fmt.Printf("%s: %s\n", k, v)
	}

	r := Record{}
	err := json.Unmarshal(req.Body, &r)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	return
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
