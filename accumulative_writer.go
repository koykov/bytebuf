package bytebuf

// AccumulativeWriter is a wrapper around Accumulative that implements IO writers interfaces.
type AccumulativeWriter struct {
	buf *Accumulative
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
