package bytebuf

import (
	"bytes"
	"net/url"
	"strconv"
	"testing"

	"github.com/koykov/byteconv"
)

func TestChain(t *testing.T) {
	t.Run("write", func(t *testing.T) {
		cb := NewChainSize(128)
		cb.Write(stage.b).WriteByte('-').
			WriteString(stage.s).WriteByte('-').
			WriteInt(stage.i).WriteByte('-').
			WriteUint(stage.u).WriteByte('-').
			WriteFloat(stage.f)

		if !bytes.Equal(cb.Bytes(), expectWS) {
			t.Error("Chain.Write*: mismatch result and expectation")
		}
	})
	t.Run("apply fn", func(t *testing.T) {
		cb := NewChainSize(128)
		cb.WriteString("foo").
			WriteApplyFnString("?q=front&p=1", func(dst, p []byte) []byte {
				p1 := url.QueryEscape(byteconv.B2S(p))
				dst = append(dst, p1...)
				return dst
			}).
			WriteString("bar")
		if cb.String() != "foo%3Fq%3Dfront%26p%3D1bar" {
			t.Error("Chain.WriteApplyFn*: mismatch result and expectation")
		}
	})
	t.Run("apply fn N", func(t *testing.T) {
		cb := NewChainSize(128)
		cb.WriteString("foo").
			WriteApplyFnNString("?q=front&p=1", func(dst, p []byte) []byte {
				p1 := url.QueryEscape(byteconv.B2S(p))
				dst = append(dst, p1...)
				return dst
			}, 3).
			WriteString("bar")
		if cb.String() != "foo%25253Fq%25253Dfront%252526p%25253D1bar" {
			t.Error("Chain.WriteApplyFn*: mismatch result and expectation")
		}
	})
	t.Run("reduce", func(t *testing.T) {
		cb := Chain{}
		cb.WriteString("0x0000").
			WriteIntBase(100, 16)
		hex := cb.Bytes()[cb.Len()-2:]
		cb.Reduce(4).Write(hex)
		if cb.String() != "0x0064" {
			t.FailNow()
		}
	})
	t.Run("format", func(t *testing.T) {
		cb := Chain{}
		cb.WriteFormat("int %d; float %f; string %s", 513, 3.1415, "foobar")
		if cb.String() != "int 513; float 3.141500; string foobar" {
			t.FailNow()
		}
	})
}

func BenchmarkChain(b *testing.B) {
	b.Run("write", func(b *testing.B) {
		cb := NewChainSize(128)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			cb.Reset().
				Write(stage.b).WriteByte('-').
				WriteString(stage.s).WriteByte('-').
				WriteInt(stage.i).WriteByte('-').
				WriteUint(stage.u).WriteByte('-').
				WriteFloat(stage.f)

			if !bytes.Equal(cb.Bytes(), expectWS) {
				b.Error("Chain.Write*: mismatch result and expectation")
			}
		}
	})
	b.Run("slice append", func(b *testing.B) {
		var buf []byte
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf = buf[:0]
			buf = append(buf, stage.b...)
			buf = append(buf, '-')
			buf = append(buf, stage.s...)
			buf = append(buf, '-')
			buf = append(buf, strconv.Itoa(int(stage.i))...)
			buf = append(buf, '-')
			buf = append(buf, strconv.Itoa(int(stage.u))...)
			buf = append(buf, '-')
			buf = append(buf, strconv.FormatFloat(stage.f, 'f', -1, 64)...)

			if !bytes.Equal(buf, expectWS) {
				b.Error("ByteArray: mismatch result and expectation")
			}
		}
	})
	b.Run("reduce", func(b *testing.B) {
		b.ReportAllocs()
		cb := Chain{}
		for i := 0; i < b.N; i++ {
			cb.Reset()
			cb.WriteString("0x0000").
				WriteIntBase(100, 16)
			hex := cb.Bytes()[cb.Len()-2:]
			cb.Reduce(4).Write(hex)
			if cb.String() != "0x0064" {
				b.FailNow()
			}
		}
	})
	b.Run("format", func(b *testing.B) {
		b.ReportAllocs()
		cb := Chain{}
		for i := 0; i < b.N; i++ {
			cb.Reset().WriteFormat("int %d; float %f; string %s", 513, 3.1415, "foobar")
		}
	})
}
