package rustine

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"time"
)

const (
	lowers  = "abcdefghijklmnopqrstuvwxyz"
	uppers  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers = "0123456789"
)

func RandomString(c int) string {
	rand.Seed(time.Now().Unix())

	b := new(bytes.Buffer)

	cs := []byte(lowers + uppers + numbers)
	for z := len(cs); b.Len() <= c; {
		ix := rand.Intn(z)
		b.WriteByte(cs[ix])
	}
	return b.String()
}

func Itob(v int) []byte {
	bs := make([]byte, 16)
	i := len(bs) - 1
	for ; v > 0; i-- {
		bs[i], v = byte(v&0xFF), v>>8
	}
	return bs[i+1:]
}

func Sum(b []byte) uint16 {
	var s, t uint16

	buf := bytes.NewReader(b)
	for buf.Len() > 0 {
		if err := binary.Read(buf, binary.BigEndian, &t); err != nil {
			break
		}
		s += t
	}
	return s
}
