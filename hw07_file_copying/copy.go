package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const bufferSize = 2048

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath) // For read access.
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := os.Create(toPath) // For read access.
	if err != nil {
		return err
	}
	defer toFile.Close()

	_, err = fromFile.Seek(offset, 0)
	if err != nil {
		return ErrOffsetExceedsFileSize
	}

	countWritesBytes := int64(0)
	for {
		b := make([]byte, bufferSize)
		n, err := fromFile.Read(b)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if countWritesBytes+bufferSize > limit && limit != 0 {
			_, err = toFile.Write(b[:countWritesBytes+limit])
			if err != nil {
				return err
			}

			return nil
		}

		if n != bufferSize {
			b = b[:n]
		}
		_, err = toFile.Write(b)
		if err != nil {
			return err
		}

		countWritesBytes += int64(n)
	}
}
