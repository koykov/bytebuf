package bytebuf

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccumulativeStakeOut(t *testing.T) {
	type step struct {
		stake bool
		write string
	}

	var stages = []struct {
		name    string
		steps   []step
		wantOff int
		wantLen int
	}{
		{
			name:    "stake before any write",
			steps:   []step{{stake: true}},
			wantOff: 0,
			wantLen: 0,
		},
		{
			name:    "hello written after stake",
			steps:   []step{{stake: true}, {write: "hello"}},
			wantOff: 0,
			wantLen: 5,
		},
		{
			name:    "restake then world",
			steps:   []step{{stake: true}, {write: "hello"}, {stake: true}, {write: "world"}},
			wantOff: 5,
			wantLen: 5,
		},
		{
			name:    "restake then !!!",
			steps:   []step{{stake: true}, {write: "hello"}, {stake: true}, {write: "world"}, {stake: true}, {write: "!!!"}},
			wantOff: 10,
			wantLen: 3,
		},
	}

	testfn := func(t *testing.T, ab *Accumulative, steps []step) {
		t.Helper()
		for _, st := range steps {
			if st.stake {
				assert.Same(t, ab, ab.StakeOut(), "StakeOut() returned wrong receiver")
			}
			if st.write != "" {
				ab.WriteString(st.write)
			}
		}
	}

	t.Run("", func(t *testing.T) {
		for _, tt := range stages {
			t.Run(tt.name, func(t *testing.T) {
				ab := NewAccumulativeSize(128)
				testfn(t, ab, tt.steps)

				assert.Equal(t, tt.wantOff, ab.StakedOffset(), "StakedOffset()")
				assert.Equal(t, tt.wantLen, ab.StakedLen(), "StakedLen()")
			})
		}
	})

	t.Run("offset", func(t *testing.T) {
		for _, tt := range stages {
			t.Run(tt.name, func(t *testing.T) {
				ab := NewAccumulativeSize(128)
				testfn(t, ab, tt.steps)

				assert.Equal(t, tt.wantOff, ab.StakedOffset(), "StakedOffset()")
			})
		}
	})

	t.Run("len", func(t *testing.T) {
		for _, tt := range stages {
			t.Run(tt.name, func(t *testing.T) {
				ab := NewAccumulativeSize(128)
				testfn(t, ab, tt.steps)

				assert.Equal(t, tt.wantLen, ab.StakedLen(), "StakedLen()")
			})
		}
	})
}

func TestAccumulativeStakedBytes(t *testing.T) {
	tests := []struct {
		name    string
		before  []string
		stake   bool
		after   []string
		want    []byte
		wantNil bool
	}{
		{
			name:  "hello staked after stake",
			stake: true,
			after: []string{"hello"},
			want:  []byte("hello"),
		},
		{
			name:   "world staked after restake",
			before: []string{"hello"},
			stake:  true,
			after:  []string{"world"},
			want:   []byte("world"),
		},
		{
			name:    "nothing staked returns nil",
			stake:   true,
			wantNil: true,
		},
		{
			name:   "fixture payload after foobar stake",
			before: []string{"foobar"},
			stake:  true,
			after:  []string{string(expectWS)},
			want:   expectWS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			for _, w := range tt.before {
				ab.WriteString(w)
			}
			if tt.stake {
				ab.StakeOut()
			}
			for _, w := range tt.after {
				ab.WriteString(w)
			}

			got := ab.StakedBytes()
			if tt.wantNil {
				assert.Nil(t, got, "StakedBytes()")
				return
			}
			assert.True(t, bytes.Equal(tt.want, got), "want %q, got %q, StakedBytes()", tt.want, got)
		})
	}
}

func TestAccumulativeStakedBytesCopy(t *testing.T) {
	tests := []struct {
		name    string
		before  []string
		after   []string
		want    []byte
		wantNil bool
	}{
		{
			name:  "copy of hello",
			after: []string{"hello"},
			want:  []byte("hello"),
		},
		{
			name:   "copy of world after restake",
			before: []string{"hello"},
			after:  []string{"world"},
			want:   []byte("world"),
		},
		{
			name:    "nothing staked returns nil",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			for _, w := range tt.before {
				ab.WriteString(w)
			}
			ab.StakeOut()
			for _, w := range tt.after {
				ab.WriteString(w)
			}

			got := ab.StakedBytesCopy()
			if tt.wantNil {
				assert.Nil(t, got, "StakedBytesCopy()")
				return
			}
			assert.True(t, bytes.Equal(tt.want, got), "want %q, got %q, StakedBytesCopy()", tt.want, got)
		})
	}

	t.Run("mutating the returned copy does not affect the buffer", func(t *testing.T) {
		ab := NewAccumulativeSize(128)
		ab.StakeOut()
		ab.WriteString("hello")

		cp := ab.StakedBytesCopy()
		cp[0] = 'X'

		assert.True(t, bytes.Equal([]byte("hello"), ab.StakedBytes()), "want %q, got %q, StakedBytes() after mutating copy", []byte("hello"), ab.StakedBytes())
	})
}

func TestAccumulativeStakedString(t *testing.T) {
	tests := []struct {
		name   string
		before []string
		after  []string
		want   string
	}{
		{
			name:  "hello staked after stake",
			after: []string{"hello"},
			want:  "hello",
		},
		{
			name:   "world staked after restake",
			before: []string{"hello"},
			after:  []string{"world"},
			want:   "world",
		},
		{
			name: "nothing staked returns empty string",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			for _, w := range tt.before {
				ab.WriteString(w)
			}
			ab.StakeOut()
			for _, w := range tt.after {
				ab.WriteString(w)
			}

			assert.Equal(t, tt.want, ab.StakedString(), "StakedString()")
		})
	}
}

func TestAccumulativeStakedStringCopy(t *testing.T) {
	tests := []struct {
		name   string
		before []string
		after  []string
		want   string
	}{
		{
			name:  "copy of hello",
			after: []string{"hello"},
			want:  "hello",
		},
		{
			name:   "copy of world after restake",
			before: []string{"hello"},
			after:  []string{"world"},
			want:   "world",
		},
		{
			name: "nothing staked returns empty string",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			for _, w := range tt.before {
				ab.WriteString(w)
			}
			ab.StakeOut()
			for _, w := range tt.after {
				ab.WriteString(w)
			}

			assert.Equal(t, tt.want, ab.StakedStringCopy(), "StakedStringCopy()")
		})
	}

	t.Run("copy survives overwriting the underlying buffer", func(t *testing.T) {
		ab := NewAccumulativeSize(128)
		ab.StakeOut()
		ab.WriteString("hello")

		sc := ab.StakedStringCopy()

		// Overwrite the region the source bytes came from.
		ab.Reset()
		ab.WriteString("XXXXX")

		assert.Equal(t, "hello", sc, "StakedStringCopy() after buffer overwrite")
	})
}

func TestAccumulativeError(t *testing.T) {
	errSentinel := errors.New("boom")

	tests := []struct {
		name        string
		writes      []string
		setErr      error
		doReset     bool
		doResetSafe bool
		want        error
	}{
		{
			name: "no writes yields nil",
			want: nil,
		},
		{
			name:   "successful writes yield nil",
			writes: []string{"hello", "world"},
			want:   nil,
		},
		{
			name:    "reset after successful writes still yields nil",
			writes:  []string{"hello"},
			doReset: true,
			want:    nil,
		},
		{
			// Reset() now clears err alongside resetting the inner buffer.
			// err is injected directly to observe this, since no plain Write* call
			// is guaranteed to produce an error with these inputs.
			name:    "reset clears a previously set error",
			setErr:  errSentinel,
			doReset: true,
			want:    nil,
		},
		{
			// ResetSafe() clears err alongside resetting the inner buffer.
			name:        "reset safe clears a previously set error",
			setErr:      errSentinel,
			doResetSafe: true,
			want:        nil,
		},
		{
			// Without a reset, err survives until one is performed.
			name:        "error preserved without reset",
			setErr:      errSentinel,
			doReset:     false,
			doResetSafe: false,
			want:        errSentinel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			if tt.setErr != nil {
				ab.err = tt.setErr
			}
			for _, w := range tt.writes {
				ab.WriteString(w)
			}
			if tt.doReset {
				ab.Reset()
			}
			if tt.doResetSafe {
				ab.ResetSafe()
			}

			assert.Equal(t, tt.want, ab.Error(), "Error()")
		})
	}
}

// The Range* methods are KNOWN BUGGY: the guard in RangeBytes/RangeString is
// inverted relative to the bounds it protects. These tests document the current
// (buggy) behavior rather than the intended behavior:
//
//   - a request that lies strictly inside the buffer returns nil/"" instead of data;
//   - a request whose end equals the buffer length returns the slice (guard falls through);
//   - a request whose end exceeds the buffer capacity, or a negative offset, panics.
//
// The buffer under test has cap == len == 10 so slice bounds are strict.
func rangeTestBuffer() (ab *Accumulative) {
	return NewAccumulative([]byte("helloworld"))
}

func TestAccumulativeRangeBytes(t *testing.T) {
	tests := []struct {
		name      string
		off       int
		length    int
		want      []byte
		wantNil   bool
		wantPanic bool
	}{
		{
			name:    "in-range start returns nil (inverted guard bug)",
			off:     0,
			length:  5,
			wantNil: true,
		},
		{
			name:    "in-range middle returns nil (inverted guard bug)",
			off:     2,
			length:  3,
			wantNil: true,
		},
		{
			name:   "reads up to buffer end returns full slice",
			off:    0,
			length: 10,
			want:   []byte("helloworld"),
		},
		{
			name:   "tail reads up to buffer end returns slice",
			off:    5,
			length: 5,
			want:   []byte("world"),
		},
		{
			name:      "end beyond buffer capacity panics",
			off:       0,
			length:    100,
			wantPanic: true,
		},
		{
			name:      "negative offset panics",
			off:       -1,
			length:    3,
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := rangeTestBuffer()
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						assert.Failf(t, "did not panic", "RangeBytes(%d, %d) did not panic, want panic", tt.off, tt.length)
					}
				}()
			}

			got := ab.RangeBytes(tt.off, tt.length)
			if tt.wantPanic {
				return
			}
			if tt.wantNil {
				assert.Nil(t, got, "RangeBytes(%d, %d)", tt.off, tt.length)
				return
			}
			assert.True(t, bytes.Equal(tt.want, got), "want %q, got %q, RangeBytes(%d, %d)", tt.want, got, tt.off, tt.length)
		})
	}
}

func TestAccumulativeRangeBytesCopy(t *testing.T) {
	tests := []struct {
		name      string
		off       int
		length    int
		want      []byte
		wantNil   bool
		wantPanic bool
	}{
		{
			name:    "in-range start returns nil (inverted guard bug)",
			off:     0,
			length:  5,
			wantNil: true,
		},
		{
			name:    "in-range middle returns nil (inverted guard bug)",
			off:     2,
			length:  3,
			wantNil: true,
		},
		{
			name:   "reads up to buffer end returns full copy",
			off:    0,
			length: 10,
			want:   []byte("helloworld"),
		},
		{
			name:   "tail reads up to buffer end returns copy",
			off:    5,
			length: 5,
			want:   []byte("world"),
		},
		{
			name:      "end beyond buffer capacity panics",
			off:       0,
			length:    100,
			wantPanic: true,
		},
		{
			name:      "negative offset panics",
			off:       -1,
			length:    3,
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := rangeTestBuffer()
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						assert.Failf(t, "did not panic", "RangeBytesCopy(%d, %d) did not panic, want panic", tt.off, tt.length)
					}
				}()
			}

			got := ab.RangeBytesCopy(tt.off, tt.length)
			if tt.wantPanic {
				return
			}
			if tt.wantNil {
				assert.Nil(t, got, "RangeBytesCopy(%d, %d)", tt.off, tt.length)
				return
			}
			assert.True(t, bytes.Equal(tt.want, got), "want %q, got %q, RangeBytesCopy(%d, %d)", tt.want, got, tt.off, tt.length)
		})
	}
}

func TestAccumulativeRangeString(t *testing.T) {
	tests := []struct {
		name      string
		off       int
		length    int
		want      string
		wantPanic bool
	}{
		{
			name:   "in-range start returns empty (inverted guard bug)",
			off:    0,
			length: 5,
			want:   "",
		},
		{
			name:   "in-range middle returns empty (inverted guard bug)",
			off:    2,
			length: 3,
			want:   "",
		},
		{
			name:   "reads up to buffer end returns full string",
			off:    0,
			length: 10,
			want:   "helloworld",
		},
		{
			name:   "tail reads up to buffer end returns string",
			off:    5,
			length: 5,
			want:   "world",
		},
		{
			name:      "end beyond buffer capacity panics",
			off:       0,
			length:    100,
			wantPanic: true,
		},
		{
			name:      "negative offset panics",
			off:       -1,
			length:    3,
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := rangeTestBuffer()
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						assert.Failf(t, "did not panic", "RangeString(%d, %d) did not panic, want panic", tt.off, tt.length)
					}
				}()
			}

			got := ab.RangeString(tt.off, tt.length)
			if tt.wantPanic {
				return
			}
			assert.Equal(t, tt.want, got, "RangeString(%d, %d)", tt.off, tt.length)
		})
	}
}

func TestAccumulativeRangeStringCopy(t *testing.T) {
	tests := []struct {
		name      string
		off       int
		length    int
		want      string
		wantPanic bool
	}{
		{
			name:   "in-range start returns empty (inverted guard bug)",
			off:    0,
			length: 5,
			want:   "",
		},
		{
			name:   "in-range middle returns empty (inverted guard bug)",
			off:    2,
			length: 3,
			want:   "",
		},
		{
			name:   "reads up to buffer end returns full copy",
			off:    0,
			length: 10,
			want:   "helloworld",
		},
		{
			name:   "tail reads up to buffer end returns copy",
			off:    5,
			length: 5,
			want:   "world",
		},
		{
			name:      "end beyond buffer capacity panics",
			off:       0,
			length:    100,
			wantPanic: true,
		},
		{
			name:      "negative offset panics",
			off:       -1,
			length:    3,
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := rangeTestBuffer()
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						assert.Failf(t, "did not panic", "RangeStringCopy(%d, %d) did not panic, want panic", tt.off, tt.length)
					}
				}()
			}

			got := ab.RangeStringCopy(tt.off, tt.length)
			if tt.wantPanic {
				return
			}
			assert.Equal(t, tt.want, got, "RangeStringCopy(%d, %d)", tt.off, tt.length)
		})
	}
}

func BenchmarkAccumulativeStaking(b *testing.B) {
	initbuf := func() *Accumulative {
		ab := NewAccumulativeSize(128)
		ab.WriteString("foobar").StakeOut().Write(stage.b)
		return ab
	}

	b.Run("stake out", func(b *testing.B) {
		ab := initbuf()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if ab.StakeOut() != ab {
				b.Error("Accumulative.StakeOut: unexpected receiver")
			}
		}
	})
	b.Run("staked bytes", func(b *testing.B) {
		ab := initbuf()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sinkBytes = ab.StakedBytes()
		}
	})
	b.Run("staked string", func(b *testing.B) {
		ab := initbuf()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sinkString = ab.StakedString()
		}
	})
	b.Run("staked offset", func(b *testing.B) {
		ab := initbuf()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sinkInt = ab.StakedOffset()
		}
	})
	b.Run("staked len", func(b *testing.B) {
		ab := initbuf()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sinkInt = ab.StakedLen()
		}
	})
	b.Run("stacked entry64", func(b *testing.B) {
		ab := initbuf()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sinkInt64 = int64(ab.StackedEntry64())
		}
	})
}
