package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func isFileEqual(fromFilePath, toFilePath string) (isEqual bool, err error) {
	fromFile, err := os.Open(fromFilePath) // For read access.
	if err != nil {
		return false, err
	}
	defer fromFile.Close()

	toFile, err := os.Open(toFilePath) // For read access.
	if err != nil {
		return false, err
	}
	defer toFile.Close()

	for {
		expBuf := make([]byte, bufferSize)
		givBuf := make([]byte, bufferSize)

		expCount, fromFileErr := fromFile.Read(expBuf)
		if fromFileErr != nil && fromFileErr != io.EOF {
			return false, err
		}

		givCount, toFileErr := toFile.Read(givBuf)
		if toFileErr != nil && toFileErr != io.EOF {
			return false, err
		}

		if fromFileErr == io.EOF && toFileErr == io.EOF {
			return true, nil
		}

		if expCount != givCount {
			return false, nil
		}

		for i, _ := range expBuf {
			if expBuf[i] != givBuf[i] {
				return false, nil
			}
		}
	}

}

type DataTestParams struct {
	FromPath string
	ToPath   string
	Offset   int64
	Limit    int64
}

func NewDataTestParams(fromPath, toPath string, offset, limit int64) *DataTestParams {
	return &DataTestParams{
		FromPath: fromPath,
		ToPath:   toPath,
		Offset:   offset,
		Limit:    limit,
	}
}

func TestCopy(t *testing.T) {
	fromFile := "./testdata/input.txt"
	toFile := "./testdata/test_input.txt"

	testParams := []*DataTestParams{
		NewDataTestParams(fromFile, toFile, 0, 0),
		NewDataTestParams(fromFile, toFile, 10, 1000),
		NewDataTestParams(fromFile, toFile, 100, 200),
		NewDataTestParams(fromFile, toFile, 213, 1),
		NewDataTestParams(fromFile, toFile, 590, 3),
		NewDataTestParams(fromFile, toFile, 600, 600),
		NewDataTestParams(fromFile, toFile, 1553, 145),
	}

	for _, params := range testParams {
		err := Copy(params.FromPath, params.ToPath, offset, limit)
		require.Nil(t, err)

		isEqual, err := isFileEqual(params.FromPath, params.ToPath)
		require.Nil(t, err)
		require.True(t, isEqual)
	}

	os.Remove(toFile)
}
