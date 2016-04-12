// friend

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

package unum

import (
	"io"
)

// UNUM-16 encodes the value v in buffer buf.
// returns number of bytes written { 1, 2 } on nil
// error.
//
// On error returns (0, e) where e is:
//    ErrorBufferOverflow  -- invalid arg b : len(b) < n
//    ErrorMaxValue        -- invalid arg v : v > 2^15
func EncodeUnum16(b []byte, v uint16) (n int, e error) {
	switch {
	case v < 0x80:
		if len(b) < 1 {
			return 0, ErrorBufferOverflow
		}
		b[0] = byte(v)
		return 1, nil
	case v < 0x8000:
		if len(b) < 2 {
			return 0, ErrorBufferOverflow
		}
		b[0] = byte(v>>8) | 0x80
		b[1] = byte(v)
		return 2, nil
	default:
		return 0, ErrorMaxValue
	}
	panic("bug - asserted unreachable")
}

// Writes UNUM-16 encoded uint value 'v' to writer 'w'.
// Returns number of bytes written 'n' (n > 0) if there are no errors.
//
// On error returns (0, e) where e is:
//    ErrorBufferOverflow  -- invalid arg b : len(b) < n
//    ErrorMaxValue        -- invalid arg v : v > 2^15
//    <other>              -- propagated io.Writer.Write error
func WriteUnum16(w io.Writer, v uint16) (n int, e error) {
	var b [Unum16Size]byte
	n0, e0 := EncodeUnum16(b[0:], v)
	if e0 != nil {
		return 0, e0
	}
	return w.Write(b[:n0])
}

// decodes UNUM-16 encoded unsigned integer value v from input buffer b.
// returns cLUW v, number of bytes n, if there are no errors.
//
// On error returns (0, 0, e) where e is:
//    ErrorBufferEOF -- invalid arg b : len(b) < n
//    ErrorInvalidBuffer   -- invalid arg b : non-conformant byte sequence
func DecodeUnum16(b []byte) (v uint16, n int, e error) {
	if len(b) == 0 {
		return 0, 0, ErrorBufferEOF
	}

	switch b[0] & 0x80 {
	case 0:
		return uint16(b[0] & 0x7f), 1, nil
	case 0x80:
		if len(b) < 2 {
			return 0, 0, ErrorInvalidBuffer
		}
		v = uint16(b[0]&0x7f)<<8 |
			uint16(b[1])
		return v, 2, nil
	}
	panic("bug - asserted unreachable")
}

// Reads UNUM-16 encoded usigned integer value from Reader r.
// Returns value v, number of bytes read n (n > 0), if there are no errors.
//
// On error returns (0, 0, e) where e is:
//    ErrorBufferEOF -- invalid arg b : len(b) < n
//    ErrorInvalidBuffer   -- invalid arg b : non-conformant byte sequence
//    <other>              -- propagated io.Reader.Read error
//
// Note that ErrorInvalidBuffer
func ReadUnum16(r io.Reader) (v uint16, n int, e error) {
	// REVU: error handling can be cleaner (i.e. io.EOF -> Underflow)
	//
	var b [Unum16Size]byte
	n, e = r.Read(b[:1])
	if e != nil {
		if e == io.EOF {
			e = ErrorBufferEOF
		}
		return
	}

	// compute expected encoded len directly
	// and read the remaining bytes (if any)
	vlen := 1 << (b[0] >> 7)
	if vlen > 1 {
		n, e = io.ReadFull(r, b[1:vlen])
		if e != nil {
			n += 1
			if e == io.EOF {
				e = ErrorInvalidBuffer
			}
			return
		}
	}
	return DecodeUnum16(b[0:])
}
