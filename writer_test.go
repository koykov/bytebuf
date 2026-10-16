package bytebuf

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
	"time"

	"github.com/koykov/x2bytes"
	"github.com/stretchr/testify/assert"
)

type mockMarshaller struct {
	data []byte
}

func (m *mockMarshaller) Size() int {
	return len(m.data)
}

func (m *mockMarshaller) MarshalTo(buf []byte) (int, error) {
	if len(buf) < len(m.data) {
		return 0, bytes.ErrTooLarge
	}
	return copy(buf, m.data), nil
}

func TestWriter(t *testing.T) {
	t.Run("write", func(t *testing.T) {
		tests := []struct {
			name     string
			data     []byte
			expected []byte
		}{
			{
				name:     "write single byte",
				data:     []byte{0x01},
				expected: []byte{0x01},
			},
			{
				name:     "write multiple bytes",
				data:     []byte("hello world"),
				expected: []byte("hello world"),
			},
			{
				name:     "write empty",
				data:     []byte{},
				expected: []byte{},
			},
			{
				name:     "write zeros",
				data:     []byte{0, 0, 0, 0},
				expected: []byte{0, 0, 0, 0},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := NewWriter(nil)
				n, err := w.Write(tt.data)

				assert.NoError(t, err)

				assert.Equal(t, len(tt.data), n, "Write() byte count")

				// bytes.Equal is kept here so that a nil buffer (empty write)
				// remains equal to an empty non-nil expected slice.
				assert.True(t, bytes.Equal(tt.expected, w.Bytes()), "Bytes() = %v, want %v", w.Bytes(), tt.expected)
			})
		}
	})

	t.Run("write string", func(t *testing.T) {
		tests := []struct {
			name     string
			str      string
			expected string
		}{
			{
				name:     "simple string",
				str:      "hello",
				expected: "hello",
			},
			{
				name:     "unicode string",
				str:      "Hello, 世界!",
				expected: "Hello, 世界!",
			},
			{
				name:     "empty string",
				str:      "",
				expected: "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := NewWriter(nil)
				n, err := w.WriteString(tt.str)

				assert.NoError(t, err)

				assert.Equal(t, len(tt.str), n, "WriteString() byte count")

				assert.Equal(t, tt.expected, w.String(), "String()")
			})
		}
	})

	t.Run("write byte", func(t *testing.T) {
		w := NewWriter(nil)

		tests := []byte{0x00, 0x01, 0xFF, 'a', 'z'}

		for _, b := range tests {
			err := w.WriteByte(b)
			assert.NoError(t, err, "WriteByte(%d)", b)
		}

		expected := tests
		assert.True(t, bytes.Equal(expected, w.Bytes()), "want %q, got %q", expected, w.Bytes())
	})

	t.Run("write at", func(t *testing.T) {
		w := NewWriter(nil)

		_, _ = w.Write([]byte("0123456789"))

		tests := []struct {
			name    string
			data    []byte
			offset  int64
			want    []byte
			wantErr bool
		}{
			{
				name:    "write at beginning",
				data:    []byte("abc"),
				offset:  0,
				want:    []byte("abc3456789"),
				wantErr: false,
			},
			{
				name:    "write in middle",
				data:    []byte("XYZ"),
				offset:  3,
				want:    []byte("abcXYZ6789"),
				wantErr: false,
			},
			{
				name:    "write at end",
				data:    []byte("END"),
				offset:  9,
				want:    []byte("abcXYZ678END"),
				wantErr: false,
			},
			{
				name:    "write beyond end",
				data:    []byte("EXTRA"),
				offset:  15,
				want:    []byte("abcXYZ678END\x00\x00\x00EXTRA"),
				wantErr: false,
			},
			{
				name:    "write negative offset",
				data:    []byte("bad"),
				offset:  -1,
				want:    nil,
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

				// w := NewWriter(nil)
				// _, _ = w.Write([]byte("0123456789"))

				n, err := w.WriteAt(tt.data, tt.offset)

				if tt.wantErr {
					assert.Error(t, err, "WriteAt() error")
				} else {
					assert.NoError(t, err, "WriteAt() error")
					assert.Equal(t, len(tt.data), n, "WriteAt() byte count")
					assert.True(t, bytes.Equal(tt.want, w.Bytes()), "After WriteAt bytes: want %q, got %q", tt.want, w.Bytes())
				}
			})
		}
	})

	t.Run("grow", func(t *testing.T) {
		tests := []struct {
			name      string
			initial   int
			growTo    int
			writeData []byte
		}{
			{
				name:      "grow to larger size",
				initial:   0,
				growTo:    1000,
				writeData: bytes.Repeat([]byte("x"), 500),
			},
			{
				name:      "grow to smaller size",
				initial:   100,
				growTo:    50,
				writeData: bytes.Repeat([]byte("x"), 30),
			},
			{
				name:      "grow zero",
				initial:   0,
				growTo:    0,
				writeData: []byte("test"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := NewWriter(nil)
				w.Grow(tt.growTo)
				w.Reset()

				initialCap := w.Cap()
				assert.GreaterOrEqual(t, initialCap, tt.growTo, "After Grow(%d) capacity", tt.growTo)

				_, _ = w.Write(tt.writeData)
				assert.Equal(t, len(tt.writeData), w.Len(), "After write length")
			})
		}
	})

	t.Run("grow delta", func(t *testing.T) {
		w := NewWriter(nil)
		initialCap := w.Cap()

		w.GrowDelta(100)
		assert.GreaterOrEqual(t, w.Cap(), initialCap+100, "After GrowDelta(100) capacity")

		w.GrowDelta(0)
		assert.GreaterOrEqual(t, w.Cap(), initialCap+100, "After GrowDelta(0) capacity decreased")
	})

	t.Run("write int", func(t *testing.T) {
		tests := []struct {
			name     string
			value    int64
			expected []byte
		}{
			{
				name:     "zero",
				value:    0,
				expected: []byte("0"),
			},
			{
				name:     "positive",
				value:    123456789,
				expected: []byte("123456789"),
			},
			{
				name:     "negative",
				value:    -123456789,
				expected: []byte("-123456789"),
			},
			{
				name:     "max int64",
				value:    math.MaxInt64,
				expected: []byte("9223372036854775807"),
			},
			{
				name:     "min int64",
				value:    math.MinInt64,
				expected: []byte("-9223372036854775808"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := NewWriter(nil)
				_, err := w.WriteInt(tt.value)

				assert.NoError(t, err)

				assert.True(t, bytes.Equal(tt.expected, w.Bytes()), "want %q, got %q", tt.expected, w.Bytes())
			})
		}
	})

	t.Run("write uint", func(t *testing.T) {
		tests := []struct {
			name     string
			value    uint64
			expected []byte
		}{
			{
				name:     "zero",
				value:    0,
				expected: []byte("0"),
			},
			{
				name:     "small",
				value:    255,
				expected: []byte("255"),
			},
			{
				name:     "max uint64",
				value:    math.MaxUint64,
				expected: []byte("18446744073709551615"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := NewWriter(nil)
				_, err := w.WriteUint(tt.value)

				assert.NoError(t, err)

				assert.True(t, bytes.Equal(tt.expected, w.Bytes()), "want %q, got %q", tt.expected, w.Bytes())
			})
		}
	})

	t.Run("write float", func(t *testing.T) {
		tests := []struct {
			name     string
			value    float64
			expected []byte
		}{
			{
				name:     "zero",
				value:    0.0,
				expected: []byte("0"),
			},
			{
				name:     "positive",
				value:    123.456,
				expected: []byte("123.456"),
			},
			{
				name:     "negative",
				value:    -123.456,
				expected: []byte("-123.456"),
			},
			{
				name:     "pi",
				value:    math.Pi,
				expected: []byte("3.141592653589793"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := NewWriter(nil)
				_, err := w.WriteFloat(tt.value)

				assert.NoError(t, err)

				assert.True(t, bytes.Equal(tt.expected, w.Bytes()), "want %q, got %q", tt.expected, w.Bytes())
			})
		}
	})

	t.Run("write bool", func(t *testing.T) {
		tests := []struct {
			name     string
			value    bool
			expected []byte
		}{
			{
				name:     "true",
				value:    true,
				expected: []byte("true"),
			},
			{
				name:     "false",
				value:    false,
				expected: []byte("false"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := NewWriter(nil)
				_, err := w.WriteBool(tt.value)

				assert.NoError(t, err)

				assert.True(t, bytes.Equal(tt.expected, w.Bytes()), "want %q, got %q", tt.expected, w.Bytes())
			})
		}
	})

	t.Run("write binary", func(t *testing.T) {
		type TestStruct struct {
			A uint32
			B uint16
			C byte
		}

		tests := []struct {
			name    string
			order   binary.ByteOrder
			value   interface{}
			expect  []byte
			wantErr bool
		}{
			{
				name:    "little endian struct",
				order:   binary.LittleEndian,
				value:   &TestStruct{A: 0x12345678, B: 0xABCD, C: 0xFF},
				expect:  []byte{0x78, 0x56, 0x34, 0x12, 0xCD, 0xAB, 0xFF},
				wantErr: false,
			},
			{
				name:    "big endian struct",
				order:   binary.BigEndian,
				value:   &TestStruct{A: 0x12345678, B: 0xABCD, C: 0xFF},
				expect:  []byte{0x12, 0x34, 0x56, 0x78, 0xAB, 0xCD, 0xFF},
				wantErr: false,
			},
			{
				name:    "invalid type",
				order:   binary.LittleEndian,
				value:   make(chan int),
				expect:  nil,
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := NewWriter(nil)
				n, err := w.WriteBinary(tt.order, tt.value)

				if tt.wantErr {
					assert.Error(t, err, "WriteBinary() error")
				} else {
					assert.NoError(t, err, "WriteBinary() error")
					assert.Equal(t, len(tt.expect), n, "WriteBinary() byte count")
					assert.True(t, bytes.Equal(tt.expect, w.Bytes()), "write binary: want %q, got %q", tt.expect, w.Bytes())
				}
			})
		}
	})

	t.Run("write X", func(t *testing.T) {
		type testStruct struct {
			Name string
		}
		tests := []struct {
			name     string
			value    interface{}
			expected string
		}{
			{
				name:     "int",
				value:    42,
				expected: "42",
			},
			{
				name:     "string",
				value:    "hello",
				expected: "hello",
			},
			{
				name:     "struct",
				value:    testStruct{"Alice"},
				expected: "{Alice}",
			},
		}
		x2bytes.RegisterToBytesFn(func(dst []byte, val any, _ ...any) ([]byte, error) {
			if x, ok := val.(testStruct); ok {
				dst = append(dst, '{')
				dst = append(dst, x.Name...)
				dst = append(dst, '}')
				return dst, nil
			}
			return dst, x2bytes.ErrUnknownType
		})

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := NewWriter(nil)
				n, err := w.WriteX(tt.value)

				assert.NoError(t, err)

				assert.Equal(t, len(tt.expected), n, "WriteX() byte count")

				assert.Equal(t, tt.expected, w.String(), "String()")
			})
		}
	})

	t.Run("write marshaller to", func(t *testing.T) {
		mock := &mockMarshaller{data: []byte("test data")}

		w := NewWriter(nil)
		n, err := w.WriteMarshallerTo(mock)

		assert.NoError(t, err)

		assert.Equal(t, mock.Size(), n, "WriteMarshallerTo() byte count")

		assert.True(t, bytes.Equal(mock.data, w.Bytes()), "WriteMarshallerTo(): want %q, got %q", mock.data, w.Bytes())
	})

	t.Run("write apply fn", func(t *testing.T) {
		addPrefix := func(dst, p []byte) []byte {
			return append(append(dst, "prefix:"...), p...)
		}

		data := []byte("test")
		expected := []byte("prefix:test")

		w := NewWriter(nil)
		n, err := w.WriteApplyFn(data, addPrefix)

		assert.NoError(t, err)

		assert.Equal(t, len(expected), n, "WriteApplyFn() byte count")

		assert.True(t, bytes.Equal(expected, w.Bytes()), "want %q, got %q", expected, w.Bytes())
	})

	t.Run("wryte apply fn N", func(t *testing.T) {
		duplicate := func(dst, p []byte) []byte {
			return append(dst, append(p, p...)...)
		}

		data := []byte("ab")
		expected := []byte("abababababababababababababababab")

		w := NewWriter(nil)
		n, err := w.WriteApplyFnN(data, duplicate, 4)

		assert.NoError(t, err)

		assert.Equal(t, 32, n, "WriteApplyFnN() byte count")

		assert.True(t, bytes.Equal(expected, w.Bytes()), "want %q, got %q", expected, w.Bytes())
	})

	t.Run("write time", func(t *testing.T) {
		layout := "%Y-%m-%d %H:%M:%S"
		tm := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)
		expected := "2024-01-15 10:30:45"

		w := NewWriter(nil)
		n, err := w.WriteTime(layout, tm)

		assert.NoError(t, err)

		assert.Equal(t, len(expected), n, "WriteTime() byte count")

		assert.Equal(t, expected, w.String(), "String()")
	})

	t.Run("write ULEB128", func(t *testing.T) {
		tests := []struct {
			name     string
			value    uint64
			expected []byte
		}{
			{
				name:     "zero",
				value:    0,
				expected: []byte{0x00},
			},
			{
				name:     "small",
				value:    127,
				expected: []byte{0x7F},
			},
			{
				name:     "medium",
				value:    128,
				expected: []byte{0x80, 0x01},
			},
			{
				name:     "large",
				value:    300,
				expected: []byte{0xAC, 0x02},
			},
			{
				name:     "max uint64",
				value:    math.MaxUint64,
				expected: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := NewWriter(nil)
				n, err := w.WriteULEB128(tt.value)

				assert.NoError(t, err)

				assert.Equal(t, len(tt.expected), n, "WriteULEB128() byte count")

				assert.True(t, bytes.Equal(tt.expected, w.Bytes()), "want %q, got %q", tt.expected, w.Bytes())
			})
		}
	})

	t.Run("write SLEB128", func(t *testing.T) {
		tests := []struct {
			name     string
			value    int64
			expected []byte
		}{
			{
				name:     "zero",
				value:    0,
				expected: []byte{0x00},
			},
			{
				name:     "positive small",
				value:    63,
				expected: []byte{0x3F},
			},
			{
				name:     "positive medium",
				value:    64,
				expected: []byte{0xC0, 0x00},
			},
			{
				name:     "negative small",
				value:    -1,
				expected: []byte{0x7F},
			},
			{
				name:     "negative medium",
				value:    -64,
				expected: []byte{0x40},
			},
			{
				name:     "negative large",
				value:    -100,
				expected: []byte{0x9C, 0x7F},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := NewWriter(nil)
				n, err := w.WriteSLEB128(tt.value)

				assert.NoError(t, err)

				assert.Equal(t, len(tt.expected), n, "WriteSLEB128() byte count")

				assert.True(t, bytes.Equal(tt.expected, w.Bytes()), "want %q, got %q", tt.expected, w.Bytes())
			})
		}
	})

	t.Run("len cap", func(t *testing.T) {
		w := NewWriter(nil)

		assert.Equal(t, 0, w.Len(), "Initial Len()")

		_, _ = w.Write([]byte("test"))
		assert.Equal(t, 4, w.Len(), "After write Len()")

		initialCap := w.Cap()
		w.Grow(100)
		assert.Greater(t, w.Cap(), initialCap, "After Grow Cap()")
	})

	t.Run("reset", func(t *testing.T) {
		w := NewWriter(nil)
		_, _ = w.Write([]byte("test data"))

		assert.NotZero(t, w.Len(), "Writer should have data before reset")

		w.Reset()

		assert.Equal(t, 0, w.Len(), "After Reset Len()")

		assert.Empty(t, w.String(), "After Reset String()")
	})

	t.Run("reset safe", func(t *testing.T) {
		w := NewWriter(nil)
		_, _ = w.Write([]byte("test data"))

		oldBytes := w.Bytes()
		w.ResetSafe()

		assert.Equal(t, 0, w.Len(), "After ResetSafe Len()")

		assert.True(t, bytes.Equal([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0}, oldBytes), "ResetSafe modified the underlying slice: want %q, got %q", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}, oldBytes)
	})

	t.Run("big data", func(t *testing.T) {
		size := 10 * 1024 * 1024
		largeData := bytes.Repeat([]byte("x"), size)

		w := NewWriter(nil)

		chunkSize := 1024 * 1024
		for i := 0; i < size; i += chunkSize {
			end := i + chunkSize
			if end > size {
				end = size
			}
			n, err := w.Write(largeData[i:end])
			assert.NoError(t, err, "Write() error at offset %d", i)
			assert.Equal(t, end-i, n, "Write() byte count at offset %d", i)
		}

		assert.Equal(t, size, w.Len(), "Final length")

		assert.True(t, bytes.Equal(largeData, w.Bytes()), "Data mismatch in large write: want %q, got %q", largeData, w.Bytes())

		t.Run("write at", func(t *testing.T) {
			size := 10 * 1024 * 1024
			w := NewWriter(nil)

			w.Grow(size)

			positions := []int64{
				0,
				int64(size / 4),
				int64(size / 2),
				int64(size - 1000),
				int64(size - 1),
			}

			for _, pos := range positions {
				data := []byte("test")
				n, err := w.WriteAt(data, pos)
				assert.NoError(t, err, "WriteAt(%d) error", pos)
				assert.Equal(t, len(data), n, "WriteAt(%d) byte count", pos)
			}

			for _, pos := range positions {
				if int(pos) >= w.Len() {
					continue
				}
				expected := []byte("test")
				actual := w.Bytes()[pos : pos+4]
				assert.True(t, bytes.Equal(expected, actual), "At position %d: want %q, got %q", pos, expected, actual)
			}
		})
	})

	t.Run("combine", func(t *testing.T) {
		w := NewWriter(nil)

		_, _ = w.Write([]byte("Data: "))
		_, _ = w.WriteString("Hello, ")
		_, _ = w.WriteRune('世')
		_, _ = w.WriteInt(12345)
		_, _ = w.WriteUint(67890)
		_, _ = w.WriteFloat(3.14159)
		_, _ = w.WriteBool(true)

		assert.NotZero(t, w.Len(), "No data written")

		data := w.Bytes()
		assert.NotEmpty(t, data, "Bytes() returned empty slice")

		saved := w.String()
		assert.NotEmpty(t, saved, "String() returned empty")

		w.Reset()
		assert.Zero(t, w.Len(), "After Reset Len() != 0")

		_, _ = w.WriteApplyFn([]byte("applied"), func(dst, p []byte) []byte {
			return append(dst, append([]byte("prefix_"), p...)...)
		})

		_, _ = w.WriteTime(time.RFC3339, time.Now())
		_, _ = w.WriteULEB128(123456)
		_, _ = w.WriteSLEB128(-123456)

		assert.NotZero(t, w.Len(), "Second batch write failed")
	})
}

func BenchmarkWriter(b *testing.B) {
	b.Run("write", func(b *testing.B) {
		w := NewWriter(nil)
		data := []byte("test data")
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			w.Reset()
			_, _ = w.Write(data)
		}
	})

	b.Run("big data", func(b *testing.B) {
		w := NewWriter(nil)
		size := 1024 * 1024
		data := bytes.Repeat([]byte("x"), size)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			w.Reset()
			_, _ = w.Write(data)
		}
	})

	b.Run("write int", func(b *testing.B) {
		w := NewWriter(nil)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w.Reset()
			_, _ = w.WriteInt(123456789)
		}
	})

	b.Run("write ULEB128", func(b *testing.B) {
		w := NewWriter(nil)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w.Reset()
			_, _ = w.WriteULEB128(123456789)
		}
	})

	b.Run("write SLEB128", func(b *testing.B) {
		w := NewWriter(nil)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			w.Reset()
			_, _ = w.WriteSLEB128(-123456789)
		}
	})

	b.Run("write apply fn", func(b *testing.B) {
		w := NewWriter(nil)
		data := []byte("test data")
		fn := func(dst, p []byte) []byte {
			return append(dst, p...)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			w.Reset()
			_, _ = w.WriteApplyFn(data, fn)
		}
	})
}
