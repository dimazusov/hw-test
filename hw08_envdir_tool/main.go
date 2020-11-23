package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var ErrCountArgs = errors.New("WrongCountArgs")

func main() {
	if len(os.Args) < 2 {
		log.Fatalln(ErrCountArgs)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	res := RunCmd(os.Args[2:], env)
	fmt.Println(res)
	os.Exit(res)
}
