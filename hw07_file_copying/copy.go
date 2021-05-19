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
	ErrEmptyFromPath         = errors.New("from path is empty")
	ErrEmptyToPath           = errors.New("to path is empty")
	ErrEqualPaths            = errors.New("from path should not be the same as to path")
	ErrNegativeLimit         = errors.New("limit must not be less than zero")
	ErrNegativeOffset        = errors.New("offset must not be less than zero")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	i, err := checkValid(fromPath, toPath, offset, limit)
	if err != nil {
		return err
	}

	size := i.Size()
	if limit == 0 || limit > size-offset {
		limit = size - offset
	}

	var r, w *os.File
	if r, err = os.OpenFile(fromPath, os.O_RDONLY, 0o600); err != nil {
		return err
	}
	defer syncClose(r)
	if _, err := r.Seek(offset, io.SeekStart); err != nil {
		return err
	}
	if w, err = os.Create(toPath); err != nil {
		return err
	}
	defer syncClose(w)

	b := pb.Start64(limit)
	br := b.NewProxyReader(r)
	defer b.Finish()
	if _, err := io.CopyN(w, br, limit); err != nil {
		return err
	}

	return nil
}

func checkValid(fromPath, toPath string, offset, limit int64) (os.FileInfo, error) {
	if fromPath == "" {
		return nil, ErrEmptyFromPath
	}
	if toPath == "" {
		return nil, ErrEmptyToPath
	}
	if fromPath == toPath {
		return nil, ErrEqualPaths
	}
	if limit < 0 {
		return nil, ErrNegativeLimit
	}
	if offset < 0 {
		return nil, ErrNegativeOffset
	}

	i, err := os.Stat(fromPath)
	if err != nil {
		return nil, err
	}
	if !i.Mode().IsRegular() {
		return nil, ErrUnsupportedFile
	}

	if offset > i.Size() {
		return nil, ErrOffsetExceedsFileSize
	}
	return i, nil
}

func syncClose(f *os.File) error {
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
