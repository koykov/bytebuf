package bytebuf

func (b *Accumulative) WriteRepeat(p []byte, count int) *Accumulative {
	b.buf.WriteRepeat(p, count)
	return b
}

func (b *Accumulative) WriteByteRepeat(p byte, count int) *Accumulative {
	b.buf.WriteByteRepeat(p, count)
	return b
}

func (b *Accumulative) WriteStringRepeat(s string, count int) *Accumulative {
	b.buf.WriteStringRepeat(s, count)
	return b
}

func (b *Accumulative) WriteIntRepeat(i int64, count int) *Accumulative {
	b.buf.WriteIntRepeat(i, count)
	return b
}

func (b *Accumulative) WriteIntBaseRepeat(i int64, base, count int) *Accumulative {
	b.buf.WriteIntBaseRepeat(i, base, count)
	return b
}

func (b *Accumulative) WriteUintRepeat(u uint64, count int) *Accumulative {
	b.buf.WriteUintRepeat(u, count)
	return b
}

func (b *Accumulative) WriteUintBaseRepeat(u uint64, base, count int) *Accumulative {
	b.buf.WriteUintBaseRepeat(u, base, count)
	return b
}

func (b *Accumulative) WriteFloatRepeat(f float64, count int) *Accumulative {
	b.buf.WriteFloatRepeat(f, count)
	return b
}

func (b *Accumulative) WriteBoolRepeat(v bool, count int) *Accumulative {
	b.buf.WriteBoolRepeat(v, count)
	return b
}

func (b *Accumulative) WriteXRepeat(x any, count int) *Accumulative {
	b.buf.WriteXRepeat(x, count)
	return b
}
