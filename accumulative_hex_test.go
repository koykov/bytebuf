package bytebuf

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccumulativeWriteHex(t *testing.T) {
	tests := []struct {
		name string
		p    []byte
		want []byte
	}{
		{name: "empty input", p: []byte{}, want: []byte{}},
		{name: "single zero byte", p: []byte{0x00}, want: []byte("00")},
		{name: "single high byte", p: []byte{0xFF}, want: []byte("ff")},
		{name: "single ascii byte", p: []byte{0x41}, want: []byte("41")},
		{name: "multiple bytes", p: []byte{0x00, 0xFF}, want: []byte("00ff")},
		{name: "odd length", p: []byte{0x01, 0x02, 0x03}, want: []byte("010203")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteHex(tt.p)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteHex(% x)", tt.want, ab.Bytes(), tt.p)
			assert.Equal(t, hex.EncodeToString(tt.p), ab.String(), "WriteHex(% x)", tt.p)
		})
	}
}

func TestAccumulativeWriteHexString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []byte
	}{
		{name: "empty input", s: "", want: []byte{}},
		{name: "zero string", s: "00", want: []byte("3030")},
		{name: "ff string", s: "ff", want: []byte("6666")},
		{name: "00ff string", s: "00ff", want: []byte("30306666")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteHexString(tt.s)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteHexString(%q)", tt.want, ab.Bytes(), tt.s)
			assert.Equal(t, hex.EncodeToString([]byte(tt.s)), ab.String(), "WriteHexString(%q)", tt.s)
		})
	}
}
