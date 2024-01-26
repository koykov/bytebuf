package bytebuf

func (b *Accumulative) WriteIf(cond bool, p []byte) *Accumulative {
	if cond {
		b.Write(p)
	}
	return b
}

func (b *Accumulative) WriteByteIf(cond bool, p byte) *Accumulative {
	if cond {
		b.WriteByte(p)
	}
	return b
}

func (b *Accumulative) WriteStringIf(cond bool, s string) *Accumulative {
	if cond {
		b.WriteString(s)
	}
	return b
}

func (b *Accumulative) WriteIntIf(cond bool, i int64) *Accumulative {
	if cond {
		b.WriteInt(i)
	}
	return b
}

func (b *Accumulative) WriteUintIf(cond bool, u uint64) *Accumulative {
	if cond {
		b.WriteUint(u)
	}
	return b
}

func (b *Accumulative) WriteFloatIf(cond bool, f float64) *Accumulative {
	if cond {
		b.WriteFloat(f)
	}
	return b
}

func (b *Accumulative) WriteBoolIf(cond bool, v bool) *Accumulative {
	if cond {
		b.WriteBool(v)
	}
	return b
}

func (b *Accumulative) WriteXIf(cond bool, x any) *Accumulative {
	if cond {
		b.WriteX(x)
	}
	return b
}

func (b *Accumulative) WriteMarshallerToIf(cond bool, m MarshallerTo) *Accumulative {
	if cond {
		b.WriteMarshallerTo(m)
	}
	return b
}

func (b *Accumulative) WriteApplyFnIf(cond bool, p []byte, fn func(dst, p []byte) []byte) *Accumulative {
	if cond {
		b.WriteApplyFn(p, fn)
	}
	return b
}

func (b *Accumulative) WriteApplyFnStringIf(cond bool, s string, fn func(dst, p []byte) []byte) *Accumulative {
	if cond {
		b.WriteApplyFnString(s, fn)
	}
	return b
}
