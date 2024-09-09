package bytebuf

// WriteN writes p to buffer N times.
// See Write.
func (b *Chain) WriteN(p []byte, n int) *Chain {
	for i := 0; i < n; i++ {
		b.Write(p)
	}
	return b
}

// WriteByteN writes p to buffer N times.
// See WriteByte.
func (b *Chain) WriteByteN(p byte, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteByte(p)
	}
	return b
}

// WriteStringN writes s to buffer N times.
// See WriteString.
func (b *Chain) WriteStringN(s string, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteString(s)
	}
	return b
}

// WriteIntN writes i to buffer N times.
// See WriteInt.
func (b *Chain) WriteIntN(i int64, n int) *Chain {
	for j := 0; j < n; j++ {
		b.WriteInt(i)
	}
	return b
}

// WriteIntBaseN writes i to buffer N times.
// See WriteIntBase.
func (b *Chain) WriteIntBaseN(i int64, base, n int) *Chain {
	for j := 0; j < n; j++ {
		b.WriteIntBase(i, base)
	}
	return b
}

// WriteUintN writes u to buffer N times.
// See WriteUint.
func (b *Chain) WriteUintN(u uint64, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteUint(u)
	}
	return b
}

// WriteUintBaseN writes u to buffer N times.
// See WriteUintBase.
func (b *Chain) WriteUintBaseN(u uint64, base, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteUintBase(u, base)
	}
	return b
}

// WriteFloatN writes f to buffer N times.
// See WriteFloat.
func (b *Chain) WriteFloatN(f float64, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteFloat(f)
	}
	return b
}

// WriteBoolN writes v to buffer N times.
// See WriteBool.
func (b *Chain) WriteBoolN(v bool, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteBool(v)
	}
	return b
}

// WriteXN writes x to buffer N times.
// See WriteX.
func (b *Chain) WriteXN(x any, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteX(x)
	}
	return b
}
