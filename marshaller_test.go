package bytebuf

import (
	"bytes"
	"testing"
)

func TestChainBuf_WriteMarshallerTo(t *testing.T) {
	cb := &ChainBuf{}
	cb.WriteMarshallerTo(stage)
	if !bytes.Equal(cb.Bytes(), expectMT) {
		t.Error("ChainBuf.WriteMarshallerTo: mismatch result and expectation")
	}
}

func TestAccumulativeBuf_WriteMarshallerTo(t *testing.T) {
	cb := &AccumulativeBuf{}
	cb.WriteStr("trash data").StakeOut().WriteMarshallerTo(stage)
	if !bytes.Equal(cb.StakedBytes(), expectMT) {
		t.Error("ChainBuf.WriteMarshallerTo: mismatch result and expectation")
	}
}

func BenchmarkChainBuf_WriteMarshallerTo(b *testing.B) {
	cb := &ChainBuf{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cb.Reset().WriteMarshallerTo(stage)
		if !bytes.Equal(cb.Bytes(), expectMT) {
			b.Error("ChainBuf.WriteMarshallerTo: mismatch result and expectation")
		}
	}
}

func BenchmarkAccumulativeBuf_WriteMarshallerTo(b *testing.B) {
	cb := &AccumulativeBuf{}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cb.Reset().WriteStr("trash data").StakeOut().WriteMarshallerTo(stage)
		if !bytes.Equal(cb.StakedBytes(), expectMT) {
			b.Error("ChainBuf.WriteMarshallerTo: mismatch result and expectation")
		}
	}
}
