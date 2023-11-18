package bytebuf

import (
	"bytes"
	"testing"
)

func TestAccumulative(t *testing.T) {
	t.Run("write", func(t *testing.T) {
		ab := &Accumulative{}
		ab.WriteString("foobar").
			StakeOut().
			Write(stage.b).WriteByte('-').
			WriteString(stage.s).WriteByte('-').
			WriteInt(stage.i).WriteByte('-').
			WriteUint(stage.u).WriteByte('-').
			WriteFloat(stage.f)

		if !bytes.Equal(ab.StakedBytes(), expectWS) {
			t.Error("Accumulative.Write*: mismatch result and expectation")
		}
	})
}

func BenchmarkAccumulative(b *testing.B) {
	b.Run("write", func(b *testing.B) {
		ab := &Accumulative{}
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
				b.Error("Accumulative.Write*: mismatch result and expectation")
			}
		}
	})
}
