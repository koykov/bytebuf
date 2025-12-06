package bytebuf

import (
	"encoding/binary"

	"github.com/koykov/x2bytes"
)

// ChainWriter is a wrapper around Chain that implements IO writers interfaces.
type ChainWriter struct {
	buf *Chain
}

func (w ChainWriter) Bytes() []byte {
	return w.buf.Bytes()
}

func (w ChainWriter) String() string {
	return w.buf.String()
}

func (w ChainWriter) Write(p []byte) (int, error) {
	w.buf.Write(p)
	return len(p), nil
}

func (w ChainWriter) WriteByte(b byte) error {
	w.buf.WriteByte(b)
	return nil
}

func (w ChainWriter) WriteString(s string) (int, error) {
	w.buf.WriteString(s)
	return len(s), nil
}

func (w ChainWriter) WriteInt(i int64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteInt(i)
	return w.buf.Len() - off, nil
}

func (w ChainWriter) WriteUint(u uint64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteUint(u)
	return w.buf.Len() - off, nil
}

func (w ChainWriter) WriteFloat(f float64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteFloat(f)
	return w.buf.Len() - off, nil
}

func (w ChainWriter) WriteBool(b bool) (int, error) {
	off := w.buf.Len()
	w.buf.WriteBool(b)
	return w.buf.Len() - off, nil
}

func (w ChainWriter) WriteBinary(order binary.ByteOrder, x any) (int, error) {
	off := w.buf.Len()
	err := binary.Write(w, order, x)
	return w.buf.Len() - off, err
}

func (w ChainWriter) WriteX(x any) (int, error) {
	off := w.buf.Len()
	var err error
	*w.buf, err = x2bytes.ToBytes(*w.buf, x)
	w.buf.WriteX(x)
	return w.buf.Len() - off, err
}

func (w ChainWriter) WriteMarshallerTo(m MarshallerTo) (int, error) {
	off := w.buf.Len()
	w.buf.WriteMarshallerTo(m)
	return w.buf.Len() - off, nil
}

func (w ChainWriter) WriteApplyFn(p []byte, fn func(dst, p []byte) []byte) (int, error) {
	off := w.buf.Len()
	w.buf.WriteApplyFn(p, fn)
	return w.buf.Len() - off, nil
}

func (w ChainWriter) WriteApplyFnN(p []byte, fn func(dst, p []byte) []byte, n int) (int, error) {
	off := w.buf.Len()
	w.buf.WriteApplyFnN(p, fn, n)
	return w.buf.Len() - off, nil
}

func (w ChainWriter) WriteApplyFnString(s string, fn func(dst, p []byte) []byte) (int, error) {
	off := w.buf.Len()
	w.buf.WriteApplyFnString(s, fn)
	return w.buf.Len() - off, nil
}

func (w ChainWriter) WriteApplyFnNString(s string, fn func(dst, p []byte) []byte, n int) (int, error) {
	off := w.buf.Len()
	w.buf.WriteApplyFnNString(s, fn, n)
	return w.buf.Len() - off, nil
}

func (w ChainWriter) Len() int {
	return w.buf.Len()
}

func (w ChainWriter) Cap() int {
	return w.buf.Cap()
}

func (w ChainWriter) Reset() {
	w.buf.Reset()
}
