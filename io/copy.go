package io

import (
	"bytes"
	"io"
	"time"
)

func CopyWithFirstReadChan(firstReadDone chan struct{}) func(io.Writer, io.Reader) (int64, error) {
	return func(dst io.Writer, src io.Reader) (written int64, err error) {
		buf := make([]byte, 2*1024)
		for {
			nr, er := src.Read(buf)
			select {
			case <-firstReadDone:
			default:
				close(firstReadDone)
				// HACK: let the close hook of the spinner
				time.Sleep(100 * time.Millisecond)
			}
			if nr > 0 {
				toWrite := bytes.Replace(buf[0:nr], []byte{'\n'}, []byte{'\r', '\n'}, -1)
				nr = len(toWrite)
				nw, ew := dst.Write(toWrite)
				if nw > 0 {
					written += int64(nw)
				}
				if ew != nil {
					err = ew
					break
				}
				if nr != nw {
					err = io.ErrShortWrite
					break
				}
			}
			if er == io.EOF {
				break
			}
			if er != nil {
				err = er
				break
			}
		}
		return written, err
	}
}
