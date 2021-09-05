package bytebuf

// Contains wrapper methods over inner ChainBuf.

// Get contents of the buffer.
func (b *AccumulativeBuffer) Bytes() []byte {
	return b.buf
}

// Get copy of the buffer.
func (b *AccumulativeBuffer) BytesCopy() []byte {
	return b.buf.BytesCopy()
}

// Get contents of the buffer as string.
func (b *AccumulativeBuffer) String() string {
	return b.buf.String()
}

// Get copy of the buffer as string.
func (b *AccumulativeBuffer) StringCopy() string {
	return b.buf.StringCopy()
}

// Write bytes to the buffer.
func (b *AccumulativeBuffer) Write(p []byte) *AccumulativeBuffer {
	b.buf.Write(p)
	return b
}

// Write single byte.
func (b *AccumulativeBuffer) WriteByte(p byte) *AccumulativeBuffer {
	b.buf.WriteByte(p)
	return b
}

// Write string to the buffer.
func (b *AccumulativeBuffer) WriteStr(s string) *AccumulativeBuffer {
	b.buf.WriteStr(s)
	return b
}

// Write integer value to the buffer.
func (b *AccumulativeBuffer) WriteInt(i int64) *AccumulativeBuffer {
	b.buf.WriteInt(i)
	return b
}

// Write unsigned integer value to the buffer.
func (b *AccumulativeBuffer) WriteUint(u uint64) *AccumulativeBuffer {
	b.buf.WriteUint(u)
	return b
}

// Write float value to the buffer.
func (b *AccumulativeBuffer) WriteFloat(f float64) *AccumulativeBuffer {
	b.buf.WriteFloat(f)
	return b
}

// Write boolean value to the buffer.
func (b *AccumulativeBuffer) WriteBool(v bool) *AccumulativeBuffer {
	b.buf.WriteBool(v)
	return b
}

func (b *AccumulativeBuffer) WriteX(v interface{}) *AccumulativeBuffer {
	b.buf.WriteX(v)
	return b
}

// Replace old to new bytes in buffer.
func (b *AccumulativeBuffer) Replace(old, new []byte, n int) *AccumulativeBuffer {
	b.buf.Replace(old, new, n)
	return b
}

// Replace old to new strings in buffer.
func (b *AccumulativeBuffer) ReplaceStr(old, new string, n int) *AccumulativeBuffer {
	b.buf.ReplaceStr(old, new, n)
	return b
}

// Replace all old to new bytes in buffer.
func (b *AccumulativeBuffer) ReplaceAll(old, new []byte) *AccumulativeBuffer {
	b.buf.ReplaceAll(old, new)
	return b
}

// Replace all old to new strings in buffer.
func (b *AccumulativeBuffer) ReplaceStrAll(old, new string) *AccumulativeBuffer {
	b.buf.ReplaceStrAll(old, new)
	return b
}

// Get length of the buffer.
func (b *AccumulativeBuffer) Len() int {
	return b.buf.Len()
}

// Get capacity of the buffer.
func (b *AccumulativeBuffer) Cap() int {
	return b.buf.Cap()
}

// Grow length of the buffer.
func (b *AccumulativeBuffer) Grow(newLen int) *AccumulativeBuffer {
	b.buf.Grow(newLen)
	return b
}

// Grow length of the buffer to actual length + delta.
//
// See Grow().
func (b *AccumulativeBuffer) GrowDelta(delta int) *AccumulativeBuffer {
	b.buf.GrowDelta(delta)
	return b
}

// Reset length of the buffer.
func (b *AccumulativeBuffer) Reset() *AccumulativeBuffer {
	b.buf.Reset()
	return b
}

