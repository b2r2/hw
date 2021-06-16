package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Traceln(args ...interface{})
}

type logger struct {
	logger *logrus.Logger
}

func (l logger) Traceln(args ...interface{}) {
	l.logger.Traceln(args...)
}

func (l logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func New(level, filename string, out io.Writer) (Logger, error) {
	log := logrus.New()

	if out != nil {
		log.SetOutput(out)
	} else if filename != "" {
		path, err := filepath.Abs(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to open logfile: %w", err)
		}
		if err := os.Mkdir(filepath.Dir(path), os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directory to logfile: %w", err)
		}

		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open logfile: %w", err)
		}
		log.SetOutput(file)
	}

	if level != "" {
		lvl, err := logrus.ParseLevel(level)
		if err != nil {
			return nil, fmt.Errorf("failed to parse log level: %w", err)
		}

		log.SetLevel(lvl)
	}

	return logger{logger: log}, nil
}
