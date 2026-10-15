package bytebuf

import (
	"bytes"
	"encoding/binary"
	"math"
	"net/url"
	"strconv"
	"testing"

	"github.com/koykov/byteconv"
	"github.com/koykov/x2bytes"
)

func TestChain(t *testing.T) {
	t.Run("write", func(t *testing.T) {
		cb := NewChainSize(128)
		cb.Write(stage.b).WriteByte('-').
			WriteString(stage.s).WriteByte('-').
			WriteInt(stage.i).WriteByte('-').
			WriteUint(stage.u).WriteByte('-').
			WriteFloat(stage.f)

		if !bytes.Equal(cb.Bytes(), expectWS) {
			t.Error("Chain.Write*: mismatch result and expectation")
		}
	})
	t.Run("apply fn", func(t *testing.T) {
		cb := NewChainSize(128)
		cb.WriteString("foo").
			WriteApplyFnString("?q=front&p=1", func(dst, p []byte) []byte {
				p1 := url.QueryEscape(byteconv.B2S(p))
				dst = append(dst, p1...)
				return dst
			}).
			WriteString("bar")
		if cb.String() != "foo%3Fq%3Dfront%26p%3D1bar" {
			t.Error("Chain.WriteApplyFn*: mismatch result and expectation")
		}
	})
	t.Run("apply fn N", func(t *testing.T) {
		cb := NewChainSize(128)
		cb.WriteString("foo").
			WriteApplyFnNString("?q=front&p=1", func(dst, p []byte) []byte {
				p1 := url.QueryEscape(byteconv.B2S(p))
				dst = append(dst, p1...)
				return dst
			}, 3).
			WriteString("bar")
		if cb.String() != "foo%25253Fq%25253Dfront%252526p%25253D1bar" {
			t.Error("Chain.WriteApplyFn*: mismatch result and expectation")
		}
	})
	t.Run("reduce", func(t *testing.T) {
		cb := Chain{}
		cb.WriteString("0x0000").
			WriteIntBase(100, 16)
		hex := cb.Bytes()[cb.Len()-2:]
		cb.Reduce(4).Write(hex)
		if cb.String() != "0x0064" {
			t.FailNow()
		}
	})
	t.Run("format", func(t *testing.T) {
		cb := Chain{}
		cb.WriteFormat("int %d; float %f; string %s", 513, 3.1415, "foobar")
		if cb.String() != "int 513; float 3.141500; string foobar" {
			t.FailNow()
		}
	})
	t.Run("bytes copy", func(t *testing.T) {
		t.Run("empty buffer", func(t *testing.T) {
			cb := NewChainSize(128)
			if cp := cb.BytesCopy(); len(cp) != 0 {
				t.Errorf("BytesCopy() = %v, want empty", cp)
			}
		})
		t.Run("copy is independent of original", func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteString("original")
			cp := cb.BytesCopy()
			if !bytes.Equal(cp, []byte("original")) {
				t.Fatalf("BytesCopy() = %q, want %q", cp, "original")
			}
			// Mutating the copy must not leak into the original buffer.
			cp[0] = 'X'
			if !bytes.Equal(cb.Bytes(), []byte("original")) {
				t.Errorf("mutating copy changed original: got %q, want %q", cb.Bytes(), "original")
			}
		})
	})
	t.Run("string copy", func(t *testing.T) {
		t.Run("empty buffer", func(t *testing.T) {
			cb := NewChainSize(128)
			if s := cb.StringCopy(); s != "" {
				t.Errorf("StringCopy() = %q, want empty", s)
			}
		})
		t.Run("returns equal string", func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteString("hello copy")
			if s := cb.StringCopy(); s != "hello copy" {
				t.Errorf("StringCopy() = %q, want %q", s, "hello copy")
			}
		})
	})
	t.Run("write uint in base", func(t *testing.T) {
		tests := []struct {
			name     string
			value    uint64
			base     int
			expected string
		}{
			{name: "base 2", value: 255, base: 2, expected: "11111111"},
			{name: "base 8", value: 255, base: 8, expected: "377"},
			{name: "base 10", value: 255, base: 10, expected: "255"},
			{name: "base 16", value: 255, base: 16, expected: "ff"},
			{name: "base 36", value: 255, base: 36, expected: "73"},
			{name: "max uint64 base 16", value: math.MaxUint64, base: 16, expected: "ffffffffffffffff"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				cb := NewChainSize(128)
				cb.WriteUintBase(tt.value, tt.base)
				if cb.String() != tt.expected {
					t.Errorf("WriteUintBase(%d, %d) = %q, want %q", tt.value, tt.base, cb.String(), tt.expected)
				}
			})
		}
	})
	t.Run("write binary", func(t *testing.T) {
		type binStruct struct {
			A uint32
			B uint16
			C byte
		}
		tests := []struct {
			name     string
			order    binary.ByteOrder
			expected []byte
		}{
			{
				name:     "little endian struct",
				order:    binary.LittleEndian,
				expected: []byte{0x78, 0x56, 0x34, 0x12, 0xCD, 0xAB, 0xFF},
			},
			{
				name:     "big endian struct",
				order:    binary.BigEndian,
				expected: []byte{0x12, 0x34, 0x56, 0x78, 0xAB, 0xCD, 0xFF},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				cb := NewChainSize(128)
				cb.WriteBinary(tt.order, &binStruct{A: 0x12345678, B: 0xABCD, C: 0xFF})
				if !bytes.Equal(cb.Bytes(), tt.expected) {
					t.Errorf("WriteBinary() = %v, want %v", cb.Bytes(), tt.expected)
				}
			})
		}
	})
	t.Run("write x", func(t *testing.T) {
		type chainXStruct struct {
			Name string
		}
		x2bytes.RegisterToBytesFn(func(dst []byte, val any, _ ...any) ([]byte, error) {
			if x, ok := val.(chainXStruct); ok {
				dst = append(dst, '[')
				dst = append(dst, x.Name...)
				dst = append(dst, ']')
				return dst, nil
			}
			return dst, x2bytes.ErrUnknownType
		})
		tests := []struct {
			name     string
			value    any
			expected string
		}{
			{name: "int", value: 42, expected: "42"},
			{name: "string", value: "hello", expected: "hello"},
			{name: "custom struct", value: chainXStruct{Name: "Alice"}, expected: "[Alice]"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				cb := NewChainSize(128)
				cb.WriteX(tt.value)
				if cb.String() != tt.expected {
					t.Errorf("WriteX(%v) = %q, want %q", tt.value, cb.String(), tt.expected)
				}
			})
		}
	})
	t.Run("to writer", func(t *testing.T) {
		cb := NewChainSize(128)
		cb.WriteString("foo")
		w := cb.ToWriter()
		if _, err := w.WriteString("bar"); err != nil {
			t.Fatalf("Writer.WriteString() error = %v", err)
		}
		if cb.String() != "foobar" {
			t.Errorf("Chain after Writer.WriteString() = %q, want %q", cb.String(), "foobar")
		}
	})
	t.Run("to reader", func(t *testing.T) {
		cb := NewChainSize(128)
		cb.WriteString("hello world")
		r := cb.ToReader()
		buf := make([]byte, 5)
		n, err := r.Read(buf)
		if err != nil {
			t.Fatalf("Reader.Read() error = %v", err)
		}
		if n != len(buf) {
			t.Fatalf("Reader.Read() read %d bytes, want %d", n, len(buf))
		}
		if !bytes.Equal(buf, []byte("hello")) {
			t.Errorf("Reader.Read() = %q, want %q", buf, "hello")
		}
	})
}

func BenchmarkChain(b *testing.B) {
	b.Run("write", func(b *testing.B) {
		cb := NewChainSize(128)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			cb.Reset().
				Write(stage.b).WriteByte('-').
				WriteString(stage.s).WriteByte('-').
				WriteInt(stage.i).WriteByte('-').
				WriteUint(stage.u).WriteByte('-').
				WriteFloat(stage.f)

			if !bytes.Equal(cb.Bytes(), expectWS) {
				b.Error("Chain.Write*: mismatch result and expectation")
			}
		}
	})
	b.Run("slice append", func(b *testing.B) {
		var buf []byte
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf = buf[:0]
			buf = append(buf, stage.b...)
			buf = append(buf, '-')
			buf = append(buf, stage.s...)
			buf = append(buf, '-')
			buf = append(buf, strconv.Itoa(int(stage.i))...)
			buf = append(buf, '-')
			buf = append(buf, strconv.Itoa(int(stage.u))...)
			buf = append(buf, '-')
			buf = append(buf, strconv.FormatFloat(stage.f, 'f', -1, 64)...)

			if !bytes.Equal(buf, expectWS) {
				b.Error("ByteArray: mismatch result and expectation")
			}
		}
	})
	b.Run("reduce", func(b *testing.B) {
		b.ReportAllocs()
		cb := Chain{}
		for i := 0; i < b.N; i++ {
			cb.Reset()
			cb.WriteString("0x0000").
				WriteIntBase(100, 16)
			hex := cb.Bytes()[cb.Len()-2:]
			cb.Reduce(4).Write(hex)
			if cb.String() != "0x0064" {
				b.FailNow()
			}
		}
	})
	b.Run("format", func(b *testing.B) {
		b.ReportAllocs()
		cb := Chain{}
		for i := 0; i < b.N; i++ {
			cb.Reset().WriteFormat("int %d; float %f; string %s", 513, 3.1415, "foobar")
		}
	})
}
