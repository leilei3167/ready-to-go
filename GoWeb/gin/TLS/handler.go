package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"time"
)

// Product 交互模型,通过binding来约束输入
type Product struct {
	Username    string    `json:"username" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Price       int       `json:"price" binding:"gte=0"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type productHandler struct {
	sync.RWMutex
	products map[string]Product
}

func newProductHandler() *productHandler {
	return &productHandler{products: make(map[string]Product)}
}

// Create 增
func (p *productHandler) Create(c *gin.Context) { //作为函数变量进行绑定?
	p.Lock()
	defer p.Unlock()
	//1.参数解析
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil { //ShouldBind出错时不会直接返回客户端错误,以便于返回我们自定义的错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//2.参数校验,看是重复等
	product.Username = c.MustGet(gin.AuthUserKey).(string) //来到这里的一定是认证过的,将product的username设置为认证后的
	if _, ok := p.products[product.Name]; ok {             //用名字从数据库中查询,如有重复直接返回
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %s already exist", product.Name)})
		return
	}

	//3.逻辑处理(模拟插入数据库)
	product.CreatedAt = time.Now()

	p.products[product.Name] = product
	log.Printf("Register product %s success", product.Name)

	// 4. 返回结果(一般来说要返回成功信息)
	c.JSON(http.StatusOK, product)
}

func (p *productHandler) Get(c *gin.Context) {
	p.Lock()
	defer p.Unlock()
	product, ok := p.products[c.Param("name")] //直接获取查询参数name对应的值
	if !ok {                                   //不存在key会返回零值,ok会为false
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("can not find product %s",
			c.Param("name")).Error()})
		return
	}
	//参数校验
	if c.MustGet(gin.AuthUserKey).(string) != product.Username { //认证的User和操作的User是否为同一个
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permission denied"})
		return
	}
	c.JSON(http.StatusOK, product)
}
