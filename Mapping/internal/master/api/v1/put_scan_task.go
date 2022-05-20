package v1

import (
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"log"
	model2 "mapping/internal/master/model"
	"mapping/internal/master/service"
	"net/http"
	"strings"
	"sync"
)

// ScanFromUpload 获取上传的IP列表文件,产生任务
func ScanFromUpload(c *gin.Context) {
	//解析参数,指定要扫描的端口
	p, err := service.CheckParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseWithErr("解析参数出错", err))
		return
	}
	//构建IP解析
	var parser *service.IPparser
	switch p.Type {
	case "file":
		file, err := p.File.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, responseWithErr("打开上传文件错误", err))
			return
		}
		defer file.Close()
		parser = service.NewIPparser(file)
	case "json":
		ips := strings.Join(p.Ip, "\n") //手动添加换行便于扫描
		ipReader := strings.NewReader(ips)
		parser = service.NewIPparser(ipReader)
	}

	//TODO:考虑生产者复用
	producer, err := service.GetDefaultProducer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseWithErr("创建producer错误", err))
		return
	}
	go func() {
		for producerError := range producer.Errors() {
			log.Println(producerError)
		}
	}()

	var count uint
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := range parser.IPChan {
			//解析好的ip发送至kafka
			a := &model2.ScanTask{
				IP:    i,
				Ports: p.Ports,
			}
			msg := sarama.ProducerMessage{Value: a, Topic: "test_10"}
			producer.Input() <- &msg
			count++
			//	log.Printf("已发送IP:%v:%v", a.IP, a.Ports)
		}
	}()
	in := make([]string, 0)
	go func() {
		defer wg.Done()
		for i := range parser.InvalidIPs {
			in = append(in, i)
		}
	}()

	parser.ReadAndParse() //开始解析ip

	wg.Wait()        //说明所有G已退出
	producer.Close() //关闭生产者

	c.JSON(http.StatusOK, model2.Response{
		InvalidIP: in,
		HasSent:   count,
	})

}

func responseWithErr(text string, err error) gin.H {

	return gin.H{
		"msg":  text,
		"err:": err.Error(),
	}

}
