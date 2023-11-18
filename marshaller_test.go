package bytebuf

import (
	"bytes"
	"testing"
)

func TestWriteMarshallerTo(t *testing.T) {
	t.Run("chainbuf", func(t *testing.T) {
		cb := &Chain{}
		cb.WriteMarshallerTo(stage)
		if !bytes.Equal(cb.Bytes(), expectMT) {
			t.Error("Chain.WriteMarshallerTo: mismatch result and expectation")
		}
	})
	t.Run("accbuf", func(t *testing.T) {
		ab := &Accumulative{}
		ab.WriteStr("trash data").StakeOut().WriteMarshallerTo(stage)
		if !bytes.Equal(ab.StakedBytes(), expectMT) {
			t.Error("Chain.WriteMarshallerTo: mismatch result and expectation")
		}
	})
}

func BenchmarkWriteMarshallerTo(b *testing.B) {
	b.Run("chainbuf", func(b *testing.B) {
		cb := &Chain{}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			cb.Reset().WriteMarshallerTo(stage)
			if !bytes.Equal(cb.Bytes(), expectMT) {
				b.Error("Chain.WriteMarshallerTo: mismatch result and expectation")
			}
		}
	})
	b.Run("accbuf", func(b *testing.B) {
		ab := &Accumulative{}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ab.Reset().WriteStr("trash data").StakeOut().WriteMarshallerTo(stage)
			if !bytes.Equal(ab.StakedBytes(), expectMT) {
				b.Error("Chain.WriteMarshallerTo: mismatch result and expectation")
			}
		}
	})
}
