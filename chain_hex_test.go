package bytebuf

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainWriteHex(t *testing.T) {
	tests := []struct {
		name string
		p    []byte
		want []byte
	}{
		{name: "empty input", p: nil, want: []byte{}},
		{name: "empty slice input", p: []byte{}, want: []byte{}},
		{name: "single byte 0x00", p: []byte{0x00}, want: []byte("00")},
		{name: "single byte 0xFF", p: []byte{0xFF}, want: []byte("ff")},
		{name: "single byte 0x41", p: []byte{0x41}, want: []byte("41")},
		{name: "two bytes 0x00 0xFF", p: []byte{0x00, 0xFF}, want: []byte("00ff")},
		{name: "odd length three bytes", p: []byte{0x01, 0x02, 0x03}, want: []byte("010203")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteHex(tt.p)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q", tt.want, cb.Bytes())
			ref := hex.EncodeToString(tt.p)
			assert.True(t, bytes.Equal([]byte(ref), cb.Bytes()), "want %q, got %q", []byte(ref), cb.Bytes())
		})
	}
}

func TestChainWriteHexString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []byte
	}{
		{name: "empty string", s: "", want: []byte{}},
		{name: "ascii 00", s: "00", want: []byte("3030")},
		{name: "ascii ff", s: "ff", want: []byte("6666")},
		{name: "ascii 00ff", s: "00ff", want: []byte("30306666")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteHexString(tt.s)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q", tt.want, cb.Bytes())
			ref := hex.EncodeToString([]byte(tt.s))
			assert.True(t, bytes.Equal([]byte(ref), cb.Bytes()), "want %q, got %q", []byte(ref), cb.Bytes())
		})
	}
}
