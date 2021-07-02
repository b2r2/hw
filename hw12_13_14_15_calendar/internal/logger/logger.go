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
	Infoln(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Errorln(args ...interface{})
	Fatal(args ...interface{})
	Traceln(args ...interface{})
	WithField(key string, value interface{}) *logrus.Entry
}

type logger struct {
	logger *logrus.Logger
}

func (l logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}

func (l logger) Traceln(args ...interface{}) {
	l.logger.Traceln(args...)
}

func (l logger) Infoln(args ...interface{}) {
	l.logger.Infoln(args...)
}

func (l logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l logger) Errorln(args ...interface{}) {
	l.logger.Errorln(args...)
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
	} else {
		log.SetOutput(os.Stdout)
	}

	if filename != "" && out != nil {
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
