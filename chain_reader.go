package bytebuf

import "io"

type ChainReader struct {
	buf *Chain
	off int64
}

func (cr *ChainReader) Read(p []byte) (n int, err error) {
	if cr.off >= int64(cr.buf.Len()) {
		err = io.EOF
		return
	}
	n = copy(p, cr.buf.Bytes()[cr.off:])
	cr.off += int64(n)
	if n < len(p) {
		err = io.EOF
	}
	return
}

func (cr *ChainReader) ReadAt(p []byte, off int64) (n int, err error) {
	// todo implement me
	return
}

func (cr *ChainReader) ReadFrom(r io.Reader) (n int64, err error) {
	// todo implement me
	return
}

func (cr *ChainReader) ReadByte() (byte, error) {
	// todo implement me
	return 0, nil
}

func (cr *ChainReader) ReadRune() (r rune, size int, err error) {
	// todo implement me
	return
}
