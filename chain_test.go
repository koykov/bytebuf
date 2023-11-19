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
		cb := &Chain{}
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
		cb := &Chain{}
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
}

func BenchmarkChain(b *testing.B) {
	b.Run("write", func(b *testing.B) {
		cb := &Chain{}
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
}
