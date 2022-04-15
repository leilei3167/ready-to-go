package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

//内置的len函数返回的是string中字节数而不是rune字符数.第i个字节不代表是第i个字符!
//因为非ASCII编码字符可能会要多个字节

func main() {
	s := `hello world`
	fmt.Println(len(s))
	fmt.Printf("%v,%#v,%v\n", s[0], string(s[0]), s[0]) //返回的是字节值
	//子字符串操作
	fmt.Println("你" + s[3:])
	//字符串不可修改
	//不变性意味着如果两个字符串共享相同的底层数据的话也是安全的，这使得复制任何长度的字符串代价是低廉的。同样，一个字符串s和对应的子字符串切片s[7:]的操作也可以安全地共享相同的内存，因此字符串切片操作代价也是低廉的。在这两种情况下都没有必要分配新的内存

	//	s[0]='d'

	//原生的字符串面值是``包括的,其中没有转义操作
	fmt.Println(`nihao\n\t`) //nihao\n\t,原生字符串内无法直接写``
	//原生字符串面值广泛用于编写正则,HTML模板,提示信息等

	//Go语言range字符串时会隐式的解码UTF-8字符串,用range可以求得字符串中字符的数量(不是字节)
	for _, v := range `hello 世界` {
		fmt.Println(string(v))
	}
	fmt.Println(utf8.RuneCountInString("hello 你好"))

	/* 字符串切片和Byte切片 */
	/* 标准库中有四个包对字符串处理尤为重要：bytes、strings、strconv和unicode包。strings包提供了许多如字符串的查询、替换、比较、截断、拆分和合并等功能。

	   bytes包也提供了很多类似功能的函数，但是针对和字符串有着相同结构的[]byte类型。因为字符串是只读的，因此逐步构建字符串会导致很多分配和复制。在这种情况下，使用bytes.Buffer类型将会更有效，稍后我们将展示。

	   strconv包提供了布尔型、整型数、浮点数和对应字符串的相互转换，还提供了双引号转义相关的转换。

	   unicode包提供了IsDigit、IsLetter、IsUpper和IsLower等类似功能，它们用于给字符分类。每个函数有一个单一的rune类型的参数，然后返回一个布尔值。而像ToUpper和ToLower之类的转换函数将用于rune字符的大小写转换。所有的这些函数都是遵循Unicode标准定义的字母、数字等分类规范。strings包也有类似的函数，它们是ToUpper和ToLower，将原始字符串的每个字符都做相应的转换，然后返回新的字符串。 */
	fmt.Println(intsToString([]int{2, 3, 4, 1}))
	
}

//后缀删除和前缀/删除
func basename(s string) string {
	// Discard last '/' and everything before.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}
	// Preserve everything before last '.'.
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}

//use strings.LastIndex
func basename2(s string) string {
	//找到 /  就返回其最后出现的索引
	slash := strings.LastIndex(s, "/") // -1 if "/" not found
	s = s[slash+1:]
	//如果最后出现了.
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		s = s[:dot]
	}
	return s
}

//每隔3个以,分割
func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	//递归调用
	return comma(s[:n-3]) + "," + s[n-3:]

}


//一个字符串是包含只读字节的数组，一旦创建，是不可变的。相比之下，一个字节slice的元素则可以自由地修改。一个[]byte(s)转换是分配了一个新的字节数组用于保存字符串数据的拷贝，然后引用这个底层的字节数组。
//将一个字节slice转换到字符串的string(b)操作则是构造一个字符串拷贝，以确保s2字符串是只读的。

/* 	s := "abc"
b := []byte(s) //先创建底层数组容纳s的拷贝,再引用形成切片
s2 := string(b)
*/
/*
 为了避免转换中不必要的内存分配，bytes包和strings同时提供了许多实用函数:
func Contains(s, substr string) bool
func Count(s, sep string) int
func Fields(s string) []string
func HasPrefix(s, prefix string) bool
func Index(s, sep string) int
func Join(a []string, sep string) string
byte包中有对应的函数

bytes包还提供了Buffer类型用于字节slice的缓存。一个Buffer开始是空的，但是随着string、byte或[]byte等类型数据的写入可以动态增长，一个bytes.Buffer变量并不需要初始化，因为零值也是有效的：

*/

func intsToString(values []int) string {
	var buf bytes.Buffer //0值可用
	buf.WriteByte('[')
	for i, v := range values {
		if i > 0 {
			buf.WriteString(",")
		}
		fmt.Fprintf(&buf, "%d", v) //把v向buf中写入
	}
	buf.WriteByte(']')
	return buf.String()
	//当向bytes.Buffer添加任意字符的UTF8编码时，最好使用bytes.Buffer的WriteRune方法，但是WriteByte方法对于写入类似'['和']'等ASCII字符则会更加有效。
}
 