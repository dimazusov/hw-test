package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dimazusov/hw-test/hw11_telnet_client/telnet"
)

func main() {
	var timeout string
	flag.StringVar(&timeout, "timeout", "10s", "timeout for connection closing")
	flag.Parse()

	var d time.Duration
	var host, port string

	var err error
	d, err = time.ParseDuration(timeout)
	if err != nil {
		log.Fatalln(err, "cannot parse timeout")
	}

	if len(os.Args) != 3 && len(os.Args) != 4 {
		log.Fatalln("wrong args count")
	}

	if len(os.Args) == 4 {
		host = os.Args[2]
		port = os.Args[3]
	} else {
		host = os.Args[1]
		port = os.Args[2]
	}

	client := telnet.NewTelnetClient(fmt.Sprintf("%s:%s", host, port), d, os.Stdin, os.Stdout)
	ctx, cancel := context.WithCancel(context.Background())

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err, "cannot connect")
	}
	defer client.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		err := client.Receive()
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()
	go func() {
		err := client.Send()
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	wg.Wait()
}
