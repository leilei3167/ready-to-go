package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

func main() {
	pipeReader, pipeWriter := io.Pipe() //构建两端
	go PipeWrite(pipeWriter)
	go PipeRead(pipeReader)
	time.Sleep(30 * time.Second)
}

func PipeWrite(w *io.PipeWriter) { //写入端
	data := []byte("hello im leilei!")
	for i := 0; i < 3; i++ {
		n, err := w.Write(data)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("写入%v字节!\n", n)

	}
	w.CloseWithError(errors.New("写入段已关闭哦")) //直接close的话读取段会接受到EOF错误
}

func PipeRead(reader *io.PipeReader) {
	buf := make([]byte, 128)
	for {
		fmt.Println("接口端开始阻塞5秒钟...")
		time.Sleep(5 * time.Second)
		fmt.Println("接收端开始接受")
		n, err := reader.Read(buf) //一直读取直到错误
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("收到字节: %d\n buf内容: %s\n", n, buf)
	}
}
