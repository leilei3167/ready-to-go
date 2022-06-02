package service

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"mapping/internal/master/model"
	"strings"
)

// CheckParams 根据不同的输入参数选择不同的解析方式,支持json,和上传文件(PostForm)
func CheckParams(c *gin.Context) (model.Param, error) {
	var p model.Param
	content := c.GetHeader("Content-Type") //根据Content-Type判断是从文件解析还是从json
	switch {
	case strings.Contains(content, "application/json"):
		p.Type = "json"
		//json请求
		var s model.ScanReq
		err := c.BindJSON(&s)
		if err != nil {
			return p, err
		}
		Ports, err := ParsePorts(s.ScanType)
		if err != nil {
			return p, err
		}
		p.Ports = Ports
		p.Ip = s.Targets
		return p, nil

	case strings.Contains(content, "multipart/form-data"):
		p.Type = "File"
		scanType := c.PostForm("scanType")
		if scanType == "" {
			return p, errors.New("scanType和file不能为空")
		}
		Ports, err := ParsePorts(scanType)
		if err != nil {
			return p, errors.Wrap(err, "解析scanType出错")
		}
		//获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			return p, errors.Wrap(err, "获取文件出错")
		}
		if !strings.HasSuffix(file.Filename, ".txt") { //只支持txt文件
			return p, errors.New("只支持.txt格式的文件!")
		}
		if file.Size >= 256*1024 {
			return p, errors.New("文件不应超过256kb")
		}
		p.Ports = Ports
		p.File = file
		return p, nil

	default:
		return p, errors.New("解析请求失败")
	}

}
