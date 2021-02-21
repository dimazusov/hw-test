package test

import (
	"github.com/dimazusov/hw-test/hw12_13_14_15_calendar/internal/integration_tests"
	"testing"
)

func TestIntegrations(t *testing.T) {
	integration_tests.RunTests(t)
}