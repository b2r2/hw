package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrEmptyPath             = errors.New("empty path")
	ErrEmptyToPath           = errors.New("the path must not be empty")
	ErrEqualPaths            = errors.New("from path should not be the same as to path")
	ErrNegativeLimit         = errors.New("limit must not be less than zero")
	ErrNegativeOffset        = errors.New("offset must not be less than zero")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := checkValid(fromPath, toPath, offset, limit); err != nil {
		return err
	}

	i, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	size := i.Size()
	if limit == 0 || limit > size-offset {
		limit = size - offset
	}

	var r, w *os.File
	if r, err = os.OpenFile(fromPath, os.O_RDONLY, 0o600); err != nil {
		fmt.Println("open")
		return err
	}
	defer syncClose(r)
	if _, err := r.Seek(offset, io.SeekStart); err != nil {
		return err
	}
	if w, err = os.Create(toPath); err != nil {
		if err := os.Mkdir(path.Dir(toPath), os.ModePerm); err != nil {
			return err
		}
		w, _ = os.Create(toPath)
	}

	b := pb.Start64(limit)
	br := b.NewProxyReader(r)
	defer b.Finish()
	if _, err := io.CopyN(w, br, limit); err != nil {
		return err
	}

	return nil
}

func checkValid(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" {
		return ErrEmptyPath
	}
	if toPath == "" {
		return ErrEmptyToPath
	}
	if fromPath == toPath {
		return ErrEqualPaths
	}
	if limit < 0 {
		return ErrNegativeLimit
	}
	if offset < 0 {
		return ErrNegativeOffset
	}

	i, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	size := i.Size()
	if size == 0 || i.IsDir() {
		return ErrUnsupportedFile
	}

	if offset > size {
		return ErrOffsetExceedsFileSize
	}
	return nil
}

func syncClose(f *os.File) {
	if err := f.Sync(); err != nil {
		log.Fatalf(err.Error())
	}
	if err := f.Close(); err != nil {
		log.Fatalf(err.Error())
	}
}
