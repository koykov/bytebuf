package bytebuf

import "github.com/koykov/x2bytes"

// Contains wrapper methods over inner ChainBuf.

// Get contents of the buffer.
func (b *AccumulativeBuf) Bytes() []byte {
	return b.buf
}

// Get copy of the buffer.
func (b *AccumulativeBuf) BytesCopy() []byte {
	return b.buf.BytesCopy()
}

// Get contents of the buffer as string.
func (b *AccumulativeBuf) String() string {
	return b.buf.String()
}

// Get copy of the buffer as string.
func (b *AccumulativeBuf) StringCopy() string {
	return b.buf.StringCopy()
}

// Write bytes to the buffer.
func (b *AccumulativeBuf) Write(p []byte) *AccumulativeBuf {
	b.buf.Write(p)
	return b
}

// Write single byte.
func (b *AccumulativeBuf) WriteByte(p byte) *AccumulativeBuf {
	b.buf.WriteByte(p)
	return b
}

// Write string to the buffer.
func (b *AccumulativeBuf) WriteStr(s string) *AccumulativeBuf {
	b.buf.WriteStr(s)
	return b
}

// Write integer value to the buffer.
func (b *AccumulativeBuf) WriteInt(i int64) *AccumulativeBuf {
	b.buf, b.err = x2bytes.IntToBytes(b.buf, i)
	return b
}

// Write unsigned integer value to the buffer.
func (b *AccumulativeBuf) WriteUint(u uint64) *AccumulativeBuf {
	b.buf, b.err = x2bytes.UintToBytes(b.buf, u)
	return b
}

// Write float value to the buffer.
func (b *AccumulativeBuf) WriteFloat(f float64) *AccumulativeBuf {
	b.buf, b.err = x2bytes.FloatToBytes(b.buf, f)
	return b
}

// Write boolean value to the buffer.
func (b *AccumulativeBuf) WriteBool(v bool) *AccumulativeBuf {
	b.buf, b.err = x2bytes.BoolToBytes(b.buf, v)
	return b
}

// Write v with arbitrary type to the buffer.
func (b *AccumulativeBuf) WriteX(x interface{}) *AccumulativeBuf {
	b.buf, b.err = x2bytes.ToBytes(b.buf, x)
	return b
}

// Marshal data of struct implemented MarshallerTo interface.
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

// Replace old to new bytes in buffer.
func (b *AccumulativeBuf) Replace(old, new []byte, n int) *AccumulativeBuf {
	b.buf.Replace(old, new, n)
	return b
}

// Replace old to new strings in buffer.
func (b *AccumulativeBuf) ReplaceStr(old, new string, n int) *AccumulativeBuf {
	b.buf.ReplaceStr(old, new, n)
	return b
}

// Replace all old to new bytes in buffer.
func (b *AccumulativeBuf) ReplaceAll(old, new []byte) *AccumulativeBuf {
	b.buf.ReplaceAll(old, new)
	return b
}

// Replace all old to new strings in buffer.
func (b *AccumulativeBuf) ReplaceStrAll(old, new string) *AccumulativeBuf {
	b.buf.ReplaceStrAll(old, new)
	return b
}

// Get length of the buffer.
func (b *AccumulativeBuf) Len() int {
	return b.buf.Len()
}

// Get capacity of the buffer.
func (b *AccumulativeBuf) Cap() int {
	return b.buf.Cap()
}

// Grow length of the buffer.
func (b *AccumulativeBuf) Grow(newLen int) *AccumulativeBuf {
	b.buf.Grow(newLen)
	return b
}

// Grow length of the buffer to actual length + delta.
//
// See Grow().
func (b *AccumulativeBuf) GrowDelta(delta int) *AccumulativeBuf {
	b.buf.GrowDelta(delta)
	return b
}

// Reset length of the buffer.
func (b *AccumulativeBuf) Reset() *AccumulativeBuf {
	b.buf.Reset()
	return b
}
