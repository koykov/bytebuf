package bytebuf

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/unicode/norm"
)

func TestChainWriteN(t *testing.T) {
	tests := []struct {
		name string
		p    []byte
		n    int
		want []byte
	}{
		{name: "zero copies", p: []byte("ab"), n: 0, want: []byte{}},
		{name: "single copy", p: []byte("ab"), n: 1, want: []byte("ab")},
		{name: "two copies", p: []byte("ab"), n: 2, want: []byte("abab")},
		{name: "three copies", p: []byte("ab"), n: 3, want: []byte("ababab")},
		{name: "negative copies", p: []byte("ab"), n: -1, want: []byte{}},
		{name: "empty slice single copy", p: []byte{}, n: 1, want: []byte{}},
		{name: "nil slice single copy", p: nil, n: 1, want: []byte{}},
		{name: "nil slice multiple copies", p: nil, n: 3, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteN(tt.p, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteN(%q, %d)", tt.want, cb.Bytes(), tt.p, tt.n)
		})
	}
}

func TestChainWriteByteN(t *testing.T) {
	tests := []struct {
		name string
		p    byte
		n    int
		want []byte
	}{
		{name: "zero copies", p: 'x', n: 0, want: []byte{}},
		{name: "single copy", p: 'x', n: 1, want: []byte("x")},
		{name: "three copies", p: 'x', n: 3, want: []byte("xxx")},
		{name: "negative copies", p: 'x', n: -1, want: []byte{}},
		{name: "zero byte three copies", p: 0x00, n: 3, want: []byte{0x00, 0x00, 0x00}},
		{name: "high byte two copies", p: 0xFF, n: 2, want: []byte{0xFF, 0xFF}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteByteN(tt.p, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteByteN(%d, %d)", tt.want, cb.Bytes(), tt.p, tt.n)
		})
	}
}

func TestChainWriteStringN(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want []byte
	}{
		{name: "zero copies", s: "ab", n: 0, want: []byte{}},
		{name: "single copy", s: "ab", n: 1, want: []byte("ab")},
		{name: "two copies", s: "ab", n: 2, want: []byte("abab")},
		{name: "three copies", s: "ab", n: 3, want: []byte("ababab")},
		{name: "negative copies", s: "ab", n: -1, want: []byte{}},
		{name: "empty string single copy", s: "", n: 1, want: []byte{}},
		{name: "empty string multiple copies", s: "", n: 3, want: []byte{}},
		{name: "unicode two copies", s: "世", n: 2, want: []byte("世世")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteStringN(tt.s, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteStringN(%q, %d)", tt.want, cb.Bytes(), tt.s, tt.n)
		})
	}
}

func TestChainWriteIntN(t *testing.T) {
	tests := []struct {
		name string
		i    int64
		n    int
		want []byte
	}{
		{name: "zero copies", i: 123, n: 0, want: []byte{}},
		{name: "single copy", i: 123, n: 1, want: []byte("123")},
		{name: "two copies", i: 123, n: 2, want: []byte("123123")},
		{name: "three copies", i: 123, n: 3, want: []byte("123123123")},
		{name: "negative copies", i: 123, n: -1, want: []byte{}},
		{name: "negative value", i: -123, n: 1, want: []byte("-123")},
		{name: "max int64", i: math.MaxInt64, n: 1, want: []byte("9223372036854775807")},
		{name: "min int64", i: math.MinInt64, n: 1, want: []byte("-9223372036854775808")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteIntN(tt.i, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteIntN(%d, %d)", tt.want, cb.Bytes(), tt.i, tt.n)
		})
	}
}

func TestChainWriteIntBaseN(t *testing.T) {
	bases := []struct {
		name string
		base int
		one  []byte
	}{
		{name: "base 2", base: 2, one: []byte("11111111")},
		{name: "base 8", base: 8, one: []byte("377")},
		{name: "base 10", base: 10, one: []byte("255")},
		{name: "base 16", base: 16, one: []byte("ff")},
		{name: "base 36", base: 36, one: []byte("73")},
	}
	for _, b := range bases {
		t.Run(b.name, func(t *testing.T) {
			tests := []struct {
				name string
				i    int64
				base int
				n    int
				want []byte
			}{
				{name: "zero copies", i: 255, base: b.base, n: 0, want: []byte{}},
				{name: "single copy", i: 255, base: b.base, n: 1, want: b.one},
				{name: "two copies", i: 255, base: b.base, n: 2, want: bytes.Repeat(b.one, 2)},
				{name: "three copies", i: 255, base: b.base, n: 3, want: bytes.Repeat(b.one, 3)},
				{name: "negative copies", i: 255, base: b.base, n: -1, want: []byte{}},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					cb := NewChainSize(128)
					cb.WriteIntBaseN(tt.i, tt.base, tt.n)
					assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteIntBaseN(%d, %d, %d)", tt.want, cb.Bytes(), tt.i, tt.base, tt.n)
				})
			}
		})
	}
}

func TestChainWriteUintN(t *testing.T) {
	tests := []struct {
		name string
		u    uint64
		n    int
		want []byte
	}{
		{name: "zero copies", u: 123, n: 0, want: []byte{}},
		{name: "single copy", u: 123, n: 1, want: []byte("123")},
		{name: "two copies", u: 123, n: 2, want: []byte("123123")},
		{name: "three copies", u: 123, n: 3, want: []byte("123123123")},
		{name: "negative copies", u: 123, n: -1, want: []byte{}},
		{name: "zero value", u: 0, n: 1, want: []byte("0")},
		{name: "max uint64", u: math.MaxUint64, n: 1, want: []byte("18446744073709551615")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteUintN(tt.u, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteUintN(%d, %d)", tt.want, cb.Bytes(), tt.u, tt.n)
		})
	}
}

func TestChainWriteUintBaseN(t *testing.T) {
	bases := []struct {
		name string
		base int
		one  []byte
	}{
		{name: "base 2", base: 2, one: []byte("11111111")},
		{name: "base 8", base: 8, one: []byte("377")},
		{name: "base 10", base: 10, one: []byte("255")},
		{name: "base 16", base: 16, one: []byte("ff")},
		{name: "base 36", base: 36, one: []byte("73")},
	}
	for _, b := range bases {
		t.Run(b.name, func(t *testing.T) {
			tests := []struct {
				name string
				u    uint64
				base int
				n    int
				want []byte
			}{
				{name: "zero copies", u: 255, base: b.base, n: 0, want: []byte{}},
				{name: "single copy", u: 255, base: b.base, n: 1, want: b.one},
				{name: "two copies", u: 255, base: b.base, n: 2, want: bytes.Repeat(b.one, 2)},
				{name: "three copies", u: 255, base: b.base, n: 3, want: bytes.Repeat(b.one, 3)},
				{name: "negative copies", u: 255, base: b.base, n: -1, want: []byte{}},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					cb := NewChainSize(128)
					cb.WriteUintBaseN(tt.u, tt.base, tt.n)
					assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteUintBaseN(%d, %d, %d)", tt.want, cb.Bytes(), tt.u, tt.base, tt.n)
				})
			}
		})
	}
}

func TestChainWriteFloatN(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		n    int
		want []byte
	}{
		{name: "zero copies", f: 1.5, n: 0, want: []byte{}},
		{name: "single copy", f: 1.5, n: 1, want: []byte("1.5")},
		{name: "two copies", f: 1.5, n: 2, want: []byte("1.51.5")},
		{name: "three copies", f: 1.5, n: 3, want: []byte("1.51.51.5")},
		{name: "negative copies", f: 1.5, n: -1, want: []byte{}},
		{name: "negative value", f: -1.5, n: 1, want: []byte("-1.5")},
		{name: "pi", f: math.Pi, n: 1, want: []byte("3.141592653589793")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteFloatN(tt.f, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteFloatN(%v, %d)", tt.want, cb.Bytes(), tt.f, tt.n)
		})
	}
}

func TestChainWriteBoolN(t *testing.T) {
	tests := []struct {
		name string
		v    bool
		n    int
		want []byte
	}{
		{name: "true zero copies", v: true, n: 0, want: []byte{}},
		{name: "true single copy", v: true, n: 1, want: []byte("true")},
		{name: "true two copies", v: true, n: 2, want: []byte("truetrue")},
		{name: "true three copies", v: true, n: 3, want: []byte("truetruetrue")},
		{name: "true negative copies", v: true, n: -1, want: []byte{}},
		{name: "false zero copies", v: false, n: 0, want: []byte{}},
		{name: "false single copy", v: false, n: 1, want: []byte("false")},
		{name: "false two copies", v: false, n: 2, want: []byte("falsefalse")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteBoolN(tt.v, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteBoolN(%v, %d)", tt.want, cb.Bytes(), tt.v, tt.n)
		})
	}
}

func TestChainWriteBinaryN(t *testing.T) {
	type testBinary struct {
		A uint32
		B uint16
		C byte
	}
	x := &testBinary{A: 0x12345678, B: 0xABCD, C: 0xFF}
	one := []byte{0x78, 0x56, 0x34, 0x12, 0xCD, 0xAB, 0xFF}

	tests := []struct {
		name string
		x    any
		n    int
		want []byte
	}{
		{name: "zero copies", x: x, n: 0, want: []byte{}},
		{name: "single copy", x: x, n: 1, want: one},
		{name: "two copies", x: x, n: 2, want: bytes.Repeat(one, 2)},
		{name: "three copies", x: x, n: 3, want: bytes.Repeat(one, 3)},
		{name: "negative copies", x: x, n: -1, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteBinaryN(binary.LittleEndian, tt.x, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteBinaryN(little endian, %+v, %d)", tt.want, cb.Bytes(), tt.x, tt.n)
		})
	}
}

func TestChainWriteXN(t *testing.T) {
	tests := []struct {
		name string
		x    any
		n    int
		want []byte
	}{
		{name: "int zero copies", x: 42, n: 0, want: []byte{}},
		{name: "int single copy", x: 42, n: 1, want: []byte("42")},
		{name: "int two copies", x: 42, n: 2, want: []byte("4242")},
		{name: "int three copies", x: 42, n: 3, want: []byte("424242")},
		{name: "int negative copies", x: 42, n: -1, want: []byte{}},
		{name: "string single copy", x: "hi", n: 1, want: []byte("hi")},
		{name: "string two copies", x: "hi", n: 2, want: []byte("hihi")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteXN(tt.x, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteXN(%v, %d)", tt.want, cb.Bytes(), tt.x, tt.n)
		})
	}
}

func TestChainWriteTimeN(t *testing.T) {
	format := "%Y-%m-%d"
	tm := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)

	tests := []struct {
		name string
		n    int
		want []byte
	}{
		{name: "zero copies", n: 0, want: []byte{}},
		{name: "single copy", n: 1, want: []byte("2024-01-15")},
		{name: "two copies", n: 2, want: []byte("2024-01-152024-01-15")},
		{name: "three copies", n: 3, want: []byte("2024-01-152024-01-152024-01-15")},
		{name: "negative copies", n: -1, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteTimeN(format, tm, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteTimeN(%q, %v, %d)", tt.want, cb.Bytes(), format, tm, tt.n)
		})
	}
}

func TestChainWriteULEB128N(t *testing.T) {
	one := []byte{0xAC, 0x02}
	max := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01}

	tests := []struct {
		name string
		v    uint64
		n    int
		want []byte
	}{
		{name: "zero copies", v: 300, n: 0, want: []byte{}},
		{name: "single copy", v: 300, n: 1, want: one},
		{name: "two copies", v: 300, n: 2, want: bytes.Repeat(one, 2)},
		{name: "three copies", v: 300, n: 3, want: bytes.Repeat(one, 3)},
		{name: "negative copies", v: 300, n: -1, want: []byte{}},
		{name: "zero value", v: 0, n: 1, want: []byte{0x00}},
		{name: "max uint64", v: math.MaxUint64, n: 1, want: max},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteULEB128N(tt.v, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteULEB128N(%d, %d)", tt.want, cb.Bytes(), tt.v, tt.n)
		})
	}
}

func TestChainWriteSLEB128N(t *testing.T) {
	one := []byte{0x9C, 0x7F}

	tests := []struct {
		name string
		v    int64
		n    int
		want []byte
	}{
		{name: "zero copies", v: -100, n: 0, want: []byte{}},
		{name: "single copy", v: -100, n: 1, want: one},
		{name: "two copies", v: -100, n: 2, want: bytes.Repeat(one, 2)},
		{name: "three copies", v: -100, n: 3, want: bytes.Repeat(one, 3)},
		{name: "negative copies", v: -100, n: -1, want: []byte{}},
		{name: "zero value", v: 0, n: 1, want: []byte{0x00}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteSLEB128N(tt.v, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteSLEB128N(%d, %d)", tt.want, cb.Bytes(), tt.v, tt.n)
		})
	}
}

func TestChainWriteNormalizeN(t *testing.T) {
	tests := []struct {
		name string
		p    []byte
		n    int
		want []byte
	}{
		{name: "zero copies", p: []byte("test"), n: 0, want: []byte{}},
		{name: "single copy", p: []byte("test"), n: 1, want: []byte("test")},
		{name: "two copies", p: []byte("test"), n: 2, want: []byte("testtest")},
		{name: "three copies", p: []byte("test"), n: 3, want: []byte("testtesttest")},
		{name: "negative copies", p: []byte("test"), n: -1, want: []byte{}},
		{name: "empty input single copy", p: []byte{}, n: 1, want: []byte{}},
		{name: "nil input multiple copies", p: nil, n: 3, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteNormalizeN(tt.p, norm.NFC, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteNormalizeN(%q, NFC, %d)", tt.want, cb.Bytes(), tt.p, tt.n)
		})
	}
}

func TestChainWriteNormalizeStringN(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want []byte
	}{
		{name: "zero copies", s: "test", n: 0, want: []byte{}},
		{name: "single copy", s: "test", n: 1, want: []byte("test")},
		{name: "two copies", s: "test", n: 2, want: []byte("testtest")},
		{name: "three copies", s: "test", n: 3, want: []byte("testtesttest")},
		{name: "negative copies", s: "test", n: -1, want: []byte{}},
		{name: "empty string single copy", s: "", n: 1, want: []byte{}},
		{name: "empty string multiple copies", s: "", n: 3, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteNormalizeStringN(tt.s, norm.NFC, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteNormalizeStringN(%q, NFC, %d)", tt.want, cb.Bytes(), tt.s, tt.n)
		})
	}
}

func TestChainWriteHexN(t *testing.T) {
	tests := []struct {
		name string
		p    []byte
		n    int
		want []byte
	}{
		{name: "zero copies", p: []byte{0x00, 0xFF}, n: 0, want: []byte{}},
		{name: "single copy", p: []byte{0x00, 0xFF}, n: 1, want: []byte("00ff")},
		{name: "two copies", p: []byte{0x00, 0xFF}, n: 2, want: []byte("00ff00ff")},
		{name: "three copies", p: []byte{0x00, 0xFF}, n: 3, want: []byte("00ff00ff00ff")},
		{name: "negative copies", p: []byte{0x00, 0xFF}, n: -1, want: []byte{}},
		{name: "empty input single copy", p: []byte{}, n: 1, want: []byte{}},
		{name: "nil input multiple copies", p: nil, n: 2, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteHexN(tt.p, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteHexN(%v, %d)", tt.want, cb.Bytes(), tt.p, tt.n)
			if tt.n > 0 {
				ref := bytes.Repeat([]byte(hex.EncodeToString(tt.p)), tt.n)
				assert.True(t, bytes.Equal(ref, cb.Bytes()), "want %q, got %q, WriteHexN hex reference (%v, %d)", ref, cb.Bytes(), tt.p, tt.n)
			}
		})
	}
}

func TestChainWriteHexStringN(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want []byte
	}{
		{name: "zero copies", s: "00ff", n: 0, want: []byte{}},
		{name: "single copy", s: "00ff", n: 1, want: []byte("30306666")},
		{name: "two copies", s: "00ff", n: 2, want: []byte("3030666630306666")},
		{name: "three copies", s: "00ff", n: 3, want: []byte("303066663030666630306666")},
		{name: "negative copies", s: "00ff", n: -1, want: []byte{}},
		{name: "empty string single copy", s: "", n: 1, want: []byte{}},
		{name: "empty string multiple copies", s: "", n: 2, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteHexStringN(tt.s, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, WriteHexStringN(%q, %d)", tt.want, cb.Bytes(), tt.s, tt.n)
			if tt.n > 0 {
				ref := bytes.Repeat([]byte(hex.EncodeToString([]byte(tt.s))), tt.n)
				assert.True(t, bytes.Equal(ref, cb.Bytes()), "want %q, got %q, WriteHexStringN hex reference (%q, %d)", ref, cb.Bytes(), tt.s, tt.n)
			}
		})
	}
}
