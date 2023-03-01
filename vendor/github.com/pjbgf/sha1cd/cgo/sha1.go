package cgo

// #include <sha1.h>
// #include <ubc_check.h>
import "C"

import (
	"crypto"
	"hash"
	"unsafe"
)

const (
	Size      = 20
	BlockSize = 64
)

func init() {
	crypto.RegisterHash(crypto.SHA1, New)
}

func New() hash.Hash {
	d := new(digest)
	d.Reset()
	return d
}

type digest struct {
	ctx C.SHA1_CTX
	h   [Size]byte
}

func (d *digest) sum() ([Size]byte, bool) {
	c := C.SHA1DCFinal((*C.uchar)(unsafe.Pointer(&d.h[0])), &d.ctx)
	if c != 0 {
		return d.h, true
	}

	return d.h, false
}

func (d *digest) Sum(in []byte) []byte {
	d0 := *d // use a copy of d to avoid race conditions.
	h, _ := d0.CollisionResistantSum(in)
	return h
}

func (d *digest) CollisionResistantSum(in []byte) ([]byte, bool) {
	d0 := *d // use a copy of d to avoid race conditions.
	h, c := d0.sum()
	return append(in, h[:]...), c
}

func (d *digest) Reset() {
	C.SHA1DCInit(&d.ctx)
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return BlockSize }

func Sum(data []byte) ([]byte, bool) {
	d := New().(*digest)
	d.Write(data)

	h, c := d.sum()
	return h[:], c
}

func (d *digest) Write(p []byte) (nn int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	data := (*C.char)(unsafe.Pointer(&p[0]))
	C.SHA1DCUpdate(&d.ctx, data, (C.size_t)(len(p)))

	return len(p), nil
}
