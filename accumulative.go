package bytebuf

import (
	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

// Accumulative is a wrapper around Chain that allows to accumulate buffer data and use only necessary part.
//
// See StakeOut and Staked* methods.
type Accumulative struct {
	buf Chain
	off int
	err error
}

// StakeOut saves current offset for further use.
func (b *Accumulative) StakeOut() *Accumulative {
	b.off = b.Len()
	return b
}

// StakedOffset returns staked offset.
func (b *Accumulative) StakedOffset() int {
	return b.off
}

// StakedLen returns length of accumulated bytes since staked offset.
func (b *Accumulative) StakedLen() int {
	return b.Len() - b.off
}

// StakedBytes returns accumulated bytes from staked offset.
func (b *Accumulative) StakedBytes() []byte {
	if b.off >= b.Len() {
		return nil
	}
	return b.buf[b.off:]
}

// StakedBytesCopy returns copy of accumulated bytes since staked offset.
func (b *Accumulative) StakedBytesCopy() []byte {
	return bytealg.Copy(b.StakedBytes())
}

// StakedString returns accumulated bytes as string.
func (b *Accumulative) StakedString() string {
	if b.off >= b.Len() {
		return ""
	}
	return b.String()[b.off:]
}

// StakedStringCopy returns copy of accumulated bytes as string.
func (b *Accumulative) StakedStringCopy() string {
	return bytealg.Copy[string](b.StakedString())
}

// RangeBytes returns buffer bytes from offset off with length len.
func (b *Accumulative) RangeBytes(off, len int) []byte {
	if off >= 0 && off+len < b.buf.Len() {
		return nil
	}
	return b.buf[off : off+len]
}

// RangeBytesCopy copies result of RangeBytes().
func (b *Accumulative) RangeBytesCopy(off, len int) []byte {
	return bytealg.Copy(b.RangeBytes(off, len))
}

// RangeString returns buffer bytes as string from offset off with length len.
func (b *Accumulative) RangeString(off, len int) string {
	if off >= 0 && off+len < b.buf.Len() {
		return ""
	}
	return fastconv.B2S(b.buf[off : off+len])
}

// RangeStringCopy copies result of RangeString().
func (b *Accumulative) RangeStringCopy(off, len int) string {
	return bytealg.Copy[string](b.RangeString(off, len))
}

// Get last error caught in Write* methods.
func (b *Accumulative) Error() error {
	return b.err
}

// ToWriter wraps buffer with class implementing IO interfaces.
func (b *Accumulative) ToWriter() *AccumulativeWriter {
	return &AccumulativeWriter{buf: b}
}
