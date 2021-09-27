package bytebuf

import "strconv"

type proto struct {
	b []byte
	s string
	i int64
	u uint64
	f float64
	l int
}

var (
	stage = &proto{
		b: []byte("bytes string"),
		s: "foobar string",
		i: 7439264324321,
		u: 9827390546234,
		f: 3.1415,
	}

	expectWS = []byte("bytes string-foobar string-7439264324321-9827390546234-3.1415")
	expectMT = []byte("bytes stringfoobar string743926432432198273905462343.1415")
)

func (p *proto) Size() int {
	if p.l == 0 {
		o := len(p.b)
		p.b = append(p.b, p.s...)
		p.b = strconv.AppendInt(p.b, p.i, 10)
		p.b = strconv.AppendUint(p.b, p.u, 10)
		p.b = strconv.AppendFloat(p.b, p.f, 'f', -1, 64)
		p.l = len(p.b)
		p.b = p.b[:o]
	}
	return p.l
}

func (p proto) MarshalTo(dst []byte) (int, error) {
	dst = append(dst[:0], p.b...)
	dst = append(dst, p.s...)
	dst = strconv.AppendInt(dst, p.i, 10)
	dst = strconv.AppendUint(dst, p.u, 10)
	dst = strconv.AppendFloat(dst, p.f, 'f', -1, 64)
	return len(dst), nil
}
