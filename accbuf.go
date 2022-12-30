package bytebuf

import (
	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

// AccumulativeBuf is a wrapper around ChainBuf that allows to accumulate buffer data and use only necessary part.
//
// See StakeOut and Staked* methods.
type AccumulativeBuf struct {
	buf ChainBuf
	off int
	err error
}

// StakeOut saves current offset for further use.
func (b *AccumulativeBuf) StakeOut() *AccumulativeBuf {
	b.off = b.Len()
	return b
}

// StakedOffset returns staked offset.
func (b *AccumulativeBuf) StakedOffset() int {
	return b.off
}

// StakedLen returns length of accumulated bytes since staked offset.
func (b *AccumulativeBuf) StakedLen() int {
	return b.Len() - b.off
}

// StakedBytes returns accumulated bytes from staked offset.
func (b *AccumulativeBuf) StakedBytes() []byte {
	if b.off >= b.Len() {
		return nil
	}
	return b.buf[b.off:]
}

// StakedBytesCopy returns copy of accumulated bytes since staked offset.
func (b *AccumulativeBuf) StakedBytesCopy() []byte {
	return bytealg.Copy(b.StakedBytes())
}

// StakedString returns accumulated bytes as string.
func (b *AccumulativeBuf) StakedString() string {
	if b.off >= b.Len() {
		return ""
	}
	return b.String()[b.off:]
}

// StakedStringCopy returns copy of accumulated bytes as string.
func (b *AccumulativeBuf) StakedStringCopy() string {
	return bytealg.CopyStr(b.StakedString())
}

// RangeBytes returns buffer bytes from offset off with length len.
func (b *AccumulativeBuf) RangeBytes(off, len int) []byte {
	if off >= 0 && off+len < b.buf.Len() {
		return nil
	}
	return b.buf[off : off+len]
}

// RangeBytesCopy copies result of RangeBytes().
func (b *AccumulativeBuf) RangeBytesCopy(off, len int) []byte {
	return bytealg.Copy(b.RangeBytes(off, len))
}

// RangeString returns buffer bytes as string from offset off with length len.
func (b *AccumulativeBuf) RangeString(off, len int) string {
	if off >= 0 && off+len < b.buf.Len() {
		return ""
	}
	return fastconv.B2S(b.buf[off : off+len])
}

// RangeStringCopy copies result of RangeString().
func (b *AccumulativeBuf) RangeStringCopy(off, len int) string {
	return bytealg.CopyStr(b.RangeString(off, len))
}

// Get last error caught in Write* methods.
func (b AccumulativeBuf) Error() error {
	return b.err
}
