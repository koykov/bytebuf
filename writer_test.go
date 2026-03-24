package bytebuf

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
	"time"

	"github.com/koykov/x2bytes"
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

				if err != nil {
					t.Errorf("Write() error = %v", err)
				}

				if n != len(tt.data) {
					t.Errorf("Write() wrote %d bytes, want %d", n, len(tt.data))
				}

				if !bytes.Equal(w.Bytes(), tt.expected) {
					t.Errorf("Bytes() = %v, want %v", w.Bytes(), tt.expected)
				}
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

				if err != nil {
					t.Errorf("WriteString() error = %v", err)
				}

				if n != len(tt.str) {
					t.Errorf("WriteString() wrote %d bytes, want %d", n, len(tt.str))
				}

				if w.String() != tt.expected {
					t.Errorf("String() = %s, want %s", w.String(), tt.expected)
				}
			})
		}
	})

	t.Run("write byte", func(t *testing.T) {
		w := NewWriter(nil)

		tests := []byte{0x00, 0x01, 0xFF, 'a', 'z'}

		for _, b := range tests {
			err := w.WriteByte(b)
			if err != nil {
				t.Errorf("WriteByte(%d) error = %v", b, err)
			}
		}

		expected := tests
		if !bytes.Equal(w.Bytes(), expected) {
			t.Errorf("Bytes() = %v, want %v", w.Bytes(), expected)
		}
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

				if (err != nil) != tt.wantErr {
					t.Errorf("WriteAt() error = %v, wantErr %v", err, tt.wantErr)
				}

				if !tt.wantErr && n != len(tt.data) {
					t.Errorf("WriteAt() wrote %d bytes, want %d", n, len(tt.data))
				}

				if !tt.wantErr && !bytes.Equal(w.Bytes(), tt.want) {
					t.Errorf("After WriteAt bytes = %s, want %s", w.String(), string(tt.want))
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
				if initialCap < tt.growTo {
					t.Errorf("After Grow(%d) capacity = %d, want at least %d", tt.growTo, initialCap, tt.growTo)
				}

				_, _ = w.Write(tt.writeData)
				if w.Len() != len(tt.writeData) {
					t.Errorf("After write length = %d, want %d", w.Len(), len(tt.writeData))
				}
			})
		}
	})

	t.Run("grow delta", func(t *testing.T) {
		w := NewWriter(nil)
		initialCap := w.Cap()

		w.GrowDelta(100)
		if w.Cap() < initialCap+100 {
			t.Errorf("After GrowDelta(100) capacity = %d, want at least %d", w.Cap(), initialCap+100)
		}

		w.GrowDelta(0)
		if w.Cap() < initialCap+100 {
			t.Errorf("After GrowDelta(0) capacity decreased")
		}
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

				if err != nil {
					t.Errorf("WriteInt() error = %v", err)
				}

				if !bytes.Equal(w.Bytes(), tt.expected) {
					t.Errorf("Bytes() = %v, want %v", w.Bytes(), tt.expected)
				}
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

				if err != nil {
					t.Errorf("WriteUint() error = %v", err)
				}

				if !bytes.Equal(w.Bytes(), tt.expected) {
					t.Errorf("Bytes() = %v, want %v", w.Bytes(), tt.expected)
				}
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

				if err != nil {
					t.Errorf("WriteFloat() error = %v", err)
				}

				if !bytes.Equal(w.Bytes(), tt.expected) {
					t.Errorf("Bytes() = %v, want %v", w.Bytes(), tt.expected)
				}
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

				if err != nil {
					t.Errorf("WriteBool() error = %v", err)
				}

				if !bytes.Equal(w.Bytes(), tt.expected) {
					t.Errorf("Bytes() = %v, want %v", w.Bytes(), tt.expected)
				}
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

				if (err != nil) != tt.wantErr {
					t.Errorf("WriteBinary() error = %v, wantErr %v", err, tt.wantErr)
				}

				if !tt.wantErr {
					if n != len(tt.expect) {
						t.Errorf("WriteBinary() wrote %d bytes, want %d", n, len(tt.expect))
					}

					if !bytes.Equal(w.Bytes(), tt.expect) {
						t.Errorf("Bytes() = %v, want %v", w.Bytes(), tt.expect)
					}
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

				if err != nil {
					t.Errorf("WriteX() error = %v", err)
				}

				if n != len(tt.expected) {
					t.Errorf("WriteX() wrote %d bytes, want %d", n, len(tt.expected))
				}

				if w.String() != tt.expected {
					t.Errorf("String() = %s, want %s", w.String(), tt.expected)
				}
			})
		}
	})

	t.Run("write marshaller to", func(t *testing.T) {
		mock := &mockMarshaller{data: []byte("test data")}

		w := NewWriter(nil)
		n, err := w.WriteMarshallerTo(mock)

		if err != nil {
			t.Errorf("WriteMarshallerTo() error = %v", err)
		}

		if n != mock.Size() {
			t.Errorf("WriteMarshallerTo() wrote %d bytes, want %d", n, mock.Size())
		}

		if !bytes.Equal(w.Bytes(), mock.data) {
			t.Errorf("Bytes() = %v, want %v", w.Bytes(), mock.data)
		}
	})

	t.Run("write apply fn", func(t *testing.T) {
		addPrefix := func(dst, p []byte) []byte {
			return append(append(dst, "prefix:"...), p...)
		}

		data := []byte("test")
		expected := []byte("prefix:test")

		w := NewWriter(nil)
		n, err := w.WriteApplyFn(data, addPrefix)

		if err != nil {
			t.Errorf("WriteApplyFn() error = %v", err)
		}

		if n != len(expected) {
			t.Errorf("WriteApplyFn() wrote %d bytes, want %d", n, len(expected))
		}

		if !bytes.Equal(w.Bytes(), expected) {
			t.Errorf("Bytes() = %v, want %v", w.Bytes(), expected)
		}
	})

	t.Run("wryte apply fn N", func(t *testing.T) {
		duplicate := func(dst, p []byte) []byte {
			return append(dst, append(p, p...)...)
		}

		data := []byte("ab")
		expected := []byte("abababababababababababababababab")

		w := NewWriter(nil)
		n, err := w.WriteApplyFnN(data, duplicate, 4)

		if err != nil {
			t.Errorf("WriteApplyFnN() error = %v", err)
		}

		if n != 32 {
			t.Errorf("WriteApplyFnN() wrote %d bytes, want 4", n)
		}

		if !bytes.Equal(w.Bytes(), expected) {
			t.Errorf("Bytes() = %v, want %v", w.Bytes(), expected)
		}
	})

	t.Run("write time", func(t *testing.T) {
		layout := "%Y-%m-%d %H:%M:%S"
		tm := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)
		expected := "2024-01-15 10:30:45"

		w := NewWriter(nil)
		n, err := w.WriteTime(layout, tm)

		if err != nil {
			t.Errorf("WriteTime() error = %v", err)
		}

		if n != len(expected) {
			t.Errorf("WriteTime() wrote %d bytes, want %d", n, len(expected))
		}

		if w.String() != expected {
			t.Errorf("String() = %s, want %s", w.String(), expected)
		}
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

				if err != nil {
					t.Errorf("WriteULEB128() error = %v", err)
				}

				if n != len(tt.expected) {
					t.Errorf("WriteULEB128() wrote %d bytes, want %d", n, len(tt.expected))
				}

				if !bytes.Equal(w.Bytes(), tt.expected) {
					t.Errorf("Bytes() = %v, want %v", w.Bytes(), tt.expected)
				}
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

				if err != nil {
					t.Errorf("WriteSLEB128() error = %v", err)
				}

				if n != len(tt.expected) {
					t.Errorf("WriteSLEB128() wrote %d bytes, want %d", n, len(tt.expected))
				}

				if !bytes.Equal(w.Bytes(), tt.expected) {
					t.Errorf("Bytes() = %v, want %v", w.Bytes(), tt.expected)
				}
			})
		}
	})

	t.Run("len cap", func(t *testing.T) {
		w := NewWriter(nil)

		if w.Len() != 0 {
			t.Errorf("Initial Len() = %d, want 0", w.Len())
		}

		_, _ = w.Write([]byte("test"))
		if w.Len() != 4 {
			t.Errorf("After write Len() = %d, want 4", w.Len())
		}

		initialCap := w.Cap()
		w.Grow(100)
		if w.Cap() <= initialCap {
			t.Errorf("After Grow Cap() = %d, should be greater than %d", w.Cap(), initialCap)
		}
	})

	t.Run("reset", func(t *testing.T) {
		w := NewWriter(nil)
		_, _ = w.Write([]byte("test data"))

		if w.Len() == 0 {
			t.Error("Writer should have data before reset")
		}

		w.Reset()

		if w.Len() != 0 {
			t.Errorf("After Reset Len() = %d, want 0", w.Len())
		}

		if w.String() != "" {
			t.Errorf("After Reset String() = %s, want empty", w.String())
		}
	})

	t.Run("reset safe", func(t *testing.T) {
		w := NewWriter(nil)
		_, _ = w.Write([]byte("test data"))

		oldBytes := w.Bytes()
		w.ResetSafe()

		if w.Len() != 0 {
			t.Errorf("After ResetSafe Len() = %d, want 0", w.Len())
		}

		if !bytes.Equal(oldBytes, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}) {
			t.Error("ResetSafe modified the underlying slice")
		}
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
			if err != nil {
				t.Errorf("Write() error at offset %d: %v", i, err)
			}
			if n != end-i {
				t.Errorf("Write() wrote %d bytes, want %d", n, end-i)
			}
		}

		if w.Len() != size {
			t.Errorf("Final length = %d, want %d", w.Len(), size)
		}

		if !bytes.Equal(w.Bytes(), largeData) {
			t.Error("Data mismatch in large write")
		}

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
				if err != nil {
					t.Errorf("WriteAt(%d) error = %v", pos, err)
				}
				if n != len(data) {
					t.Errorf("WriteAt(%d) wrote %d bytes, want %d", pos, n, len(data))
				}
			}

			for _, pos := range positions {
				if int(pos) >= w.Len() {
					continue
				}
				expected := []byte("test")
				actual := w.Bytes()[pos : pos+4]
				if !bytes.Equal(actual, expected) {
					t.Errorf("At position %d got %v, want %v", pos, actual, expected)
				}
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

		if w.Len() == 0 {
			t.Error("No data written")
		}

		data := w.Bytes()
		if len(data) == 0 {
			t.Error("Bytes() returned empty slice")
		}

		saved := w.String()
		if saved == "" {
			t.Error("String() returned empty")
		}

		w.Reset()
		if w.Len() != 0 {
			t.Error("After Reset Len() != 0")
		}

		_, _ = w.WriteApplyFn([]byte("applied"), func(dst, p []byte) []byte {
			return append(dst, append([]byte("prefix_"), p...)...)
		})

		_, _ = w.WriteTime(time.RFC3339, time.Now())
		_, _ = w.WriteULEB128(123456)
		_, _ = w.WriteSLEB128(-123456)

		if w.Len() == 0 {
			t.Error("Second batch write failed")
		}
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
