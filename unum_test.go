// friend!

package unum_test

import (
	"math/rand"
	"testing"
	"testing/quick"
	"unum"
)

func BenchmarkRandGenForReference(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// encode
		v := uint64(rand.Int63n(int64(unum.Unum64MaxValue)))
		if v < 0 {
			b.Errorf("dummy test - just insure loop is not optimized away!\n")
		}
	}
}

func BenchmarkEncodeUnum64(b *testing.B) {
	//	v := uint64(rand.Int63n(int64(unum.Unum64MaxValue)))
	var vb0 [8]byte
	vb := vb0[:]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// encode
		v := uint64(rand.Int63n(int64(unum.Unum64MaxValue)))
		_, e := unum.EncodeUint(vb, v)
		if e != nil {
			b.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
		}
	}
}

func BenchmarkDecodeUnum64(b *testing.B) {
	v := uint64(rand.Int63n(int64(unum.Unum64MaxValue)))
	var vb0 [8]byte
	vb := vb0[:]
	_, e := unum.EncodeUint(vb, v)
	if e != nil {
		b.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// encode
		_, _, e := unum.DecodeUint(vb)
		if e != nil {
			b.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
		}
	}
}

func TestCodec(t *testing.T) {
	f := func(v uint64) bool {
		if v >= unum.Unum64MaxValue {
			return true // ignore max values
		}
		var b0 [8]byte
		b := b0[:]

		// encode
		_, e := unum.EncodeUint(b, v)
		if e != nil {
			t.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
		}

		// decode
		v0, _, e := unum.DecodeUint(b)
		if e != nil {
			t.Errorf("error decoding - v:%d - e:%s\n", v, e.Error())
		}
		if v0 != v {
			t.Errorf("BUG - v:%d - v0%d\n", v, v0)
		}
		return true
	}
	if e := quick.Check(f, nil); e != nil {
		t.Error(e)
	}
}
