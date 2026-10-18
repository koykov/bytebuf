package bytebuf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainReplace(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		old  []byte
		new  []byte
		n    int
		want []byte
	}{
		{name: "n zero no replacements", in: []byte("foo foo foo"), old: []byte("foo"), new: []byte("bar"), n: 0, want: []byte("foo foo foo")},
		{name: "n one first occurrence only", in: []byte("foo foo foo"), old: []byte("foo"), new: []byte("bar"), n: 1, want: []byte("bar foo foo")},
		{name: "n two first two occurrences", in: []byte("foo foo foo"), old: []byte("foo"), new: []byte("bar"), n: 2, want: []byte("bar bar foo")},
		{name: "n negative replaces all", in: []byte("foo foo foo"), old: []byte("foo"), new: []byte("bar"), n: -1, want: []byte("bar bar bar")},
		{name: "n greater than occurrences replaces all", in: []byte("foo foo foo"), old: []byte("foo"), new: []byte("bar"), n: 10, want: []byte("bar bar bar")},
		{name: "old not found no change", in: []byte("hello world"), old: []byte("xyz"), new: []byte("abc"), n: -1, want: []byte("hello world")},
		{name: "old equals entire buffer replaced", in: []byte("abc"), old: []byte("abc"), new: []byte("xyz"), n: 1, want: []byte("xyz")},
		{name: "old longer than buffer no change", in: []byte("ab"), old: []byte("abc"), new: []byte("x"), n: -1, want: []byte("ab")},
		{name: "new longer than old grows", in: []byte("a-b"), old: []byte("-"), new: []byte("---"), n: -1, want: []byte("a---b")},
		{name: "new shorter than old shrinks", in: []byte("a---b"), old: []byte("---"), new: []byte("-"), n: -1, want: []byte("a-b")},
		{name: "empty buffer no change", in: []byte(""), old: []byte("a"), new: []byte("b"), n: -1, want: []byte("")},
		{name: "overlapping pattern aa to aaa", in: []byte("aaaa"), old: []byte("aa"), new: []byte("aaa"), n: -1, want: []byte("aaaaaa")},
		{name: "single byte replaced everywhere", in: []byte("a.b.c"), old: []byte("."), new: []byte("_"), n: -1, want: []byte("a_b_c")},
		{name: "empty new deletes occurrences", in: []byte("aXbXc"), old: []byte("X"), new: []byte(""), n: -1, want: []byte("abc")},
		{name: "growth beyond capacity reallocates", in: bytes.Repeat([]byte("a"), 100), old: []byte("a"), new: []byte("bb"), n: -1, want: bytes.Repeat([]byte("b"), 200)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.Write(tt.in).Replace(tt.old, tt.new, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, Replace(%q, %q, %q, %d)", tt.want, cb.Bytes(), tt.in, tt.old, tt.new, tt.n)
		})
	}
}

func TestChainReplaceString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		old  string
		new  string
		n    int
		want []byte
	}{
		{name: "n zero no replacements", in: "foo foo foo", old: "foo", new: "bar", n: 0, want: []byte("foo foo foo")},
		{name: "n one first occurrence only", in: "foo foo foo", old: "foo", new: "bar", n: 1, want: []byte("bar foo foo")},
		{name: "n two first two occurrences", in: "foo foo foo", old: "foo", new: "bar", n: 2, want: []byte("bar bar foo")},
		{name: "n negative replaces all", in: "foo foo foo", old: "foo", new: "bar", n: -1, want: []byte("bar bar bar")},
		{name: "n greater than occurrences replaces all", in: "foo foo foo", old: "foo", new: "bar", n: 5, want: []byte("bar bar bar")},
		{name: "old not found no change", in: "hello world", old: "xyz", new: "abc", n: -1, want: []byte("hello world")},
		{name: "old equals entire buffer replaced", in: "abc", old: "abc", new: "xyz", n: 1, want: []byte("xyz")},
		{name: "old longer than buffer no change", in: "ab", old: "abc", new: "x", n: -1, want: []byte("ab")},
		{name: "new longer than old grows", in: "a-b", old: "-", new: "---", n: -1, want: []byte("a---b")},
		{name: "new shorter than old shrinks", in: "a---b", old: "---", new: "-", n: -1, want: []byte("a-b")},
		{name: "empty buffer no change", in: "", old: "a", new: "b", n: -1, want: []byte("")},
		{name: "overlapping pattern aa to aaa", in: "aaaa", old: "aa", new: "aaa", n: -1, want: []byte("aaaaaa")},
		{name: "unicode multibyte replaced", in: "héllo", old: "é", new: "e", n: -1, want: []byte("hello")},
		{name: "empty new deletes occurrences", in: "aXbXc", old: "X", new: "", n: -1, want: []byte("abc")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteString(tt.in).ReplaceString(tt.old, tt.new, tt.n)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, ReplaceString(%q, %q, %q, %d)", tt.want, cb.Bytes(), tt.in, tt.old, tt.new, tt.n)
		})
	}
}

func TestChainReplaceAll(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		old  []byte
		new  []byte
		want []byte
	}{
		{name: "replaces all occurrences", in: []byte("foo foo foo"), old: []byte("foo"), new: []byte("bar"), want: []byte("bar bar bar")},
		{name: "single occurrence replaced", in: []byte("a.b"), old: []byte("."), new: []byte("_"), want: []byte("a_b")},
		{name: "old not found no change", in: []byte("hello world"), old: []byte("xyz"), new: []byte("abc"), want: []byte("hello world")},
		{name: "empty buffer no change", in: []byte(""), old: []byte("a"), new: []byte("b"), want: []byte("")},
		{name: "old equals entire buffer replaced", in: []byte("abc"), old: []byte("abc"), new: []byte("xyz"), want: []byte("xyz")},
		{name: "old longer than buffer no change", in: []byte("ab"), old: []byte("abc"), new: []byte("x"), want: []byte("ab")},
		{name: "new longer than old grows", in: []byte("a-b"), old: []byte("-"), new: []byte("--"), want: []byte("a--b")},
		{name: "new shorter than old shrinks", in: []byte("a--b"), old: []byte("--"), new: []byte("-"), want: []byte("a-b")},
		{name: "overlapping pattern aa to aaa", in: []byte("aaaa"), old: []byte("aa"), new: []byte("aaa"), want: []byte("aaaaaa")},
		{name: "empty new deletes occurrences", in: []byte("aXbXc"), old: []byte("X"), new: []byte(""), want: []byte("abc")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.Write(tt.in).ReplaceAll(tt.old, tt.new)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, ReplaceAll(%q, %q, %q)", tt.want, cb.Bytes(), tt.in, tt.old, tt.new)
		})
	}
}

func TestChainReplaceStringAll(t *testing.T) {
	tests := []struct {
		name string
		in   string
		old  string
		new  string
		want []byte
	}{
		{name: "replaces all occurrences", in: "foo foo foo", old: "foo", new: "bar", want: []byte("bar bar bar")},
		{name: "single occurrence replaced", in: "a.b", old: ".", new: "_", want: []byte("a_b")},
		{name: "old not found no change", in: "hello world", old: "xyz", new: "abc", want: []byte("hello world")},
		{name: "empty buffer no change", in: "", old: "a", new: "b", want: []byte("")},
		{name: "old equals entire buffer replaced", in: "abc", old: "abc", new: "xyz", want: []byte("xyz")},
		{name: "old longer than buffer no change", in: "ab", old: "abc", new: "x", want: []byte("ab")},
		{name: "new longer than old grows", in: "a-b", old: "-", new: "--", want: []byte("a--b")},
		{name: "new shorter than old shrinks", in: "a--b", old: "--", new: "-", want: []byte("a-b")},
		{name: "overlapping pattern aa to aaa", in: "aaaa", old: "aa", new: "aaa", want: []byte("aaaaaa")},
		{name: "unicode multibyte replaced", in: "héllo", old: "é", new: "e", want: []byte("hello")},
		{name: "empty new deletes occurrences", in: "aXbXc", old: "X", new: "", want: []byte("abc")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChainSize(128)
			cb.WriteString(tt.in).ReplaceStringAll(tt.old, tt.new)
			assert.True(t, bytes.Equal(tt.want, cb.Bytes()), "want %q, got %q, ReplaceStringAll(%q, %q, %q)", tt.want, cb.Bytes(), tt.in, tt.old, tt.new)
		})
	}
}
