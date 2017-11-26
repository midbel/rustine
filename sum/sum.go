package sum

import (
  "bytes"
  "encoding/binary"
)

func Sum1071(bs []byte) uint16 {
	r := bytes.NewBuffer(bs)
	if len(bs)%2 != 0 {
		r.WriteByte(0)
	}

	var s uint32
	for {
		var v uint16
		if err := binary.Read(r, binary.BigEndian, &v); err != nil {
			break
		}
		s += uint32(v)
	}
	for i := s >> 16; i > 0; i = s >> 16 {
		s = (s & 0xffff) | (s >> 16)
	}
	return uint16(s)
}
