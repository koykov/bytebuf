package bytebuf

import (
	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

// A wrapper around ChainBuf that allows to accumulate buffer data and use only necessary part.
//
// See StakeOut and Staked* methods.
type AccumulativeBuffer struct {
	buf ChainBuf
	off int
}

// Stake out current offset for further use.
func (b *AccumulativeBuffer) StakeOut() *AccumulativeBuffer {
	b.off = b.Len()
	return b
}

// Get staked offset.
func (b *AccumulativeBuffer) StakedOffset() int {
	return b.off
}

// Get accumulated bytes from staked offset.
func (b *AccumulativeBuffer) StakedBytes() []byte {
	if b.off >= b.Len() {
		return nil
	}
	return b.buf[b.off:]
}

// Get copy of accumulated bytes from staked offset.
func (b *AccumulativeBuffer) StakedBytesCopy() []byte {
	return bytealg.Copy(b.StakedBytes())
}

// Get accumulated bytes as string.
func (b *AccumulativeBuffer) StakedString() string {
	if b.off >= b.Len() {
		return ""
	}
	return b.String()[b.off:]
}

// Get copy of accumulated bytes as string.
func (b *AccumulativeBuffer) StakedStringCopy() string {
	return bytealg.CopyStr(b.StakedString())
}

// Get buffer bytes from offset off with length len.
func (b *AccumulativeBuffer) RangeBytes(off, len int) []byte {
	if off >= 0 && off+len < b.buf.Len() {
		return nil
	}
	return b.buf[off : off+len]
}

// Copy version of RangeBytes().
func (b *AccumulativeBuffer) RangeBytesCopy(off, len int) []byte {
	return bytealg.Copy(b.RangeBytes(off, len))
}

// Get buffer bytes as string from offset off with length len.
func (b *AccumulativeBuffer) RangeString(off, len int) string {
	if off >= 0 && off+len < b.buf.Len() {
		return ""
	}
	return fastconv.B2S(b.buf[off : off+len])
}

// Copy version of RangeString().
func (b *AccumulativeBuffer) RangeStringCopy(off, len int) string {
	return bytealg.CopyStr(b.RangeString(off, len))
}
