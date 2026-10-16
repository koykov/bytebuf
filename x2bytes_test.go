package bytebuf

import (
	"bytes"
	"testing"

	"github.com/koykov/x2bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChainToBytes(t *testing.T) {
	chainData := NewChain([]byte("chain data"))
	emptyChain := NewChain(nil)

	tests := []struct {
		name    string
		dst     []byte
		val     any
		want    []byte
		wantErr bool
	}{
		{name: "chain with data", dst: nil, val: chainData, want: []byte("chain data"), wantErr: false},
		{name: "empty chain", dst: nil, val: emptyChain, want: []byte{}, wantErr: false},
		{name: "appends to non-empty dst", dst: []byte("prefix:"), val: chainData, want: []byte("prefix:chain data"), wantErr: false},
		{name: "non-matching type", dst: nil, val: []byte("plain bytes"), want: nil, wantErr: true},
		{name: "nil value", dst: nil, val: nil, want: nil, wantErr: true},
		{name: "non-matching type keeps dst", dst: []byte("keep:"), val: 42, want: []byte("keep:"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChainToBytes(tt.dst, tt.val)
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, x2bytes.ErrUnknownType, err)
			} else {
				require.NoError(t, err)
			}
			// ChainToBytes returns a nil slice for an empty chain while the table
			// expects []byte{}, so keep bytes.Equal (nil == empty) semantics.
			require.True(t, bytes.Equal(tt.want, got), "got %q, want %q", got, tt.want)
		})
	}
}

func TestAccumulativeToBytes(t *testing.T) {
	accData := NewAccumulative([]byte("accumulative data"))
	emptyAcc := NewAccumulative(nil)

	tests := []struct {
		name    string
		dst     []byte
		val     any
		want    []byte
		wantErr bool
	}{
		{name: "accumulative with data", dst: nil, val: accData, want: []byte("accumulative data"), wantErr: false},
		{name: "empty accumulative", dst: nil, val: emptyAcc, want: []byte{}, wantErr: false},
		{name: "appends to non-empty dst", dst: []byte("prefix:"), val: accData, want: []byte("prefix:accumulative data"), wantErr: false},
		{name: "non-matching type", dst: nil, val: []byte("plain bytes"), want: nil, wantErr: true},
		{name: "nil value", dst: nil, val: nil, want: nil, wantErr: true},
		{name: "non-matching type keeps dst", dst: []byte("keep:"), val: 42, want: []byte("keep:"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AccumulativeToBytes(tt.dst, tt.val)
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, x2bytes.ErrUnknownType, err)
			} else {
				require.NoError(t, err)
			}
			// AccumulativeToBytes returns a nil slice for an empty buffer while the
			// table expects []byte{}, so keep bytes.Equal (nil == empty) semantics.
			require.True(t, bytes.Equal(tt.want, got), "got %q, want %q", got, tt.want)
		})
	}
}
