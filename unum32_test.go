// friend!

package unum_test

import (
	"math/rand"
	"testing"
	"testing/quick"
	"unum"
)

func BenchmarkEncodeUnum32(b *testing.B) {
	v := uint32(rand.Int31n(int32(unum.Unum32ValueBound)))
	var vb0 [unum.Unum32Size]byte
	vb := vb0[:]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, e := unum.EncodeUnum32(vb, v)
		if e != nil {
			b.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
		}
	}
}

func BenchmarkDecodeUnum32(b *testing.B) {
	v := uint32(rand.Int31n(int32(unum.Unum32ValueBound)))
	var vb0 [unum.Unum32Size]byte
	vb := vb0[:]
	_, e := unum.EncodeUnum32(vb, v)
	if e != nil {
		b.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, e := unum.DecodeUnum32(vb)
		if e != nil {
			b.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
		}
	}
}

func TestCodecUnum32(t *testing.T) {
	f := func(v uint32) bool {
		if v >= unum.Unum32ValueBound {
			return true // ignore max values
		}
		var b0 [unum.Unum32Size]byte
		b := b0[:]

		// encode
		_, e := unum.EncodeUnum32(b, v)
		if e != nil {
			t.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
		}

		// decode
		v0, _, e := unum.DecodeUnum32(b)
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
