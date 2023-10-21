package bytebuf

// AccBufWriter is a wrapper around AccBuf that implements io.Writer.
type AccBufWriter struct {
	AccBuf *AccumulativeBuf
}

func (b AccBufWriter) Write(p []byte) (int, error) {
	b.AccBuf.Write(p)
	return len(p), nil
}
