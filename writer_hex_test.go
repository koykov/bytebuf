package bytebuf

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriterWriteHex(t *testing.T) {
	tests := []struct {
		name string
		p    []byte
		want []byte
	}{
		{name: "empty input", p: []byte{}, want: nil},
		{name: "single zero byte", p: []byte{0x00}, want: []byte("00")},
		{name: "single high byte", p: []byte{0xFF}, want: []byte("ff")},
		{name: "single ascii byte", p: []byte{0x41}, want: []byte("41")},
		{name: "multiple bytes", p: []byte{0x00, 0xFF}, want: []byte("00ff")},
		{name: "odd length", p: []byte{0x01, 0x02, 0x03}, want: []byte("010203")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWriter(nil)
			n, err := w.WriteHex(tt.p)
			require.NoError(t, err, "WriteHex(% x)", tt.p)

			assert.True(t, bytes.Equal(tt.want, w.Bytes()), "WriteHex(% x): want %q, got %q", tt.p, tt.want, w.Bytes())
			assert.Equal(t, hex.EncodeToString(tt.p), w.String(), "WriteHex(% x)", tt.p)
			assert.Equal(t, len(hex.EncodeToString(tt.p)), n, "WriteHex(% x) returned n", tt.p)
		})
	}
}

func TestWriterWriteHexString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []byte
	}{
		{name: "empty input", s: "", want: nil},
		{name: "zero string", s: "00", want: []byte("3030")},
		{name: "ff string", s: "ff", want: []byte("6666")},
		{name: "00ff string", s: "00ff", want: []byte("30306666")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWriter(nil)
			n, err := w.WriteHexString(tt.s)
			require.NoError(t, err, "WriteHexString(%q)", tt.s)

			assert.True(t, bytes.Equal(tt.want, w.Bytes()), "WriteHexString(%q): want %q, got %q", tt.s, tt.want, w.Bytes())
			assert.Equal(t, hex.EncodeToString([]byte(tt.s)), w.String(), "WriteHexString(%q)", tt.s)
			assert.Equal(t, len(hex.EncodeToString([]byte(tt.s))), n, "WriteHexString(%q) returned n", tt.s)
		})
	}
}
