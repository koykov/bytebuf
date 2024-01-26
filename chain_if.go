package bytebuf

func (b *Chain) WriteIf(cond bool, p []byte) *Chain {
	if cond {
		b.Write(p)
	}
	return b
}

func (b *Chain) WriteByteIf(cond bool, p byte) *Chain {
	if cond {
		b.WriteByte(p)
	}
	return b
}

func (b *Chain) WriteStringIf(cond bool, s string) *Chain {
	if cond {
		b.WriteString(s)
	}
	return b
}

func (b *Chain) WriteIntIf(cond bool, i int64) *Chain {
	if cond {
		b.WriteInt(i)
	}
	return b
}

func (b *Chain) WriteUintIf(cond bool, u uint64) *Chain {
	if cond {
		b.WriteUint(u)
	}
	return b
}

func (b *Chain) WriteFloatIf(cond bool, f float64) *Chain {
	if cond {
		b.WriteFloat(f)
	}
	return b
}

func (b *Chain) WriteBoolIf(cond bool, v bool) *Chain {
	if cond {
		b.WriteBool(v)
	}
	return b
}

func (b *Chain) WriteXIf(cond bool, x any) *Chain {
	if cond {
		b.WriteX(x)
	}
	return b
}

func (b *Chain) WriteMarshallerToIf(cond bool, m MarshallerTo) *Chain {
	if cond {
		b.WriteMarshallerTo(m)
	}
	return b
}

func (b *Chain) WriteApplyFnIf(cond bool, p []byte, fn func(dst, p []byte) []byte) *Chain {
	if cond {
		b.WriteApplyFn(p, fn)
	}
	return b
}

func (b *Chain) WriteApplyFnStringIf(cond bool, s string, fn func(dst, p []byte) []byte) *Chain {
	if cond {
		b.WriteApplyFnString(s, fn)
	}
	return b
}
