package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	env, err := ReadDir("./testdata/env")
	require.Nil(t, err)
	require.NotNil(t, env)
	require.Equal(t, env["BAR"], "bar")
	require.Equal(t, env["FOO"], "   foo\nwith new line")
	require.Equal(t, env["HELLO"], `"hello"`)
	require.Equal(t, env["UNSET"], "")
}
