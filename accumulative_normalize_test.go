package bytebuf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/unicode/norm"
)

func TestAccumulativeWriteNormalize(t *testing.T) {
	tests := []struct {
		name string
		p    []byte
		form norm.Form
		want []byte
	}{
		{
			name: "empty input NFC",
			p:    nil,
			form: norm.NFC,
			want: []byte{},
		},
		{
			name: "empty input NFD",
			p:    nil,
			form: norm.NFD,
			want: []byte{},
		},
		{
			name: "empty input NFKC",
			p:    nil,
			form: norm.NFKC,
			want: []byte{},
		},
		{
			name: "empty input NFKD",
			p:    nil,
			form: norm.NFKD,
			want: []byte{},
		},
		{
			name: "ascii test NFC",
			p:    []byte("test"),
			form: norm.NFC,
			want: []byte("test"),
		},
		{
			name: "ascii test NFD",
			p:    []byte("test"),
			form: norm.NFD,
			want: []byte("test"),
		},
		{
			name: "ascii test NFKC",
			p:    []byte("test"),
			form: norm.NFKC,
			want: []byte("test"),
		},
		{
			name: "ascii test NFKD",
			p:    []byte("test"),
			form: norm.NFKD,
			want: []byte("test"),
		},
		{
			name: "combining marks composed by NFC",
			p:    []byte("e\u0301"),
			form: norm.NFC,
			want: []byte("\u00E9"),
		},
		{
			name: "combining marks preserved by NFD",
			p:    []byte("e\u0301"),
			form: norm.NFD,
			want: []byte("e\u0301"),
		},
		{
			name: "precomposed decomposed by NFD",
			p:    []byte("\u00E9"),
			form: norm.NFD,
			want: []byte("e\u0301"),
		},
		{
			name: "precomposed unchanged by NFC",
			p:    []byte("\u00E9"),
			form: norm.NFC,
			want: []byte("\u00E9"),
		},
		{
			name: "ligature expanded by NFKC",
			p:    []byte("\uFB03"),
			form: norm.NFKC,
			want: []byte("ffi"),
		},
		{
			name: "ligature expanded by NFKD",
			p:    []byte("\uFB03"),
			form: norm.NFKD,
			want: []byte("ffi"),
		},
		{
			name: "ligature unchanged by NFC",
			p:    []byte("\uFB03"),
			form: norm.NFC,
			want: []byte("\uFB03"),
		},
		{
			name: "ligature unchanged by NFD",
			p:    []byte("\uFB03"),
			form: norm.NFD,
			want: []byte("\uFB03"),
		},
		{
			name: "mixed precomposed and ligature NFKC",
			p:    []byte("\u00E9\uFB03"),
			form: norm.NFKC,
			want: []byte("\u00E9ffi"),
		},
		{
			name: "invalid UTF-8 passed through NFC",
			p:    []byte{0xff, 0xfe, 'a'},
			form: norm.NFC,
			want: []byte{0xff, 0xfe, 'a'},
		},
		{
			name: "invalid UTF-8 passed through NFKC",
			p:    []byte{0xff, 0xfe, 'a'},
			form: norm.NFKC,
			want: []byte{0xff, 0xfe, 'a'},
		},
		{
			name: "truncated multibyte passed through NFD",
			p:    []byte{0xe2, 0x82},
			form: norm.NFD,
			want: []byte{0xe2, 0x82},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteNormalize(tt.p, tt.form)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteNormalize(%q, %d)", tt.want, ab.Bytes(), tt.p, tt.form)
		})
	}
}

func TestAccumulativeWriteNormalizeString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		form norm.Form
		want []byte
	}{
		{
			name: "empty string NFC",
			s:    "",
			form: norm.NFC,
			want: []byte{},
		},
		{
			name: "empty string NFD",
			s:    "",
			form: norm.NFD,
			want: []byte{},
		},
		{
			name: "empty string NFKC",
			s:    "",
			form: norm.NFKC,
			want: []byte{},
		},
		{
			name: "empty string NFKD",
			s:    "",
			form: norm.NFKD,
			want: []byte{},
		},
		{
			name: "ascii test NFC",
			s:    "test",
			form: norm.NFC,
			want: []byte("test"),
		},
		{
			name: "ascii test NFD",
			s:    "test",
			form: norm.NFD,
			want: []byte("test"),
		},
		{
			name: "ascii test NFKC",
			s:    "test",
			form: norm.NFKC,
			want: []byte("test"),
		},
		{
			name: "ascii test NFKD",
			s:    "test",
			form: norm.NFKD,
			want: []byte("test"),
		},
		{
			name: "combining marks composed by NFC",
			s:    "e\u0301",
			form: norm.NFC,
			want: []byte("\u00E9"),
		},
		{
			name: "precomposed decomposed by NFD",
			s:    "\u00E9",
			form: norm.NFD,
			want: []byte("e\u0301"),
		},
		{
			name: "ligature expanded by NFKC",
			s:    "\uFB03",
			form: norm.NFKC,
			want: []byte("ffi"),
		},
		{
			name: "ligature expanded by NFKD",
			s:    "\uFB03",
			form: norm.NFKD,
			want: []byte("ffi"),
		},
		{
			name: "ligature unchanged by NFC",
			s:    "\uFB03",
			form: norm.NFC,
			want: []byte("\uFB03"),
		},
		{
			name: "mixed precomposed and ligature NFKD",
			s:    "\u00E9\uFB03",
			form: norm.NFKD,
			want: []byte("e\u0301ffi"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteNormalizeString(tt.s, tt.form)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteNormalizeString(%q, %d)", tt.want, ab.Bytes(), tt.s, tt.form)
		})
	}
}
