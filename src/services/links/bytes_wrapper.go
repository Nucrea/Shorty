package links

import "io"

type BytesWriterCloser struct {
	io.Writer
}

func (BytesWriterCloser) Close() error {
	return nil
}
