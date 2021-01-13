package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/dimazusov/hw-test/hw11_telnet_client/telnet"
)

func main() {
	var timeout string
	flag.StringVar(&timeout, "timeout", "", "timeout for connection closing")
	flag.Parse()

	var d time.Duration
	var host, port string

	if timeout != "" {
		if len(os.Args) != 4 {
			log.Fatalln("wrong args count")
		}

		var err error
		d, err = time.ParseDuration(timeout)
		if err != nil {
			log.Fatalln("cannot parse timeout")
		}

		host = os.Args[2]
		port = os.Args[3]
	} else {
		if len(os.Args) != 3 {
			log.Fatalln("wrong args count")
		}

		host = os.Args[1]
		port = os.Args[2]
	}

	fmt.Println(fmt.Sprintf("%s:%s", host, port))
	fmt.Println(d)
	client := telnet.NewTelnetClient(fmt.Sprintf("%s:%s", host, port), d, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatalln(err, "cannot connect")
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		client.Receive()
		wg.Done()
	}()
	go func() {
		client.Send()
		wg.Done()
	}()
	wg.Wait()
}