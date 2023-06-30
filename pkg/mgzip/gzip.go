package mgzip

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

const _encoding string = "gzip"

type GZipper struct {
	b *bytes.Buffer
}

func (w *GZipper) Compress(data []byte) ([]byte, error) {
	w.b.Reset()
	gz, err := gzip.NewWriterLevel(w.b, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	_, err = gz.Write(data)
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}
	return w.b.Bytes(), err
}
func (w *GZipper) Decompress(data []byte) ([]byte, error) {
	r, _ := gzip.NewReader(bytes.NewReader(data))
	defer r.Close()

	w.b.Reset()
	_, err := w.b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return w.b.Bytes(), nil
}
func (w *GZipper) GetEncoding() string {
	return _encoding
}

func NewGZipper() Zipper {
	var buf bytes.Buffer
	w := &GZipper{b: &buf}
	return w
}
