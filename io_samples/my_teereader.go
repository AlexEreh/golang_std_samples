package main

import "io"

type teeReader struct {
	r io.Reader
	w io.Writer
}

func (r teeReader) Read(data []byte) (n int, err error) {
	n, err = r.r.Read(data)
	if n <= 0 {
		return n, err
	}
	{
		n, err := r.w.Write(data[:n])
		if err != nil {
			return n, err
		}
	}
	return n, err
}

func MyTeeReader(r io.Reader, w io.Writer) io.Reader {
	return teeReader{
		r: r,
		w: w,
	}
}
