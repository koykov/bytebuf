package bytebuf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteMarshallerTo(t *testing.T) {
	t.Run("chain", func(t *testing.T) {
		cb := NewAccumulativeSize(1024)
		cb.WriteMarshallerTo(stage)
		assert.True(t, bytes.Equal(expectMT, cb.Bytes()), "Chain.WriteMarshallerTo: mismatch result and expectation, want %q, got %q", expectMT, cb.Bytes())
	})
	t.Run("accumulative", func(t *testing.T) {
		ab := &Accumulative{}
		ab.WriteString("trash data").StakeOut().WriteMarshallerTo(stage)
		assert.True(t, bytes.Equal(expectMT, ab.StakedBytes()), "Chain.WriteMarshallerTo: mismatch result and expectation, want %q, got %q", expectMT, ab.StakedBytes())
	})
}

func BenchmarkWriteMarshallerTo(b *testing.B) {
	b.Run("chain", func(b *testing.B) {
		cb := NewAccumulativeSize(1024)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			cb.Reset().WriteMarshallerTo(stage)
			if !bytes.Equal(cb.Bytes(), expectMT) {
				b.Error("Chain.WriteMarshallerTo: mismatch result and expectation")
			}
		}
	})
	b.Run("accumulative", func(b *testing.B) {
		ab := &Accumulative{}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ab.Reset().WriteString("trash data").StakeOut().WriteMarshallerTo(stage)
			if !bytes.Equal(ab.StakedBytes(), expectMT) {
				b.Error("Chain.WriteMarshallerTo: mismatch result and expectation")
			}
		}
	})
}
