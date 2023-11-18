package bytebuf

import (
	"github.com/koykov/x2bytes"
)

// Contains wrapper methods over inner Chain.

// Bytes returns contents of the buffer.
func (b *Accumulative) Bytes() []byte {
	return b.buf
}

// BytesCopy returns copy of buffer contents.
func (b *Accumulative) BytesCopy() []byte {
	return b.buf.BytesCopy()
}

// Get contents of the buffer as string.
func (b *Accumulative) String() string {
	return b.buf.String()
}

// StringCopy returns copy of the buffer contents as string.
func (b *Accumulative) StringCopy() string {
	return b.buf.StringCopy()
}

// Write bytes to the buffer.
func (b *Accumulative) Write(p []byte) *Accumulative {
	b.buf.Write(p)
	return b
}

// WriteByte writes single byte.
func (b *Accumulative) WriteByte(p byte) *Accumulative {
	b.buf.WriteByte(p)
	return b
}

// WriteStr writes string to the buffer.
// DEPRECATED: use WriteString() instead.
func (b *Accumulative) WriteStr(s string) *Accumulative {
	return b.WriteString(s)
}

// WriteString writes string to the buffer.
func (b *Accumulative) WriteString(s string) *Accumulative {
	b.buf.WriteString(s)
	return b
}

// WriteInt writes integer value to the buffer.
func (b *Accumulative) WriteInt(i int64) *Accumulative {
	b.buf, b.err = x2bytes.IntToBytes(b.buf, i)
	return b
}

// WriteUint writes unsigned integer value to the buffer.
func (b *Accumulative) WriteUint(u uint64) *Accumulative {
	b.buf, b.err = x2bytes.UintToBytes(b.buf, u)
	return b
}

// WriteFloat writes float value to the buffer.
func (b *Accumulative) WriteFloat(f float64) *Accumulative {
	b.buf, b.err = x2bytes.FloatToBytes(b.buf, f)
	return b
}

// WriteBool writes boolean value to the buffer.
func (b *Accumulative) WriteBool(v bool) *Accumulative {
	b.buf, b.err = x2bytes.BoolToBytes(b.buf, v)
	return b
}

// WriteX write v with arbitrary type to the buffer.
func (b *Accumulative) WriteX(x interface{}) *Accumulative {
	b.buf, b.err = x2bytes.ToBytes(b.buf, x)
	return b
}

// WriteMarshallerTo marshal data of struct implemented MarshallerTo interface and write it to the buffer.
func (b *Accumulative) WriteMarshallerTo(m MarshallerTo) *Accumulative {
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
func (b *Accumulative) WriteApplyFn(p []byte, fn func(dst, p []byte) []byte) *Accumulative {
	b.buf.WriteApplyFn(p, fn)
	return b
}

// WriteApplyFnStr applies fn to s and write result to the buffer.
// DEPRECATED: use WriteApplyFnString() instead.
func (b *Accumulative) WriteApplyFnStr(s string, fn func(dst, p []byte) []byte) *Accumulative {
	return b.WriteApplyFnString(s, fn)
}

// WriteApplyFnString applies fn to s and write result to the buffer.
func (b *Accumulative) WriteApplyFnString(s string, fn func(dst, p []byte) []byte) *Accumulative {
	b.buf.WriteApplyFnString(s, fn)
	return b
}

// Replace replaces old bytes to new in buffer.
func (b *Accumulative) Replace(old, new []byte, n int) *Accumulative {
	b.buf.Replace(old, new, n)
	return b
}

// ReplaceStr replace old to new substrings in buffer.
// DEPRECATED: use ReplaceString() instead.
func (b *Accumulative) ReplaceStr(old, new string, n int) *Accumulative {
	return b.ReplaceString(old, new, n)
}

// ReplaceString replace old to new substrings in buffer.
func (b *Accumulative) ReplaceString(old, new string, n int) *Accumulative {
	b.buf.ReplaceString(old, new, n)
	return b
}

// ReplaceAll replace all occurrences of old bytes to new in buffer.
func (b *Accumulative) ReplaceAll(old, new []byte) *Accumulative {
	b.buf.ReplaceAll(old, new)
	return b
}

// ReplaceStrAll replaces all occurrences of old substrings to new in buffer.
// DEPRECATED: use ReplaceStringAll() instead.
func (b *Accumulative) ReplaceStrAll(old, new string) *Accumulative {
	return b.ReplaceStringAll(old, new)
}

// ReplaceStringAll replaces all occurrences of old substrings to new in buffer.
func (b *Accumulative) ReplaceStringAll(old, new string) *Accumulative {
	b.buf.ReplaceStringAll(old, new)
	return b
}

func (b *Accumulative) Len() int {
	return b.buf.Len()
}

func (b *Accumulative) Cap() int {
	return b.buf.Cap()
}

// Grow increases length of the buffer.
func (b *Accumulative) Grow(newLen int) *Accumulative {
	b.buf.Grow(newLen)
	return b
}

// GrowDelta increases length of the buffer to actual length + delta.
//
// See Grow().
func (b *Accumulative) GrowDelta(delta int) *Accumulative {
	b.buf.GrowDelta(delta)
	return b
}

func (b *Accumulative) Reset() *Accumulative {
	b.buf.Reset()
	return b
}
