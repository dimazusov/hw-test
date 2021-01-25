package main

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

var (
	ErrMessageOpenFile              = "open file error"
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
		curBufSize := int64(bufferSize)
		if countWritesBytes+curBufSize > limit && limit != 0 {
			curBufSize = limit - countWritesBytes
		}

		n, err := io.CopyN(toFile, fromFile, curBufSize)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return errors.Wrap(err, ErrMessageWriteFile)
		}
		countWritesBytes += n

		if countWritesBytes >= limit && limit != 0 {
			return nil
		}
	}
}
