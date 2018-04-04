package rw

import (
	"crypto/rand"
	"io"
)

type empty struct{}

func Rand() io.Reader {
	return &empty{}
}

func (e *empty) Read(bs []byte) (int, error) {
	n, err := io.ReadFull(rand.Reader, bs)
	return int(n), err
}

type bytes struct {
	buffer []byte
}

func Alea(n int) io.Reader {
	bs := make([]byte, n)
	io.ReadFull(rand.Reader, bs)

	return &bytes{bs}
}

func Zero(n int) io.Reader {
	return &bytes{make([]byte, n)}
}

func (b *bytes) Read(bs []byte) (int, error) {
	return copy(bs, b.buffer), nil
}
