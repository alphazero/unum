// friend!

// The MIT License (MIT)
//
// Copyright (c) 2016 Joubin Muhammad Houshyar
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
