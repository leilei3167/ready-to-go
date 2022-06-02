package master

import (
	"github.com/gin-gonic/gin"
	"log"
	v1 "mapping/internal/master/api/v1"
	"mapping/internal/pkg/db"
	"mapping/internal/pkg/mid"
	"time"
)

type Master struct {
	//web控制
	Addr string
	E    *gin.Engine
}

func NewMaster(addr string, mongoaddr string) *Master {
	var M Master
	M.E = gin.Default()
	M.Addr = addr
	db.InitDB(mongoaddr)
	return &M
}

func (m *Master) initHanders() {
	//初始化数据库

	//中间件
	m.E.Use(mid.RateLimitMiddleware(time.Second, 5, 1))
	m.E.POST("/scantask", v1.ScanFromUpload)
	m.E.GET("/ipinfo", v1.GetResult)

}

func (m *Master) Run() {
	m.initHanders()
	log.Fatal(m.E.Run(m.Addr))
}
