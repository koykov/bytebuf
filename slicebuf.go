package bytebuf

import (
	"unicode/utf8"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

// SliceBuf is a buffer implementation around bytes slice.
// Is a replacement of bytes.Buffer when needs to convert byte slice to buffer object.
type SliceBuf []byte

func (b *SliceBuf) Bytes() []byte {
	return *b
}

func (b *SliceBuf) String() string {
	return fastconv.B2S(*b)
}

func (b *SliceBuf) Len() int {
	return len(*b)
}

func (b *SliceBuf) Cap() int {
	return cap(*b)
}

func (b *SliceBuf) Truncate(n int) {
	if n < 0 {
		*b = (*b)[:0]
	} else if n > len(*b) {
		// do noting
	} else {
		*b = (*b)[:n]
	}
}

func (b *SliceBuf) Reset() {
	*b = (*b)[:0]
}

func (b *SliceBuf) Grow(n int) {
	*b = bytealg.Grow(*b, n)
}

func (b *SliceBuf) Write(p []byte) (n int, err error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *SliceBuf) WriteString(s string) (n int, err error) {
	return b.Write(fastconv.S2B(s))
}

func (b *SliceBuf) WriteByte(c byte) error {
	*b = append(*b, c)
	return nil
}

func (b *SliceBuf) WriteRune(r rune) (n int, err error) {
	if uint32(r) < utf8.RuneSelf {
		_ = b.WriteByte(byte(r))
		n = 1
		return
	}
	off := len(*b)
	*b = bytealg.GrowDelta(*b, utf8.UTFMax)
	n = utf8.EncodeRune((*b)[off:off+utf8.UTFMax], r)
	*b = (*b)[:off+n]
	return
}

// todo: implement rest of bytes.Buffer methods.
