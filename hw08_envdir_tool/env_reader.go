package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Environment map[string]string

const trimString = " \t"
const tSymbols = "\x00"

var ErrReading = errors.New("read error")
var ErrOpen = errors.New("open error")
var ErrScan = errors.New("scan error")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := Environment{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, ErrReading)
	}

	for _, f := range files {
		file, err := os.Open(dir + "/" + f.Name())
		if err != nil {
			return nil, errors.Wrap(err, ErrOpen)
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
			return nil, errors.Wrap(err, ErrScan)
		}

		file.Close()
	}

	return env, nil
}
