package sum

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"testing"
)

func TestFletcher16(t *testing.T) {
	data := []struct {
		Value string
		Want  string
	}{
		{Value: "0102", Want: "0403"},
		{Value: "6162636465", Want: "c8f0"},
		{Value: "616263646566", Want: "2057"},
		{Value: "6162636465666768", Want: "0627"},
	}
	for i, d := range data {
		v, _ := hex.DecodeString(d.Value)
		w, _ := hex.DecodeString(d.Want)

		r := Fletcher16(v)
		bs := make([]byte, binary.Size(r))
		binary.BigEndian.PutUint16(bs, r)

		if !bytes.Equal(bs, w) {
			t.Errorf("%d) want %x, got %x", i+1, w, bs)
		}
	}
}

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
