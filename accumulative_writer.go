package bytebuf

import (
	"encoding/binary"
	"time"
)

// AccumulativeWriter is a wrapper around Accumulative that implements IO writers interfaces.
type AccumulativeWriter struct {
	buf *Accumulative
}

func (w AccumulativeWriter) Bytes() []byte {
	return w.buf.Bytes()
}

func (w AccumulativeWriter) String() string {
	return w.buf.String()
}

func (w AccumulativeWriter) Write(p []byte) (int, error) {
	w.buf.Write(p)
	return len(p), nil
}

func (w AccumulativeWriter) WriteByte(b byte) error {
	w.buf.WriteByte(b)
	return nil
}

func (w AccumulativeWriter) WriteString(s string) (int, error) {
	w.buf.WriteString(s)
	return len(s), nil
}

func (w AccumulativeWriter) WriteInt(i int64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteInt(i)
	return w.buf.Len() - off, nil
}

func (w AccumulativeWriter) WriteUint(u uint64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteUint(u)
	return w.buf.Len() - off, nil
}

func (w AccumulativeWriter) WriteFloat(f float64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteFloat(f)
	return w.buf.Len() - off, nil
}

func (w AccumulativeWriter) WriteBool(b bool) (int, error) {
	off := w.buf.Len()
	w.buf.WriteBool(b)
	return w.buf.Len() - off, nil
}

func (w AccumulativeWriter) WriteBinary(order binary.ByteOrder, x any) (int, error) {
	off := w.buf.Len()
	err := binary.Write(w, order, x)
	return w.buf.Len() - off, err
}

func (w AccumulativeWriter) WriteX(x any) (int, error) {
	off := w.buf.Len()
	w.buf.WriteX(x)
	return w.buf.Len() - off, w.buf.Error()
}

func (w AccumulativeWriter) WriteMarshallerTo(m MarshallerTo) (int, error) {
	off := w.buf.Len()
	w.buf.WriteMarshallerTo(m)
	return w.buf.Len() - off, nil
}

func (w AccumulativeWriter) WriteApplyFn(p []byte, fn func(dst, p []byte) []byte) (int, error) {
	off := w.buf.Len()
	w.buf.WriteApplyFn(p, fn)
	return w.buf.Len() - off, nil
}

func (w AccumulativeWriter) WriteApplyFnN(p []byte, fn func(dst, p []byte) []byte, n int) (int, error) {
	off := w.buf.Len()
	w.buf.WriteApplyFnN(p, fn, n)
	return w.buf.Len() - off, nil
}

func (w AccumulativeWriter) WriteApplyFnString(s string, fn func(dst, p []byte) []byte) (int, error) {
	off := w.buf.Len()
	w.buf.WriteApplyFnString(s, fn)
	return w.buf.Len() - off, nil
}

func (w AccumulativeWriter) WriteApplyFnNString(s string, fn func(dst, p []byte) []byte, n int) (int, error) {
	off := w.buf.Len()
	w.buf.WriteApplyFnNString(s, fn, n)
	return w.buf.Len() - off, nil
}

func (w AccumulativeWriter) WriteTime(format string, t time.Time) (int, error) {
	off := w.buf.Len()
	w.buf.WriteTime(format, t)
	return w.buf.Len() - off, nil
}

func (w AccumulativeWriter) WriteULEB128(v uint64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteULEB128(v)
	return w.buf.Len() - off, nil
}

func (w AccumulativeWriter) WriteSLEB128(v int64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteSLEB128(v)
	return w.buf.Len() - off, nil
}

func (w AccumulativeWriter) Len() int {
	return w.buf.Len()
}

func (w AccumulativeWriter) Cap() int {
	return w.buf.Cap()
}

func (w AccumulativeWriter) Reset() {
	w.buf.Reset()
}
