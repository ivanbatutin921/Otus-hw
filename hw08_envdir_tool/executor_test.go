package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("empty env", func(t *testing.T) {
		cmdArds := []string{"", "", "/bin/bash", "./testdata/echo.sh", "1", "1"}
		resultCode := RunCmd(cmdArds, nil)
		require.Equal(t, resultCode, 0)
	})
	t.Run("Incorrect script", func(t *testing.T) {
		cmdArds := []string{"", "", "/bin/bash", "./testdata/echo.s", "1", "1"}
		resultCode := RunCmd(cmdArds, nil)
		require.Equal(t, 127, resultCode)
	})
}
