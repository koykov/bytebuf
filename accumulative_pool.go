package bytebuf

import "sync"

var ap = sync.Pool{
	New: func() any {
		return &Accumulative{}
	},
}

// AcquireAccumulative gets Accumulative buffer instance from the pool.
func AcquireAccumulative() *Accumulative {
	return ap.Get().(*Accumulative)
}

// ReleaseAccumulative puts Accumulative buffer instance to the pool.
func ReleaseAccumulative(buf *Accumulative) {
	buf.Reset()
	ap.Put(buf)
}

var _, _ = AcquireAccumulative, ReleaseAccumulative
