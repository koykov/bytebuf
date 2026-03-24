package bytebuf

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/koykov/x2bytes"
)

type Writer interface {
	io.Writer
	io.StringWriter
	io.ByteWriter
	io.WriterAt
	Grow(newLen int)
	GrowDelta(delta int)
	WriteInt(i int64) (int, error)
	WriteUint(u uint64) (int, error)
	WriteFloat(f float64) (int, error)
	WriteBool(b bool) (int, error)
	WriteBinary(order binary.ByteOrder, x any) (int, error)
	WriteX(x any) (int, error)
	WriteMarshallerTo(m MarshallerTo) (int, error)
	WriteApplyFn(p []byte, fn func(dst, p []byte) []byte) (int, error)
	WriteApplyFnN(p []byte, fn func(dst, p []byte) []byte, n int) (int, error)
	WriteApplyFnString(s string, fn func(dst, p []byte) []byte) (int, error)
	WriteApplyFnNString(s string, fn func(dst, p []byte) []byte, n int) (int, error)
	WriteTime(format string, t time.Time) (int, error)
	WriteULEB128(v uint64) (int, error)
	WriteSLEB128(v int64) (int, error)
	Len() int
	Cap() int
	Bytes() []byte
	String() string
	Reset()
	ResetSafe()
}

type writer struct {
	buf *Chain
}

func (w writer) Grow(newLen int) {
	w.buf.Grow(newLen)
}

func (w writer) GrowDelta(delta int) {
	w.buf.GrowDelta(delta)
}

func (w writer) Bytes() []byte {
	return w.buf.Bytes()
}

func (w writer) String() string {
	return w.buf.String()
}

func (w writer) Write(p []byte) (int, error) {
	w.buf.Write(p)
	return len(p), nil
}

func (w writer) WriteByte(b byte) error {
	w.buf.WriteByte(b)
	return nil
}

func (w writer) WriteString(s string) (int, error) {
	w.buf.WriteString(s)
	return len(s), nil
}

func (w writer) WriteAt(p []byte, off int64) (int, error) {
	w.buf.Grow(int(off))
	w.buf.Write(p)
	return len(p), nil
}

func (w writer) WriteInt(i int64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteInt(i)
	return w.buf.Len() - off, nil
}

func (w writer) WriteUint(u uint64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteUint(u)
	return w.buf.Len() - off, nil
}

func (w writer) WriteFloat(f float64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteFloat(f)
	return w.buf.Len() - off, nil
}

func (w writer) WriteBool(b bool) (int, error) {
	off := w.buf.Len()
	w.buf.WriteBool(b)
	return w.buf.Len() - off, nil
}

func (w writer) WriteBinary(order binary.ByteOrder, x any) (int, error) {
	off := w.buf.Len()
	err := binary.Write(w, order, x)
	return w.buf.Len() - off, err
}

func (w writer) WriteX(x any) (int, error) {
	off := w.buf.Len()
	var err error
	*w.buf, err = x2bytes.ToBytes(*w.buf, x)
	w.buf.WriteX(x)
	return w.buf.Len() - off, err
}

func (w writer) WriteMarshallerTo(m MarshallerTo) (int, error) {
	off := w.buf.Len()
	w.buf.WriteMarshallerTo(m)
	return w.buf.Len() - off, nil
}

func (w writer) WriteApplyFn(p []byte, fn func(dst, p []byte) []byte) (int, error) {
	off := w.buf.Len()
	w.buf.WriteApplyFn(p, fn)
	return w.buf.Len() - off, nil
}

func (w writer) WriteApplyFnN(p []byte, fn func(dst, p []byte) []byte, n int) (int, error) {
	off := w.buf.Len()
	w.buf.WriteApplyFnN(p, fn, n)
	return w.buf.Len() - off, nil
}

func (w writer) WriteApplyFnString(s string, fn func(dst, p []byte) []byte) (int, error) {
	off := w.buf.Len()
	w.buf.WriteApplyFnString(s, fn)
	return w.buf.Len() - off, nil
}

func (w writer) WriteApplyFnNString(s string, fn func(dst, p []byte) []byte, n int) (int, error) {
	off := w.buf.Len()
	w.buf.WriteApplyFnNString(s, fn, n)
	return w.buf.Len() - off, nil
}

func (w writer) WriteTime(format string, t time.Time) (int, error) {
	off := w.buf.Len()
	w.buf.WriteTime(format, t)
	return w.buf.Len() - off, nil
}

func (w writer) WriteULEB128(v uint64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteULEB128(v)
	return w.buf.Len() - off, nil
}

func (w writer) WriteSLEB128(v int64) (int, error) {
	off := w.buf.Len()
	w.buf.WriteSLEB128(v)
	return w.buf.Len() - off, nil
}

func (w writer) Len() int {
	return w.buf.Len()
}

func (w writer) Cap() int {
	return w.buf.Cap()
}

func (w writer) Reset() {
	w.buf.Reset()
}

func (w writer) ResetSafe() {
	w.buf.ResetSafe()
}
