package models

/* 此包对DB的句柄进行统一管理 */
import (
	"fmt"
	"log"

	//老版的驱动
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/leilei3167/ready-to-go/GoWeb/Blog2/pkg/setting"
)

//句柄
var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

func init() {
	//先从配置文件中获取数据库的配置,私有项,因为不希望数据库相关的内容
	//被随意调用
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)
	//Cfg在setting包中已创建,这里直接调用
	//Section会在没有分区时创建,而Getsection没有分区会返回错误
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}
	//获取链接数据库必须的信息
	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()

	//根据配置文件链接数据库
	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	if err != nil {
		log.Println("连接数据库失败:", err)
	}
	//默认的表,输入auth 就会返回blog_auth
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}
	//数据库的配置 链接数等
	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	
}

//关闭数据库
func CloseDB() {
	defer db.Close()
}
