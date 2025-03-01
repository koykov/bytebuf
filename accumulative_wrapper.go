package bytebuf

import (
	"strconv"
	"time"

	"github.com/koykov/clock"
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

// WriteRune writes rune to the buffer.
func (b *Accumulative) WriteRune(r rune) *Accumulative {
	b.buf.WriteRune(r)
	return b
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

// WriteIntBase writes integer value in given base to the buffer.
func (b *Accumulative) WriteIntBase(i int64, base int) *Accumulative {
	b.buf = strconv.AppendInt(b.buf, i, base)
	return b
}

// WriteUint writes unsigned integer value to the buffer.
func (b *Accumulative) WriteUint(u uint64) *Accumulative {
	b.buf, b.err = x2bytes.UintToBytes(b.buf, u)
	return b
}

// WriteUintBase writes unsigned integer value in given base to the buffer.
func (b *Accumulative) WriteUintBase(u uint64, base int) *Accumulative {
	b.buf = strconv.AppendUint(b.buf, u, base)
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

// WriteFormat writes formatted string to the buffer.
func (b *Accumulative) WriteFormat(format string, args ...any) *Accumulative {
	b.buf.WriteFormat(format, args...)
	return b
}

// WriteX write v with arbitrary type to the buffer.
func (b *Accumulative) WriteX(x any) *Accumulative {
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

// WriteApplyFnN applies fn to p N times and write result to the buffer.
func (b *Accumulative) WriteApplyFnN(p []byte, fn func(dst, p []byte) []byte, n int) *Accumulative {
	b.buf.WriteApplyFnN(p, fn, n)
	return b
}

// WriteApplyFnString applies fn to s and write result to the buffer.
func (b *Accumulative) WriteApplyFnString(s string, fn func(dst, p []byte) []byte) *Accumulative {
	b.buf.WriteApplyFnString(s, fn)
	return b
}

// WriteApplyFnNString applies fn to s and write result to the buffer.
func (b *Accumulative) WriteApplyFnNString(s string, fn func(dst, p []byte) []byte, n int) *Accumulative {
	b.buf.WriteApplyFnNString(s, fn, n)
	return b
}

// WriteTime writes time t in given format to the buffer.
func (b *Accumulative) WriteTime(format string, t time.Time) *Accumulative {
	b.buf, b.err = clock.AppendFormat(b.buf, format, t)
	return b
}

// Replace replaces old bytes to new in buffer.
func (b *Accumulative) Replace(old, new []byte, n int) *Accumulative {
	b.buf.Replace(old, new, n)
	return b
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

func (b *Accumulative) Reduce(delta int) *Accumulative {
	b.buf.Reduce(delta)
	return b
}
