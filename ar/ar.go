package ar

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"
)

var (
	magic = []byte("!<arch>")
	feed  = []byte{0x60, 0x0A}
)

type Header struct {
	Name    string
	Uid     int
	Gid     int
	Mode    int
	Length  int
	ModTime time.Time
}

type Reader struct {
	inner io.Reader
	hdr   *Header
	err   error
}

func NewReader(r io.Reader) (*Reader, error) {
	rs := bufio.NewReader(r)
	bs, err := rs.Peek(len(magic))
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(bs, magic) {
		return nil, fmt.Errorf("not an ar archive")
	}
	if _, err := rs.Discard(len(bs) + 1); err != nil {
		return nil, err
	}
	return &Reader{inner: rs}, nil
}

func (r *Reader) Next() (*Header, error) {
	var h Header
	if r.err != nil {
		return nil, r.err
	}
	if err := readFilename(r.inner, &h); err != nil {
		r.err = err
		return nil, err
	}
	if err := readModTime(r.inner, &h); err != nil {
		r.err = err
		return nil, err
	}
	if err := readFileInfos(r.inner, &h); err != nil {
		r.err = err
		return nil, err
	}
	bs := make([]byte, len(feed))
	if _, err := r.inner.Read(bs); err != nil || !bytes.Equal(bs, feed) {
		return nil, err
	}
	r.hdr = &h
	return r.hdr, r.err
}

func (r *Reader) Read(bs []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.hdr == nil {
		_, err := r.Next()
		if err != nil {
			return 0, err
		}
	}
	return io.ReadAtLeast(r.inner, bs, r.hdr.Length)
}

func readFilename(r io.Reader, h *Header) error {
	bs, err := readHeaderField(r, 16)
	if err != nil {
		return err
	}
	h.Name = string(bs)
	return nil
}

func readModTime(r io.Reader, h *Header) error {
	bs, err := readHeaderField(r, 12)
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(string(bs), 0, 64)
	if err != nil {
		return err
	}
	h.ModTime = time.Unix(i, 0)
	return nil
}

func readFileInfos(r io.Reader, h *Header) error {
	if bs, err := readHeaderField(r, 6); err != nil {
		return err
	} else {
		i, err := strconv.ParseInt(string(bs), 0, 64)
		if err != nil {
			return err
		}
		h.Uid = int(i)
	}
	if bs, err := readHeaderField(r, 6); err != nil {
		return err
	} else {
		i, err := strconv.ParseInt(string(bs), 0, 64)
		if err != nil {
			return err
		}
		h.Gid = int(i)
	}
	if bs, err := readHeaderField(r, 8); err != nil {
		return err
	} else {
		i, err := strconv.ParseInt(string(bs), 0, 64)
		if err != nil {
			return err
		}
		h.Mode = int(i)
	}
	if bs, err := readHeaderField(r, 10); err != nil {
		return err
	} else {
		i, err := strconv.ParseInt(string(bs), 0, 64)
		if err != nil {
			return err
		}
		h.Length = int(i)
	}
	return nil
}

func readHeaderField(r io.Reader, n int) ([]byte, error) {
	bs := make([]byte, n)
	if _, err := io.ReadFull(r, bs); err != nil {
		return nil, err
	}
	return bytes.TrimSpace(bs), nil
}
