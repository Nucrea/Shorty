package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

func New(outputFile string) (Logger, error) {
	writers := []io.Writer{}
	writers = append(writers, os.Stdout)

	if outputFile != "" {
		file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err != nil {
			return nil, err
		}
		writers = append(writers, file)
	}

	writer := io.MultiWriter(writers...)
	l := zerolog.New(writer).
		Level(zerolog.DebugLevel).
		With().Timestamp().
		Logger()

	return &logger{
		zeroLogger: &l,
	}, nil
}
