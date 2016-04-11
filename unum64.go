// friend

package unum

import (
	"io"
)

// UNUM-64 encodes the value v in buffer buf.
// returns number of bytes written { 1, 2, 4, 8 } on nil
// error.
//
// On error returns (0, e) where e is:
//    ErrorBufferOverflow  -- invalid arg b : len(b) < n
//    ErrorMaxValue        -- invalid arg v : v > 2^62
func EncodeUnum64(b []byte, v uint64) (n int, e error) {
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
func WriteUnum64(w io.Writer, v uint64) (n int, e error) {
	var b [8]byte
	n0, e0 := EncodeUnum64(b[0:], v)
	if e0 != nil {
		return 0, e0
	}
	return w.Write(b[:n0])
}

// decodes UNUM-64 encoded unsigned integer value v from input buffer b.
// returns cLUW v, number of bytes n, if there are no errors.
//
// On error returns (0, 0, e) where e is:
//    ErrorBufferEOF -- invalid arg b : len(b) < n
//    ErrorInvalidBuffer   -- invalid arg b : non-conformant byte sequence
func DecodeUnum64(b []byte) (v uint64, n int, e error) {
	if len(b) == 0 {
		return 0, 0, ErrorBufferEOF
	}

	switch b[0] & 0xc0 {
	case 0:
		return uint64(b[0] & 0x3f), 1, nil
	case 0x40:
		if len(b) < 2 {
			return 0, 0, ErrorInvalidBuffer
		}
		v = uint64(b[0]&0x3f)<<8 |
			uint64(b[1])
		return v, 2, nil
	case 0x80:
		if len(b) < 4 {
			return 0, 0, ErrorInvalidBuffer
		}
		v = uint64(b[0]&0x3f)<<24 |
			uint64(b[1])<<16 |
			uint64(b[2])<<8 |
			uint64(b[3])
		return v, 4, nil
	case 0xc0:
		if len(b) < 8 {
			return 0, 0, ErrorInvalidBuffer
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

// Reads UNUM-64 encoded usigned integer value from Reader r.
// Returns value v, number of bytes read n (n > 0), if there are no errors.
//
// On error returns (0, 0, e) where e is:
//    ErrorBufferEOF -- invalid arg b : len(b) < n
//    ErrorInvalidBuffer   -- invalid arg b : non-conformant byte sequence
//    <other>              -- propagated io.Reader.Read error
//
// Note that ErrorInvalidBuffer
func ReadUint(r io.Reader) (v uint64, n int, e error) {
	// REVU: error handling can be cleaner (i.e. io.EOF -> Underflow)
	//
	var b [8]byte
	n, e = r.Read(b[:1])
	if e != nil {
		if e == io.EOF {
			e = ErrorBufferEOF
		}
		return
	}

	// compute expected encoded len directly
	// and read the remaining bytes (if any)
	vlen := 1 << (b[0] >> 6)
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
	return DecodeUnum64(b[0:])
}
