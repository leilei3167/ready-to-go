package main

/* 关键点是创建相应的文件Create
通过io.Copy将http请求的Body写入到create创建的文件
要注意验证下载是否正常
使用os.Stat获取下载文件 并判断其字节数是否和response的
contetLength一致
*/
import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	cli := http.Client{}

	req, _ := http.NewRequest("GET", "http://localhost:8080/download", nil)

	res, _ := cli.Do(req)
	defer res.Body.Close()

	file, _ := os.Create("downloadfile.json")

	n, err := io.Copy(file, res.Body)
	defer file.Close()
	fileInfo, _ := os.Stat("downloadfile.json")
	//应该也可以判断Copy写入的字节数和Contentlength的值
	if err != nil || fileInfo.Size() != res.ContentLength {
		fmt.Println("下载出错!", err)
	}

	fmt.Println("下载成功:", n, "字节")

}
