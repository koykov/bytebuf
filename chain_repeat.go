package bytebuf

func (b *Chain) WriteRepeat(p []byte, count int) *Chain {
	for i := 0; i < count; i++ {
		b.Write(p)
	}
	return b
}

func (b *Chain) WriteByteRepeat(p byte, count int) *Chain {
	for i := 0; i < count; i++ {
		b.WriteByte(p)
	}
	return b
}

func (b *Chain) WriteStringRepeat(s string, count int) *Chain {
	for i := 0; i < count; i++ {
		b.WriteString(s)
	}
	return b
}

func (b *Chain) WriteIntRepeat(i int64, count int) *Chain {
	for j := 0; j < count; j++ {
		b.WriteInt(i)
	}
	return b
}

func (b *Chain) WriteIntBaseRepeat(i int64, base, count int) *Chain {
	for j := 0; j < count; j++ {
		b.WriteIntBase(i, base)
	}
	return b
}

func (b *Chain) WriteUintRepeat(u uint64, count int) *Chain {
	for i := 0; i < count; i++ {
		b.WriteUint(u)
	}
	return b
}

func (b *Chain) WriteUintBaseRepeat(u uint64, base, count int) *Chain {
	for i := 0; i < count; i++ {
		b.WriteUintBase(u, base)
	}
	return b
}

func (b *Chain) WriteFloatRepeat(f float64, count int) *Chain {
	for i := 0; i < count; i++ {
		b.WriteFloat(f)
	}
	return b
}

func (b *Chain) WriteBoolRepeat(v bool, count int) *Chain {
	for i := 0; i < count; i++ {
		b.WriteBool(v)
	}
	return b
}

func (b *Chain) WriteXRepeat(x any, count int) *Chain {
	for i := 0; i < count; i++ {
		b.WriteX(x)
	}
	return b
}
