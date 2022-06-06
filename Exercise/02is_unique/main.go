package main

import (
	"strings"
)

func main() {

}

func IsUnique(s string) bool {
	//Count用于计算一个字符串包含另一个字符串的数量,""时会返回unicode数量+1的值
	if strings.Count(s, "") > 3000 { //不要输入大于3000的字符串
		return false
	}

	for _, v := range s {
		if v > 127 { //128以后的ascii无法在键盘打印
			return false
		}
		if strings.Count(s, string(v)) > 1 { //如果有多个,则返回false
			return false
		}
	}
	return true

}

func IsUnique2(s string) bool {
	if strings.Count(s, "") > 3000 {
		return false
	}

	for k, v := range s {
		if strings.Index(s, string(v)) != k { //v在s中第一次出现的位置,如果不与索引相等,说明有重复
			return false
		}
	}
	return true
}
