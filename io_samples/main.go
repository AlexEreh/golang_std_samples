package main

import (
	"fmt"
	"io"
)

func main() {
	defaultPipeUsage()
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
