package main

import (
	"log"
	"os"
	"time"
)

func main() {
	var address string
	var timeout time.Duration

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatal(err, "cannot connect")
	}
	defer client.Close()

	go func() {
		client.Receive()
	}()
	go func() {
		client.Send()
	}()
}
