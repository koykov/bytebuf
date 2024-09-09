package bytebuf

func (b *Chain) WriteN(p []byte, n int) *Chain {
	for i := 0; i < n; i++ {
		b.Write(p)
	}
	return b
}

func (b *Chain) WriteByteN(p byte, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteByte(p)
	}
	return b
}

func (b *Chain) WriteStringN(s string, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteString(s)
	}
	return b
}

func (b *Chain) WriteIntN(i int64, n int) *Chain {
	for j := 0; j < n; j++ {
		b.WriteInt(i)
	}
	return b
}

func (b *Chain) WriteIntBaseN(i int64, base, n int) *Chain {
	for j := 0; j < n; j++ {
		b.WriteIntBase(i, base)
	}
	return b
}

func (b *Chain) WriteUintN(u uint64, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteUint(u)
	}
	return b
}

func (b *Chain) WriteUintBaseN(u uint64, base, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteUintBase(u, base)
	}
	return b
}

func (b *Chain) WriteFloatN(f float64, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteFloat(f)
	}
	return b
}

func (b *Chain) WriteBoolN(v bool, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteBool(v)
	}
	return b
}

func (b *Chain) WriteXN(x any, n int) *Chain {
	for i := 0; i < n; i++ {
		b.WriteX(x)
	}
	return b
}
