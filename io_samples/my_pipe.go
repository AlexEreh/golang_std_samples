package main

import (
	"errors"
	"io"
	"sync"
)

var errClosed = errors.New("closed pipe")
var errInsufficientLength = errors.New("insufficient length")

// pipeStruct без какого-либо декорирования и всяких отдельных ридеров, врайтеров, ибо самописное (да, мне лень)
type pipeStruct struct {
	// writeMutex нужен на случай, если метод Write будут вызывать одновременно сразу несколько горутин.
	// Следовательно, надо будет обкладывать запись, ибо это будет место для потенциального Data Race
	writeMutex sync.Mutex
	// Нужен для того, чтобы удостовериться, что мы закрыли канал единожды.
	// Примечание: закрытие закрытого канала вызовет панику.
	doneOnce sync.Once
	// В этом канале нам важно лишь два состояния - он закрыт и не закрыт.
	// Примечание: чтение из закрытого канала нам вернёт дефолтное значение типа, то есть пустую структуру.
	doneCh chan struct{}
	// Канал в который записывается последнее прочтённое число байт
	readCh  chan int
	writeCh chan []byte
}

func (p *pipeStruct) Read(dst []byte) (n int, err error) {
	select {
	case <-p.doneCh:
		return 0, errClosed
	default:
	}

	select {
	case <-p.doneCh:
		return 0, errClosed
	case bw := <-p.writeCh:
		nr := copy(dst, bw)
		p.readCh <- nr
		return nr, nil
	}
}

func (p *pipeStruct) Write(dst []byte) (n int, err error) {
	select {
	case <-p.doneCh:
		return 0, errClosed
	default:
		p.writeMutex.Lock()
		defer p.writeMutex.Unlock()
	}

	if len(dst) <= 0 {
		return 0, errInsufficientLength
	}

	bruh := func() (int, error) {
		select {
		case <-p.doneCh:
			return 0, errClosed
		case p.writeCh <- dst:
			// Место дедлока? Write будет ждать чтения из ридера пайпы
			lastReadBytes := <-p.readCh
			dst = dst[lastReadBytes:]
			n += lastReadBytes
		}
		return n, nil
	}
	n, err = bruh()
	if err != nil {
		return n, err
	}
	for len(dst) > 0 {
		n, err = bruh()
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

// Close закрывает врайтер
func (p *pipeStruct) Close() error {
	p.doneOnce.Do(func() {
		close(p.doneCh)
	})
	return nil
}

func MyPipe() (io.Reader, io.WriteCloser) {
	p := pipeStruct{
		writeMutex: sync.Mutex{},
		doneOnce:   sync.Once{},
		doneCh:     make(chan struct{}),
		readCh:     make(chan int),
		writeCh:    make(chan []byte),
	}
	return &p, &p
}
