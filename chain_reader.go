package bytebuf

import "io"

type ChainReader struct {
	buf *Chain
	i   int64
}

func (cr *ChainReader) Read(p []byte) (n int, err error) {
	if n = cr.buf.Len() - int(cr.i); n < len(p) {
		copy(p, cr.buf.Bytes()[cr.i:])
		err = io.EOF
		cr.i += int64(n)
		return
	}
	n = len(p)
	copy(p, cr.buf.Bytes()[:n])
	cr.i += int64(n)
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
