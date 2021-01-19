package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{}
	env["BAR"] = "b"
	env["FOO"] = "bar"

	cmd := []string{"echo", "$BAR"}

	returnCode := RunCmd(cmd, env)
	require.Equal(t, returnCode, 0)
}
