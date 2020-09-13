package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	currentTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatalln(err)
	}

	const timeOutputFormat = "2006-01-02 15:04:05"
	exactTime := time.Now()

	fmt.Println("current time: ", currentTime.Format(timeOutputFormat))
	fmt.Println("exact time: ", exactTime.Format(timeOutputFormat))
}
