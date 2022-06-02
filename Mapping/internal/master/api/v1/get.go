package v1

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
	"mapping/internal/pkg/db"
	"net/http"
)

// GetResult 根据IP查询单个
func GetResult(c *gin.Context) {
	//接收查询参数
	var req string
	req = c.Query("ip")
	if req == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "ip不能为空!"})
		return
	}
	log.Println(req)
	//查询
	coll := db.MgoClient.Database("result").Collection("IP")

	filter := bson.D{{"ip", req}}
	var r db.ScanResult
	err := coll.FindOne(context.Background(), filter).Decode(&r)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{
				"msg": "没有记录",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, responseWithErr("内部错误", err))
		return
	}
	c.JSON(http.StatusOK, r)
	return
}

//根据端口返回多个
