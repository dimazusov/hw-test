package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	exactTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatalln(err)
	}

	const timeOutputFormat = "2006-01-02 15:04:05 +0000 UTC"
	currentTime := time.Now()

	fmt.Println("current time:", currentTime.Format(timeOutputFormat))
	fmt.Println("exact time:", exactTime.Format(timeOutputFormat))
}
