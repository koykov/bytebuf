# Byte buffers

Collection of various byte buffer implementations.

## Chain buffer

A wrapper around `[]byte` slice that supports chain write methods. Example of usage:
```go
import "github.com/koykov/bytebuf"

var buf bytebuf.Chain
b := buf.Write([]byte("foo")).
    WriteString("bar").
    WriteByte('@').
    WriteInt(-123).
    WriteUint(123).
    WriteFloat(3.1415).
    WriteBool(true).
    WriteRune('X').
    WriteFormat("pre%dpost", 15)
    WriteX(1.23456).
    Bytes()
println(string(b)) // foobar@-1231233.1415trueXpre15post1.23456
```

The same operations may be proceeded using `append()`/`AppendInt()`/... functions around `[]byte` slice, but `Chain`
provides handy API for that.

Chain buffer may be reused using internal pool. To get instance of the buffer from the pool call `AcquireChain` function.
To put buffer back to the pool call `ReleaseChain` function.

## Accumulative buffer

A wrapper around `Chain` buffer with "staked" features. The main idea is to accumulate (bufferize) various data and
provide handy access to chunks of buffered data. Example of usage:
```go
import "github.com/koykov/bytebuf"

var (
    buf bytebuf.Chain
    msg protobuf.ExampleMessage
)
chunk0 := buf.WriteMarshallerTo(&msg).Bytes()
chunk1 := buf.StakeOut().Write([]byte("foo")).
    WriteString("bar").
    WriteByte('@').
    WriteInt(-123).
    WriteUint(123).
    WriteFloat(3.1415).
    WriteBool(true).
    WriteRune('X').
    WriteFormat("pre%dpost", 15)
    WriteX(1.23456).
    StakedBytes()
println(string(chunk0)) // h�����,����?�ihttps
println(string(chunk1)) // foobar@-1231233.1415trueXpre15post1.23456
```
Thus, one buffer may be used to bufferize multiple data and hence reduce pointers in application.

Accumulative buffer may be reused using internal pool. To get instance of the buffer from the pool call `AcquireAccumulative`
function. To put buffer back to the pool call `ReleaseAccumulative` function.

# Pools

Both `Chain` and `Accumulative` buffers are using internal pools to reduce memory allocations.

Usage:
```go
cbuf := bytebuf.AcquireChain()
defer bytebuf.ReleaseChain(cbuf)
cbuf.WriteString("foo")
...

abuf := bytebuf.AcquireAccumulative()
defer bytebuf.ReleaseAccumulative(abuf)
abuf.WriteString("foo")
...
```
