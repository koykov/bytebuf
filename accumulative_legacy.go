package bytebuf

// AccumulativeBuf is an alias of Accumulative type.
// DEPRECATED: use Accumulative instead. Kept for backward compatibility.
type AccumulativeBuf = Accumulative

// AccBufWriter is an alias of AccumulativeWriter type.
// DEPRECATED: use AccumulativeWriter instead. Kept for backward compatibility.
type AccBufWriter = AccumulativeWriter

// WriteStr writes string to the buffer.
// DEPRECATED: use WriteString() instead.
func (b *Accumulative) WriteStr(s string) *Accumulative {
	return b.WriteString(s)
}

// WriteApplyFnStr applies fn to s and write result to the buffer.
// DEPRECATED: use WriteApplyFnString() instead.
func (b *Accumulative) WriteApplyFnStr(s string, fn func(dst, p []byte) []byte) *Accumulative {
	return b.WriteApplyFnString(s, fn)
}

// ReplaceStr replace old to new substrings in buffer.
// DEPRECATED: use ReplaceString() instead.
func (b *Accumulative) ReplaceStr(old, new string, n int) *Accumulative {
	return b.ReplaceString(old, new, n)
}

// ReplaceStrAll replaces all occurrences of old substrings to new in buffer.
// DEPRECATED: use ReplaceStringAll() instead.
func (b *Accumulative) ReplaceStrAll(old, new string) *Accumulative {
	return b.ReplaceStringAll(old, new)
}
