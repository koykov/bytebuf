package bytebuf

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// failingReader returns its data on the first Read call and the configured
// error on every subsequent call, mimicking a stream that fails mid-read.
type failingReader struct {
	data []byte
	err  error
	read bool
}

func (r *failingReader) Read(p []byte) (int, error) {
	if r.read {
		return 0, r.err
	}
	r.read = true
	return copy(p, r.data), nil
}

func TestWriterReadFrom(t *testing.T) {
	errRead := errors.New("read error")

	tests := []struct {
		name    string
		reader  io.Reader
		want    []byte
		wantErr bool
	}{
		{
			name:    "bytes.Reader with short data",
			reader:  bytes.NewReader([]byte("hello")),
			want:    []byte("hello"),
			wantErr: false,
		},
		{
			name:    "bytes.Reader with large data",
			reader:  bytes.NewReader(bytes.Repeat([]byte("x"), 600)),
			want:    bytes.Repeat([]byte("x"), 600),
			wantErr: false,
		},
		{
			name:    "strings.Reader",
			reader:  strings.NewReader("hello world"),
			want:    []byte("hello world"),
			wantErr: false,
		},
		{
			name:    "empty reader",
			reader:  bytes.NewReader(nil),
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "reader errors mid-stream",
			reader:  &failingReader{data: []byte("partial data"), err: errRead},
			want:    []byte("partial data"),
			wantErr: true,
		},
		{
			name:    "reader returns exactly bytes.MinRead",
			reader:  bytes.NewReader(bytes.Repeat([]byte("y"), bytes.MinRead)),
			want:    bytes.Repeat([]byte("y"), bytes.MinRead),
			wantErr: false,
		},
		{
			name:    "reader returns more than bytes.MinRead",
			reader:  bytes.NewReader(bytes.Repeat([]byte("z"), bytes.MinRead*2+13)),
			want:    bytes.Repeat([]byte("z"), bytes.MinRead*2+13),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWriter(nil)
			n, err := w.ReadFrom(tt.reader)

			if tt.wantErr {
				require.Error(t, err, "ReadFrom() error")
			} else {
				require.NoError(t, err, "ReadFrom() error")
			}

			if tt.wantErr {
				assert.ErrorIs(t, err, errRead, "ReadFrom() error")
			}

			assert.Equal(t, int64(len(tt.want)), n, "ReadFrom() bytes read")

			assert.True(t, bytes.Equal(tt.want, w.Bytes()), "Bytes(): want %q, got %q", tt.want, w.Bytes())
		})
	}
}
