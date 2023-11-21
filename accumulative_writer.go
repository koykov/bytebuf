package bytebuf

// AccumulativeWriter is a wrapper around Accumulative that implements IO writers interfaces.
type AccumulativeWriter struct {
	buf *Accumulative
}

func (b AccumulativeWriter) Bytes() []byte {
	return b.buf.Bytes()
}

func (b AccumulativeWriter) String() string {
	return b.buf.String()
}

func (b AccumulativeWriter) Write(p []byte) (int, error) {
	b.buf.Write(p)
	return len(p), nil
}

func (b AccumulativeWriter) WriteByte(p byte) error {
	b.buf.WriteByte(p)
	return nil
}

func (b AccumulativeWriter) WriteString(s string) (int, error) {
	b.buf.WriteString(s)
	return len(s), nil
}

func (b AccumulativeWriter) WriteInt(i int64) (int, error) {
	off := b.buf.Len()
	b.buf.WriteInt(i)
	return b.buf.Len() - off, nil
}

func (b AccumulativeWriter) WriteUint(u uint64) (int, error) {
	off := b.buf.Len()
	b.buf.WriteUint(u)
	return b.buf.Len() - off, nil
}

func (b AccumulativeWriter) WriteFloat(f float64) (int, error) {
	off := b.buf.Len()
	b.buf.WriteFloat(f)
	return b.buf.Len() - off, nil
}

func (b AccumulativeWriter) WriteBool(b_ bool) (int, error) {
	off := b.buf.Len()
	b.buf.WriteBool(b_)
	return b.buf.Len() - off, nil
}

func (b AccumulativeWriter) WriteX(x any) (int, error) {
	off := b.buf.Len()
	b.buf.WriteX(x)
	return b.buf.Len() - off, b.buf.Error()
}

func (b AccumulativeWriter) WriteMarshallerTo(m MarshallerTo) (int, error) {
	off := b.buf.Len()
	b.buf.WriteMarshallerTo(m)
	return b.buf.Len() - off, nil
}

func (b AccumulativeWriter) WriteApplyFn(p []byte, fn func(dst, p []byte) []byte) (int, error) {
	off := b.buf.Len()
	b.buf.WriteApplyFn(p, fn)
	return b.buf.Len() - off, nil
}

func (b AccumulativeWriter) WriteApplyFnString(s string, fn func(dst, p []byte) []byte) (int, error) {
	off := b.buf.Len()
	b.buf.WriteApplyFnString(s, fn)
	return b.buf.Len() - off, nil
}

func (b AccumulativeWriter) Len() int {
	return b.buf.Len()
}

func (b AccumulativeWriter) Cap() int {
	return b.buf.Cap()
}

func (b AccumulativeWriter) Reset() {
	b.buf.Reset()
}
