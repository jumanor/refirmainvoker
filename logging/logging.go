package logging

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	path, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}
	pathFileLog := filepath.Join(path, "main.log")

	file, err := os.OpenFile(pathFileLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	var writers []io.Writer
	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	writers = append(writers, file)
	mw := io.MultiWriter(writers...)
	log = zerolog.New(mw).With().Caller().Timestamp().Logger()
}
func Log() *zerolog.Logger {
	return &log
}
