package bytebuf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/unicode/norm"
)

func TestChainWriteNormalize(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		form norm.Form
		want []byte
	}{
		{name: "empty input NFC", in: []byte(""), form: norm.NFC, want: []byte("")},
		{name: "empty input NFD", in: []byte(""), form: norm.NFD, want: []byte("")},
		{name: "empty input NFKC", in: []byte(""), form: norm.NFKC, want: []byte("")},
		{name: "empty input NFKD", in: []byte(""), form: norm.NFKD, want: []byte("")},

		{name: "ascii unchanged NFC", in: []byte("test"), form: norm.NFC, want: []byte("test")},
		{name: "ascii unchanged NFD", in: []byte("test"), form: norm.NFD, want: []byte("test")},
		{name: "ascii unchanged NFKC", in: []byte("test"), form: norm.NFKC, want: []byte("test")},
		{name: "ascii unchanged NFKD", in: []byte("test"), form: norm.NFKD, want: []byte("test")},

		{name: "combining mark composed by NFC", in: []byte("e\u0301"), form: norm.NFC, want: []byte("\u00e9")},
		{name: "combining mark preserved by NFD", in: []byte("e\u0301"), form: norm.NFD, want: []byte("e\u0301")},
		{name: "combining mark composed by NFKC", in: []byte("e\u0301"), form: norm.NFKC, want: []byte("\u00e9")},
		{name: "combining mark preserved by NFKD", in: []byte("e\u0301"), form: norm.NFKD, want: []byte("e\u0301")},

		{name: "precomposed preserved by NFC", in: []byte("\u00e9"), form: norm.NFC, want: []byte("\u00e9")},
		{name: "precomposed decomposed by NFD", in: []byte("\u00e9"), form: norm.NFD, want: []byte("e\u0301")},
		{name: "precomposed preserved by NFKC", in: []byte("\u00e9"), form: norm.NFKC, want: []byte("\u00e9")},
		{name: "precomposed decomposed by NFKD", in: []byte("\u00e9"), form: norm.NFKD, want: []byte("e\u0301")},

		{name: "ligature preserved by NFC", in: []byte("\uFB03"), form: norm.NFC, want: []byte("\uFB03")},
		{name: "ligature preserved by NFD", in: []byte("\uFB03"), form: norm.NFD, want: []byte("\uFB03")},
		{name: "ligature expanded by NFKC", in: []byte("\uFB03"), form: norm.NFKC, want: []byte("ffi")},
		{name: "ligature expanded by NFKD", in: []byte("\uFB03"), form: norm.NFKD, want: []byte("ffi")},
		{name: "ligature within ascii expanded by NFKC", in: []byte("a\uFB03b"), form: norm.NFKC, want: []byte("affib")},

		{name: "invalid utf8 preserved by NFC", in: []byte{0xff, 0xfe}, form: norm.NFC, want: []byte{0xff, 0xfe}},
		{name: "invalid utf8 preserved by NFD", in: []byte{0xff, 0xfe}, form: norm.NFD, want: []byte{0xff, 0xfe}},
		{name: "invalid utf8 preserved by NFKC", in: []byte{0xff, 0xfe}, form: norm.NFKC, want: []byte{0xff, 0xfe}},
		{name: "invalid utf8 preserved by NFKD", in: []byte{0xff, 0xfe}, form: norm.NFKD, want: []byte{0xff, 0xfe}},
		{name: "trailing invalid byte preserved", in: []byte("abc\xff"), form: norm.NFKC, want: []byte("abc\xff")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteNormalize(tt.in, tt.form)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteNormalize", tt.want, cb.Bytes())
		})
	}
}

func TestChainWriteNormalizeString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		form norm.Form
		want []byte
	}{
		{name: "empty string NFC", in: "", form: norm.NFC, want: []byte("")},
		{name: "empty string NFD", in: "", form: norm.NFD, want: []byte("")},
		{name: "empty string NFKC", in: "", form: norm.NFKC, want: []byte("")},
		{name: "empty string NFKD", in: "", form: norm.NFKD, want: []byte("")},

		{name: "ascii unchanged NFC", in: "test", form: norm.NFC, want: []byte("test")},
		{name: "ascii unchanged NFD", in: "test", form: norm.NFD, want: []byte("test")},
		{name: "ascii unchanged NFKC", in: "test", form: norm.NFKC, want: []byte("test")},
		{name: "ascii unchanged NFKD", in: "test", form: norm.NFKD, want: []byte("test")},

		{name: "combining mark composed by NFC", in: "e\u0301", form: norm.NFC, want: []byte("\u00e9")},
		{name: "combining mark preserved by NFD", in: "e\u0301", form: norm.NFD, want: []byte("e\u0301")},
		{name: "precomposed decomposed by NFD", in: "\u00e9", form: norm.NFD, want: []byte("e\u0301")},

		{name: "ligature expanded by NFKC", in: "\uFB03", form: norm.NFKC, want: []byte("ffi")},
		{name: "ligature expanded by NFKD", in: "\uFB03", form: norm.NFKD, want: []byte("ffi")},
		{name: "ligature preserved by NFC", in: "\uFB03", form: norm.NFC, want: []byte("\uFB03")},
		{name: "ligature within ascii expanded by NFKC", in: "a\uFB03b", form: norm.NFKC, want: []byte("affib")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteNormalizeString(tt.in, tt.form)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteNormalizeString", tt.want, cb.Bytes())
		})
	}
}
