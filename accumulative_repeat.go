package bytebuf

// WriteN writes p to buffer N times.
// See Write.
func (b *Accumulative) WriteN(p []byte, n int) *Accumulative {
	b.buf.WriteN(p, n)
	return b
}

// WriteByteN writes p to buffer N times.
// See WriteByte.
func (b *Accumulative) WriteByteN(p byte, n int) *Accumulative {
	b.buf.WriteByteN(p, n)
	return b
}

// WriteStringN writes s to buffer N times.
// See WriteString.
func (b *Accumulative) WriteStringN(s string, n int) *Accumulative {
	b.buf.WriteStringN(s, n)
	return b
}

// WriteIntN writes i to buffer N times.
// See WriteInt.
func (b *Accumulative) WriteIntN(i int64, n int) *Accumulative {
	b.buf.WriteIntN(i, n)
	return b
}

// WriteIntBaseN writes i to buffer N times.
// See WriteIntBase.
func (b *Accumulative) WriteIntBaseN(i int64, base, n int) *Accumulative {
	b.buf.WriteIntBaseN(i, base, n)
	return b
}

// WriteUintN writes u to buffer N times.
// See WriteUint.
func (b *Accumulative) WriteUintN(u uint64, n int) *Accumulative {
	b.buf.WriteUintN(u, n)
	return b
}

// WriteUintBaseN writes u to buffer N times.
// See WriteUintBase.
func (b *Accumulative) WriteUintBaseN(u uint64, base, n int) *Accumulative {
	b.buf.WriteUintBaseN(u, base, n)
	return b
}

// WriteFloatN writes f to buffer N times.
// See WriteFloat.
func (b *Accumulative) WriteFloatN(f float64, n int) *Accumulative {
	b.buf.WriteFloatN(f, n)
	return b
}

// WriteBoolN writes v to buffer N times.
// See WriteBool.
func (b *Accumulative) WriteBoolN(v bool, n int) *Accumulative {
	b.buf.WriteBoolN(v, n)
	return b
}

// WriteXN writes x to buffer N times.
// See WriteX.
func (b *Accumulative) WriteXN(x any, n int) *Accumulative {
	b.buf.WriteXN(x, n)
	return b
}
