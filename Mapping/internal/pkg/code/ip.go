package code

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

// ParseMixIP 实现对多IP地址的解析和验证,返回一个ip的切片,如果输入不合法,返回nil切片
//支持混合输入,如:8.8.8.8,127.0.0.1,9.9.9.9-9.9.9.255,8.8.8.8/24
func ParseMixIP(ip string) []string {
	var res []string
	//以 , 分隔
	ips := strings.Split(ip, ",")
	for _, ip := range ips {
		r := ParseIP(ip)
		if r == nil {
			return nil
		}
		res = append(res, r...)
	}
	return res
}

// ParseIP 对单个输入进行解析验证
func ParseIP(ip string) []string {
	switch {
	// 扫描/8时,只扫网关和随机IP,避免扫描过多IP;仍然会产生六十多万ip,可以考虑直接禁止使用;
	case strings.HasSuffix(ip, "/8"):
		return parseIP8(ip)
	//解析CIDR格式
	case strings.Contains(ip, "/"):
		return parseIP2(ip)
	//解析范围IP格式 192.168.1.1-192.168.1.100
	case strings.Contains(ip, "-"):
		return parseIP1(ip)
	//处理单个ip
	default:
		testIP := net.ParseIP(ip)
		if testIP == nil {
			return nil
		}
		return []string{ip}
	}
}

//解析范围  192.168.111.1-192.168.112.255
func parseIP1(ip string) []string {
	IPRange := strings.Split(ip, "-")
	testIP := net.ParseIP(IPRange[0])
	var AllIP []string //TODO:根据范围预先分配内存,避免频繁扩容
	if len(IPRange[1]) < 4 {
		Range, err := strconv.Atoi(IPRange[1])
		if testIP == nil || Range > 255 || err != nil {
			return nil
		}
		SplitIP := strings.Split(IPRange[0], ".")
		ip1, err1 := strconv.Atoi(SplitIP[3])
		ip2, err2 := strconv.Atoi(IPRange[1])
		PrefixIP := strings.Join(SplitIP[0:3], ".")
		if ip1 > ip2 || err1 != nil || err2 != nil {
			return nil
		}
		for i := ip1; i <= ip2; i++ {
			AllIP = append(AllIP, PrefixIP+"."+strconv.Itoa(i))
		}
	} else {
		SplitIP1 := strings.Split(IPRange[0], ".")
		SplitIP2 := strings.Split(IPRange[1], ".")
		if len(SplitIP1) != 4 || len(SplitIP2) != 4 {
			return nil
		}
		start, end := [4]int{}, [4]int{}
		for i := 0; i < 4; i++ {
			ip1, err1 := strconv.Atoi(SplitIP1[i])
			ip2, err2 := strconv.Atoi(SplitIP2[i])
			if ip1 > ip2 || err1 != nil || err2 != nil {
				return nil
			}
			start[i], end[i] = ip1, ip2
		}
		startNum := start[0]<<24 | start[1]<<16 | start[2]<<8 | start[3]
		endNum := end[0]<<24 | end[1]<<16 | end[2]<<8 | end[3]
		for num := startNum; num <= endNum; num++ {
			ip := strconv.Itoa((num>>24)&0xff) + "." + strconv.Itoa((num>>16)&0xff) + "." + strconv.Itoa((num>>8)&0xff) + "." + strconv.Itoa((num)&0xff)
			AllIP = append(AllIP, ip)
		}
	}

	return AllIP //去重
}

func parseIP2(host string) (hosts []string) {
	_, ipNet, err := net.ParseCIDR(host) //验证是否是合法CIDR地址
	if err != nil {
		return
	}
	hosts = parseIP1(ipRange(ipNet))
	return
}

func ipRange(c *net.IPNet) string {
	start := c.IP.String()
	mask := c.Mask
	bcst := make(net.IP, len(c.IP))
	copy(bcst, c.IP)
	for i := 0; i < len(mask); i++ {
		ipIdx := len(bcst) - i - 1
		bcst[ipIdx] = c.IP[ipIdx] | ^mask[len(mask)-i-1]
	}
	end := bcst.String()
	return fmt.Sprintf("%s-%s", start, end) //返回用-表示的ip段,192.168.1.0-192.168.255.255
}

func parseIP8(ip string) []string {
	realIP := ip[:len(ip)-2]      //去除/8
	testIP := net.ParseIP(realIP) //验证
	if testIP == nil {
		return nil
	}

	IPrange := strings.Split(ip, ".")[0] //以.分隔后 只取第0位
	var AllIP []string
	for a := 0; a <= 255; a++ {
		for b := 0; b <= 255; b++ {
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPrange, a, b, 1))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPrange, a, b, 2))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPrange, a, b, 4))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPrange, a, b, 5))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPrange, a, b, RandInt(6, 55)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPrange, a, b, RandInt(56, 100)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPrange, a, b, RandInt(101, 150)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPrange, a, b, RandInt(151, 200)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPrange, a, b, RandInt(201, 253)))
			AllIP = append(AllIP, fmt.Sprintf("%s.%d.%d.%d", IPrange, a, b, 254))
		}
	}
	return AllIP
}

func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

// RemoveDuplicateIPs 去重,影响性能
func RemoveDuplicateIPs(old []string) []string {
	var result []string
	temp := map[string]struct{}{}
	for _, item := range old {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
