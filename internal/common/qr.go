package common

import (
	"bytes"
	"encoding/base64"
	"io"

	qrcode "github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type BytesWriterCloser struct {
	io.Writer
}

func (BytesWriterCloser) Close() error {
	return nil
}

func NewQRBase64(value string) (string, error) {
	qrc, err := qrcode.New(value)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(nil)
	wc := BytesWriterCloser{buf}

	wr := standard.NewWithWriter(wc, standard.WithQRWidth(20))
	if err := qrc.Save(wr); err != nil {
		return "", err
	}

	qrBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return qrBase64, err
}
