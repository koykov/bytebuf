package bytebuf

import "github.com/koykov/x2bytes"

func init() {
	x2bytes.RegisterToBytesFn(ChainBufToBytes)
	x2bytes.RegisterToBytesFn(AccumulativeBufToBytes)
}

// Conversion to bytes function.
func ChainBufToBytes(dst []byte, val interface{}) ([]byte, error) {
	if b, ok := val.(*ChainBuf); ok {
		dst = append(dst, *b...)
		return dst, nil
	}
	return dst, x2bytes.ErrUnknownType
}

// Conversion to bytes function.
func AccumulativeBufToBytes(dst []byte, val interface{}) ([]byte, error) {
	if b, ok := val.(*AccumulativeBuf); ok {
		dst = append(dst, b.Bytes()...)
		return dst, nil
	}
	return dst, x2bytes.ErrUnknownType
}
