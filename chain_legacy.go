package bytebuf

import "github.com/koykov/byteconv"

// ChainBuf is an alias of Chain type.
// DEPRECATED: use Chain instead. Kept for backward compatibility.
type ChainBuf = Chain

// WriteStr writes string to the buffer.
// DEPRECATED: use WriteString() instead.
func (b *Chain) WriteStr(s string) *Chain {
	*b = append(*b, s...)
	return b
}

// WriteApplyFnStr applies fn to s and write result to the buffer.
// DEPRECATED: use WriteApplyFnString() instead.
func (b *Chain) WriteApplyFnStr(s string, fn func(dst, p []byte) []byte) *Chain {
	*b = fn(*b, byteconv.S2B(s))
	return b
}

// ReplaceStr replace old to new substrings in buffer.
// DEPRECATED: use ReplaceString() instead.
func (b *Chain) ReplaceStr(old, new string, n int) *Chain {
	return b.ReplaceString(old, new, n)
}

// ReplaceStrAll replaces all occurrences of old substrings to new in buffer.
// DEPRECATED: use ReplaceStringAll() instead.
func (b *Chain) ReplaceStrAll(old, new string) *Chain {
	return b.ReplaceStringAll(old, new)
}
