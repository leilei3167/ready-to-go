package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/leilei3167/ready-to-go/DBoperation/GORM/model"
)

//Model代表想要在哪个表操作,如果参数是一个实例并且其有主键,则
//相当于通过其主键查找到再修改,只修改一行

func Update() {
	//1.更新所有,Save会保存所有字段 即使字段是零值
	user := model.User{}
	model.Db.Find(&user, "name=?", "leilei")
	user.Name = "leileiNew"
	user.Age = 19
	model.Db.Save(&user)
	//1.1更新单个列,必须指定Where
	model.Db.Model(&model.User{}).Where("name=?", "leileiNew").Update("Age", 222)

	//2.更新多列,用map或结构体,Updates
	model.Db.Model(&model.User{}).Where("name=?", "leileiNew").Updates(map[string]interface{}{
		"name": "leileiNNNNN",
		"age":  "1111113",
		"job":  "golang",
	})
	//如果在更新时,只想更新某些字段,加上Select,或者Omit忽略某些字段
	model.Db.Model(&model.User{}).Where("id between ? and ?", 1, 5).Updates(map[string]interface{}{
		"name": "leilei111",
		"age":  "11233",
		"job":  "go5",
	})
	//默认情况下会阻止全局更新,会返回ErrMissingWhereClause错误
	//想要全局更新必须加条件,或者使用原生语句
	res := model.Db.Exec("update users set name =?", "zhutou")
	fmt.Println("受影响:", res.RowsAffected)

	//3.用sql表达式更新
	model.Db.Model(&model.User{}).Where("id between ? and ?", 1, 10).Update("age", gorm.Expr("age*?+?", 2, 3))

}

//可添加钩子函数,在更新前更新后 触发
//BeforeSave BeforeUpdate ...

func main() {
	Update()
}
