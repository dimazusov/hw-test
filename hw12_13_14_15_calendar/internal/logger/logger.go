package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var ErrLevelNotExists = errors.New("level not exists")

const (
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errLevel   = "err"
)

var logLevels = map[string]int{
	debugLevel: 1,
	infoLevel:  2,
	warnLevel:  3,
	errLevel:   4,
}

type logger struct {
	mu    sync.Mutex
	file  *os.File
	level int
}

type Logger interface {
	Debug(data interface{}) error
	Info(data interface{}) error
	Warn(data interface{}) error
	Error(data interface{}) error
	Close() error
}

func New(filePath string, level string) (Logger, error) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create logger")
	}

	levelVal, ok := logLevels[level]
	if !ok {
		return nil, errors.Wrap(ErrLevelNotExists, "cannot create logger")
	}

	return &logger{
		file:  f,
		level: levelVal,
	}, nil
}

func (m *logger) Debug(data interface{}) error {
	return m.write(data, debugLevel)
}

func (m *logger) Info(data interface{}) error {
	return m.write(data, infoLevel)
}

func (m *logger) Warn(data interface{}) error {
	return m.write(data, warnLevel)
}

func (m *logger) Error(data interface{}) error {
	return m.write(data, errLevel)
}

func (m *logger) Close() error {
	err := m.file.Close()
	if err != nil {
		return errors.Wrap(err, "cannot close file")
	}

	return nil
}

func (m *logger) write(data interface{}, level string) error {
	if !m.isNeedLogForLevel(level) {
		return nil
	}

	t := time.Now().Format(time.RFC822)
	m.mu.Lock()
	_, err := m.file.WriteString(fmt.Sprintf("%s %s: %v\n", t, level, data))
	m.mu.Unlock()

	if err != nil {
		return errors.Wrap(err, "cannot write")
	}

	return nil
}

func (m *logger) isNeedLogForLevel(loggableLevel string) bool {
	return logLevels[loggableLevel] >= m.level
}
