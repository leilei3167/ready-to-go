package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	a := make([]int, 0)
	s := strings.Split("1,2,3,45,6,1-47", ",")
	fmt.Printf("s以逗号分割后s的值:%#v\n", s)
	for _, v := range s {
		//判断是否是范围表达
		if strings.Contains(v, "-") {
			ports, err := Scope(v)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Scope后的:%#v", ports)
			a = append(a, ports...)
		} else {
			//不是范围表达转换为单个数字
			x, err := strconv.Atoi(v)
			if err != nil {
				log.Fatal("输入不合法,端口间以逗号分隔,范围以开始-结束处表示", err)
			}
			a = append(a, x)
		}

	}
	a = delrep(a)
	fmt.Printf("a去重后的值:%#v\n", a)
}

//处理范围端口
func Scope(v string) ([]int, error) {
	//以-来切分开始和结束端口
	tri := strings.SplitN(v, "-", 2) //切分成开始和结束2个部分
	log.Printf("SplitN:%#v", tri)
	//将开始和结束解析为int,出错说明输入有误
	start, err := strconv.Atoi(tri[0])
	if err != nil {
		return nil, err
	}
	end, err := strconv.Atoi(tri[1])
	if err != nil {
		return nil, err
	}
	//以start和end来构建ports切片
	n := end - start + 1
	if n < 1 {
		return nil, fmt.Errorf("结束端口%v小于开始端口%v", start, end)
	}
	log.Printf("start:%#v,end:%#v,n:%v", start, end, n)
	ports := make([]int, n)
	for i := 0; i < n; i++ {
		ports[i] = start
		start++
	}
	return ports, nil

}

//字符串去重
func delrep(v []int) []int {
	tmp := make(map[int]struct{})
	result := make([]int, 0)

	for _, v := range v {
		_, ok := tmp[v]
		if ok {
			continue
		} else {
			tmp[v] = struct{}{}
			result = append(result, v)
		}

	}
	return result
}
