package bytebuf

import (
	"reflect"
	"unsafe"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
	"github.com/koykov/x2bytes"
)

// ChainBuf is a primitive byte buffer with chain call support.
type ChainBuf []byte

// Bytes returns contents of the buffer.
func (b *ChainBuf) Bytes() []byte {
	return *b
}

// BytesCopy returns copy of the buffer.
func (b *ChainBuf) BytesCopy() []byte {
	return bytealg.Copy(*b)
}

// Get contents of the buffer as string.
func (b *ChainBuf) String() string {
	return fastconv.B2S(*b)
}

// StringCopy returns copy of the buffer contents as string.
func (b *ChainBuf) StringCopy() string {
	return bytealg.Copy[string](fastconv.B2S(*b))
}

// Write bytes to the buffer.
func (b *ChainBuf) Write(p []byte) *ChainBuf {
	*b = append(*b, p...)
	return b
}

// WriteByte writes single byte.
func (b *ChainBuf) WriteByte(p byte) *ChainBuf {
	*b = append(*b, p)
	return b
}

// WriteStr writes string to the buffer.
// DEPRECATED: use WriteString() instead.
func (b *ChainBuf) WriteStr(s string) *ChainBuf {
	*b = append(*b, s...)
	return b
}

// WriteString writes string to the buffer.
func (b *ChainBuf) WriteString(s string) *ChainBuf {
	*b = append(*b, s...)
	return b
}

// WriteInt writes integer value to the buffer.
func (b *ChainBuf) WriteInt(i int64) *ChainBuf {
	*b, _ = x2bytes.IntToBytes(*b, i)
	return b
}

// WriteUint writes unsigned integer value to the buffer.
func (b *ChainBuf) WriteUint(u uint64) *ChainBuf {
	*b, _ = x2bytes.UintToBytes(*b, u)
	return b
}

// WriteFloat writes float value to the buffer.
func (b *ChainBuf) WriteFloat(f float64) *ChainBuf {
	*b, _ = x2bytes.FloatToBytes(*b, f)
	return b
}

// WriteBool writes boolean value to the buffer.
func (b *ChainBuf) WriteBool(v bool) *ChainBuf {
	*b, _ = x2bytes.BoolToBytes(*b, v)
	return b
}

// WriteX write x with arbitrary type to the buffer.
func (b *ChainBuf) WriteX(x interface{}) *ChainBuf {
	*b, _ = x2bytes.ToBytes(*b, x)
	return b
}

// WriteMarshallerTo marshal data of struct implemented MarshallerTo interface and write it to the buffer.
func (b *ChainBuf) WriteMarshallerTo(m MarshallerTo) *ChainBuf {
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
func (b *ChainBuf) WriteApplyFn(p []byte, fn func(dst, p []byte) []byte) *ChainBuf {
	*b = fn(*b, p)
	return b
}

// WriteApplyFnStr applies fn to s and write result to the buffer.
// DEPRECATED: use WriteApplyFnString() instead.
func (b *ChainBuf) WriteApplyFnStr(s string, fn func(dst, p []byte) []byte) *ChainBuf {
	*b = fn(*b, fastconv.S2B(s))
	return b
}

// WriteApplyFnString applies fn to s and write result to the buffer.
func (b *ChainBuf) WriteApplyFnString(s string, fn func(dst, p []byte) []byte) *ChainBuf {
	*b = fn(*b, fastconv.S2B(s))
	return b
}

// Replace replaces old bytes to new in buffer.
func (b *ChainBuf) Replace(old, new []byte, n int) *ChainBuf {
	if b.Len() == 0 || n == 0 {
		return b
	}
	var i, at, c int
	// Use the same byte buffer to make replacement and avoid alloc.
	dst := (*b)[b.Len():]
	for {
		if i = bytealg.IndexAt(*b, old, at); i < 0 || c == n {
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

// ReplaceStr replace old to new substrings in buffer.
// DEPRECATED: use ReplaceString() instead.
func (b *ChainBuf) ReplaceStr(old, new string, n int) *ChainBuf {
	return b.ReplaceString(old, new, n)
}

// ReplaceString replace old to new substrings in buffer.
func (b *ChainBuf) ReplaceString(old, new string, n int) *ChainBuf {
	return b.Replace(fastconv.S2B(old), fastconv.S2B(new), n)
}

// ReplaceAll replace all occurrences of old bytes to new in buffer.
func (b *ChainBuf) ReplaceAll(old, new []byte) *ChainBuf {
	return b.Replace(old, new, -1)
}

// ReplaceStrAll replaces all occurrences of old substrings to new in buffer.
// DEPRECATED: use ReplaceStringAll() instead.
func (b *ChainBuf) ReplaceStrAll(old, new string) *ChainBuf {
	return b.ReplaceStringAll(old, new)
}

// ReplaceStringAll replaces all occurrences of old substrings to new in buffer.
func (b *ChainBuf) ReplaceStringAll(old, new string) *ChainBuf {
	return b.Replace(fastconv.S2B(old), fastconv.S2B(new), -1)
}

func (b *ChainBuf) Len() int {
	return len(*b)
}

func (b *ChainBuf) Cap() int {
	return cap(*b)
}

// Grow increases length of the buffer.
func (b *ChainBuf) Grow(newLen int) *ChainBuf {
	if newLen <= 0 {
		return b
	}
	// Get buffer header.
	h := *(*reflect.SliceHeader)(unsafe.Pointer(b))
	if newLen < h.Cap {
		// Just increase header's length if capacity allows
		h.Len = newLen
		// .. and restore the buffer from the header.
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
func (b *ChainBuf) GrowDelta(delta int) *ChainBuf {
	return b.Grow(b.Len() + delta)
}

func (b *ChainBuf) Reset() *ChainBuf {
	*b = (*b)[:0]
	return b
}
