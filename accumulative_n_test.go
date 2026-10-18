package bytebuf

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/unicode/norm"
)

func TestAccumulativeWriteN(t *testing.T) {
	tests := []struct {
		name    string
		p       []byte
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", p: []byte("ab"), n: 0, want: []byte{}},
		{name: "n=1 single copy", p: []byte("ab"), n: 1, want: []byte("ab")},
		{name: "n=2 two copies", p: []byte("ab"), n: 2, want: []byte("abab")},
		{name: "n=3 three copies", p: []byte("ab"), n: 3, want: []byte("ababab")},
		{name: "negative n writes nothing", p: []byte("ab"), n: -1, want: []byte{}},
		{name: "empty input with n=1 writes nothing", p: nil, n: 1, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteN(tt.p, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteN(%q, %d)", tt.want, ab.Bytes(), tt.p, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteByteN(t *testing.T) {
	tests := []struct {
		name    string
		p       byte
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", p: 'x', n: 0, want: []byte{}},
		{name: "n=1 single byte", p: 'x', n: 1, want: []byte("x")},
		{name: "n=3 three bytes", p: 'x', n: 3, want: []byte("xxx")},
		{name: "negative n writes nothing", p: 'x', n: -1, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteByteN(tt.p, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteByteN(%q, %d)", tt.want, ab.Bytes(), tt.p, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteStringN(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", s: "ab", n: 0, want: []byte{}},
		{name: "n=1 single copy", s: "ab", n: 1, want: []byte("ab")},
		{name: "n=2 two copies", s: "ab", n: 2, want: []byte("abab")},
		{name: "n=3 three copies", s: "ab", n: 3, want: []byte("ababab")},
		{name: "negative n writes nothing", s: "ab", n: -1, want: []byte{}},
		{name: "empty string with n=1 writes nothing", s: "", n: 1, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteStringN(tt.s, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteStringN(%q, %d)", tt.want, ab.Bytes(), tt.s, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteIntN(t *testing.T) {
	tests := []struct {
		name    string
		i       int64
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", i: 123, n: 0, want: []byte{}},
		{name: "n=1 single copy", i: 123, n: 1, want: []byte("123")},
		{name: "n=2 two copies", i: 123, n: 2, want: []byte("123123")},
		{name: "n=3 three copies", i: 123, n: 3, want: []byte("123123123")},
		{name: "negative n writes nothing", i: 123, n: -1, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteIntN(tt.i, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteIntN(%d, %d)", tt.want, ab.Bytes(), tt.i, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteIntBaseN(t *testing.T) {
	tests := []struct {
		name    string
		i       int64
		base    int
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "base 16 n=0 writes nothing", i: 255, base: 16, n: 0, want: []byte{}},
		{name: "base 16 n=1 single copy", i: 255, base: 16, n: 1, want: []byte("ff")},
		{name: "base 16 n=2 two copies", i: 255, base: 16, n: 2, want: []byte("ffff")},
		{name: "base 16 n=3 three copies", i: 255, base: 16, n: 3, want: []byte("ffffff")},
		{name: "base 16 negative n writes nothing", i: 255, base: 16, n: -1, want: []byte{}},
		{name: "base 2 n=1", i: 255, base: 2, n: 1, want: []byte("11111111")},
		{name: "base 8 n=1", i: 255, base: 8, n: 1, want: []byte("377")},
		{name: "base 10 n=1", i: 255, base: 10, n: 1, want: []byte("255")},
		{name: "base 36 n=1", i: 255, base: 36, n: 1, want: []byte("73")},
		{name: "negative value base 10", i: -255, base: 10, n: 1, want: []byte("-255")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteIntBaseN(tt.i, tt.base, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteIntBaseN(%d, %d, %d)", tt.want, ab.Bytes(), tt.i, tt.base, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteUintN(t *testing.T) {
	tests := []struct {
		name    string
		u       uint64
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", u: 123, n: 0, want: []byte{}},
		{name: "n=1 single copy", u: 123, n: 1, want: []byte("123")},
		{name: "n=2 two copies", u: 123, n: 2, want: []byte("123123")},
		{name: "n=3 three copies", u: 123, n: 3, want: []byte("123123123")},
		{name: "negative n writes nothing", u: 123, n: -1, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteUintN(tt.u, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteUintN(%d, %d)", tt.want, ab.Bytes(), tt.u, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteUintBaseN(t *testing.T) {
	tests := []struct {
		name    string
		u       uint64
		base    int
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "base 16 n=0 writes nothing", u: 255, base: 16, n: 0, want: []byte{}},
		{name: "base 16 n=1 single copy", u: 255, base: 16, n: 1, want: []byte("ff")},
		{name: "base 16 n=2 two copies", u: 255, base: 16, n: 2, want: []byte("ffff")},
		{name: "base 16 n=3 three copies", u: 255, base: 16, n: 3, want: []byte("ffffff")},
		{name: "base 16 negative n writes nothing", u: 255, base: 16, n: -1, want: []byte{}},
		{name: "base 2 n=1", u: 255, base: 2, n: 1, want: []byte("11111111")},
		{name: "base 8 n=1", u: 255, base: 8, n: 1, want: []byte("377")},
		{name: "base 10 n=1", u: 255, base: 10, n: 1, want: []byte("255")},
		{name: "base 36 n=1", u: 255, base: 36, n: 1, want: []byte("73")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteUintBaseN(tt.u, tt.base, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteUintBaseN(%d, %d, %d)", tt.want, ab.Bytes(), tt.u, tt.base, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteFloatN(t *testing.T) {
	tests := []struct {
		name    string
		f       float64
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", f: 1.5, n: 0, want: []byte{}},
		{name: "n=1 single copy", f: 1.5, n: 1, want: []byte("1.5")},
		{name: "n=2 two copies", f: 1.5, n: 2, want: []byte("1.51.5")},
		{name: "n=3 three copies", f: 1.5, n: 3, want: []byte("1.51.51.5")},
		{name: "negative n writes nothing", f: 1.5, n: -1, want: []byte{}},
		{name: "pi n=1", f: math.Pi, n: 1, want: []byte("3.141592653589793")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteFloatN(tt.f, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteFloatN(%v, %d)", tt.want, ab.Bytes(), tt.f, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteBoolN(t *testing.T) {
	tests := []struct {
		name    string
		v       bool
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "true n=0 writes nothing", v: true, n: 0, want: []byte{}},
		{name: "true n=1 single copy", v: true, n: 1, want: []byte("true")},
		{name: "true n=2 two copies", v: true, n: 2, want: []byte("truetrue")},
		{name: "true n=3 three copies", v: true, n: 3, want: []byte("truetruetrue")},
		{name: "true negative n writes nothing", v: true, n: -1, want: []byte{}},
		{name: "false n=0 writes nothing", v: false, n: 0, want: []byte{}},
		{name: "false n=1 single copy", v: false, n: 1, want: []byte("false")},
		{name: "false n=2 two copies", v: false, n: 2, want: []byte("falsefalse")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteBoolN(tt.v, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteBoolN(%v, %d)", tt.want, ab.Bytes(), tt.v, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteBinaryN(t *testing.T) {
	type binStruct struct {
		A uint16
		B uint32
	}
	base := []byte{0x02, 0x01, 0x06, 0x05, 0x04, 0x03}

	tests := []struct {
		name    string
		order   binary.ByteOrder
		x       any
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", order: binary.LittleEndian, x: binStruct{A: 0x0102, B: 0x03040506}, n: 0, want: []byte{}},
		{name: "n=1 single copy", order: binary.LittleEndian, x: binStruct{A: 0x0102, B: 0x03040506}, n: 1, want: bytes.Repeat(base, 1)},
		{name: "n=2 two copies", order: binary.LittleEndian, x: binStruct{A: 0x0102, B: 0x03040506}, n: 2, want: bytes.Repeat(base, 2)},
		{name: "n=3 three copies", order: binary.LittleEndian, x: binStruct{A: 0x0102, B: 0x03040506}, n: 3, want: bytes.Repeat(base, 3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteBinaryN(tt.order, tt.x, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteBinaryN(%d)", tt.want, ab.Bytes(), tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteXN(t *testing.T) {
	tests := []struct {
		name    string
		x       any
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", x: 123, n: 0, want: []byte{}},
		{name: "n=1 single copy", x: 123, n: 1, want: []byte("123")},
		{name: "n=2 two copies", x: 123, n: 2, want: []byte("123123")},
		{name: "n=3 three copies", x: 123, n: 3, want: []byte("123123123")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteXN(tt.x, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteXN(%v, %d)", tt.want, ab.Bytes(), tt.x, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteTimeN(t *testing.T) {
	tm := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)
	base := []byte("2024-01-15")

	tests := []struct {
		name    string
		format  string
		t       time.Time
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", format: "%Y-%m-%d", t: tm, n: 0, want: []byte{}},
		{name: "n=1 single copy", format: "%Y-%m-%d", t: tm, n: 1, want: bytes.Repeat(base, 1)},
		{name: "n=2 two copies", format: "%Y-%m-%d", t: tm, n: 2, want: bytes.Repeat(base, 2)},
		{name: "n=3 three copies", format: "%Y-%m-%d", t: tm, n: 3, want: bytes.Repeat(base, 3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteTimeN(tt.format, tt.t, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteTimeN(%q, %d)", tt.want, ab.Bytes(), tt.format, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteULEB128N(t *testing.T) {
	// 300 encodes to two bytes: 0xAC 0x02.
	base := []byte{0xAC, 0x02}

	tests := []struct {
		name    string
		v       uint64
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", v: 300, n: 0, want: []byte{}},
		{name: "n=1 single copy", v: 300, n: 1, want: bytes.Repeat(base, 1)},
		{name: "n=2 two copies", v: 300, n: 2, want: bytes.Repeat(base, 2)},
		{name: "n=3 three copies", v: 300, n: 3, want: bytes.Repeat(base, 3)},
		{name: "negative n writes nothing", v: 300, n: -1, want: []byte{}},
		{name: "max uint64 n=1", v: math.MaxUint64, n: 1, want: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteULEB128N(tt.v, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteULEB128N(%d, %d)", tt.want, ab.Bytes(), tt.v, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteSLEB128N(t *testing.T) {
	// -100 encodes to two bytes: 0x9C 0x7F.
	base := []byte{0x9C, 0x7F}

	tests := []struct {
		name    string
		v       int64
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", v: -100, n: 0, want: []byte{}},
		{name: "n=1 single copy", v: -100, n: 1, want: bytes.Repeat(base, 1)},
		{name: "n=2 two copies", v: -100, n: 2, want: bytes.Repeat(base, 2)},
		{name: "n=3 three copies", v: -100, n: 3, want: bytes.Repeat(base, 3)},
		{name: "negative n writes nothing", v: -100, n: -1, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteSLEB128N(tt.v, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteSLEB128N(%d, %d)", tt.want, ab.Bytes(), tt.v, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteNormalizeN(t *testing.T) {
	tests := []struct {
		name    string
		p       []byte
		form    norm.Form
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", p: []byte("test"), form: norm.NFC, n: 0, want: []byte{}},
		{name: "n=1 single copy", p: []byte("test"), form: norm.NFC, n: 1, want: []byte("test")},
		{name: "n=2 two copies", p: []byte("test"), form: norm.NFC, n: 2, want: []byte("testtest")},
		{name: "n=3 three copies", p: []byte("test"), form: norm.NFC, n: 3, want: []byte("testtesttest")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteNormalizeN(tt.p, tt.form, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteNormalizeN(%q, %d)", tt.want, ab.Bytes(), tt.p, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteNormalizeStringN(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		form    norm.Form
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", s: "test", form: norm.NFC, n: 0, want: []byte{}},
		{name: "n=1 single copy", s: "test", form: norm.NFC, n: 1, want: []byte("test")},
		{name: "n=2 two copies", s: "test", form: norm.NFC, n: 2, want: []byte("testtest")},
		{name: "n=3 three copies", s: "test", form: norm.NFC, n: 3, want: []byte("testtesttest")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteNormalizeStringN(tt.s, tt.form, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteNormalizeStringN(%q, %d)", tt.want, ab.Bytes(), tt.s, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteHexN(t *testing.T) {
	// hex(0x00, 0xFF) == "00ff".
	base := []byte("00ff")

	tests := []struct {
		name    string
		p       []byte
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", p: []byte{0x00, 0xFF}, n: 0, want: []byte{}},
		{name: "n=1 single copy", p: []byte{0x00, 0xFF}, n: 1, want: bytes.Repeat(base, 1)},
		{name: "n=2 two copies", p: []byte{0x00, 0xFF}, n: 2, want: bytes.Repeat(base, 2)},
		{name: "n=3 three copies", p: []byte{0x00, 0xFF}, n: 3, want: bytes.Repeat(base, 3)},
		{name: "negative n writes nothing", p: []byte{0x00, 0xFF}, n: -1, want: []byte{}},
		{name: "empty input n=1 writes nothing", p: nil, n: 1, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteHexN(tt.p, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteHexN(% x, %d)", tt.want, ab.Bytes(), tt.p, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}

func TestAccumulativeWriteHexStringN(t *testing.T) {
	// hex("00ff") == hex of the ASCII bytes 0x30 0x30 0x66 0x66 == "30306666".
	base := []byte("30306666")

	tests := []struct {
		name    string
		s       string
		n       int
		want    []byte
		wantErr bool
	}{
		{name: "n=0 writes nothing", s: "00ff", n: 0, want: []byte{}},
		{name: "n=1 single copy", s: "00ff", n: 1, want: bytes.Repeat(base, 1)},
		{name: "n=2 two copies", s: "00ff", n: 2, want: bytes.Repeat(base, 2)},
		{name: "n=3 three copies", s: "00ff", n: 3, want: bytes.Repeat(base, 3)},
		{name: "negative n writes nothing", s: "00ff", n: -1, want: []byte{}},
		{name: "empty input n=1 writes nothing", s: "", n: 1, want: []byte{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteHexStringN(tt.s, tt.n)

			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, WriteHexStringN(%q, %d)", tt.want, ab.Bytes(), tt.s, tt.n)
			if tt.wantErr {
				assert.Error(t, ab.Error())
			}
		})
	}
}
