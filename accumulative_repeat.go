package bytebuf

func (b *Accumulative) WriteN(p []byte, n int) *Accumulative {
	b.buf.WriteN(p, n)
	return b
}

func (b *Accumulative) WriteByteN(p byte, n int) *Accumulative {
	b.buf.WriteByteN(p, n)
	return b
}

func (b *Accumulative) WriteStringN(s string, n int) *Accumulative {
	b.buf.WriteStringN(s, n)
	return b
}

func (b *Accumulative) WriteIntN(i int64, n int) *Accumulative {
	b.buf.WriteIntN(i, n)
	return b
}

func (b *Accumulative) WriteIntBaseN(i int64, base, n int) *Accumulative {
	b.buf.WriteIntBaseN(i, base, n)
	return b
}

func (b *Accumulative) WriteUintN(u uint64, n int) *Accumulative {
	b.buf.WriteUintN(u, n)
	return b
}

func (b *Accumulative) WriteUintBaseN(u uint64, base, n int) *Accumulative {
	b.buf.WriteUintBaseN(u, base, n)
	return b
}

func (b *Accumulative) WriteFloatN(f float64, n int) *Accumulative {
	b.buf.WriteFloatN(f, n)
	return b
}

func (b *Accumulative) WriteBoolN(v bool, n int) *Accumulative {
	b.buf.WriteBoolN(v, n)
	return b
}

func (b *Accumulative) WriteXN(x any, n int) *Accumulative {
	b.buf.WriteXN(x, n)
	return b
}
