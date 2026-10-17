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

func TestChainWriteIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		p    []byte
		want []byte
	}{
		{name: "cond true writes bytes", cond: true, p: []byte("hello world"), want: []byte("hello world")},
		{name: "cond true writes single byte", cond: true, p: []byte{0x01}, want: []byte{0x01}},
		{name: "cond true writes empty slice", cond: true, p: []byte{}, want: []byte{}},
		{name: "cond true writes nil slice", cond: true, p: nil, want: []byte{}},
		{name: "cond false leaves buffer empty", cond: false, p: []byte("hello world"), want: []byte{}},
		{name: "cond false with empty slice leaves buffer empty", cond: false, p: []byte{}, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteIf(tt.cond, tt.p)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteIf(%v, %q)", tt.want, cb.Bytes(), tt.cond, tt.p)
		})
	}
}

func TestChainWriteByteIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		p    byte
		want []byte
	}{
		{name: "cond true writes byte", cond: true, p: 'a', want: []byte("a")},
		{name: "cond true writes zero byte", cond: true, p: 0x00, want: []byte{0x00}},
		{name: "cond true writes max byte", cond: true, p: 0xFF, want: []byte{0xFF}},
		{name: "cond false leaves buffer empty", cond: false, p: 'a', want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteByteIf(tt.cond, tt.p)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteByteIf(%v, %d)", tt.want, cb.Bytes(), tt.cond, tt.p)
		})
	}
}

func TestChainWriteStringIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		s    string
		want []byte
	}{
		{name: "cond true writes string", cond: true, s: "hello world", want: []byte("hello world")},
		{name: "cond true writes empty string", cond: true, s: "", want: []byte{}},
		{name: "cond true writes unicode string", cond: true, s: "Hello, 世界", want: []byte("Hello, 世界")},
		{name: "cond false leaves buffer empty", cond: false, s: "hello world", want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteStringIf(tt.cond, tt.s)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteStringIf(%v, %q)", tt.want, cb.Bytes(), tt.cond, tt.s)
		})
	}
}

func TestChainWriteIntIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		i    int64
		want []byte
	}{
		{name: "cond true writes zero", cond: true, i: 0, want: []byte("0")},
		{name: "cond true writes positive", cond: true, i: 123456789, want: []byte("123456789")},
		{name: "cond true writes negative", cond: true, i: -123456789, want: []byte("-123456789")},
		{name: "cond true writes max int64", cond: true, i: math.MaxInt64, want: []byte("9223372036854775807")},
		{name: "cond true writes min int64", cond: true, i: math.MinInt64, want: []byte("-9223372036854775808")},
		{name: "cond false leaves buffer empty", cond: false, i: 42, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteIntIf(tt.cond, tt.i)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteIntIf(%v, %d)", tt.want, cb.Bytes(), tt.cond, tt.i)
		})
	}
}

func TestChainWriteIntBaseIf(t *testing.T) {
	tests := []struct {
		name      string
		cond      bool
		i         int64
		base      int
		want      []byte
		wantPanic bool
	}{
		{name: "cond true base 2", cond: true, i: 100, base: 2, want: []byte("1100100")},
		{name: "cond true base 8", cond: true, i: 100, base: 8, want: []byte("144")},
		{name: "cond true base 10", cond: true, i: 100, base: 10, want: []byte("100")},
		{name: "cond true base 16", cond: true, i: 100, base: 16, want: []byte("64")},
		{name: "cond true base 36", cond: true, i: 100, base: 36, want: []byte("2s")},
		{name: "cond true base 16 negative", cond: true, i: -100, base: 16, want: []byte("-64")},
		{name: "cond true base 2 zero", cond: true, i: 0, base: 2, want: []byte("0")},
		{name: "cond false leaves buffer empty", cond: false, i: 100, base: 16, want: []byte{}},
		{name: "invalid base 0 panics", cond: true, i: 100, base: 0, wantPanic: true},
		{name: "invalid base 1 panics", cond: true, i: 100, base: 1, wantPanic: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			if tt.wantPanic {
				defer func() {
					assert.NotNil(t, recover(), "Chain.WriteIntBaseIf(%v, %d, %d) did not panic, want panic", tt.cond, tt.i, tt.base)
				}()
			}
			cb.WriteIntBaseIf(tt.cond, tt.i, tt.base)
			if tt.wantPanic {
				return
			}
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteIntBaseIf(%v, %d, %d)", tt.want, cb.Bytes(), tt.cond, tt.i, tt.base)
		})
	}
}

func TestChainWriteUintIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		u    uint64
		want []byte
	}{
		{name: "cond true writes zero", cond: true, u: 0, want: []byte("0")},
		{name: "cond true writes small", cond: true, u: 255, want: []byte("255")},
		{name: "cond true writes max uint64", cond: true, u: math.MaxUint64, want: []byte("18446744073709551615")},
		{name: "cond false leaves buffer empty", cond: false, u: 255, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteUintIf(tt.cond, tt.u)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteUintIf(%v, %d)", tt.want, cb.Bytes(), tt.cond, tt.u)
		})
	}
}

func TestChainWriteUintBaseIf(t *testing.T) {
	tests := []struct {
		name      string
		cond      bool
		u         uint64
		base      int
		want      []byte
		wantPanic bool
	}{
		{name: "cond true base 2", cond: true, u: 100, base: 2, want: []byte("1100100")},
		{name: "cond true base 8", cond: true, u: 100, base: 8, want: []byte("144")},
		{name: "cond true base 10", cond: true, u: 100, base: 10, want: []byte("100")},
		{name: "cond true base 16", cond: true, u: 100, base: 16, want: []byte("64")},
		{name: "cond true base 36", cond: true, u: 100, base: 36, want: []byte("2s")},
		{name: "cond true base 16 max uint64", cond: true, u: math.MaxUint64, base: 16, want: []byte("ffffffffffffffff")},
		{name: "cond true base 2 zero", cond: true, u: 0, base: 2, want: []byte("0")},
		{name: "cond false leaves buffer empty", cond: false, u: 100, base: 16, want: []byte{}},
		{name: "invalid base 0 panics", cond: true, u: 100, base: 0, wantPanic: true},
		{name: "invalid base 1 panics", cond: true, u: 100, base: 1, wantPanic: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			if tt.wantPanic {
				defer func() {
					assert.NotNil(t, recover(), "Chain.WriteUintBaseIf(%v, %d, %d) did not panic, want panic", tt.cond, tt.u, tt.base)
				}()
			}
			cb.WriteUintBaseIf(tt.cond, tt.u, tt.base)
			if tt.wantPanic {
				return
			}
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteUintBaseIf(%v, %d, %d)", tt.want, cb.Bytes(), tt.cond, tt.u, tt.base)
		})
	}
}

func TestChainWriteFloatIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		f    float64
		want []byte
	}{
		{name: "cond true writes zero", cond: true, f: 0.0, want: []byte("0")},
		{name: "cond true writes positive", cond: true, f: 123.456, want: []byte("123.456")},
		{name: "cond true writes negative", cond: true, f: -123.456, want: []byte("-123.456")},
		{name: "cond true writes pi", cond: true, f: math.Pi, want: []byte("3.141592653589793")},
		{name: "cond false leaves buffer empty", cond: false, f: math.Pi, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteFloatIf(tt.cond, tt.f)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteFloatIf(%v, %v)", tt.want, cb.Bytes(), tt.cond, tt.f)
		})
	}
}

func TestChainWriteBoolIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		v    bool
		want []byte
	}{
		{name: "cond true writes true", cond: true, v: true, want: []byte("true")},
		{name: "cond true writes false", cond: true, v: false, want: []byte("false")},
		{name: "cond false with true leaves buffer empty", cond: false, v: true, want: []byte{}},
		{name: "cond false with false leaves buffer empty", cond: false, v: false, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteBoolIf(tt.cond, tt.v)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteBoolIf(%v, %v)", tt.want, cb.Bytes(), tt.cond, tt.v)
		})
	}
}

func TestChainWriteBinaryIf(t *testing.T) {
	type testStruct struct {
		A uint32
		B uint16
		C byte
	}
	value := &testStruct{A: 0x12345678, B: 0xABCD, C: 0xFF}

	tests := []struct {
		name  string
		cond  bool
		order binary.ByteOrder
		x     any
		want  []byte
	}{
		{
			name:  "cond true little endian",
			cond:  true,
			order: binary.LittleEndian,
			x:     value,
			want:  []byte{0x78, 0x56, 0x34, 0x12, 0xCD, 0xAB, 0xFF},
		},
		{
			name:  "cond true big endian",
			cond:  true,
			order: binary.BigEndian,
			x:     value,
			want:  []byte{0x12, 0x34, 0x56, 0x78, 0xAB, 0xCD, 0xFF},
		},
		{
			name:  "cond false leaves buffer empty",
			cond:  false,
			order: binary.LittleEndian,
			x:     value,
			want:  []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteBinaryIf(tt.cond, tt.order, tt.x)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteBinaryIf(%v, %v, %v)", tt.want, cb.Bytes(), tt.cond, tt.order, tt.x)
		})
	}
}

func TestChainWriteXIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		x    any
		want []byte
	}{
		{name: "cond true writes int", cond: true, x: 42, want: []byte("42")},
		{name: "cond true writes negative int64", cond: true, x: int64(-7), want: []byte("-7")},
		{name: "cond true writes string", cond: true, x: "hello", want: []byte("hello")},
		{name: "cond true writes bytes", cond: true, x: []byte("bytes"), want: []byte("bytes")},
		{name: "cond true writes empty string", cond: true, x: "", want: []byte{}},
		{name: "cond false leaves buffer empty", cond: false, x: 42, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteXIf(tt.cond, tt.x)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteXIf(%v, %v)", tt.want, cb.Bytes(), tt.cond, tt.x)
		})
	}
}

func TestChainWriteMarshallerToIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		m    MarshallerTo
		want []byte
	}{
		{name: "cond true writes marshalled proto", cond: true, m: stage, want: expectMT},
		{name: "cond true with nil marshaller leaves buffer empty", cond: true, m: nil, want: []byte{}},
		{name: "cond false leaves buffer empty", cond: false, m: stage, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteMarshallerToIf(tt.cond, tt.m)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteMarshallerToIf(%v, %T)", tt.want, cb.Bytes(), tt.cond, tt.m)
		})
	}
}

func TestChainWriteApplyFnIf(t *testing.T) {
	addPrefix := func(dst, p []byte) []byte {
		return append(append(dst, "prefix:"...), p...)
	}

	tests := []struct {
		name string
		cond bool
		p    []byte
		want []byte
	}{
		{name: "cond true applies fn", cond: true, p: []byte("test"), want: []byte("prefix:test")},
		{name: "cond true applies fn to empty input", cond: true, p: []byte{}, want: []byte("prefix:")},
		{name: "cond false leaves buffer empty", cond: false, p: []byte("test"), want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteApplyFnIf(tt.cond, tt.p, addPrefix)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteApplyFnIf(%v, %q)", tt.want, cb.Bytes(), tt.cond, tt.p)
		})
	}
}

func TestChainWriteApplyFnNIf(t *testing.T) {
	duplicate := func(dst, p []byte) []byte {
		return append(dst, append(p, p...)...)
	}

	tests := []struct {
		name string
		cond bool
		p    []byte
		n    int
		want []byte
	}{
		{name: "cond true n=3 expands input", cond: true, p: []byte("ab"), n: 3, want: bytes.Repeat([]byte("ab"), 8)},
		{name: "cond true n=1 duplicates once", cond: true, p: []byte("ab"), n: 1, want: []byte("abab")},
		{name: "cond true n=0 writes nothing", cond: true, p: []byte("ab"), n: 0, want: []byte{}},
		{name: "cond false leaves buffer empty", cond: false, p: []byte("ab"), n: 3, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteApplyFnNIf(tt.cond, tt.p, duplicate, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteApplyFnNIf(%v, %q, fn, %d)", tt.want, cb.Bytes(), tt.cond, tt.p, tt.n)
		})
	}
}

func TestChainWriteApplyFnStringIf(t *testing.T) {
	addPrefix := func(dst, p []byte) []byte {
		return append(append(dst, "prefix:"...), p...)
	}

	tests := []struct {
		name string
		cond bool
		s    string
		want []byte
	}{
		{name: "cond true applies fn", cond: true, s: "test", want: []byte("prefix:test")},
		{name: "cond true applies fn to empty string", cond: true, s: "", want: []byte("prefix:")},
		{name: "cond false leaves buffer empty", cond: false, s: "test", want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteApplyFnStringIf(tt.cond, tt.s, addPrefix)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteApplyFnStringIf(%v, %q)", tt.want, cb.Bytes(), tt.cond, tt.s)
		})
	}
}

func TestChainWriteApplyFnNStringIf(t *testing.T) {
	duplicate := func(dst, p []byte) []byte {
		return append(dst, append(p, p...)...)
	}

	tests := []struct {
		name string
		cond bool
		s    string
		n    int
		want []byte
	}{
		{name: "cond true n=3 expands input", cond: true, s: "ab", n: 3, want: bytes.Repeat([]byte("ab"), 8)},
		{name: "cond true n=0 writes nothing", cond: true, s: "ab", n: 0, want: []byte{}},
		{name: "cond false leaves buffer empty", cond: false, s: "ab", n: 3, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteApplyFnNStringIf(tt.cond, tt.s, duplicate, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteApplyFnNStringIf(%v, %q, fn, %d)", tt.want, cb.Bytes(), tt.cond, tt.s, tt.n)
		})
	}
}

func TestChainWriteTimeIf(t *testing.T) {
	layout := "%Y-%m-%d %H:%M:%S"
	tm := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)

	tests := []struct {
		name string
		cond bool
		want []byte
	}{
		{name: "cond true writes formatted time", cond: true, want: []byte("2024-01-15 10:30:45")},
		{name: "cond false leaves buffer empty", cond: false, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteTimeIf(tt.cond, layout, tm)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteTimeIf(%v, %q, %v)", tt.want, cb.Bytes(), tt.cond, layout, tm)
		})
	}
}

func TestChainWriteULEB128If(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		v    uint64
		want []byte
	}{
		{name: "cond true zero", cond: true, v: 0, want: []byte{0x00}},
		{name: "cond true boundary 127", cond: true, v: 127, want: []byte{0x7F}},
		{name: "cond true boundary 128", cond: true, v: 128, want: []byte{0x80, 0x01}},
		{name: "cond true 300", cond: true, v: 300, want: []byte{0xAC, 0x02}},
		{
			name: "cond true max uint64",
			cond: true,
			v:    math.MaxUint64,
			want: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01},
		},
		{name: "cond false leaves buffer empty", cond: false, v: 300, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteULEB128If(tt.cond, tt.v)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteULEB128If(%v, %d)", tt.want, cb.Bytes(), tt.cond, tt.v)
		})
	}
}

func TestChainWriteSLEB128If(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		v    int64
		want []byte
	}{
		{name: "cond true zero", cond: true, v: 0, want: []byte{0x00}},
		{name: "cond true 63", cond: true, v: 63, want: []byte{0x3F}},
		{name: "cond true 64", cond: true, v: 64, want: []byte{0xC0, 0x00}},
		{name: "cond true -1", cond: true, v: -1, want: []byte{0x7F}},
		{name: "cond true -64", cond: true, v: -64, want: []byte{0x40}},
		{name: "cond true -100", cond: true, v: -100, want: []byte{0x9C, 0x7F}},
		{name: "cond false leaves buffer empty", cond: false, v: -100, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteSLEB128If(tt.cond, tt.v)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteSLEB128If(%v, %d)", tt.want, cb.Bytes(), tt.cond, tt.v)
		})
	}
}

func TestChainWriteNormalizeIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		p    []byte
		form norm.Form
		want []byte
	}{
		{name: "NFC keeps ligature", cond: true, p: []byte("\uFB03"), form: norm.NFC, want: []byte("\uFB03")},
		{name: "NFD keeps ligature", cond: true, p: []byte("\uFB03"), form: norm.NFD, want: []byte("\uFB03")},
		{name: "NFKC expands ligature", cond: true, p: []byte("\uFB03"), form: norm.NFKC, want: []byte("ffi")},
		{name: "NFKD expands ligature", cond: true, p: []byte("\uFB03"), form: norm.NFKD, want: []byte("ffi")},
		{name: "NFC composes combining marks", cond: true, p: []byte("e\u0301"), form: norm.NFC, want: []byte("\u00e9")},
		{name: "NFD decomposes precomposed", cond: true, p: []byte("\u00e9"), form: norm.NFD, want: []byte("e\u0301")},
		{name: "empty input", cond: true, p: []byte{}, form: norm.NFKC, want: []byte{}},
		{name: "cond false leaves buffer empty", cond: false, p: []byte("\uFB03"), form: norm.NFKC, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteNormalizeIf(tt.cond, tt.p, tt.form)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteNormalizeIf(%v, %q, %d)", tt.want, cb.Bytes(), tt.cond, tt.p, tt.form)
		})
	}
}

func TestChainWriteNormalizeStringIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		s    string
		form norm.Form
		want []byte
	}{
		{name: "NFC keeps ligature", cond: true, s: "\uFB03", form: norm.NFC, want: []byte("\uFB03")},
		{name: "NFD keeps ligature", cond: true, s: "\uFB03", form: norm.NFD, want: []byte("\uFB03")},
		{name: "NFKC expands ligature", cond: true, s: "\uFB03", form: norm.NFKC, want: []byte("ffi")},
		{name: "NFKD expands ligature", cond: true, s: "\uFB03", form: norm.NFKD, want: []byte("ffi")},
		{name: "NFC composes combining marks", cond: true, s: "e\u0301", form: norm.NFC, want: []byte("\u00e9")},
		{name: "NFD decomposes precomposed", cond: true, s: "\u00e9", form: norm.NFD, want: []byte("e\u0301")},
		{name: "empty input", cond: true, s: "", form: norm.NFKC, want: []byte{}},
		{name: "cond false leaves buffer empty", cond: false, s: "\uFB03", form: norm.NFKC, want: []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteNormalizeStringIf(tt.cond, tt.s, tt.form)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Chain.WriteNormalizeStringIf(%v, %q, %d)", tt.want, cb.Bytes(), tt.cond, tt.s, tt.form)
		})
	}
}

func TestChainWriteHexIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		p    []byte
	}{
		{name: "cond true empty input", cond: true, p: []byte{}},
		{name: "cond true nil input", cond: true, p: nil},
		{name: "cond true single zero byte", cond: true, p: []byte{0x00}},
		{name: "cond true single max byte", cond: true, p: []byte{0xFF}},
		{name: "cond true ascii", cond: true, p: []byte("abc")},
		{name: "cond true odd length", cond: true, p: []byte{0x01, 0x02, 0x03}},
		{name: "cond true unicode", cond: true, p: []byte("世界")},
		{name: "cond false leaves buffer empty", cond: false, p: []byte("abc")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := []byte{}
			if tt.cond {
				want = []byte(hex.EncodeToString(tt.p))
			}
			cb := NewChainSize(128)
			cb.WriteHexIf(tt.cond, tt.p)
			assert.True(t, bytes.Equal(want, cb.Bytes()), "want %q, got %q, Chain.WriteHexIf(%v, %q)", want, cb.Bytes(), tt.cond, tt.p)
		})
	}
}

func TestChainWriteHexStringIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		s    string
	}{
		{name: "cond true empty string", cond: true, s: ""},
		{name: "cond true single zero byte", cond: true, s: "\x00"},
		{name: "cond true single max byte", cond: true, s: "\xff"},
		{name: "cond true ascii", cond: true, s: "abc"},
		{name: "cond true odd length", cond: true, s: "\x01\x02\x03"},
		{name: "cond true unicode", cond: true, s: "世界"},
		{name: "cond false leaves buffer empty", cond: false, s: "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := []byte{}
			if tt.cond {
				want = []byte(hex.EncodeToString([]byte(tt.s)))
			}
			cb := NewChainSize(128)
			cb.WriteHexStringIf(tt.cond, tt.s)
			assert.True(t, bytes.Equal(want, cb.Bytes()), "want %q, got %q, Chain.WriteHexStringIf(%v, %q)", want, cb.Bytes(), tt.cond, tt.s)
		})
	}
}
