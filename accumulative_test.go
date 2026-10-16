package bytebuf

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
	"time"

	"github.com/koykov/x2bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccumulative(t *testing.T) {
	t.Run("write", func(t *testing.T) {
		ab := NewAccumulativeSize(128)
		ab.WriteString("foobar").
			StakeOut().
			Write(stage.b).WriteByte('-').
			WriteString(stage.s).WriteByte('-').
			WriteInt(stage.i).WriteByte('-').
			WriteUint(stage.u).WriteByte('-').
			WriteFloat(stage.f)

		assert.True(t, bytes.Equal(expectWS, ab.StakedBytes()), "want %q, got %q", expectWS, ab.StakedBytes())
	})

	t.Run("bytes copy", func(t *testing.T) {
		t.Run("empty buffer", func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			cp := ab.BytesCopy()
			assert.Equal(t, 0, len(cp))
		})
		t.Run("copy is independent of original", func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteString("original")
			cp := ab.BytesCopy()
			require.True(t, bytes.Equal([]byte("original"), cp), "want %q, got %q", []byte("original"), cp)
			cp[0] = 'X'
			assert.True(t, bytes.Equal([]byte("original"), ab.Bytes()), "want %q, got %q", []byte("original"), ab.Bytes())
		})
	})

	t.Run("string copy", func(t *testing.T) {
		t.Run("empty buffer", func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			s := ab.StringCopy()
			assert.Equal(t, "", s)
		})
		t.Run("returns equal string", func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteString("hello copy")
			s := ab.StringCopy()
			assert.Equal(t, "hello copy", s)
		})
	})

	t.Run("write rune", func(t *testing.T) {
		tests := []struct {
			name     string
			r        rune
			expected string
		}{
			{name: "ascii", r: 'a', expected: "a"},
			{name: "three byte rune", r: '世', expected: "世"},
			{name: "four byte rune", r: '😀', expected: "😀"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ab := NewAccumulativeSize(16)
				ab.WriteRune(tt.r)
				assert.Equal(t, tt.expected, ab.String())
			})
		}
	})

	t.Run("write int in base", func(t *testing.T) {
		tests := []struct {
			name     string
			value    int64
			base     int
			expected string
		}{
			{name: "base 2", value: 255, base: 2, expected: "11111111"},
			{name: "base 8", value: 255, base: 8, expected: "377"},
			{name: "base 10", value: 255, base: 10, expected: "255"},
			{name: "base 16", value: 255, base: 16, expected: "ff"},
			{name: "base 36", value: 255, base: 36, expected: "73"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ab := NewAccumulativeSize(16)
				ab.WriteIntBase(tt.value, tt.base)
				assert.Equal(t, tt.expected, ab.String())
			})
		}
	})

	t.Run("write uint in base", func(t *testing.T) {
		tests := []struct {
			name     string
			value    uint64
			base     int
			expected string
		}{
			{name: "base 16", value: 255, base: 16, expected: "ff"},
			{name: "base 2", value: 255, base: 2, expected: "11111111"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ab := NewAccumulativeSize(16)
				ab.WriteUintBase(tt.value, tt.base)
				assert.Equal(t, tt.expected, ab.String())
			})
		}
	})

	t.Run("write bool", func(t *testing.T) {
		tests := []struct {
			name     string
			value    bool
			expected string
		}{
			{name: "true", value: true, expected: "true"},
			{name: "false", value: false, expected: "false"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ab := NewAccumulativeSize(16)
				ab.WriteBool(tt.value)
				assert.Equal(t, tt.expected, ab.String())
			})
		}
	})

	t.Run("write format", func(t *testing.T) {
		tests := []struct {
			name     string
			format   string
			args     []any
			expected string
		}{
			{name: "int", format: "int %d", args: []any{42}, expected: "int 42"},
			{name: "float", format: "float %f", args: []any{3.14}, expected: "float 3.140000"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ab := NewAccumulativeSize(32)
				ab.WriteFormat(tt.format, tt.args...)
				assert.Equal(t, tt.expected, ab.String())
			})
		}
	})

	t.Run("write binary", func(t *testing.T) {
		type binStruct struct {
			A uint32
			B uint16
		}
		tests := []struct {
			name     string
			order    binary.ByteOrder
			expected []byte
		}{
			{
				name:     "little endian struct",
				order:    binary.LittleEndian,
				expected: []byte{0x78, 0x56, 0x34, 0x12, 0xCD, 0xAB},
			},
			{
				name:     "big endian struct",
				order:    binary.BigEndian,
				expected: []byte{0x12, 0x34, 0x56, 0x78, 0xAB, 0xCD},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ab := NewAccumulativeSize(16)
				ab.WriteBinary(tt.order, binStruct{A: 0x12345678, B: 0xABCD})
				assert.True(t, bytes.Equal(tt.expected, ab.Bytes()), "want %q, got %q", tt.expected, ab.Bytes())
			})
		}
	})

	t.Run("write x", func(t *testing.T) {
		type accXStruct struct {
			Name string
		}
		x2bytes.RegisterToBytesFn(func(dst []byte, val any, _ ...any) ([]byte, error) {
			if x, ok := val.(accXStruct); ok {
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
			{name: "custom struct", value: accXStruct{Name: "Alice"}, expected: "[Alice]"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ab := NewAccumulativeSize(32)
				ab.WriteX(tt.value)
				assert.Equal(t, tt.expected, ab.String())
			})
		}
	})

	t.Run("write time", func(t *testing.T) {
		ab := NewAccumulativeSize(32)
		ab.WriteTime("%Y-%m-%d", time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))
		assert.Equal(t, "2024-01-15", ab.String())
	})

	t.Run("write uleb128", func(t *testing.T) {
		tests := []struct {
			name     string
			value    uint64
			expected []byte
		}{
			{name: "zero", value: 0, expected: []byte{0x00}},
			{name: "small", value: 127, expected: []byte{0x7F}},
			{name: "medium", value: 128, expected: []byte{0x80, 0x01}},
			{name: "large", value: 300, expected: []byte{0xAC, 0x02}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ab := NewAccumulativeSize(16)
				ab.WriteULEB128(tt.value)
				assert.True(t, bytes.Equal(tt.expected, ab.Bytes()), "want %q, got %q", tt.expected, ab.Bytes())
			})
		}
	})

	t.Run("write sleb128", func(t *testing.T) {
		tests := []struct {
			name     string
			value    int64
			expected []byte
		}{
			{name: "zero", value: 0, expected: []byte{0x00}},
			{name: "negative one", value: -1, expected: []byte{0x7F}},
			{name: "positive small", value: 63, expected: []byte{0x3F}},
			{name: "negative small", value: -64, expected: []byte{0x40}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ab := NewAccumulativeSize(16)
				ab.WriteSLEB128(tt.value)
				assert.True(t, bytes.Equal(tt.expected, ab.Bytes()), "want %q, got %q", tt.expected, ab.Bytes())
			})
		}
	})

	t.Run("cap", func(t *testing.T) {
		ab := NewAccumulativeSize(0)
		assert.GreaterOrEqual(t, ab.Cap(), 0)
		ab.WriteString("hello")
		assert.GreaterOrEqual(t, ab.Cap(), ab.Len())
	})

	t.Run("grow", func(t *testing.T) {
		ab := NewAccumulativeSize(0)
		ab.Grow(1000)
		assert.GreaterOrEqual(t, ab.Cap(), 1000)
		// Grow with a non-positive value must be a no-op and must not panic.
		ab.Grow(0)
	})

	t.Run("reset", func(t *testing.T) {
		ab := NewAccumulativeSize(64)
		ab.WriteString("test data")
		require.NotZero(t, ab.Len())
		ab.Reset()
		assert.Equal(t, 0, ab.Len())
	})

	t.Run("reset safe", func(t *testing.T) {
		ab := NewAccumulativeSize(64)
		ab.WriteString("test data")
		old := ab.Bytes()
		ab.ResetSafe()
		assert.Equal(t, 0, ab.Len())
		for i, b := range old {
			require.Zero(t, b, "ResetSafe: underlying byte %d, want 0", i)
		}
	})

	t.Run("write apply fn", func(t *testing.T) {
		addPrefix := func(dst, p []byte) []byte {
			return append(append(dst, "prefix:"...), p...)
		}
		ab := NewAccumulativeSize(64)
		ab.WriteApplyFn([]byte("test"), addPrefix)
		assert.Equal(t, "prefix:test", ab.String())

		t.Run("empty data still calls fn", func(t *testing.T) {
			called := false
			ab := NewAccumulativeSize(16)
			ab.WriteApplyFn(nil, func(dst, p []byte) []byte {
				called = true
				return append(dst, 'X')
			})
			assert.True(t, called, "WriteApplyFn() did not call fn for empty data")
			assert.Equal(t, "X", ab.String())
		})
	})

	t.Run("write apply fn n", func(t *testing.T) {
		duplicate := func(dst, p []byte) []byte {
			return append(dst, append(p, p...)...)
		}
		ab := NewAccumulativeSize(64)
		ab.WriteApplyFnN([]byte("ab"), duplicate, 3)
		// Each iteration feeds the previous result back in, so n=3 yields the
		// original input repeated 8 times (see Chain.WriteApplyFnN semantics).
		want := bytes.Repeat([]byte("ab"), 8)
		assert.True(t, bytes.Equal(want, ab.Bytes()), "want %q, got %q", want, ab.Bytes())
	})

	t.Run("write apply fn string", func(t *testing.T) {
		addPrefix := func(dst, p []byte) []byte {
			return append(append(dst, "prefix:"...), p...)
		}
		ab := NewAccumulativeSize(64)
		ab.WriteApplyFnString("test", addPrefix)
		assert.Equal(t, "prefix:test", ab.String())
	})

	t.Run("write apply fn n string", func(t *testing.T) {
		duplicate := func(dst, p []byte) []byte {
			return append(dst, append(p, p...)...)
		}
		ab := NewAccumulativeSize(64)
		ab.WriteApplyFnNString("ab", duplicate, 3)
		want := bytes.Repeat([]byte("ab"), 8)
		assert.True(t, bytes.Equal(want, ab.Bytes()), "want %q, got %q", want, ab.Bytes())
	})

	t.Run("error", func(t *testing.T) {
		ab := NewAccumulativeSize(64)
		ab.WriteString("foo").WriteInt(42).WriteBool(true)
		assert.NoError(t, ab.Error())
		ab.Reset()
		assert.NoError(t, ab.Error())
	})

	t.Run("to writer", func(t *testing.T) {
		ab := NewAccumulativeSize(64)
		ab.WriteString("foo")
		w := ab.ToWriter()
		_, err := w.WriteString("bar")
		require.NoError(t, err)
		assert.Equal(t, "foobar", ab.String())
	})

	t.Run("to reader", func(t *testing.T) {
		ab := NewAccumulativeSize(64)
		ab.WriteString("hello world")
		r := ab.ToReader()
		data, err := io.ReadAll(r)
		require.NoError(t, err)
		assert.True(t, bytes.Equal([]byte("hello world"), data), "want %q, got %q", []byte("hello world"), data)
	})
}

func BenchmarkAccumulative(b *testing.B) {
	b.Run("write", func(b *testing.B) {
		ab := NewAccumulativeSize(128)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ab.Reset().
				WriteString("foobar").
				StakeOut().
				Write(stage.b).WriteByte('-').
				WriteString(stage.s).WriteByte('-').
				WriteInt(stage.i).WriteByte('-').
				WriteUint(stage.u).WriteByte('-').
				WriteFloat(stage.f)

			if !bytes.Equal(ab.StakedBytes(), expectWS) {
				b.Error("Accumulative.Write*: mismatch result and expectation")
			}
		}
	})
}
