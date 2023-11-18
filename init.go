package bytebuf

import "github.com/koykov/x2bytes"

func init() {
	x2bytes.RegisterToBytesFn(ChainBufToBytes)
	x2bytes.RegisterToBytesFn(AccumulativeBufToBytes)
}

// ChainBufToBytes registers x2bytes conversion function accepts Chain.
func ChainBufToBytes(dst []byte, val interface{}) ([]byte, error) {
	if b, ok := val.(*Chain); ok {
		dst = append(dst, *b...)
		return dst, nil
	}
	return dst, x2bytes.ErrUnknownType
}

// AccumulativeBufToBytes registers x2bytes conversion function accepts Accumulative.
func AccumulativeBufToBytes(dst []byte, val interface{}) ([]byte, error) {
	if b, ok := val.(*Accumulative); ok {
		dst = append(dst, b.Bytes()...)
		return dst, nil
	}
	return dst, x2bytes.ErrUnknownType
}
