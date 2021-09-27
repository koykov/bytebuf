package bytebuf

import (
	"bytes"
	"strconv"
	"testing"
)

func TestChainBuf_Write(t *testing.T) {
	cb := &ChainBuf{}
	cb.Write(stage.b).WriteByte('-').
		WriteStr(stage.s).WriteByte('-').
		WriteInt(stage.i).WriteByte('-').
		WriteUint(stage.u).WriteByte('-').
		WriteFloat(stage.f)

	if !bytes.Equal(cb.Bytes(), expectWS) {
		t.Error("ChainBuf.Write*: mismatch result and expectation")
	}
}

func BenchmarkChainBuf_Write(b *testing.B) {
	cb := &ChainBuf{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cb.Reset().
			Write(stage.b).WriteByte('-').
			WriteStr(stage.s).WriteByte('-').
			WriteInt(stage.i).WriteByte('-').
			WriteUint(stage.u).WriteByte('-').
			WriteFloat(stage.f)

		if !bytes.Equal(cb.Bytes(), expectWS) {
			b.Error("ChainBuf.Write*: mismatch result and expectation")
		}
	}
}

func BenchmarkByteSlice(b *testing.B) {
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
}
