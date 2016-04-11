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
// Note that the scheme is uniformly byte-aligned, not word-aligned.
//
//      -- bit layout for UNUM-64 and UNUM-32 unsigned values
//
//      0          2                            *
//      [ tag-bits | variable length value rep. ]
//
//
// The 2-bits of the tag determine the value-range and physical length of the image:
//
//      encoding: UNUM-32
//      tag        | bytes | range
//      -----------+-------+------------------------------------------------
//      00         | 1     | uint: (0, 2^6]
//      -----------+-------+------------------------------------------------
//      01         | 2     | uint: (2^6, 2^14]
//      -----------+-------+------------------------------------------------
//      10         | 3     | uint: (2^14, 2^22]
//      -----------+-------+------------------------------------------------
//      11         | 4     | uint: (2^22, 2^30]
//
//      encoding: UNUM-64
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
//
//      -- bit layout for UNUM-16 unsigned values
//
//      0          1                            *
//      [ tag-bit  | variable length value rep. ]
//
//      encoding: UNUM-16
//      tag        | bytes | range
//      -----------+-------+------------------------------------------------
//      0          | 1     | uint: (0, 2^7]
//      -----------+-------+------------------------------------------------
//      1          | 2     | uint: (2^7, 2^15]
//      -----------+-------+------------------------------------------------
//
package unum

import (
	"fmt"
	//	"goerror" // REVU: use goerror
)

// REVU: these are not max values, they are upper bounds. (0, bound]
// unsigned int max values
const (
	Unum64ValueBound = uint64(0x4000000000000000)
	Unum32ValueBound = uint32(0x40000000)
	Unum16ValueBound = uint32(0x8000)
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
	ErrorBufferOverflow = fmt.Errorf("unum.ErrorBufferOverflow")
	ErrorMaxValue       = fmt.Errorf("unum.ErrorMaxValue")
)

// Decode/Read errors
var (
	ErrorBufferEOF     = fmt.Errorf("unum.ErrorBufferEOF")
	ErrorInvalidBuffer = fmt.Errorf("unum.ErrorInvalidBuffer")
)
