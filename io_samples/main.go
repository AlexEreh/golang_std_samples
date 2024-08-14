package main

import (
	"fmt"
	"io"
	"time"
)

func main() {
	defaultPipeUsage()
}

func defaultPipeUsage() {
	reader, writer := MyPipe()
	go func() {
		for i := 0; i < 5; i++ {
			_, _ = writer.Write([]byte("abc"))
			time.Sleep(1 * time.Second)
		}
		_ = writer.Close()
	}()
	bytes, _ := io.ReadAll(reader)
	fmt.Println(string(bytes))
}
