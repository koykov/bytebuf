package bytebuf

import (
	"github.com/koykov/x2bytes"
)

// Contains wrapper methods over inner ChainBuf.

// Bytes returns contents of the buffer.
func (b *AccumulativeBuf) Bytes() []byte {
	return b.buf
}

// BytesCopy returns copy of buffer contents.
func (b *AccumulativeBuf) BytesCopy() []byte {
	return b.buf.BytesCopy()
}

// Get contents of the buffer as string.
func (b *AccumulativeBuf) String() string {
	return b.buf.String()
}

// StringCopy returns copy of the buffer contents as string.
func (b *AccumulativeBuf) StringCopy() string {
	return b.buf.StringCopy()
}

// Write bytes to the buffer.
func (b *AccumulativeBuf) Write(p []byte) *AccumulativeBuf {
	b.buf.Write(p)
	return b
}

// WriteByte writes single byte.
func (b *AccumulativeBuf) WriteByte(p byte) *AccumulativeBuf {
	b.buf.WriteByte(p)
	return b
}

// WriteStr writes string to the buffer.
// DEPRECATED: use WriteString() instead.
func (b *AccumulativeBuf) WriteStr(s string) *AccumulativeBuf {
	return b.WriteString(s)
}

// WriteString writes string to the buffer.
func (b *AccumulativeBuf) WriteString(s string) *AccumulativeBuf {
	b.buf.WriteString(s)
	return b
}

// WriteInt writes integer value to the buffer.
func (b *AccumulativeBuf) WriteInt(i int64) *AccumulativeBuf {
	b.buf, b.err = x2bytes.IntToBytes(b.buf, i)
	return b
}

// WriteUint writes unsigned integer value to the buffer.
func (b *AccumulativeBuf) WriteUint(u uint64) *AccumulativeBuf {
	b.buf, b.err = x2bytes.UintToBytes(b.buf, u)
	return b
}

// WriteFloat writes float value to the buffer.
func (b *AccumulativeBuf) WriteFloat(f float64) *AccumulativeBuf {
	b.buf, b.err = x2bytes.FloatToBytes(b.buf, f)
	return b
}

// WriteBool writes boolean value to the buffer.
func (b *AccumulativeBuf) WriteBool(v bool) *AccumulativeBuf {
	b.buf, b.err = x2bytes.BoolToBytes(b.buf, v)
	return b
}

// WriteX write v with arbitrary type to the buffer.
func (b *AccumulativeBuf) WriteX(x interface{}) *AccumulativeBuf {
	b.buf, b.err = x2bytes.ToBytes(b.buf, x)
	return b
}

// WriteMarshallerTo marshal data of struct implemented MarshallerTo interface and write it to the buffer.
func (b *AccumulativeBuf) WriteMarshallerTo(m MarshallerTo) *AccumulativeBuf {
	if m == nil {
		return b
	}
	n := b.Len()
	d := m.Size()
	b.GrowDelta(d)
	_, b.err = m.MarshalTo(b.buf[n:])
	return b
}

// WriteApplyFn applies fn to p and write result to the buffer.
func (b *AccumulativeBuf) WriteApplyFn(p []byte, fn func(dst, p []byte) []byte) *AccumulativeBuf {
	b.buf.WriteApplyFn(p, fn)
	return b
}

// WriteApplyFnStr applies fn to s and write result to the buffer.
// DEPRECATED: use WriteApplyFnString() instead.
func (b *AccumulativeBuf) WriteApplyFnStr(s string, fn func(dst, p []byte) []byte) *AccumulativeBuf {
	return b.WriteApplyFnString(s, fn)
}

// WriteApplyFnString applies fn to s and write result to the buffer.
func (b *AccumulativeBuf) WriteApplyFnString(s string, fn func(dst, p []byte) []byte) *AccumulativeBuf {
	b.buf.WriteApplyFnString(s, fn)
	return b
}

// Replace replaces old bytes to new in buffer.
func (b *AccumulativeBuf) Replace(old, new []byte, n int) *AccumulativeBuf {
	b.buf.Replace(old, new, n)
	return b
}

// ReplaceStr replace old to new substrings in buffer.
// DEPRECATED: use ReplaceString() instead.
func (b *AccumulativeBuf) ReplaceStr(old, new string, n int) *AccumulativeBuf {
	return b.ReplaceString(old, new, n)
}

// ReplaceString replace old to new substrings in buffer.
func (b *AccumulativeBuf) ReplaceString(old, new string, n int) *AccumulativeBuf {
	b.buf.ReplaceString(old, new, n)
	return b
}

// ReplaceAll replace all occurrences of old bytes to new in buffer.
func (b *AccumulativeBuf) ReplaceAll(old, new []byte) *AccumulativeBuf {
	b.buf.ReplaceAll(old, new)
	return b
}

// ReplaceStrAll replaces all occurrences of old substrings to new in buffer.
// DEPRECATED: use ReplaceStringAll() instead.
func (b *AccumulativeBuf) ReplaceStrAll(old, new string) *AccumulativeBuf {
	return b.ReplaceStringAll(old, new)
}

// ReplaceStringAll replaces all occurrences of old substrings to new in buffer.
func (b *AccumulativeBuf) ReplaceStringAll(old, new string) *AccumulativeBuf {
	b.buf.ReplaceStringAll(old, new)
	return b
}

func (b *AccumulativeBuf) Len() int {
	return b.buf.Len()
}

func (b *AccumulativeBuf) Cap() int {
	return b.buf.Cap()
}

// Grow increases length of the buffer.
func (b *AccumulativeBuf) Grow(newLen int) *AccumulativeBuf {
	b.buf.Grow(newLen)
	return b
}

// GrowDelta increases length of the buffer to actual length + delta.
//
// See Grow().
func (b *AccumulativeBuf) GrowDelta(delta int) *AccumulativeBuf {
	b.buf.GrowDelta(delta)
	return b
}

func (b *AccumulativeBuf) Reset() *AccumulativeBuf {
	b.buf.Reset()
	return b
}
