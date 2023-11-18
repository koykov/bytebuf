package bytebuf

import (
	"bytes"
	"net/url"
	"strconv"
	"testing"

	"github.com/koykov/fastconv"
)

func TestChainBuf(t *testing.T) {
	t.Run("cb write", func(t *testing.T) {
		cb := &ChainBuf{}
		cb.Write(stage.b).WriteByte('-').
			WriteString(stage.s).WriteByte('-').
			WriteInt(stage.i).WriteByte('-').
			WriteUint(stage.u).WriteByte('-').
			WriteFloat(stage.f)

		if !bytes.Equal(cb.Bytes(), expectWS) {
			t.Error("ChainBuf.Write*: mismatch result and expectation")
		}
	})
	t.Run("apply fn", func(t *testing.T) {
		cb := &ChainBuf{}
		cb.WriteString("foo").
			WriteApplyFnString("?q=front&p=1", func(dst, p []byte) []byte {
				p1 := url.QueryEscape(fastconv.B2S(p))
				dst = append(dst, p1...)
				return dst
			}).
			WriteString("bar")
		if cb.String() != "foo%3Fq%3Dfront%26p%3D1bar" {
			t.Error("ChainBuf.WriteApplyFn*: mismatch result and expectation")
		}
	})
}

func BenchmarkChainBuf(b *testing.B) {
	b.Run("cb write", func(b *testing.B) {
		cb := &ChainBuf{}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			cb.Reset().
				Write(stage.b).WriteByte('-').
				WriteString(stage.s).WriteByte('-').
				WriteInt(stage.i).WriteByte('-').
				WriteUint(stage.u).WriteByte('-').
				WriteFloat(stage.f)

			if !bytes.Equal(cb.Bytes(), expectWS) {
				b.Error("ChainBuf.Write*: mismatch result and expectation")
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
