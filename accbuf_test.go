package bytebuf

import (
	"bytes"
	"testing"
)

func TestAccumulativeBuf_Write(t *testing.T) {
	ab := &AccumulativeBuf{}
	ab.WriteStr("foobar").
		StakeOut().
		Write(stage.b).WriteByte('-').
		WriteStr(stage.s).WriteByte('-').
		WriteInt(stage.i).WriteByte('-').
		WriteUint(stage.u).WriteByte('-').
		WriteFloat(stage.f)

	if !bytes.Equal(ab.StakedBytes(), expectWS) {
		t.Error("AccumulativeBuf.Write*: mismatch result and expectation")
	}
}

func BenchmarkAccumulativeBuf_Write(b *testing.B) {
	ab := &AccumulativeBuf{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ab.Reset().
			WriteStr("foobar").
			StakeOut().
			Write(stage.b).WriteByte('-').
			WriteStr(stage.s).WriteByte('-').
			WriteInt(stage.i).WriteByte('-').
			WriteUint(stage.u).WriteByte('-').
			WriteFloat(stage.f)

		if !bytes.Equal(ab.StakedBytes(), expectWS) {
			b.Error("AccumulativeBuf.Write*: mismatch result and expectation")
		}
	}
}
