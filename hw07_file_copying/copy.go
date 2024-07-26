package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

var bytesToCopy int64

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	fromFileInfo, err := fromFile.Stat()
	if err != nil {
		return err
	}

	if fromFileInfo.Mode()&os.ModeDevice != 0 {
		return ErrUnsupportedFile
	}

	if offset > fromFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	if limit == 0 {
		bytesToCopy = fromFileInfo.Size() - offset
	} else {
		bytesToCopy = limit
	}

	bar := pb.StartNew(int(bytesToCopy))

	_, err = io.CopyN(toFile, io.NewSectionReader(fromFile, offset, bytesToCopy), bytesToCopy)
	if err != nil {
		if err == io.EOF {
			// EOF is expected when we reach the end of the file
			err = nil
		} else {
			return err
		}
	}

	bar.Finish()

	return nil
}
