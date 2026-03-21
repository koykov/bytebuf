package bytebuf

import (
	"errors"
	"io"
	"unicode/utf8"

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
	io.Seeker
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
	if cr.off >= int64(len(cr.buf)) {
		err = io.EOF
		return
	}
	if b := cr.buf[cr.off]; b < utf8.RuneSelf {
		r = rune(b)
		size = 1
		return
	}
	r, size = utf8.DecodeRune(cr.buf[cr.off:])
	cr.off += int64(size)
	return
}

func (cr *reader) UnreadRune() error {
	if cr.off <= 0 {
		return ErrOutOfRange
	}
	r, size := utf8.DecodeLastRune(cr.buf[:cr.off])
	if r == utf8.RuneError {
		return ErrBadRune
	}
	cr.off -= int64(size)
	return nil
}

func (cr *reader) Seek(off int64, whence int) (int64, error) {
	var pos int64
	switch whence {
	case io.SeekStart:
		pos = off
	case io.SeekCurrent:
		pos = cr.off + off
	case io.SeekEnd:
		pos = int64(len(cr.buf)) + off
	default:
		return 0, ErrUnknownWhence
	}
	if pos < 0 {
		return 0, ErrNegativePosition
	}
	cr.off = pos
	return 0, nil
}

func (cr *reader) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var (
	ErrNegativeOffset   = errors.New("negative offset")
	ErrNegativePosition = errors.New("negative position")
	ErrOutOfRange       = errors.New("out of range")
	ErrBadRune          = errors.New("bad rune")
	ErrUnknownWhence    = errors.New("unknown whence")
)
