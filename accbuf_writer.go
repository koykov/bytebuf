package bytebuf

// AccBufWriter is a wrapper around AccBuf that implements io.Writer.
type AccBufWriter struct {
	AccBuf *AccumulativeBuf
}

func (b AccBufWriter) Write(p []byte) (int, error) {
	b.AccBuf.Write(p)
	return len(p), nil
}

func (b AccBufWriter) WriteByte(p byte) error {
	b.AccBuf.WriteByte(p)
	return nil
}

func (b AccBufWriter) WriteString(s string) (int, error) {
	b.AccBuf.WriteStr(s)
	return len(s), nil
}
