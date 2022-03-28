package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	/* 字符串的常用操作 */
	//1.判断是否包含
	fmt.Println(strings.Contains("leieliehahah", "yes")) //false
	fmt.Println(strings.Contains("leieliehahah", "lei")) //true
	//2.字符串链接
	s := []string{"l", "e", "i"}
	fmt.Println(strings.Join(s, ".")) //l.e.i
	//3.查找出现的位置,没有返回-1
	fmt.Println(strings.Index("hello world!", "!")) //11
	fmt.Println(strings.Index("hello world!", "?")) //-1
	//4.重复指定字符串count次
	fmt.Println("haha" + strings.Repeat("nni", 5))
	//5.替换指定字符,-1表示不限制替换次数
	fmt.Println(strings.Replace("hahaleinihao", "haha", "wu", -1))
	//6.分割字符串
	fmt.Printf("%q\n", strings.Split("a man a plan a canal pananma", "a"))
	fmt.Printf("%q\n", strings.Split(" xyz ", ""))
	//7.在头尾去除指定字符串
	fmt.Printf("%q\n", strings.Trim(" !!! Achtung !!! ", "! "))
	fmt.Printf("%q\n", strings.Trim("a man a plan a canal pananma", "a"))
	//8.去掉空格,并按空格分割,返回切片
	fmt.Printf("Fields are: %q\n", strings.Fields("  foo bar  baz   "))

	/* 字符串的转换 */
	//1.Append函数将指定值转换为字符串后加入到字节数组中
	a := make([]byte, 0)
	a = strconv.AppendBool(a, false)
	a = strconv.AppendInt(a, 12321321, 10) //10进制
	a = strconv.AppendQuote(a, "abcdefg")
	a = strconv.AppendQuoteRune(a, '单')
	fmt.Printf("%s\n", string(a))
	//2.Format将其他类型转换为字符串
	a1 := strconv.FormatBool(false)
	b := strconv.FormatFloat(123.23, 'g', 12, 64)
	c := strconv.FormatInt(1234, 10)
	d := strconv.FormatUint(12345, 10)
	e := strconv.Itoa(1023)
	fmt.Println(a1, b, c, d, e)

	//3.Parse函数将字符串转换为其他类型
	c1, err := strconv.ParseInt("1234", 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("val is %s ,type:%d\n", c1, c1)
}
