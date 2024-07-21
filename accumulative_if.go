package bytebuf

import "time"

// Contains conditional write methods.

func (b *Accumulative) WriteIf(cond bool, p []byte) *Accumulative {
	b.buf.WriteIf(cond, p)
	return b
}

func (b *Accumulative) WriteByteIf(cond bool, p byte) *Accumulative {
	b.buf.WriteByteIf(cond, p)
	return b
}

func (b *Accumulative) WriteStringIf(cond bool, s string) *Accumulative {
	b.buf.WriteStringIf(cond, s)
	return b
}

func (b *Accumulative) WriteIntIf(cond bool, i int64) *Accumulative {
	b.buf.WriteIntIf(cond, i)
	return b
}

func (b *Accumulative) WriteUintIf(cond bool, u uint64) *Accumulative {
	b.buf.WriteUintIf(cond, u)
	return b
}

func (b *Accumulative) WriteFloatIf(cond bool, f float64) *Accumulative {
	b.buf.WriteFloatIf(cond, f)
	return b
}

func (b *Accumulative) WriteBoolIf(cond bool, v bool) *Accumulative {
	b.buf.WriteBoolIf(cond, v)
	return b
}

func (b *Accumulative) WriteXIf(cond bool, x any) *Accumulative {
	b.buf.WriteXIf(cond, x)
	return b
}

func (b *Accumulative) WriteMarshallerToIf(cond bool, m MarshallerTo) *Accumulative {
	b.buf.WriteMarshallerToIf(cond, m)
	return b
}

func (b *Accumulative) WriteApplyFnIf(cond bool, p []byte, fn func(dst, p []byte) []byte) *Accumulative {
	b.buf.WriteApplyFnIf(cond, p, fn)
	return b
}

func (b *Accumulative) WriteApplyFnNIf(cond bool, p []byte, fn func(dst, p []byte) []byte, n int) *Accumulative {
	b.buf.WriteApplyFnNIf(cond, p, fn, n)
	return b
}

func (b *Accumulative) WriteApplyFnStringIf(cond bool, s string, fn func(dst, p []byte) []byte) *Accumulative {
	b.buf.WriteApplyFnStringIf(cond, s, fn)
	return b
}

func (b *Accumulative) WriteApplyFnNStringIf(cond bool, s string, fn func(dst, p []byte) []byte, n int) *Accumulative {
	b.buf.WriteApplyFnNStringIf(cond, s, fn, n)
	return b
}

func (b *Accumulative) WriteTimeIf(cond bool, format string, t time.Time) *Accumulative {
	b.buf.WriteTimeIf(cond, format, t)
	return b
}
