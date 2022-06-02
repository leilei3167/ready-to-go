package service

import (
	"github.com/pkg/errors"
	"mapping/internal/pkg/code"
	"sort"
	"strconv"
	"strings"
)

// ParsePorts PortsParser 实现解析 22,24,25,67-78,6000-7000或者top100的输入,解析成一个[]int的切片
//支持单个输入:22 支持内置模式:top1000,top100,all 支持范围输入:22,23,24-56
//对于范围值起始值必须小于结束值,否则报错
func ParsePorts(s string) ([]int, error) {
	switch {
	case s == "top100":
		return code.TOP100, nil
	case s == "top1000":
		return code.TOP1000, nil
	case s == "all":
		return code.All, nil
		//return nil, errors.New("ParsePorts:暂时不支持all扫描")
	case s == "":
		return nil, errors.New("ParsePorts:无效的scantype")
	default:
		p := strings.Split(s, ",")
		var res []int
		for _, input := range p {
			input = strings.TrimSpace(input)
			if strings.Contains(input, "-") { //-1 ["","1"]
				r, err := transScope(input)
				if err != nil {
					return nil, errors.Wrap(err, "ParsePorts")
				}
				res = append(res, r...)
			} else {
				//转换为单个数字,并且范围在1到65535之间
				port, err := strconv.Atoi(input)
				if err != nil {
					return nil, errors.Wrap(err, "ParsePorts")
				} else if !isPort(port) {
					return nil, errors.New("port必须是在1-65535之间的整数")
				}
				res = append(res, port)
			}
		}
		res = RmDuplicatePorts(res)
		return res, nil
	}

}

func transScope(v string) ([]int, error) {
	tri := strings.SplitN(v, "-", 2)
	//将开始和结束解析为int,出错说明输入有误
	start, err := strconv.Atoi(tri[0])
	if err != nil {
		return nil, err
	} else if !isPort(start) {
		return nil, errors.New("port必须是在1-65535之间的整数")
	}
	end, err := strconv.Atoi(tri[1])
	if err != nil {
		return nil, err
	} else if !isPort(end) {
		return nil, errors.New("port必须是在1-65535之间的整数")
	}
	//以start和end来构建ports切片
	n := end - start + 1
	if n < 1 {
		return nil, errors.New("起始端口必须小于结束端口")
	}
	//log.Printf("start:%#v,end:%#v,n:%v", start, end, n)
	ports := make([]int, n)
	for i := 0; i < n; i++ {
		ports[i] = start
		start++
	}
	return ports, nil
}

func isPort(p int) bool {
	if p > 65535 || p < 1 {
		return false
	}
	return true
}

// RmDuplicatePorts int切片去重并递增排序
func RmDuplicatePorts(v []int) []int {
	tmp := make(map[int]struct{})
	result := make([]int, 0, len(v))
	for _, p := range v {
		if _, ok := tmp[p]; !ok {
			tmp[p] = struct{}{}
			result = append(result, p)
		}
	}
	sort.Ints(result)
	return result
}
