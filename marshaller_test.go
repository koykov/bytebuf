package bytebuf

import (
	"bytes"
	"testing"
)

func TestWriteMarshallerTo(t *testing.T) {
	t.Run("chainbuf", func(t *testing.T) {
		cb := &ChainBuf{}
		cb.WriteMarshallerTo(stage)
		if !bytes.Equal(cb.Bytes(), expectMT) {
			t.Error("ChainBuf.WriteMarshallerTo: mismatch result and expectation")
		}
	})
	t.Run("accbuf", func(t *testing.T) {
		ab := &AccumulativeBuf{}
		ab.WriteStr("trash data").StakeOut().WriteMarshallerTo(stage)
		if !bytes.Equal(ab.StakedBytes(), expectMT) {
			t.Error("ChainBuf.WriteMarshallerTo: mismatch result and expectation")
		}
	})
}

func BenchmarkWriteMarshallerTo(b *testing.B) {
	b.Run("chainbuf", func(b *testing.B) {
		cb := &ChainBuf{}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			cb.Reset().WriteMarshallerTo(stage)
			if !bytes.Equal(cb.Bytes(), expectMT) {
				b.Error("ChainBuf.WriteMarshallerTo: mismatch result and expectation")
			}
		}
	})
	b.Run("accbuf", func(b *testing.B) {
		ab := &AccumulativeBuf{}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ab.Reset().WriteStr("trash data").StakeOut().WriteMarshallerTo(stage)
			if !bytes.Equal(ab.StakedBytes(), expectMT) {
				b.Error("ChainBuf.WriteMarshallerTo: mismatch result and expectation")
			}
		}
	})
}
