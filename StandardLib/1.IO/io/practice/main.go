package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	file, _ := os.Open("./test.txt")
	a := make([]byte, 55)
	n, err := ReadAtleast(file, a, 51)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("读取到:", n)

	/*
	   //指定len和不指定len的区别
	   b := make([]byte, 0, 12)
	   	c := make([]byte, 12)
	   	fmt.Printf("b:%#v,len(b):%d,cap(b):%d\nc:%#v,len(c):%d,cap(c):%d\n\n", b, len(b), cap(b), c, len(c), cap(c)) */

}

//io中的WriteString函数
func WriteString(w io.Writer, s string) (n int, err error) {
	if ws, ok := w.(io.StringWriter); ok {
		return ws.WriteString(s)
	}
	return w.Write([]byte(s)) //将string转换为byte后用Write

}

//ReadAtleast
func ReadAtleast(r io.Reader, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, errors.New("缓冲区小于要读取的字节数")

	}
	for n < min && err == nil {
		var nn int
		nn, err = r.Read(buf[n:]) //最多读取len(buf[n:])个字节
		n = n + nn                //直到读取到err==EOF 或者n>=min
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = errors.New("未知错误")
	}
	return
}

//仿造copy
func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	//首先判断两个参数有没有实现ReaderFrom或者WriterTo
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src) //很少有实现以上这两个方法的
	}

	if buf == nil { //没有缓冲区就创建一个
		size := 32 * 1024
		//src是否是limitedReader,如果是以它的限定为准
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf) //从数据流读取len(buf)个字节,放入缓冲区
		if nr > 0 {             //读取到数据则写入
			nw, ew := dst.Write(buf[0:nr]) //只取有数据的部分
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("无效的写入")
				}
			}
			written = written + int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = errors.New("读写字节数不一致") //读的数不等于写入的数据,中断
				break
			}
			if er != nil {
				if er != io.EOF {
					err = er
				}
				break
			}

		}
	}
	return written, err
}

//ReadAll
func ReadAll(r io.Reader) ([]byte, error) {
	b := make([]byte, 512)
	for {
		if len(b) == cap(b) { //容量满
			b = append(b, 0)[:len(b)] //容量得到了扩充

		}
		n, err := r.Read(b[len(b):cap(b)]) //只往未写入的值的部分读
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}

	}

}
