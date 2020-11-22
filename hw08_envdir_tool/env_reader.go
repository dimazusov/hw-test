package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]string

const trimString = " \t"
const tSymbols = "\x00"

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := Environment{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		file, err := os.Open(dir+"/"+f.Name())
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		ok := scanner.Scan()
		if ok {
			s := strings.Join(strings.Split(scanner.Text(), tSymbols), " ")
			env[f.Name()] = strings.Trim(s, trimString)
		} else {
			env[f.Name()] = ""
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		file.Close()
	}

	return env, nil
}
