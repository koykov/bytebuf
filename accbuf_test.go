package bytebuf

import (
	"bytes"
	"testing"
)

func TestAccumulativeBuf(t *testing.T) {
	t.Run("ab write", func(t *testing.T) {
		ab := &AccumulativeBuf{}
		ab.WriteString("foobar").
			StakeOut().
			Write(stage.b).WriteByte('-').
			WriteString(stage.s).WriteByte('-').
			WriteInt(stage.i).WriteByte('-').
			WriteUint(stage.u).WriteByte('-').
			WriteFloat(stage.f)

		if !bytes.Equal(ab.StakedBytes(), expectWS) {
			t.Error("AccumulativeBuf.Write*: mismatch result and expectation")
		}
	})
}

func BenchmarkAccumulativeBuf(b *testing.B) {
	b.Run("ab write", func(b *testing.B) {
		ab := &AccumulativeBuf{}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ab.Reset().
				WriteString("foobar").
				StakeOut().
				Write(stage.b).WriteByte('-').
				WriteString(stage.s).WriteByte('-').
				WriteInt(stage.i).WriteByte('-').
				WriteUint(stage.u).WriteByte('-').
				WriteFloat(stage.f)

			if !bytes.Equal(ab.StakedBytes(), expectWS) {
				b.Error("AccumulativeBuf.Write*: mismatch result and expectation")
			}
		}
	})
}
