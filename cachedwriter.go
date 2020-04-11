package lzfse

import (
	"io"
)

type cachedWriter struct {
	w   io.Writer
	buf []byte
}

func newCachedWriter(w io.Writer) *cachedWriter {
	return &cachedWriter{
		w:   w,
		buf: make([]byte, 0, 1024),
	}
}

func (cw *cachedWriter) Write(b []byte) (int, error) {
	n, err := cw.w.Write(b)
	if n > 0 {
		cw.buf = append(cw.buf, b[:n]...)
	}
	return n, err
}

func (cw *cachedWriter) ReadRelativeToEnd(b []byte, offset int64) (copied int, err error) {
	copied = copy(b, cw.buf[int64(len(cw.buf))-offset:])
	return
}
