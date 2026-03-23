package bytebuf

import (
	"bytes"
	"io"
	"testing"
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

				if (err != nil) != tt.wantErr {
					t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if n != len(tt.want) {
					t.Errorf("Read() read %d bytes, want %d", n, len(tt.want))
				}

				if !bytes.Equal(buf[:n], tt.want) {
					t.Errorf("Read() got %v, want %v", buf[:n], tt.want)
				}
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

				if err != tt.wantErr {
					t.Errorf("ReadAt() error = %v, wantErr %v", err, tt.wantErr)
				}

				if n != len(tt.want) {
					t.Errorf("ReadAt() read %d bytes, want %d", n, len(tt.want))
				}

				if !bytes.Equal(buf[:n], tt.want) {
					t.Errorf("ReadAt() got %v, want %v", buf[:n], tt.want)
				}
			})
		}
	})

	t.Run("read byte", func(t *testing.T) {
		data := []byte("abc")
		r := NewReader(data)

		for i, expected := range data {
			b, err := r.ReadByte()
			if err != nil {
				t.Errorf("ReadByte() at index %d error = %v", i, err)
			}
			if b != expected {
				t.Errorf("ReadByte() at index %d got %c, want %c", i, b, expected)
			}
		}

		_, err := r.ReadByte()
		if err != io.EOF {
			t.Errorf("ReadByte() after EOF error = %v, want io.EOF", err)
		}
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
					if err != nil {
						t.Fatalf("Failed to read: %v", err)
					}
				}

				var err error
				for i := 0; i < tt.unreadCount; i++ {
					err = r.UnreadByte()
				}

				if (err != nil) != tt.wantErr {
					t.Errorf("UnreadByte() error = %v, wantErr %v", err, tt.wantErr)
				}

				if !tt.wantErr {
					buf := make([]byte, len(tt.wantAfter))
					n, _ := r.Read(buf)
					if !bytes.Equal(buf[:n], tt.wantAfter) {
						t.Errorf("After unread got %v, want %v", buf[:n], tt.wantAfter)
					}
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
					if err != nil {
						t.Errorf("ReadRune() at index %d error = %v", i, err)
					}
					if ch != expectedRune {
						t.Errorf("ReadRune() got %c, want %c", ch, expectedRune)
					}
					if size != tt.expectedSizes[i] {
						t.Errorf("ReadRune() size = %d, want %d", size, tt.expectedSizes[i])
					}
				}

				_, _, err := r.ReadRune()
				if err != io.EOF {
					t.Errorf("ReadRune() after EOF error = %v, want io.EOF", err)
				}
			})
		}
	})

	t.Run("unread rune", func(t *testing.T) {
		data := "Hello, 世界"
		r := NewReader([]byte(data))

		first, _, _ := r.ReadRune()
		if first != 'H' {
			t.Errorf("First rune = %c, want H", first)
		}

		err := r.UnreadRune()
		if err != nil {
			t.Errorf("UnreadRune() error = %v", err)
		}

		again, _, _ := r.ReadRune()
		if again != 'H' {
			t.Errorf("After unread got %c, want H", again)
		}

		_ = r.UnreadRune()
		err = r.UnreadRune()
		if err == nil {
			t.Error("Second UnreadRune() should fail")
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
				if (err != nil) != tt.wantErr {
					t.Errorf("Seek() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !tt.wantErr && pos != tt.wantPos {
					t.Errorf("Seek() position = %d, want %d", pos, tt.wantPos)
				}

				if tt.readWant != nil {
					buf := make([]byte, len(tt.readWant))
					n, _ := r.Read(buf)
					if !bytes.Equal(buf[:n], tt.readWant) {
						t.Errorf("After seek read = %v, want %v", buf[:n], tt.readWant)
					}
				}
			})
		}
	})

	t.Run("combine", func(t *testing.T) {
		data := []byte("Hello, 世界! 🎉")
		r := NewReader(data)

		buf := make([]byte, 7)
		n, _ := r.Read(buf)
		if string(buf[:n]) != "Hello, " {
			t.Errorf("Read() got %s, want 'Hello, '", buf[:n])
		}

		ch, size, _ := r.ReadRune()
		if ch != '世' || size != 3 {
			t.Errorf("ReadRune() got %c (size %d), want 世 (size 3)", ch, size)
		}

		err := r.UnreadRune()
		if err != nil {
			t.Errorf("UnreadRune() error = %v", err)
		}

		ch, size, _ = r.ReadRune()
		if ch != '世' || size != 3 {
			t.Errorf("After unread got %c (size %d), want 世 (size 3)", ch, size)
		}

		ch, size, _ = r.ReadRune()
		if ch != '界' || size != 3 {
			t.Errorf("ReadRune() got %c (size %d), want 界 (size 3)", ch, size)
		}

		pos, _ := r.Seek(0, io.SeekStart)
		if pos != 0 {
			t.Errorf("Seek to start got %d", pos)
		}

		buf2 := make([]byte, 7)
		n, _ = r.ReadAt(buf2, 0)
		if string(buf2[:n]) != "Hello, " {
			t.Errorf("ReadAt() got %s, want 'Hello, '", buf2[:n])
		}

		currentPos, _ := r.Seek(0, io.SeekCurrent)
		if currentPos != 0 {
			t.Errorf("After ReadAt position = %d, want 0", currentPos)
		}
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

			if err != nil && offset+readSize < size {
				t.Errorf("Read() at offset %d error = %v", offset, err)
			}

			if n != readSize {
				t.Errorf("Read() at offset %d read %d bytes, want %d", offset, n, readSize)
			}

			for i := 0; i < n; i++ {
				if buf[i] != largeData[offset+i] {
					t.Errorf("Data mismatch at offset %d/%d: got %d, want %d",
						i, offset+i, buf[i], largeData[offset+i])
					break
				}
			}
		}

		buf := make([]byte, 1)
		_, err := r.Read(buf)
		if err != io.EOF {
			t.Errorf("After reading all data, Read() error = %v, want io.EOF", err)
		}

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

					if err != nil && pos+100 < int64(size) {
						t.Errorf("ReadAt(%d) error = %v", pos, err)
					}

					expectedSize := 100
					if pos+100 > int64(size) {
						expectedSize = size - int(pos)
					}

					if n != expectedSize {
						t.Errorf("ReadAt(%d) read %d bytes, want %d", pos, n, expectedSize)
					}

					for i := 0; i < n; i++ {
						if buf[i] != largeData[pos+int64(i)] {
							t.Errorf("ReadAt(%d) data mismatch at offset %d", pos, pos+int64(i))
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
}
