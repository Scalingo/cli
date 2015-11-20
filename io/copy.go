package io

import (
	"io"
	"net"
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
				nw, ew := dst.Write(buf[0:nr])
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

func CopyWithTimeout(timeout time.Duration) func(io.Writer, io.Reader) (int64, error) {
	return func(dst io.Writer, src io.Reader) (written int64, err error) {
		buf := make([]byte, 2*1024)
		for {
			if sock, ok := dst.(net.Conn); ok {
				sock.SetReadDeadline(time.Now().Add(timeout))
			}
			nr, er := src.Read(buf)
			if sock, ok := dst.(net.Conn); ok {
				sock.SetReadDeadline(time.Now().Add(time.Hour))
			}
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
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
