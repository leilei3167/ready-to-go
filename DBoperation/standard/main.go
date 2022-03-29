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
func GetPosts(limit int) (posts []Post, err error) {
	//查询多个结果
	rows, err := Db.Query("select * from posts limit $1 ", limit)
	if err != nil {
		return
	}
	//逐行扫描
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Author, &post.Content)
		if err != nil {
			return
		}
		posts = append(posts, post)

	}
	rows.Close()
	return

}

//获取一条结果
func GetPost(id int) (a Post, err error) {

	err = Db.QueryRow("select * from posts where id = $1", id).Scan(&a.Id, &a.Content, &a.Author)
	if err != nil {
		return
	}
	return

}

//插入
func (post *Post) Create() (err error) {
	statment := "insert into posts (content,author) values($1,$2) returning id"
	stmt, err := Db.Prepare(statment)
	if err != nil {
		return
	}
	defer stmt.Close()
	//执行并获取returning的结果
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	return
}

//更新
func (post *Post) Update() (err error) {
	_, err = Db.Exec("update posts set content=$2,author=$3,where id =$1",
		post.Id, post.Content, post.Author)
	return

}

//删除
func (post *Post) Delete() (err error) {
	_, err = Db.Exec("delete from posts where id = $1", post.Id)
	return

}

func main() {
	var a Post
	a.Author = "leilei"
	a.Content = "nihaonihao"
	//	a.Create()
	fmt.Println(a)
	res, _ := GetPosts(3)
	fmt.Println(res)

	fmt.Println(GetPost(1))

}
