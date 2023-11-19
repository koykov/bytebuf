package bytebuf

import "github.com/koykov/x2bytes"

func init() {
	x2bytes.RegisterToBytesFn(ChainToBytes)
	x2bytes.RegisterToBytesFn(AccumulativeToBytes)
}

// ChainToBytes registers x2bytes conversion function accepts Chain.
func ChainToBytes(dst []byte, val interface{}) ([]byte, error) {
	if b, ok := val.(*Chain); ok {
		dst = append(dst, *b...)
		return dst, nil
	}
	return dst, x2bytes.ErrUnknownType
}

// AccumulativeToBytes registers x2bytes conversion function accepts Accumulative.
func AccumulativeToBytes(dst []byte, val interface{}) ([]byte, error) {
	if b, ok := val.(*Accumulative); ok {
		dst = append(dst, b.Bytes()...)
		return dst, nil
	}
	return dst, x2bytes.ErrUnknownType
}
