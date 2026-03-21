package bytebuf

import (
	"io"

	"github.com/koykov/simd/memcpy"
)

type Reader interface {
	io.Reader
	io.ReaderAt
	io.ReaderFrom
	io.ByteReader
	io.ByteScanner
	io.RuneReader
	io.RuneScanner
}

type reader struct {
	buf []byte
	off int64
}

func (cr *reader) Read(p []byte) (n int, err error) {
	if cr.off >= int64(len(cr.buf)) {
		err = io.EOF
		return
	}
	n = min(len(p), len(cr.buf)-int(cr.off))
	memcpy.Copy(p, cr.buf[cr.off:cr.off+int64(n)])
	cr.off += int64(n)
	if n < len(p) {
		err = io.EOF
	}
	return
}

func (cr *reader) ReadAt(p []byte, off int64) (n int, err error) {
	// todo implement me
	return
}

func (cr *reader) ReadFrom(r io.Reader) (n int64, err error) {
	// todo implement me
	return
}

func (cr *reader) ReadByte() (byte, error) {
	// todo implement me
	return 0, nil
}

func (cr *reader) UnreadByte() error {
	// todo implement me
	return nil
}

func (cr *reader) ReadRune() (r rune, size int, err error) {
	// todo implement me
	return
}

func (cr *reader) UnreadRune() error {
	// todo implement me
	return nil
}
