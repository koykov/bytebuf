package bytebuf

import "sync"

var cp = sync.Pool{
	New: func() any {
		return &Chain{}
	},
}

// AcquireChain gets Chain buffer instance from the pool.
func AcquireChain() *Chain {
	return cp.Get().(*Chain)
}

// ReleaseChain puts Chain buffer instance to the pool.
func ReleaseChain(buf *Chain) {
	buf.Reset()
	cp.Put(buf)
}

var _, _ = AcquireChain, ReleaseChain
