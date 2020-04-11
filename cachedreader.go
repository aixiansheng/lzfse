package lzfse

import (
	"io"
)

type cachedReader struct {
	r   io.Reader
	buf []byte
}

func newCachedReader(r io.Reader) *cachedReader {
	return &cachedReader{
		r:   r,
		buf: make([]byte, 0, 1024),
	}
}

func (cr *cachedReader) Read(b []byte) (int, error) {
	n, err := cr.r.Read(b)
	if n > 0 {
		cr.buf = append(cr.buf, b[:n]...)
	}
	return n, err
}

func (cr *cachedReader) Bytes() []byte {
	return cr.buf
}
