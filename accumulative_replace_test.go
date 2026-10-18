package bytebuf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccumulativeReplace(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		old     []byte
		new     []byte
		n       int
		want    []byte
	}{
		{name: "n=0 no replacements", initial: "foo foo foo", old: []byte("foo"), new: []byte("bar"), n: 0, want: []byte("foo foo foo")},
		{name: "n=1 first occurrence only", initial: "foo foo foo", old: []byte("foo"), new: []byte("bar"), n: 1, want: []byte("bar foo foo")},
		{name: "n=2 first two occurrences", initial: "foo foo foo", old: []byte("foo"), new: []byte("bar"), n: 2, want: []byte("bar bar foo")},
		{name: "n=-1 all occurrences", initial: "foo foo foo", old: []byte("foo"), new: []byte("bar"), n: -1, want: []byte("bar bar bar")},
		{name: "n greater than occurrences", initial: "foo foo foo", old: []byte("foo"), new: []byte("bar"), n: 10, want: []byte("bar bar bar")},
		{name: "old not found", initial: "foo foo foo", old: []byte("xyz"), new: []byte("bar"), n: -1, want: []byte("foo foo foo")},
		{name: "old equals entire buffer", initial: "hello", old: []byte("hello"), new: []byte("world"), n: -1, want: []byte("world")},
		{name: "old longer than buffer", initial: "hi", old: []byte("hello"), new: []byte("world"), n: -1, want: []byte("hi")},
		{name: "new longer than old grows buffer", initial: "a-b-c", old: []byte("-"), new: []byte("--"), n: -1, want: []byte("a--b--c")},
		{name: "new shorter than old shrinks buffer", initial: "a---b---c", old: []byte("---"), new: []byte("-"), n: -1, want: []byte("a-b-c")},
		{name: "empty buffer", initial: "", old: []byte("foo"), new: []byte("bar"), n: -1, want: []byte("")},
		{name: "overlapping pattern non-overlapping matches", initial: "aaaa", old: []byte("aa"), new: []byte("aaa"), n: -1, want: []byte("aaaaaa")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteString(tt.initial).Replace(tt.old, tt.new, tt.n)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, Accumulative.Replace(%q, %q, %d) on %q", tt.want, ab.Bytes(), tt.old, tt.new, tt.n, tt.initial)
		})
	}
}

func TestAccumulativeReplaceString(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		old     string
		new     string
		n       int
		want    []byte
	}{
		{name: "n=0 no replacements", initial: "foo foo foo", old: "foo", new: "bar", n: 0, want: []byte("foo foo foo")},
		{name: "foo to bar n=1", initial: "foo foo foo", old: "foo", new: "bar", n: 1, want: []byte("bar foo foo")},
		{name: "n=2 first two occurrences", initial: "foo foo foo", old: "foo", new: "bar", n: 2, want: []byte("bar bar foo")},
		{name: "foo to bar n=-1", initial: "foo foo foo", old: "foo", new: "bar", n: -1, want: []byte("bar bar bar")},
		{name: "n greater than occurrences", initial: "foo foo foo", old: "foo", new: "bar", n: 10, want: []byte("bar bar bar")},
		{name: "old not found", initial: "foo foo foo", old: "xyz", new: "bar", n: -1, want: []byte("foo foo foo")},
		{name: "old equals entire buffer", initial: "hello", old: "hello", new: "world", n: -1, want: []byte("world")},
		{name: "old longer than buffer", initial: "hi", old: "hello", new: "world", n: -1, want: []byte("hi")},
		{name: "new longer than old grows buffer", initial: "a-b-c", old: "-", new: "--", n: -1, want: []byte("a--b--c")},
		{name: "new shorter than old shrinks buffer", initial: "a---b---c", old: "---", new: "-", n: -1, want: []byte("a-b-c")},
		{name: "empty buffer", initial: "", old: "foo", new: "bar", n: -1, want: []byte("")},
		{name: "overlapping pattern non-overlapping matches", initial: "aaaa", old: "aa", new: "aaa", n: -1, want: []byte("aaaaaa")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteString(tt.initial).ReplaceString(tt.old, tt.new, tt.n)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, Accumulative.ReplaceString(%q, %q, %d) on %q", tt.want, ab.Bytes(), tt.old, tt.new, tt.n, tt.initial)
		})
	}
}

func TestAccumulativeReplaceAll(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		old     []byte
		new     []byte
		want    []byte
	}{
		{name: "replaces all occurrences", initial: "foo foo foo", old: []byte("foo"), new: []byte("bar"), want: []byte("bar bar bar")},
		{name: "old not found", initial: "foo foo foo", old: []byte("xyz"), new: []byte("bar"), want: []byte("foo foo foo")},
		{name: "empty buffer", initial: "", old: []byte("foo"), new: []byte("bar"), want: []byte("")},
		{name: "grows buffer", initial: "a-b-c", old: []byte("-"), new: []byte("--"), want: []byte("a--b--c")},
		{name: "shrinks buffer", initial: "a---b---c", old: []byte("---"), new: []byte("-"), want: []byte("a-b-c")},
		{name: "old equals entire buffer", initial: "hello", old: []byte("hello"), new: []byte("world"), want: []byte("world")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteString(tt.initial).ReplaceAll(tt.old, tt.new)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, Accumulative.ReplaceAll(%q, %q) on %q", tt.want, ab.Bytes(), tt.old, tt.new, tt.initial)
		})
	}
}

func TestAccumulativeReplaceStringAll(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		old     string
		new     string
		want    []byte
	}{
		{name: "replaces all occurrences", initial: "foo foo foo", old: "foo", new: "bar", want: []byte("bar bar bar")},
		{name: "old not found", initial: "foo foo foo", old: "xyz", new: "bar", want: []byte("foo foo foo")},
		{name: "empty buffer", initial: "", old: "foo", new: "bar", want: []byte("")},
		{name: "grows buffer", initial: "a-b-c", old: "-", new: "--", want: []byte("a--b--c")},
		{name: "shrinks buffer", initial: "a---b---c", old: "---", new: "-", want: []byte("a-b-c")},
		{name: "old equals entire buffer", initial: "hello", old: "hello", new: "world", want: []byte("world")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := NewAccumulativeSize(128)
			ab.WriteString(tt.initial).ReplaceStringAll(tt.old, tt.new)
			assert.True(t, bytes.Equal(tt.want, ab.Bytes()), "want %q, got %q, Accumulative.ReplaceStringAll(%q, %q) on %q", tt.want, ab.Bytes(), tt.old, tt.new, tt.initial)
		})
	}
}
