package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	path, err := os.MkdirTemp("", "myenv")
	require.NoError(t, err)
	defer os.RemoveAll(path)

	t.Run("dir is empty", func(t *testing.T) {
		_, err = ReadDir(path)
		require.NoError(t, err)
	})

	t.Run("filename with '=' ", func(t *testing.T) {
		filename := path + "/te=st"
		_, err = os.Create(filename)
		require.NoError(t, err)
		_, err = ReadDir(path)
		require.ErrorIs(t, err, ErrWrongFileName)
	})

	t.Run("dir path is empty", func(t *testing.T) {
		dir := ""
		_, err = ReadDir(dir)
		require.ErrorIs(t, err, ErrEmptyDirectoryPath)
	})
}
