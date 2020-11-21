package main

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

var (
	ErrMessageOpenFile              = "open file error"
	ErrMessageReadFile              = "read file error"
	ErrMessageWriteFile             = "write file error"
	ErrMessageOffsetExceedsFileSize = "offset exceeds file size"
)

const bufferSize = 2048

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath) // For read access.
	if err != nil {
		return errors.Wrap(err, ErrMessageOpenFile)
	}
	defer fromFile.Close()

	toFile, err := os.Create(toPath) // For read access.
	if err != nil {
		return errors.Wrap(err, ErrMessageOpenFile)
	}
	defer toFile.Close()

	_, err = fromFile.Seek(offset, 0)
	if err != nil {
		return errors.Wrap(err, ErrMessageOffsetExceedsFileSize)
	}

	countWritesBytes := int64(0)
	for {
		b := make([]byte, bufferSize)
		n, err := fromFile.Read(b)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return errors.Wrap(err, ErrMessageReadFile)
		}

		if countWritesBytes+bufferSize > limit && limit != 0 {
			_, err = toFile.Write(b[:countWritesBytes+limit])
			if err != nil {
				return errors.Wrap(err, ErrMessageWriteFile)
			}

			return nil
		}

		if n != bufferSize {
			b = b[:n]
		}
		_, err = toFile.Write(b)
		if err != nil {
			return errors.Wrap(err, ErrMessageWriteFile)
		}

		countWritesBytes += int64(n)
	}
}
