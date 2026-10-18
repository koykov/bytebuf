package bytebuf

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	t.Run("read", func(t *testing.T) {
		tests := []struct {
			name     string
			data     []byte
			readSize int
			want     []byte
			wantErr  bool
		}{
			{
				name:     "read all data",
				data:     []byte("hello world"),
				readSize: 11,
				want:     []byte("hello world"),
				wantErr:  false,
			},
			{
				name:     "read partial data",
				data:     []byte("hello world"),
				readSize: 5,
				want:     []byte("hello"),
				wantErr:  false,
			},
			{
				name:     "read more than available",
				data:     []byte("hi"),
				readSize: 10,
				want:     []byte("hi"),
				wantErr:  true, // EOF
			},
			{
				name:     "read empty buffer",
				data:     []byte{},
				readSize: 10,
				want:     []byte{},
				wantErr:  true, // EOF
			},
			{
				name:     "read zero bytes",
				data:     []byte("test"),
				readSize: 0,
				want:     []byte{},
				wantErr:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := NewReader(tt.data)
				buf := make([]byte, tt.readSize)
				n, err := r.Read(buf)

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}

				assert.Equal(t, len(tt.want), n)
				assert.True(t, bytes.Equal(tt.want, buf[:n]), "want %q, got %q", tt.want, buf[:n])
			})
		}
	})

	t.Run("read extra", func(t *testing.T) {
		tests := []struct {
			name    string
			data    []byte
			offset  int64
			size    int
			want    []byte
			wantErr error
		}{
			{
				name:    "offset equals length with non-empty p",
				data:    []byte("0123456789"),
				offset:  10,
				size:    4,
				want:    []byte{},
				wantErr: io.EOF,
			},
			{
				name:    "offset equals length with empty p",
				data:    []byte("0123456789"),
				offset:  10,
				size:    0,
				want:    []byte{},
				wantErr: nil,
			},
			{
				name:    "offset length-1 with p larger than remaining",
				data:    []byte("0123456789"),
				offset:  9,
				size:    3,
				want:    []byte("9"),
				wantErr: io.EOF,
			},
			{
				name:    "offset zero on empty buffer",
				data:    []byte{},
				offset:  0,
				size:    4,
				want:    []byte{},
				wantErr: io.EOF,
			},
			{
				name:    "negative offset",
				data:    []byte("0123456789"),
				offset:  -1,
				size:    2,
				want:    nil,
				wantErr: ErrNegativeOffset,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := NewReader(tt.data)
				buf := make([]byte, tt.size)
				n, err := r.ReadAt(buf, tt.offset)

				assert.Equal(t, tt.wantErr, err)
				assert.Equal(t, len(tt.want), n)
				assert.True(t, bytes.Equal(tt.want, buf[:n]), "ReadAt() got %v, want %v", buf[:n], tt.want)
			})
		}
	})

	t.Run("read at", func(t *testing.T) {
		data := []byte("0123456789")
		r := NewReader(data)

		tests := []struct {
			name    string
			offset  int64
			size    int
			want    []byte
			wantErr error
		}{
			{
				name:    "read from beginning",
				offset:  0,
				size:    5,
				want:    []byte("01234"),
				wantErr: nil,
			},
			{
				name:    "read from middle",
				offset:  5,
				size:    3,
				want:    []byte("567"),
				wantErr: nil,
			},
			{
				name:    "read from end",
				offset:  8,
				size:    2,
				want:    []byte("89"),
				wantErr: nil,
			},
			{
				name:    "read beyond end",
				offset:  9,
				size:    2,
				want:    []byte("9"),
				wantErr: io.EOF,
			},
			{
				name:    "read from negative offset",
				offset:  -1,
				size:    2,
				want:    nil,
				wantErr: ErrNegativeOffset,
			},
			{
				name:    "read zero bytes",
				offset:  5,
				size:    0,
				want:    []byte{},
				wantErr: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				buf := make([]byte, tt.size)
				n, err := r.ReadAt(buf, tt.offset)

				assert.Equal(t, tt.wantErr, err)
				assert.Equal(t, len(tt.want), n)
				assert.True(t, bytes.Equal(tt.want, buf[:n]), "ReadAt() got %v, want %v", buf[:n], tt.want)
			})
		}
	})

	t.Run("read byte", func(t *testing.T) {
		data := []byte("abc")
		r := NewReader(data)

		for i, expected := range data {
			b, err := r.ReadByte()
			assert.NoError(t, err, "ReadByte() at index %d", i)
			assert.Equal(t, expected, b, "ReadByte() at index %d", i)
		}

		_, err := r.ReadByte()
		assert.Equal(t, io.EOF, err)
	})

	t.Run("unread byte", func(t *testing.T) {
		tests := []struct {
			name        string
			data        []byte
			readCount   int
			unreadCount int
			wantAfter   []byte
			wantErr     bool
		}{
			{
				name:        "unread one byte",
				data:        []byte("abc"),
				readCount:   2,
				unreadCount: 1,
				wantAfter:   []byte("bc"),
				wantErr:     false,
			},
			{
				name:        "unread multiple times",
				data:        []byte("abc"),
				readCount:   1,
				unreadCount: 2,
				wantAfter:   []byte("abc"),
				wantErr:     true, // second unread should fail
			},
			{
				name:        "unread at beginning",
				data:        []byte("abc"),
				readCount:   0,
				unreadCount: 1,
				wantAfter:   []byte("abc"),
				wantErr:     true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := NewReader(tt.data)

				for i := 0; i < tt.readCount; i++ {
					_, err := r.ReadByte()
					require.NoError(t, err, "Failed to read")
				}

				var err error
				for i := 0; i < tt.unreadCount; i++ {
					err = r.UnreadByte()
				}

				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}

				if !tt.wantErr {
					buf := make([]byte, len(tt.wantAfter))
					n, _ := r.Read(buf)
					assert.True(t, bytes.Equal(tt.wantAfter, buf[:n]), "want %q, got %q", tt.wantAfter, buf[:n])
				}
			})
		}
	})

	t.Run("read rune", func(t *testing.T) {
		tests := []struct {
			name          string
			data          string
			expected      []rune
			expectedSizes []int
		}{
			{
				name:          "ascii only",
				data:          "abc",
				expected:      []rune{'a', 'b', 'c'},
				expectedSizes: []int{1, 1, 1},
			},
			{
				name:          "unicode characters",
				data:          "Hello, 世界",
				expected:      []rune{'H', 'e', 'l', 'l', 'o', ',', ' ', '世', '界'},
				expectedSizes: []int{1, 1, 1, 1, 1, 1, 1, 3, 3},
			},
			{
				name:          "emoji",
				data:          "😀🎉",
				expected:      []rune{'😀', '🎉'},
				expectedSizes: []int{4, 4},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := NewReader([]byte(tt.data))

				for i, expectedRune := range tt.expected {
					ch, size, err := r.ReadRune()
					assert.NoError(t, err, "ReadRune() at index %d", i)
					assert.Equal(t, expectedRune, ch)
					assert.Equal(t, tt.expectedSizes[i], size)
				}

				_, _, err := r.ReadRune()
				assert.Equal(t, io.EOF, err)
			})
		}
	})

	t.Run("unread rune", func(t *testing.T) {
		data := "Hello, 世界"
		r := NewReader([]byte(data))

		first, _, _ := r.ReadRune()
		assert.Equal(t, 'H', first)

		err := r.UnreadRune()
		assert.NoError(t, err)

		again, _, _ := r.ReadRune()
		assert.Equal(t, 'H', again)

		_ = r.UnreadRune()
		err = r.UnreadRune()
		assert.Error(t, err, "Second UnreadRune() should fail")
	})

	t.Run("unread rune extra", func(t *testing.T) {
		tests := []struct {
			name    string
			data    []byte
			reads   int
			unreads []error
		}{
			{
				name:    "unread after invalid utf-8 start byte",
				data:    []byte{0xFF},
				reads:   1,
				unreads: []error{ErrBadRune},
			},
			{
				name:    "unread without prior read",
				data:    []byte("abc"),
				reads:   0,
				unreads: []error{ErrOutOfRange},
			},
			{
				name:    "unread twice in a row",
				data:    []byte("abc"),
				reads:   1,
				unreads: []error{nil, ErrOutOfRange},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := NewReader(tt.data)

				for i := 0; i < tt.reads; i++ {
					if _, _, err := r.ReadRune(); err != nil {
						require.NoError(t, err, "ReadRune() #%d", i+1)
					}
				}

				for i, wantErr := range tt.unreads {
					err := r.UnreadRune()
					assert.Equal(t, wantErr, err, "UnreadRune() #%d", i+1)
				}
			})
		}
	})

	t.Run("seek", func(t *testing.T) {
		data := []byte("0123456789")

		tests := []struct {
			name      string
			preoffset int64
			offset    int64
			whence    int
			wantPos   int64
			wantErr   bool
			readWant  []byte
		}{
			{
				name:     "seek from start",
				offset:   5,
				whence:   io.SeekStart,
				wantPos:  5,
				wantErr:  false,
				readWant: []byte("56789"),
			},
			{
				name:      "seek from current",
				preoffset: 5,
				offset:    2,
				whence:    io.SeekCurrent,
				wantPos:   7,
				wantErr:   false,
				readWant:  []byte("789"),
			},
			{
				name:     "seek from end",
				offset:   -3,
				whence:   io.SeekEnd,
				wantPos:  7,
				wantErr:  false,
				readWant: []byte("789"),
			},
			{
				name:     "seek beyond beginning",
				offset:   -5,
				whence:   io.SeekStart,
				wantPos:  0,
				wantErr:  true,
				readWant: nil,
			},
			{
				name:     "seek beyond end",
				offset:   15,
				whence:   io.SeekStart,
				wantPos:  10,
				wantErr:  false,
				readWant: []byte{},
			},
			{
				name:     "invalid whence",
				offset:   0,
				whence:   100,
				wantPos:  0,
				wantErr:  true,
				readWant: nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := NewReader(data)
				if tt.preoffset > 0 {
					buf := make([]byte, tt.preoffset)
					_, _ = r.Read(buf)
					_ = buf
				}

				pos, err := r.Seek(tt.offset, tt.whence)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}

				if !tt.wantErr {
					assert.Equal(t, tt.wantPos, pos)
				}

				if tt.readWant != nil {
					buf := make([]byte, len(tt.readWant))
					n, _ := r.Read(buf)
					assert.True(t, bytes.Equal(tt.readWant, buf[:n]), "want %q, got %q", tt.readWant, buf[:n])
				}
			})
		}
	})

	t.Run("seek extra", func(t *testing.T) {
		tests := []struct {
			name    string
			data    []byte
			preRead int
			offset  int64
			whence  int
			wantPos int64
			wantErr error
		}{
			{
				name:    "seek current zero is a no-op",
				data:    []byte("0123456789"),
				preRead: 3,
				offset:  0,
				whence:  io.SeekCurrent,
				wantPos: 3,
				wantErr: nil,
			},
			{
				name:    "seek end positive overshoot clamps to end",
				data:    []byte("0123456789"),
				offset:  5,
				whence:  io.SeekEnd,
				wantPos: 10,
				wantErr: nil,
			},
			{
				name:    "seek current past end clamps to end",
				data:    []byte("0123456789"),
				preRead: 8,
				offset:  5,
				whence:  io.SeekCurrent,
				wantPos: 10,
				wantErr: nil,
			},
			{
				name:    "seek to exact end",
				data:    []byte("0123456789"),
				offset:  10,
				whence:  io.SeekStart,
				wantPos: 10,
				wantErr: nil,
			},
			{
				name:    "negative resulting position",
				data:    []byte("0123456789"),
				offset:  -5,
				whence:  io.SeekStart,
				wantPos: 0,
				wantErr: ErrNegativePosition,
			},
			{
				name:    "seek end negative overshoot",
				data:    []byte("0123456789"),
				offset:  -15,
				whence:  io.SeekEnd,
				wantPos: 0,
				wantErr: ErrNegativePosition,
			},
			{
				name:    "invalid whence",
				data:    []byte("0123456789"),
				offset:  0,
				whence:  100,
				wantPos: 0,
				wantErr: ErrUnknownWhence,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := NewReader(tt.data)
				if tt.preRead > 0 {
					if _, err := r.Read(make([]byte, tt.preRead)); err != nil {
						require.NoError(t, err, "pre-read %d bytes", tt.preRead)
					}
				}

				pos, err := r.Seek(tt.offset, tt.whence)
				assert.Equal(t, tt.wantErr, err)
				assert.Equal(t, tt.wantPos, pos)

				if tt.wantErr == nil {
					cur, err := r.Seek(0, io.SeekCurrent)
					require.NoError(t, err, "Seek(0, SeekCurrent)")
					assert.Equal(t, tt.wantPos, cur, "reader offset after Seek")
				}
			})
		}
	})

	t.Run("combine", func(t *testing.T) {
		data := []byte("Hello, 世界! 🎉")
		r := NewReader(data)

		buf := make([]byte, 7)
		n, _ := r.Read(buf)
		assert.Equal(t, "Hello, ", string(buf[:n]))

		ch, size, _ := r.ReadRune()
		assert.Equal(t, '世', ch)
		assert.Equal(t, 3, size)

		err := r.UnreadRune()
		assert.NoError(t, err)

		ch, size, _ = r.ReadRune()
		assert.Equal(t, '世', ch)
		assert.Equal(t, 3, size)

		ch, size, _ = r.ReadRune()
		assert.Equal(t, '界', ch)
		assert.Equal(t, 3, size)

		pos, _ := r.Seek(0, io.SeekStart)
		assert.Equal(t, int64(0), pos)

		buf2 := make([]byte, 7)
		n, _ = r.ReadAt(buf2, 0)
		assert.Equal(t, "Hello, ", string(buf2[:n]))

		currentPos, _ := r.Seek(0, io.SeekCurrent)
		assert.Equal(t, int64(0), currentPos)
	})

	t.Run("big data", func(t *testing.T) {
		size := 10 * 1024 * 1024
		largeData := make([]byte, size)

		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		r := NewReader(largeData)

		chunkSize := 1024 * 1024 // 1 MB chunks
		for offset := 0; offset < size; offset += chunkSize {
			readSize := chunkSize
			if offset+chunkSize > size {
				readSize = size - offset
			}

			buf := make([]byte, readSize)
			n, err := r.Read(buf)

			if offset+readSize < size {
				assert.NoError(t, err, "Read() at offset %d", offset)
			}

			assert.Equal(t, readSize, n, "Read() at offset %d", offset)

			for i := 0; i < n; i++ {
				if buf[i] != largeData[offset+i] {
					assert.Equal(t, largeData[offset+i], buf[i],
						"Data mismatch at offset %d/%d", i, offset+i)
					break
				}
			}
		}

		buf := make([]byte, 1)
		_, err := r.Read(buf)
		assert.Equal(t, io.EOF, err, "after reading all data")

		t.Run("random access", func(t *testing.T) {
			size := 50 * 1024 * 1024 // 50 MB
			largeData := make([]byte, size)

			for i := range largeData {
				largeData[i] = byte(i % 256)
			}

			r := NewReader(largeData)

			positions := []int64{
				0,
				1,
				1000,
				1000000,
				int64(size / 2),
				int64(size - 1000),
				int64(size - 1),
			}

			for _, pos := range positions {
				t.Run("ReadAt position", func(t *testing.T) {
					buf := make([]byte, 100)
					n, err := r.ReadAt(buf, pos)

					if pos+100 < int64(size) {
						assert.NoError(t, err, "ReadAt(%d)", pos)
					}

					expectedSize := 100
					if pos+100 > int64(size) {
						expectedSize = size - int(pos)
					}

					assert.Equal(t, expectedSize, n, "ReadAt(%d)", pos)

					for i := 0; i < n; i++ {
						if buf[i] != largeData[pos+int64(i)] {
							assert.Equal(t, largeData[pos+int64(i)], buf[i],
								"ReadAt(%d) data mismatch at offset %d", pos, pos+int64(i))
							break
						}
					}
				})
			}
		})
	})
}

func BenchmarkReader(b *testing.B) {
	b.Run("read", func(b *testing.B) {
		size := 10 * 1024 * 1024 // 10 MB
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i % 256)
		}

		b.ReportAllocs()
		b.ResetTimer()
		b.SetBytes(int64(size))
		for i := 0; i < b.N; i++ {
			r := NewReader(data)
			buf := make([]byte, 1024)
			for {
				_, err := r.Read(buf)
				if err == io.EOF {
					break
				}
			}
		}
	})

	b.Run("read at", func(b *testing.B) {
		b.ReportAllocs()
		size := 10 * 1024 * 1024 // 10 MB
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i % 256)
		}

		r := NewReader(data)
		buf := make([]byte, 1024)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			offset := int64(i % size)
			_, _ = r.ReadAt(buf, offset)
		}
	})

	b.Run("read byte", func(b *testing.B) {
		size := 1024
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i % 256)
		}

		r := NewReader(data)
		b.ReportAllocs()
		b.ResetTimer()
		b.SetBytes(int64(size))
		for i := 0; i < b.N; i++ {
			_, _ = r.Seek(0, io.SeekStart)
			var sum byte
			for j := 0; j < size; j++ {
				c, err := r.ReadByte()
				if err != nil {
					b.FailNow()
				}
				sum += c
			}
			sinkByte = sum
		}
	})

	b.Run("read rune", func(b *testing.B) {
		size := 1024
		data := make([]byte, size)
		for i := range data {
			data[i] = byte('a' + i%26)
		}

		r := NewReader(data)
		b.ReportAllocs()
		b.ResetTimer()
		b.SetBytes(int64(size))
		for i := 0; i < b.N; i++ {
			_, _ = r.Seek(0, io.SeekStart)
			var sum rune
			for j := 0; j < size; j++ {
				c, _, err := r.ReadRune()
				if err != nil {
					b.FailNow()
				}
				sum += c
			}
			sinkRune = sum
		}
	})

	b.Run("seek", func(b *testing.B) {
		size := 10 * 1024 * 1024 // 10 MB
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i % 256)
		}

		r := NewReader(data)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pos, err := r.Seek(int64(i%size), io.SeekStart)
			if err != nil {
				b.FailNow()
			}
			sinkInt64 = pos
		}
	})
}
