package sum

import (
  "bytes"
  "encoding/binary"
	"encoding/hex"
	"testing"
)

func TestSum1071(t *testing.T) {
	data := []struct {
		Value string
		Want  string
	}{
    {Value: "0001f203f4f5f6f7", Want: "ddf2"},
    {Value: "010003f2f5f4f7f6", Want: "f2dd"},
    {Value: "0001f200", Want: "f201"},
  }
	for i, d := range data {
		v, _ := hex.DecodeString(d.Value)
		w, _ := hex.DecodeString(d.Want)

    r := Sum1071(v)
    bs := make([]byte, binary.Size(r))
    binary.BigEndian.PutUint16(bs, r)
    if !bytes.Equal(bs, w) {
      t.Errorf("%d) want %x, got %x", i+1, w, bs)
    }
	}
}
