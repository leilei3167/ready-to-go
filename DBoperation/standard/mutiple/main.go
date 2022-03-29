package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

//一个帖子有多条评论
type Post struct {
	Id       int
	Content  string
	Author   string
	Comments []Comment
}

type Comment struct {
	Id      int
	Content string
	Author  string
	Post    *Post
}

var Db *sql.DB

//初始化数据库
func init() {
	dsn := "postgresql://root:8888@localhost:5432/postgres?sslmode=disable"
	var err error
	Db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("链接数据库失败!", err)
	}
	if Db.Ping() == nil {
		fmt.Println("数据库初始化成功")

	}
}

//创建评论
func (comment *Comment) CreateCom() (err error) {
	if comment.Content == "" {
		err = errors.New("Post not found")
		return
	}
	err = Db.QueryRow("insert into comments (content,author,post_id)values($1,$2,$3) returning id",
		comment.Content, comment.Author, comment.Post.Id).Scan(&comment.Id)
	return

}

//根据id获取Post
func GetPost(id int) (post Post, err error) {
	post.Comments = []Comment{}
	err = Db.QueryRow("select id,content,author from posts where id=$1",
		id).Scan(&post.Id, &post.Author, &post.Content)

	//一个帖子可能有多个comments
	rows, err := Db.Query("select id,content,author from comments")
	if err != nil {
		return
	}
	for rows.Next() {
		//这些comment所属的帖子都相同
		comment := Comment{Post: &post}
		err = rows.Scan(&comment.Id, &comment.Author, &comment.Content)
		if err != nil {
			return
		}
		post.Comments = append(post.Comments, comment)

	}
	rows.Close()
	return

}

func (post *Post) Create() (err error) {
	err = Db.QueryRow("insert into posts (content,author)values($1,$2) returning id",
		post.Content, post.Author).Scan(&post.Id)
	return

}

func main() {

}
