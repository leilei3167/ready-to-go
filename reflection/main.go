package main

import (
	"fmt"
	"reflect"
)

/*
go的反射包可以在运行时检查参数的type,最常用的地方就是结构体的转换,如某个函数根据传入的不同结构体
生成不同的sql
自带的reflect包实现了运行时进行反射的功能,可以帮助识别一个空接口类型的变量的底层具体类型和值

*/

/*
经过反射后interface{}类型的变量的底层具体类型由reflect.Type表示，底层值由reflect.Value表示。
reflect包里有两个函数reflect.TypeOf() 和reflect.ValueOf()
分别能将interface{}类型的变量转换为reflect.Type和reflect.Value。
这两种类型是创建我们的SQL生成器函数的基础。
*/

/*
反射三法则:
	1.从接口值可以反射出反射对象。
我们能够吧Go中的接口类型变量转换成反射对象，
上面提到的reflect.TypeOf和 reflect.ValueOf 就是完成的这种转换

	2.从反射对象可反射出接口值。
指的是我们能把反射类型的变量再转换回到接口类型

	3.要修改反射对象，其值必须可设置
与反射值是否可以被更改有关
*/

/*
反射获取结构体字段的方法:
reflect.StructField类型提供的各种方法可以用于获取结构体内的字段的类型属性,StructField可以通过
reflect.Type对象提供的几种方式拿到(Field)


*/

type order struct {
	ordId      int `json:"orider_id"`
	customerId int
}

func createQuery(q interface{}) {
	t := reflect.TypeOf(q)  //获取底层type类型,表示的是接口的实际类型(main.order)
	v := reflect.ValueOf(q) //获取底层value
	k := t.Kind()           //Kind表示的是type所属的种类,即main.order是一个struct类型
	//就像map[string]string的Kind就会是map
	fmt.Printf("Type底层:%#v 普通表示:%v \n ", t, t)
	fmt.Printf("Kind:%#v 普通表示:%v \n", k, k)
	fmt.Printf("Value:%#v 普通表示:%v \n ", v, v)

	if t.Kind() != reflect.Struct {
		panic("预期外的类型!")
	}

	fmt.Println()
	fmt.Println()
	/*
		Type的NumField可以获取到结构体的字段数量,如果不是结构体会直接pannic
		Type的Field方法可以获取到结构体的第几个字段的structField类型,从而得到其类型和值或者tag

	*/
	for i := 0; i < t.NumField(); i++ {
		fmt.Println("FieldName:", t.Field(i).Name, "FiledType:", t.Field(i).Type,
			"FiledValue:", v.Field(i), "FieldTag:", t.Field(i).Tag)
	}
	/*
		Value类型的Int(),String()等方法可以将Value类型转换为实际类型的值

	*/
}

type employee struct {
	name    string
	id      int
	address string
	salary  int
	country string
}

//一个适用于任何结构体的sql生成器
func CreateSQL(q interface{}) string {
	t := reflect.TypeOf(q)  //获取空接口的底层类型
	v := reflect.ValueOf(q) //获取其底层值
	if v.Kind() != reflect.Struct {
		panic("错误的类型")
	}

	tableName := t.Name() //将结构体的类型名字作为sql的表名!
	sql := fmt.Sprintf("INSERT INTO %s", tableName)
	columns := "("
	values := "VALUES ("
	//遍历每一个字段
	for i := 0; i < v.NumField(); i++ {
		// 注意reflect.Value 也实现了NumField,Kind这些方法
		// 这里的v.Field(i).Kind()等价于t.Field(i).Type.Kind()
		switch v.Field(i).Kind() {
		case reflect.Int:
			if i == 0 { //第一个值
				columns += fmt.Sprintf("%s", t.Field(i).Name)
				values += fmt.Sprintf("%d", v.Field(i).Int())
			} else { //不是第一个值,加逗号
				columns += fmt.Sprintf(", %s", t.Field(i).Name)
				values += fmt.Sprintf(", %d", v.Field(i).Int())
			}
		case reflect.String:
			if i == 0 {
				columns += fmt.Sprintf("%s", t.Field(i).Name)
				values += fmt.Sprintf("'%s'", v.Field(i).String())
			} else {
				columns += fmt.Sprintf(", %s", t.Field(i).Name)
				values += fmt.Sprintf(", '%s'", v.Field(i).String())
			}

		}
	}
	//添加右括号
	columns += ");"
	values += ");"
	sql += columns + values
	fmt.Println(sql)
	return sql
}

func main() {
	o := order{
		ordId:      456,
		customerId: 56,
	}
	//createQuery(o)
	e := employee{
		name:    "leieli",
		id:      123,
		address: "dsahic",
		salary:  55,
		country: "cxjzo",
	}
	CreateSQL(o)
	CreateSQL(e)

}
