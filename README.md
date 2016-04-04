![image](./resources/unum-logo.png)

####stat
    star date         04 04 16

    specification     John Gustafson | "Right-Sizing Precision" (presentation) 
    api               stable 
    documentation     inlined godoc

##about

John Gustafson's `unum` numeric value representation scheme is a tagged value encoding scheme supporting variable physical storage of values. 

This implementation (as of now) only supports unsigned integer values.

### spec

The scheme is **byte-aligned**, not word-aligned.

     -- bit layout for unsigned integer values
     
     0          2                            *    
     [ tag-bits | variable length value rep. ]


The 2-bits of the tag determine the value-range and physical length of the image:


     tag        | bytes | range
     -----------+-------+------------------------------------------------
     00         | 1     | uint: (0, 2^6] 
     -----------+-------+------------------------------------------------
     01         | 2     | uint: (2^6, 2^14] 
     -----------+-------+------------------------------------------------
     10         | 4     | uint: (2^14, 2^30] 
     -----------+-------+------------------------------------------------
     11         | 8     | uint: (2^30, 2^62] 

### usage

####`encode`

**using byte array**

	var value []uint64 = { .. }         
    var b []byte = ..           // provided by you
    
    // encoding to a byte buffer
    var offset int
    for _, v := range values {
        n, e := unum.EncodeUint(b[offset:], v)
        if e != nil {
            /* if e is ErrorBufferOverflo you could resize the buffer here */
            break
        }
        offset += n
    }

**using io.Writer **
 
    var w Writer = ..           // provider by you
    // encoding to a Writer
	for _, v := range values {
		_, e := unum.WriteUint(w, v)
		if e != nil {
			log.Fatalf("err - %s - value:%0x", e.Error(), v)
		}
	}
	 
####`decode`

**using byte array**
 
    var b []byte = ..           // unum encoded provided by you
    
    // decoding from a byte buffer
    var offset int
    for {
        v, n, e := unum.DecodeUint(b[offset:])
        if e != nil {
            break 
        }
        offset += n
    }
 