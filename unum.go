// friend

// package implements the 'unum' tagged value binary encoding of unsigned integers.
// Schme per John Gustafson, "Right-Sizing Precision"(presentation) , March 2013.
//
// Gustafson's `unum` numeric value representation scheme is a tagged value encoding
// scheme supporting variable physical storage of numeric values. The idea is pretty
// straight forward: typically values used hardly reach the value type's maximum
// value and space is wasted. By tagging the physical image of the value with an
// indication of the size category (see below) we can use less bytes than the full
// type (e.g. uint64 at 8 bytes) would require.
//
// Given that unlike Gustafson's we're not discussing hardware here, and we use
// byte buffers, this scheme is a good match for shuffling sequence of numbers in I/O.
//
// This implementation (as of now) only supports unsigned integer values.
//
// Note that the scheme is byte-aligned, not word-aligned.
//
//      -- bit layout for unsigned values
//      0          2                            *
//      [ tag-bits | variable length value rep. ]
//
//
// The 2-bits of the tag determine the value-range and physical length of the image:
//
//
//      tag        | bytes | range
//      -----------+-------+------------------------------------------------
//      00         | 1     | uint: (0, 2^6]
//      -----------+-------+------------------------------------------------
//      01         | 2     | uint: (2^6, 2^14]
//      -----------+-------+------------------------------------------------
//      10         | 4     | uint: (2^14, 2^30]
//      -----------+-------+------------------------------------------------
//      11         | 8     | uint: (2^30, 2^62]
//
package unum

import (
	"fmt"
	"io"
)

// unsigned int value masks
const (
	uint6mask  = uint8(0x3f)
	uint14mask = uint16(0x3fff)
	uint30mask = uint32(0x3fffffff)
	uint62mask = uint64(0x3fffffffffffffff)
)

// Encode/Write errors
var (
	ErrorBufferOverflow = fmt.Errorf("buffer overflow")
	ErrorMaxValue       = fmt.Errorf("value exceeds max value")
)

// unum encodes the value v in buffer buf.
// returns number of bytes written { 1, 2, 4, 8 } on nil
// error.
//
// On error returns (0, e) where e is:
//    ErrorBufferOverflow  -- invalid arg b : len(b) < n
//    ErrorMaxValue        -- invalid arg v : v > 2^62
func EncodeUint(b []byte, v uint64) (n int, e error) {
	switch {
	case v < 0x40:
		if len(b) < 1 {
			return 0, ErrorBufferOverflow
		}
		b[0] = uint8(v)
		return 1, nil
	case v < 0x4000:
		if len(b) < 2 {
			return 0, ErrorBufferOverflow
		}
		b[0] = uint8((uint14mask&uint16(v))>>8) | 0x40
		b[1] = uint8(v)
		return 2, nil
	case v < 0x40000000:
		if len(b) < 4 {
			return 0, ErrorBufferOverflow
		}
		b[0] = uint8((uint30mask&uint32(v))>>24) | 0x80
		b[1] = uint8((v & 0xff0000) >> 16)
		b[2] = uint8((v & 0xff00) >> 8)
		b[3] = uint8(v & 0xff)
		return 4, nil
	case v < 0x4000000000000000:
		if len(b) < 8 {
			return 0, ErrorBufferOverflow
		}
		b[0] = uint8((uint62mask&v)>>56) | 0xc0
		b[1] = uint8((v & 0xff000000000000) >> 48)
		b[2] = uint8((v & 0xff0000000000) >> 40)
		b[3] = uint8((v & 0xff00000000) >> 32)
		b[4] = uint8((v & 0xff000000) >> 24)
		b[5] = uint8((v & 0xff0000) >> 16)
		b[6] = uint8((v & 0xff00) >> 8)
		b[7] = uint8(v & 0xff)
		return 8, nil
	default:
		return 0, ErrorMaxValue
	}
	panic("bug - asserted unreachable")
}

// Writes encoded uint value 'v' to writer 'w'.
// Returns number of bytes written 'n' (n > 0) if there are no errors.
//
// On error returns (0, e) where e is:
//    ErrorBufferOverflow  -- invalid arg b : len(b) < n
//    ErrorMaxValue        -- invalid arg v : v > 2^62
//    <other>              -- propagated io.Writer.Write error
func WriteUint(w io.Writer, v uint64) (n int, e error) {
	var b [8]byte
	n0, e0 := EncodeUint(b[0:], v)
	if e0 != nil {
		return 0, e0
	}
	return w.Write(b[:n0])
}

// Decode/Read errors
var (
	ErrorBufferUnderflow = fmt.Errorf("buffer underflow")
	ErrorInvalidBuffer   = fmt.Errorf("invalid buffer")
)

// decodes unum encoded unsigned integer value v from input buffer b.
// returns cLUW v, number of bytes n, if there are no errors.
//
// On error returns (0, 0, e) where e is:
//    ErrorBufferUnderflow -- invalid arg b : len(b) < n
//    ErrorInvalidBuffer   -- invalid arg b : non-conformant byte sequence
func DecodeUint(b []byte) (v uint64, n int, e error) {
	if len(b) == 0 {
		return 0, 0, ErrorBufferUnderflow
	}
	switch b[0] & 0xc0 {
	case 0:
		return uint64(b[0] & 0x3f), 1, nil
	case 0x40:
		if len(b) < 2 {
			return 0, 0, ErrorBufferUnderflow
		}
		v = uint64(b[0]&0x3f)<<8 |
			uint64(b[1])
		return v, 2, nil
	case 0x80:
		if len(b) < 4 {
			return 0, 0, ErrorBufferUnderflow
		}
		v = uint64(b[0]&0x3f)<<24 |
			uint64(b[1])<<16 |
			uint64(b[2])<<8 |
			uint64(b[3])
		return v, 4, nil
	case 0xc0:
		if len(b) < 8 {
			return 0, 0, ErrorBufferUnderflow
		}
		v = uint64(b[0]&0x3f)<<56 |
			uint64(b[1])<<48 |
			uint64(b[2])<<40 |
			uint64(b[3])<<32 |
			uint64(b[4])<<24 |
			uint64(b[5])<<16 |
			uint64(b[6])<<8 |
			uint64(b[7])
		return v, 8, nil
	}
	panic("bug - asserted unreachable")
}

// Reads unum encoded usigned integer value from Reader r.
// Returns value v, number of bytes read n (n > 0), if there are no errors.
//
// On error returns (0, 0, e) where e is:
//    ErrorBufferUnderflow -- invalid arg b : len(b) < n
//    ErrorInvalidBuffer   -- invalid arg b : non-conformant byte sequence
//    <other>              -- propagated io.Reader.Read error
func ReadUint(r io.Reader) (v uint64, n int, e error) {
	var b [8]byte
	n, e = r.Read(b[:1])
	if e != nil {
		return
	}

	vlen := (b[0] & 0xc0) >> 6
	if vlen > 1 {
		n, e = io.ReadFull(r, b[1:vlen])
		if e != nil {
			n += 1
			return
		}
	}
	return DecodeUint(b[0:])
}
