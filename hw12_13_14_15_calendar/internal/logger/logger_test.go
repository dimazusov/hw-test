package logger

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	filePath := "logfile.txt"
	os.Remove(filePath)

	logger, err := New(filePath, debugLevel)
	require.Nil(t, err)
	require.NotNil(t, logger)

	logger.Info("test, data")

	file, err := os.Open(filePath)
	require.Nil(t, err)

	b, err := ioutil.ReadAll(file)
	require.Equal(t, "01 Jan 21 17:57 MSK debug: test, data", string(bytes.TrimRight(b,"\n")))
}
