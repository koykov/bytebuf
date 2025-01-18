package bytebuf

import (
	"strconv"
	"time"
	"unsafe"

	"github.com/koykov/bytealg"
	"github.com/koykov/byteconv"
	"github.com/koykov/clock"
	"github.com/koykov/x2bytes"
)

// Chain is a primitive byte buffer with chain call support.
type Chain []byte

// NewChain creates and initializes a new chain buffer instance using buf as its initial contents.
func NewChain(buf []byte) *Chain {
	cb := Chain(buf)
	return &cb
}

// NewChainSize creates new chain buffer and initializes byte slice as an internal buffer.
func NewChainSize(size int) *Chain {
	if size < 0 {
		return nil
	}
	buf := make([]byte, 0, size)
	return NewChain(buf)
}

// Bytes returns contents of the buffer.
func (b *Chain) Bytes() []byte {
	return *b
}

// BytesCopy returns copy of the buffer.
func (b *Chain) BytesCopy() []byte {
	return append([]byte(nil), *b...)
}

// Get contents of the buffer as string.
func (b *Chain) String() string {
	return byteconv.B2S(*b)
}

// StringCopy returns copy of the buffer contents as string.
func (b *Chain) StringCopy() string {
	return byteconv.B2S(b.BytesCopy())
}

// Write bytes to the buffer.
func (b *Chain) Write(p []byte) *Chain {
	*b = append(*b, p...)
	return b
}

// WriteByte writes single byte.
func (b *Chain) WriteByte(p byte) *Chain {
	*b = append(*b, p)
	return b
}

// WriteString writes string to the buffer.
func (b *Chain) WriteString(s string) *Chain {
	*b = append(*b, s...)
	return b
}

// WriteInt writes integer value to the buffer.
func (b *Chain) WriteInt(i int64) *Chain {
	*b, _ = x2bytes.IntToBytes(*b, i)
	return b
}

// WriteIntBase writes integer value in given base to the buffer.
func (b *Chain) WriteIntBase(i int64, base int) *Chain {
	*b = strconv.AppendInt(*b, i, base)
	return b
}

// WriteUint writes unsigned integer value to the buffer.
func (b *Chain) WriteUint(u uint64) *Chain {
	*b, _ = x2bytes.UintToBytes(*b, u)
	return b
}

// WriteUintBase writes unsigned integer value in given base to the buffer.
func (b *Chain) WriteUintBase(u uint64, base int) *Chain {
	*b = strconv.AppendUint(*b, u, base)
	return b
}

// WriteFloat writes float value to the buffer.
func (b *Chain) WriteFloat(f float64) *Chain {
	*b, _ = x2bytes.FloatToBytes(*b, f)
	return b
}

// WriteBool writes boolean value to the buffer.
func (b *Chain) WriteBool(v bool) *Chain {
	*b, _ = x2bytes.BoolToBytes(*b, v)
	return b
}

// WriteX write x with arbitrary type to the buffer.
func (b *Chain) WriteX(x any) *Chain {
	*b, _ = x2bytes.ToBytes(*b, x)
	return b
}

// WriteMarshallerTo marshal data of struct implemented MarshallerTo interface and write it to the buffer.
func (b *Chain) WriteMarshallerTo(m MarshallerTo) *Chain {
	if m == nil {
		return b
	}
	n := b.Len()
	d := m.Size()
	b.GrowDelta(d)
	_, _ = m.MarshalTo(b.Bytes()[n:])
	return b
}

// WriteApplyFn applies fn to p and write result to the buffer.
func (b *Chain) WriteApplyFn(p []byte, fn func(dst, p []byte) []byte) *Chain {
	*b = fn(*b, p)
	return b
}

// WriteApplyFnN applies fn to p N times and write result to the buffer.
func (b *Chain) WriteApplyFnN(p []byte, fn func(dst, p []byte) []byte, n int) *Chain {
	off := b.Len()
	var poff int
	for i := 0; i < n; i++ {
		poff = b.Len()
		*b = fn(*b, p)
		p = (*b)[poff:]
	}
	*b = append((*b)[:off], (*b)[poff:]...)
	return b
}

// WriteApplyFnString applies fn to s and write result to the buffer.
func (b *Chain) WriteApplyFnString(s string, fn func(dst, p []byte) []byte) *Chain {
	*b = fn(*b, byteconv.S2B(s))
	return b
}

// WriteApplyFnNString applies fn to s N times and write result to the buffer.
func (b *Chain) WriteApplyFnNString(s string, fn func(dst, p []byte) []byte, n int) *Chain {
	return b.WriteApplyFnN(byteconv.S2B(s), fn, n)
}

// WriteTime writes time t in given format to the buffer.
func (b *Chain) WriteTime(format string, t time.Time) *Chain {
	*b, _ = clock.AppendFormat(*b, format, t)
	return b
}

// Replace replaces old bytes to new in buffer.
func (b *Chain) Replace(old, new []byte, n int) *Chain {
	if b.Len() == 0 || n == 0 {
		return b
	}
	var i, at, c int
	// Use the same byte buffer to make replacement and avoid alloc.
	dst := (*b)[b.Len():]
	for {
		if i = bytealg.IndexAtBytes(*b, old, at); i < 0 || c == n {
			dst = append(dst, (*b)[at:]...)
			break
		}
		dst = append(dst, (*b)[at:i]...)
		dst = append(dst, new...)
		at = i + len(old)
		c++
	}
	// Move result to the beginning of buffer.
	b.Reset().Write(dst)
	return b
}

// ReplaceString replace old to new substrings in buffer.
func (b *Chain) ReplaceString(old, new string, n int) *Chain {
	return b.Replace(byteconv.S2B(old), byteconv.S2B(new), n)
}

// ReplaceAll replace all occurrences of old bytes to new in buffer.
func (b *Chain) ReplaceAll(old, new []byte) *Chain {
	return b.Replace(old, new, -1)
}

// ReplaceStringAll replaces all occurrences of old substrings to new in buffer.
func (b *Chain) ReplaceStringAll(old, new string) *Chain {
	return b.Replace(byteconv.S2B(old), byteconv.S2B(new), -1)
}

func (b *Chain) Len() int {
	return len(*b)
}

func (b *Chain) Cap() int {
	return cap(*b)
}

// Grow increases length of the buffer.
func (b *Chain) Grow(newLen int) *Chain {
	if newLen <= 0 {
		return b
	}
	// Get buffer header.
	h := *(*byteconv.SliceHeader)(unsafe.Pointer(b))
	if newLen < h.Cap {
		// Just increase header's length if capacity allows
		h.Len = newLen
		// ... and restore the buffer from the header.
		*b = *(*[]byte)(unsafe.Pointer(&h))
	} else {
		// Append necessary space.
		*b = append(*b, make([]byte, newLen-b.Len())...)
	}
	return b
}

// GrowDelta increases length of the buffer to actual length + delta.
//
// See Grow().
func (b *Chain) GrowDelta(delta int) *Chain {
	return b.Grow(b.Len() + delta)
}

func (b *Chain) Reset() *Chain {
	*b = (*b)[:0]
	return b
}

// Reduce decreases buffer length to delta.
func (b *Chain) Reduce(delta int) *Chain {
	if delta <= 0 {
		return b
	}
	if len(*b) < delta {
		return b.Reset()
	}
	*b = (*b)[:len(*b)-delta]
	return b
}
