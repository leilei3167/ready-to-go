package main

import (
	"bufio"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/upload", GetFile)

	r.Run()
}

func GetFile(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusInternalServerError, "读取上传文件出错!%v", err)
		return
	}

	f, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "读取上传文件出错!%v", err)
		return
	}
	defer f.Close()

	s := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		a := scanner.Text()
		isCDIR := strings.Contains(a, "/")
		switch isCDIR {
		case true:
			ip, _, err := net.ParseCIDR(a)
			if err != nil {
				s = append(s, err.Error())
				continue
			}
			s = append(s, ip.String()+"cidr")
		default: //不是网段的话
			if ip := net.ParseIP(a); ip != nil {
				s = append(s, ip.String())
			}
		}

	}

	if err := scanner.Err(); err != nil {
		c.String(http.StatusInternalServerError, "扫描出错!%v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": s,
	})
	return

}
