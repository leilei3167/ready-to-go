package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

//创建和数据库对应的结构体
type Post struct {
	Id      int
	Content string
	Author  string
}

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "postgresql://root:8888@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = Db.Ping()
	if err != nil {
		log.Fatal("连接失败:", err)
	} else {

		fmt.Println("连接成功")
	}

}

func main() {

}
