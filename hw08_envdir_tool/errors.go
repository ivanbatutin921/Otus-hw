package main

import "errors"

var (
	ErrWrongFileName      = errors.New("filename shouldn't contain '=' symbol")
	ErrEmptyDirectoryPath = errors.New("directory path should be set")
)
