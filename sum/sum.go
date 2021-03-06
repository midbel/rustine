package sum

import (
	"bytes"
	"encoding/binary"
)

const (
	mod16 = 255
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

func Sum1071Bis(bs []byte) uint16 {
	var (
		answer uint16
		sum   uint64
	)
	r := bytes.NewReader(bs)
	for r.Len() > 0 {
		var word uint16
		binary.Read(r, binary.LittleEndian, &word)
		sum += uint64(word)
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	answer = uint16(sum)
	return ^answer
}

func Fletcher16(bs []byte) uint16 {
	var s1, s2 uint16
	for i := 0; i < len(bs); i++ {
		s1 = (s1 + uint16(bs[i])) % mod16
		s2 = (s2 + s1) % mod16
	}
	return (s2 << 8) | uint16(s1)
}

func Fletcher32(bs []byte) uint32 {
	return 0
}

func Fletcher64(bs []byte) uint64 {
	return 0
}
