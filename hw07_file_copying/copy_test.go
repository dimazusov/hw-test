package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type DataTestParams struct {
	FromPath     string
	ToPath       string
	ExpectedPath string
	Offset       int64
	Limit        int64
}

func NewDataTestParams(fromPath, toPath, expectedPath string, offset, limit int64) *DataTestParams {
	return &DataTestParams{
		FromPath:     fromPath,
		ToPath:       toPath,
		ExpectedPath: expectedPath,
		Offset:       offset,
		Limit:        limit,
	}
}

func TestCopy(t *testing.T) {
	fromFile := "./testdata/input.txt"
	toFile := "./testdata/test_input.txt"

	testParams := []*DataTestParams{
		NewDataTestParams(fromFile, toFile, "./testdata/out_offset0_limit0.txt", 0, 0),
		NewDataTestParams(fromFile, toFile, "./testdata/out_offset0_limit10.txt", 0, 10),
		NewDataTestParams(fromFile, toFile, "./testdata/out_offset0_limit1000.txt", 0, 1000),
		NewDataTestParams(fromFile, toFile, "./testdata/out_offset0_limit10000.txt", 0, 10000),
		NewDataTestParams(fromFile, toFile, "./testdata/out_offset100_limit1000.txt", 100, 1000),
		NewDataTestParams(fromFile, toFile, "./testdata/out_offset6000_limit1000.txt", 6000, 1000),
	}

	for _, params := range testParams {
		err := Copy(params.FromPath, params.ToPath, params.Offset, params.Limit)
		require.Nil(t, err)

		isEqual, err := isFileEqual(params.ExpectedPath, params.ToPath)
		require.Nil(t, err)
		require.True(t, isEqual)
	}

	os.Remove(toFile)
}

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

		_, fromFileErr := fromFile.Read(expBuf)
		if fromFileErr != nil && fromFileErr != io.EOF {
			return false, err
		}

		_, toFileErr := toFile.Read(givBuf)
		if toFileErr != nil && toFileErr != io.EOF {
			return false, err
		}

		if fromFileErr == io.EOF && toFileErr == io.EOF {
			return true, nil
		}

		for i, _ := range expBuf {
			if expBuf[i] != givBuf[i] {
				return false, nil
			}
		}
	}

}