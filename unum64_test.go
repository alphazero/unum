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
	//	"bytes"
	"math/rand"
	"testing"
	"testing/quick"
	"unum"
)

func BenchmarkEncodeUnum64(b *testing.B) {
	v := uint64(rand.Int63n(int64(unum.Unum64ValueBound)))
	var vb0 [unum.Unum64Size]byte
	vb := vb0[:]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, e := unum.EncodeUnum64(vb, v)
		if e != nil {
			b.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
		}
	}
}

func BenchmarkDecodeUnum64(b *testing.B) {
	v := uint64(rand.Int63n(int64(unum.Unum64ValueBound)))
	var vb0 [unum.Unum64Size]byte
	vb := vb0[:]
	_, e := unum.EncodeUnum64(vb, v)
	if e != nil {
		b.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, e := unum.DecodeUnum64(vb)
		if e != nil {
			b.Errorf("error encoding - v:%d - e:%s\n", v, e.Error())
		}
	}
}

func TestCodecUnum64(t *testing.T) {
	f := func(v uint64) bool {
		var errorExpected bool
		if v >= unum.Unum64ValueBound {
			errorExpected = true
		}
		var b0 [unum.Unum64Size]byte
		b := b0[:]

		// encode
		n, e := unum.EncodeUnum64(b, v)
		if e != nil {
			if !errorExpected {
				t.Errorf("unexpected error encoding - v:%d - e:%s\n", v, e.Error())
			} else if e != unum.ErrorMaxValue {
				t.Errorf("expected ErrorMaxValue econding - v:%d - have:%s\n", v, e.Error())
			}
			return true
		} else if errorExpected {
			t.Errorf("expected error encoding - v:%d\n", v, e.Error())
		}

		// check encoding size
		size := unum.Unum64Size
		if n > size {
			t.Errorf("encoding size exceeds %d for v:%d - have:%d\n", size, v, n)
		}

		// decode
		v0, _, e := unum.DecodeUnum64(b)
		if e != nil {
			t.Errorf("error decoding - v:%d - e:%s\n", v, e.Error())
		}
		if v0 != v {
			t.Errorf("BUG - v:%08x - v0:%08x\n", v, v0)
		}
		return true
	}
	if e := quick.Check(f, nil); e != nil {
		t.Error(e)
	}
}
