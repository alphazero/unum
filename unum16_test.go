// friend!

package unum_test

import (
	"math/rand"
	"testing"
	"testing/quick"
	"unum"
)

func BenchmarkEncodeUnum16(b *testing.B) {
	v := uint16(rand.Intn(int(unum.Unum16ValueBound)))
	var vb0 [unum.Unum16Size]byte
	vb := vb0[:]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// encode
		//		v := uint16(rand.Intn(int(unum.Unum16ValueBound)))
		_, e := unum.EncodeUnum16(vb, v)
		if e != nil {
			b.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
		}
	}
}

func BenchmarkDecodeUnum16(b *testing.B) {
	v := uint16(rand.Intn(int(unum.Unum16ValueBound)))
	var vb0 [unum.Unum16Size]byte
	vb := vb0[:]
	_, e := unum.EncodeUnum16(vb, v)
	if e != nil {
		b.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// encode
		_, _, e := unum.DecodeUnum16(vb)
		if e != nil {
			b.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
		}
	}
}

func TestCodecUnum16(t *testing.T) {
	f := func(v uint16) bool {
		if v >= unum.Unum16ValueBound {
			return true // ignore max values
		}
		var b0 [unum.Unum16Size]byte
		b := b0[:]

		// encode
		_, e := unum.EncodeUnum16(b, v)
		if e != nil {
			t.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
		}

		// decode
		v0, _, e := unum.DecodeUnum16(b)
		if e != nil {
			t.Errorf("error decoding - v:%d - e:%s\n", v, e.Error())
		}
		if v0 != v {
			t.Errorf("BUG - v:%04x - v0:%04x\n", v, v0)
		}
		return true
	}
	if e := quick.Check(f, nil); e != nil {
		t.Error(e)
	}
}
