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

func TestAccumulativeWriteIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		p    []byte
		want []byte
	}{
		{name: "cond true writes bytes", cond: true, p: []byte("foobar"), want: []byte("foobar")},
		{name: "cond true writes empty", cond: true, p: []byte{}, want: nil},
		{name: "cond true writes nil", cond: true, p: nil, want: nil},
		{name: "cond false leaves empty", cond: false, p: []byte("foobar"), want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteIf(tt.cond, tt.p)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteIf(%v, %q)", tt.want, ab.Bytes(), tt.cond, tt.p)
		})
	}
}

func TestAccumulativeWriteByteIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		p    byte
		want []byte
	}{
		{name: "cond true writes byte", cond: true, p: 'a', want: []byte{'a'}},
		{name: "cond true writes zero byte", cond: true, p: 0x00, want: []byte{0x00}},
		{name: "cond true writes high byte", cond: true, p: 0xFF, want: []byte{0xFF}},
		{name: "cond false leaves empty", cond: false, p: 'a', want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteByteIf(tt.cond, tt.p)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteByteIf(%v, %d)", tt.want, ab.Bytes(), tt.cond, tt.p)
		})
	}
}

func TestAccumulativeWriteStringIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		s    string
		want []byte
	}{
		{name: "cond true writes string", cond: true, s: "foobar string", want: []byte("foobar string")},
		{name: "cond true writes empty string", cond: true, s: "", want: nil},
		{name: "cond true writes unicode string", cond: true, s: "Hello, 世界!", want: []byte("Hello, 世界!")},
		{name: "cond false leaves empty", cond: false, s: "foobar string", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteStringIf(tt.cond, tt.s)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteStringIf(%v, %q)", tt.want, ab.Bytes(), tt.cond, tt.s)
		})
	}
}

func TestAccumulativeWriteIntIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		i    int64
		want []byte
	}{
		{name: "zero cond true", cond: true, i: 0, want: []byte("0")},
		{name: "positive cond true", cond: true, i: 123456789, want: []byte("123456789")},
		{name: "negative cond true", cond: true, i: -123456789, want: []byte("-123456789")},
		{name: "max int64 cond true", cond: true, i: math.MaxInt64, want: []byte("9223372036854775807")},
		{name: "min int64 cond true", cond: true, i: math.MinInt64, want: []byte("-9223372036854775808")},
		{name: "zero cond false", cond: false, i: 0, want: nil},
		{name: "positive cond false", cond: false, i: 123456789, want: nil},
		{name: "negative cond false", cond: false, i: -123456789, want: nil},
		{name: "max int64 cond false", cond: false, i: math.MaxInt64, want: nil},
		{name: "min int64 cond false", cond: false, i: math.MinInt64, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteIntIf(tt.cond, tt.i)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteIntIf(%v, %d)", tt.want, ab.Bytes(), tt.cond, tt.i)
		})
	}
}

func TestAccumulativeWriteUintIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		u    uint64
		want []byte
	}{
		{name: "zero cond true", cond: true, u: 0, want: []byte("0")},
		{name: "small cond true", cond: true, u: 255, want: []byte("255")},
		{name: "max uint64 cond true", cond: true, u: math.MaxUint64, want: []byte("18446744073709551615")},
		{name: "zero cond false", cond: false, u: 0, want: nil},
		{name: "small cond false", cond: false, u: 255, want: nil},
		{name: "max uint64 cond false", cond: false, u: math.MaxUint64, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteUintIf(tt.cond, tt.u)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteUintIf(%v, %d)", tt.want, ab.Bytes(), tt.cond, tt.u)
		})
	}
}

func TestAccumulativeWriteFloatIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		f    float64
		want []byte
	}{
		{name: "zero cond true", cond: true, f: 0.0, want: []byte("0")},
		{name: "positive cond true", cond: true, f: 123.456, want: []byte("123.456")},
		{name: "negative cond true", cond: true, f: -123.456, want: []byte("-123.456")},
		{name: "pi cond true", cond: true, f: math.Pi, want: []byte("3.141592653589793")},
		{name: "zero cond false", cond: false, f: 0.0, want: nil},
		{name: "positive cond false", cond: false, f: 123.456, want: nil},
		{name: "negative cond false", cond: false, f: -123.456, want: nil},
		{name: "pi cond false", cond: false, f: math.Pi, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteFloatIf(tt.cond, tt.f)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteFloatIf(%v, %v)", tt.want, ab.Bytes(), tt.cond, tt.f)
		})
	}
}

func TestAccumulativeWriteBoolIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		v    bool
		want []byte
	}{
		{name: "cond true writes true", cond: true, v: true, want: []byte("true")},
		{name: "cond true writes false", cond: true, v: false, want: []byte("false")},
		{name: "cond false with true value leaves empty", cond: false, v: true, want: nil},
		{name: "cond false with false value leaves empty", cond: false, v: false, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteBoolIf(tt.cond, tt.v)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteBoolIf(%v, %v)", tt.want, ab.Bytes(), tt.cond, tt.v)
		})
	}
}

func TestAccumulativeWriteBinaryIf(t *testing.T) {
	type binStruct struct {
		A uint32
		B uint16
		C byte
	}

	tests := []struct {
		name  string
		cond  bool
		order binary.ByteOrder
		value any
		want  []byte
	}{
		{
			name:  "little endian cond true",
			cond:  true,
			order: binary.LittleEndian,
			value: &binStruct{A: 0x12345678, B: 0xABCD, C: 0xFF},
			want:  []byte{0x78, 0x56, 0x34, 0x12, 0xCD, 0xAB, 0xFF},
		},
		{
			name:  "big endian cond true",
			cond:  true,
			order: binary.BigEndian,
			value: &binStruct{A: 0x12345678, B: 0xABCD, C: 0xFF},
			want:  []byte{0x12, 0x34, 0x56, 0x78, 0xAB, 0xCD, 0xFF},
		},
		{
			name:  "little endian cond false",
			cond:  false,
			order: binary.LittleEndian,
			value: &binStruct{A: 0x12345678, B: 0xABCD, C: 0xFF},
			want:  nil,
		},
		{
			name:  "big endian cond false",
			cond:  false,
			order: binary.BigEndian,
			value: &binStruct{A: 0x12345678, B: 0xABCD, C: 0xFF},
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteBinaryIf(tt.cond, tt.order, tt.value)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteBinaryIf(%v, %v, %T)", tt.want, ab.Bytes(), tt.cond, tt.order, tt.value)
		})
	}
}

func TestAccumulativeWriteXIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		x    any
		want []byte
	}{
		{name: "int cond true", cond: true, x: 42, want: []byte("42")},
		{name: "string cond true", cond: true, x: "hello", want: []byte("hello")},
		{name: "int cond false", cond: false, x: 42, want: nil},
		{name: "string cond false", cond: false, x: "hello", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteXIf(tt.cond, tt.x)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteXIf(%v, %v)", tt.want, ab.Bytes(), tt.cond, tt.x)
		})
	}
}

func TestAccumulativeWriteMarshallerToIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		m    MarshallerTo
		want []byte
	}{
		{name: "cond true writes marshalled data", cond: true, m: stage, want: expectMT},
		{name: "cond false leaves empty", cond: false, m: stage, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteMarshallerToIf(tt.cond, tt.m)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteMarshallerToIf(%v, %T)", tt.want, ab.Bytes(), tt.cond, tt.m)
		})
	}
}

func TestAccumulativeWriteApplyFnIf(t *testing.T) {
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
		{name: "cond true applies fn to empty", cond: true, p: []byte{}, want: []byte("prefix:")},
		{name: "cond false leaves empty", cond: false, p: []byte("test"), want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteApplyFnIf(tt.cond, tt.p, addPrefix)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteApplyFnIf(%v, %q)", tt.want, ab.Bytes(), tt.cond, tt.p)
		})
	}
}

func TestAccumulativeWriteApplyFnNIf(t *testing.T) {
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
		{name: "cond true applies fn n times", cond: true, p: []byte("ab"), n: 3, want: []byte("abababababababab")},
		{name: "cond true applies fn zero times", cond: true, p: []byte("ab"), n: 0, want: nil},
		{name: "cond false leaves empty", cond: false, p: []byte("ab"), n: 3, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteApplyFnNIf(tt.cond, tt.p, duplicate, tt.n)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteApplyFnNIf(%v, %q, fn, %d)", tt.want, ab.Bytes(), tt.cond, tt.p, tt.n)
		})
	}
}

func TestAccumulativeWriteApplyFnStringIf(t *testing.T) {
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
		{name: "cond true applies fn to empty", cond: true, s: "", want: []byte("prefix:")},
		{name: "cond false leaves empty", cond: false, s: "test", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteApplyFnStringIf(tt.cond, tt.s, addPrefix)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteApplyFnStringIf(%v, %q)", tt.want, ab.Bytes(), tt.cond, tt.s)
		})
	}
}

func TestAccumulativeWriteApplyFnNStringIf(t *testing.T) {
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
		{name: "cond true applies fn n times", cond: true, s: "ab", n: 3, want: []byte("abababababababab")},
		{name: "cond true applies fn zero times", cond: true, s: "ab", n: 0, want: nil},
		{name: "cond false leaves empty", cond: false, s: "ab", n: 3, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteApplyFnNStringIf(tt.cond, tt.s, duplicate, tt.n)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteApplyFnNStringIf(%v, %q, fn, %d)", tt.want, ab.Bytes(), tt.cond, tt.s, tt.n)
		})
	}
}

func TestAccumulativeWriteTimeIf(t *testing.T) {
	tm := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)

	tests := []struct {
		name   string
		cond   bool
		format string
		t      time.Time
		want   []byte
	}{
		{name: "cond true formats time", cond: true, format: "%Y-%m-%d %H:%M:%S", t: tm, want: []byte("2024-01-15 10:30:45")},
		{name: "cond false leaves empty", cond: false, format: "%Y-%m-%d %H:%M:%S", t: tm, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteTimeIf(tt.cond, tt.format, tt.t)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteTimeIf(%v, %q, %v)", tt.want, ab.Bytes(), tt.cond, tt.format, tt.t)
		})
	}
}

func TestAccumulativeWriteULEB128If(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		v    uint64
		want []byte
	}{
		{name: "zero cond true", cond: true, v: 0, want: []byte{0x00}},
		{name: "127 cond true", cond: true, v: 127, want: []byte{0x7F}},
		{name: "128 cond true", cond: true, v: 128, want: []byte{0x80, 0x01}},
		{name: "300 cond true", cond: true, v: 300, want: []byte{0xAC, 0x02}},
		{name: "max uint64 cond true", cond: true, v: math.MaxUint64, want: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01}},
		{name: "zero cond false", cond: false, v: 0, want: nil},
		{name: "127 cond false", cond: false, v: 127, want: nil},
		{name: "128 cond false", cond: false, v: 128, want: nil},
		{name: "300 cond false", cond: false, v: 300, want: nil},
		{name: "max uint64 cond false", cond: false, v: math.MaxUint64, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteULEB128If(tt.cond, tt.v)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteULEB128If(%v, %d)", tt.want, ab.Bytes(), tt.cond, tt.v)
		})
	}
}

func TestAccumulativeWriteSLEB128If(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		v    int64
		want []byte
	}{
		{name: "zero cond true", cond: true, v: 0, want: []byte{0x00}},
		{name: "63 cond true", cond: true, v: 63, want: []byte{0x3F}},
		{name: "64 cond true", cond: true, v: 64, want: []byte{0xC0, 0x00}},
		{name: "-1 cond true", cond: true, v: -1, want: []byte{0x7F}},
		{name: "-64 cond true", cond: true, v: -64, want: []byte{0x40}},
		{name: "-100 cond true", cond: true, v: -100, want: []byte{0x9C, 0x7F}},
		{name: "zero cond false", cond: false, v: 0, want: nil},
		{name: "63 cond false", cond: false, v: 63, want: nil},
		{name: "64 cond false", cond: false, v: 64, want: nil},
		{name: "-1 cond false", cond: false, v: -1, want: nil},
		{name: "-64 cond false", cond: false, v: -64, want: nil},
		{name: "-100 cond false", cond: false, v: -100, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteSLEB128If(tt.cond, tt.v)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteSLEB128If(%v, %d)", tt.want, ab.Bytes(), tt.cond, tt.v)
		})
	}
}

func TestAccumulativeWriteNormalizeIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		p    []byte
		form norm.Form
		want []byte
	}{
		{name: "NFC cond true", cond: true, p: []byte("\u00E9\uFB03"), form: norm.NFC, want: []byte("\u00E9\uFB03")},
		{name: "NFD cond true", cond: true, p: []byte("\u00E9\uFB03"), form: norm.NFD, want: []byte("e\u0301\uFB03")},
		{name: "NFKC cond true", cond: true, p: []byte("\u00E9\uFB03"), form: norm.NFKC, want: []byte("\u00E9ffi")},
		{name: "NFKD cond true", cond: true, p: []byte("\u00E9\uFB03"), form: norm.NFKD, want: []byte("e\u0301ffi")},
		{name: "NFKC cond true expands ligature", cond: true, p: []byte("\uFB03"), form: norm.NFKC, want: []byte("ffi")},
		{name: "NFKD cond true expands ligature", cond: true, p: []byte("\uFB03"), form: norm.NFKD, want: []byte("ffi")},
		{name: "NFC cond false", cond: false, p: []byte("\u00E9\uFB03"), form: norm.NFC, want: nil},
		{name: "NFD cond false", cond: false, p: []byte("\u00E9\uFB03"), form: norm.NFD, want: nil},
		{name: "NFKC cond false", cond: false, p: []byte("\u00E9\uFB03"), form: norm.NFKC, want: nil},
		{name: "NFKD cond false", cond: false, p: []byte("\u00E9\uFB03"), form: norm.NFKD, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteNormalizeIf(tt.cond, tt.p, tt.form)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteNormalizeIf(%v, %q, %d)", tt.want, ab.Bytes(), tt.cond, tt.p, tt.form)
		})
	}
}

func TestAccumulativeWriteNormalizeStringIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		s    string
		form norm.Form
		want []byte
	}{
		{name: "NFC cond true", cond: true, s: "\u00E9\uFB03", form: norm.NFC, want: []byte("\u00E9\uFB03")},
		{name: "NFD cond true", cond: true, s: "\u00E9\uFB03", form: norm.NFD, want: []byte("e\u0301\uFB03")},
		{name: "NFKC cond true", cond: true, s: "\u00E9\uFB03", form: norm.NFKC, want: []byte("\u00E9ffi")},
		{name: "NFKD cond true", cond: true, s: "\u00E9\uFB03", form: norm.NFKD, want: []byte("e\u0301ffi")},
		{name: "NFKC cond true expands ligature", cond: true, s: "\uFB03", form: norm.NFKC, want: []byte("ffi")},
		{name: "NFKD cond true expands ligature", cond: true, s: "\uFB03", form: norm.NFKD, want: []byte("ffi")},
		{name: "NFC cond false", cond: false, s: "\u00E9\uFB03", form: norm.NFC, want: nil},
		{name: "NFD cond false", cond: false, s: "\u00E9\uFB03", form: norm.NFD, want: nil},
		{name: "NFKC cond false", cond: false, s: "\u00E9\uFB03", form: norm.NFKC, want: nil},
		{name: "NFKD cond false", cond: false, s: "\u00E9\uFB03", form: norm.NFKD, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteNormalizeStringIf(tt.cond, tt.s, tt.form)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteNormalizeStringIf(%v, %q, %d)", tt.want, ab.Bytes(), tt.cond, tt.s, tt.form)
		})
	}
}

func TestAccumulativeWriteHexIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		p    []byte
		want []byte
	}{
		{name: "empty cond true", cond: true, p: []byte{}, want: nil},
		{name: "single byte cond true", cond: true, p: []byte{0x41}, want: []byte("41")},
		{name: "zero byte cond true", cond: true, p: []byte{0x00}, want: []byte("00")},
		{name: "high byte cond true", cond: true, p: []byte{0xFF}, want: []byte("ff")},
		{name: "multiple bytes cond true", cond: true, p: []byte{0x01, 0x02, 0xAB}, want: []byte("0102ab")},
		{name: "empty cond false", cond: false, p: []byte{}, want: nil},
		{name: "single byte cond false", cond: false, p: []byte{0x41}, want: nil},
		{name: "zero byte cond false", cond: false, p: []byte{0x00}, want: nil},
		{name: "high byte cond false", cond: false, p: []byte{0xFF}, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteHexIf(tt.cond, tt.p)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteHexIf(%v, %q)", tt.want, ab.Bytes(), tt.cond, tt.p)
			if tt.cond {
				assert.Equal(t, hex.EncodeToString(tt.p), ab.String(), "WriteHexIf hex reference (%q)", tt.p)
			}
		})
	}
}

func TestAccumulativeWriteHexStringIf(t *testing.T) {
	tests := []struct {
		name string
		cond bool
		s    string
		want []byte
	}{
		{name: "empty cond true", cond: true, s: "", want: nil},
		{name: "single byte cond true", cond: true, s: "A", want: []byte("41")},
		{name: "zero byte cond true", cond: true, s: "\x00", want: []byte("00")},
		{name: "high byte cond true", cond: true, s: "\xff", want: []byte("ff")},
		{name: "multiple bytes cond true", cond: true, s: "\x01\x02\xab", want: []byte("0102ab")},
		{name: "empty cond false", cond: false, s: "", want: nil},
		{name: "single byte cond false", cond: false, s: "A", want: nil},
		{name: "zero byte cond false", cond: false, s: "\x00", want: nil},
		{name: "high byte cond false", cond: false, s: "\xff", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteHexStringIf(tt.cond, tt.s)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteHexStringIf(%v, %q)", tt.want, ab.Bytes(), tt.cond, tt.s)
			if tt.cond {
				assert.Equal(t, hex.EncodeToString([]byte(tt.s)), ab.String(), "WriteHexStringIf hex reference (%q)", tt.s)
			}
		})
	}
}
