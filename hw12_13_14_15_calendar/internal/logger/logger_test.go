package logger

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	filePath := "logfile.txt"
	os.Remove(filePath)

	logger, err := New(filePath, debugLevel)
	require.Nil(t, err)
	require.NotNil(t, logger)

	err = logger.Info("test, data")
	require.Nil(t, err)

	file, err := os.Open(filePath)
	require.Nil(t, err)

	b, err := ioutil.ReadAll(file)
	require.Nil(t, err)

	formatedTime := time.Now().Format(time.RFC822)
	loggerString := fmt.Sprintf("%s info: test, data", formatedTime)
	require.Equal(t, loggerString, string(bytes.TrimRight(b, "\n")))
}
