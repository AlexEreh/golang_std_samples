package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
	"sync"
)

type sizeWriter struct {
	size int64
}

func (w *sizeWriter) Write(data []byte) (n int, err error) {
	w.size += int64(len(data))
	return len(data), nil
}

func (w *sizeWriter) Size() int64 {
	return w.size
}

func Size(r io.Reader) int64 {
	w := &sizeWriter{}
	_, _ = io.Copy(w, r)
	return w.Size()
}

///

func main() {
	//teeReaderVariant()
	//multiWriterVariant1()
	multiWriterVariant2()
	//defaultPipeUsage()
}

func multiWriterVariant1() {
	r := strings.NewReader("abcde")

	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()
	r3, w3 := io.Pipe()

	go func() {
		mw := io.MultiWriter(w1, w2, w3)

		_, _ = io.Copy(mw, r)

		_ = w1.Close()
		_ = w2.Close()
		_ = w3.Close()
	}()
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		size := Size(r1)
		fmt.Println(size)
		wg.Done()
	}()
	go func() {
		b2, _ := io.ReadAll(r2)
		fmt.Println(string(b2))
		wg.Done()
	}()
	go func() {
		h := sha256.New()
		b3, _ := io.ReadAll(r3)
		h.Write(b3)
		fmt.Println(fmt.Sprintf("%x", h.Sum(nil)))
		wg.Done()
	}()
	wg.Wait()
}

func multiWriterVariant2() {
	r := strings.NewReader("abcde")
	sw := sizeWriter{}
	h := sha256.New()
	r1, w1 := io.Pipe()
	go func() {
		mw := io.MultiWriter(&sw, h, w1)
		_, _ = io.Copy(mw, r)
		_ = w1.Close()
	}()
	fmt.Printf("Size pre read: %d\n", sw.Size())
	fmt.Printf("Hash pre read: %x\n", h.Sum(nil))
	b, _ := io.ReadAll(r1)
	fmt.Printf("Result: %s\n", string(b))
	fmt.Printf("Size: %d\n", sw.Size())
	fmt.Printf("Hash: %x\n", h.Sum(nil))
}

func teeReaderVariant() {
	r := strings.NewReader("abcde")

	sw := sizeWriter{}

	r1 := io.TeeReader(r, &sw)

	h := sha256.New()

	r2 := io.TeeReader(r1, h)

	fmt.Println("Size pre read: ", sw.Size())
	fmt.Println("Hash pre read: ", h.Sum(nil))

	data, _ := io.ReadAll(r2)

	fmt.Println(data)
	fmt.Println("Size post read: ", sw.Size())
	fmt.Println("Hash post read", h.Sum(nil))
}

func defaultPipeUsage() {
	reader, writer := MyPipe()
	go func() {
		for i := 0; i < 5; i++ {
			_, _ = writer.Write([]byte("abc"))
		}
		_ = writer.Close()
	}()
	bytes, _ := io.ReadAll(reader)
	fmt.Println(string(bytes))
}
