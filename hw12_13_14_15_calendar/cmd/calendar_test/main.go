package main

import (
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/integration_tests"
	"testing"
)

func main() {
	integration_tests.RunTests(&testing.T{})
}