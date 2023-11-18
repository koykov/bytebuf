package bytebuf

// AccBufWriter is a wrapper around AccBuf that implements IO writers interfaces.
type AccBufWriter struct {
	AccBuf *Accumulative
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
	b.AccBuf.WriteString(s)
	return len(s), nil
}
