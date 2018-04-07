package ar

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
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

type Writer struct {
	inner io.Writer
}

func NewWriter(w io.Writer) *Writer {
	w.Write(magic)
	w.Write([]byte{0x0A})
	return &Writer{w}
}

func (w *Writer) WriteHeader(h *Header) error {
	buf := new(bytes.Buffer)
	if err := writeHeaderField(buf, path.Base(h.Name)+"/", 16); err != nil {
		return err
	}
	t := h.ModTime.Unix()
	if err := writeHeaderField(buf, strconv.FormatInt(t, 10), 12); err != nil {
		return err
	}
	if err := writeHeaderField(buf, strconv.FormatInt(int64(h.Uid), 10), 6); err != nil {
		return err
	}
	if err := writeHeaderField(buf, strconv.FormatInt(int64(h.Gid), 10), 6); err != nil {
		return err
	}
	if err := writeHeaderField(buf, strconv.FormatInt(int64(h.Mode), 8), 8); err != nil {
		return err
	}
	if err := writeHeaderField(buf, strconv.FormatInt(int64(h.Length), 10), 10); err != nil {
		return err
	}
	buf.Write(feed)
	_, err := io.Copy(w.inner, buf)
	return err
}

func writeHeaderField(w io.Writer, s string, n int) error {
	v := strings.Repeat(" ", n-len(s))
	_, err := io.WriteString(w, s+v)
	return err
}

func (w *Writer) Write(bs []byte) (int, error) {
	vs := make([]byte, len(bs))
	copy(vs, bs)
	if len(bs)%2 == 1 {
		vs = append(vs, 0x0A)
	}
	n, err := w.inner.Write(vs)
	if err != nil {
		return n, err
	}
	return len(bs), err
}

func (w *Writer) Close() error {
	_, err := w.inner.Write([]byte{0x0A})
	return err
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
	size := r.hdr.Length
	if size%2 == 1 {
		size++
	}
	vs := make([]byte, size)
	n, err := io.ReadFull(r.inner, vs)
	if err != nil {
		r.err = err
	}
	r.hdr = nil
	if size%2 == 1 {
		n--
	}
	return copy(bs, vs[:n]), r.err
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
		return fmt.Errorf("time: %s", err)
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
			return fmt.Errorf("uid: %s", err)
		}
		h.Uid = int(i)
	}
	if bs, err := readHeaderField(r, 6); err != nil {
		return err
	} else {
		i, err := strconv.ParseInt(string(bs), 0, 64)
		if err != nil {
			return fmt.Errorf("gid: %s", err)
		}
		h.Gid = int(i)
	}
	if bs, err := readHeaderField(r, 8); err != nil {
		return err
	} else {
		i, err := strconv.ParseInt(string(bs), 0, 64)
		if err != nil {
			return fmt.Errorf("mode: %s", err)
		}
		h.Mode = int(i)
	}
	if bs, err := readHeaderField(r, 10); err != nil {
		return err
	} else {
		i, err := strconv.ParseInt(string(bs), 0, 64)
		if err != nil {
			return fmt.Errorf("length: %s", err)
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
