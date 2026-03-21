package bytebuf

import (
	"errors"
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
	n = cr.min(len(p), len(cr.buf)-int(cr.off))
	memcpy.Copy(p, cr.buf[cr.off:cr.off+int64(n)])
	cr.off += int64(n)
	if n < len(p) {
		err = io.EOF
	}
	return
}

func (cr *reader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		err = ErrNegativeOffset
		return
	}
	if off > int64(len(cr.buf)) {
		err = io.EOF
		return
	}
	n = cr.min(len(p), len(cr.buf)-int(off))
	memcpy.Copy(p, cr.buf[off:])
	if n < len(p) {
		err = io.EOF
	}
	return
}

func (cr *reader) ReadFrom(r io.Reader) (n int64, err error) {
	// todo implement me
	return
}

func (cr *reader) ReadByte() (b byte, err error) {
	if cr.off >= int64(len(cr.buf)) {
		err = io.EOF
		return
	}
	b = cr.buf[cr.off]
	cr.off++
	return
}

func (cr *reader) UnreadByte() error {
	if cr.off <= 0 {
		return ErrOutOfRange
	}
	cr.off--
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

func (cr *reader) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var (
	ErrNegativeOffset = errors.New("negative offset")
	ErrOutOfRange     = errors.New("out of range")
)
