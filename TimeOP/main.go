package main

import (
	"fmt"
	"time"
)

/*  type Time struct {
    wall uint64
    ext  int64
    loc *Location
}*/
//https://mp.weixin.qq.com/s?__biz=MzI2MDA1MTcxMg==&mid=2648467756&idx=1&sn=0f38ab0b7b7ced36e360f70d78684404&chksm=f2474f43c530c6557ea1f7568d65db3b3a1ef714f18c388215d7eab3a34ceb68675f0799909f&cur_album_id=1506050738668486658&scene=190#rd
func main() {
	/* 1.基础操作 */
	//获取当前时间,时间戳
	fmt.Println(time.Now())            //2022-04-08 11:39:33.0670155 +0800 CST m=+0.001586501
	fmt.Println(time.Now().Unix())     //1649389195
	fmt.Println(time.Now().UnixNano()) //1649389271141104900

	//返回当前年月日时分秒、星期几、一年中的第几天等操作
	n := time.Now()

	year, month, day := n.Date()
	fmt.Printf("年%v月%v日%v\n", year, month, day)
	h, m, s := n.Clock()
	fmt.Printf("时%v分%v秒%v\n", h, m, s) //可使用n的hour等方法单独获取时等
	fmt.Println(n.Weekday())           //星期几
	fmt.Println(n.Year())
	fmt.Println(n.YearDay()) //一年的第几天

	//格式化时间为string,go中的时间模板是2006-01-02 15:04:05 即2006 12345
	fmt.Println(n.Format("2006-01-02 15:04:05"))
	fmt.Println(n.Format("2006/01/02 15:04:05"))
	fmt.Println(n.Format("15:04:05 2006.1.2"))

	/* 2.时间戳和日期字符串 */
	//时间戳转换为时间
	now := time.Now()
	layout := "2006-01-02 15:04:05"
	t := time.Unix(now.Unix(), 0) // 参数分别是：秒数,纳秒数
	fmt.Println(t.Format(layout))

	//将某个日期字符串解析成Time类型
	t1, _ := time.ParseInLocation("2006-01-02 3:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)
	fmt.Println(t1)

	/* 3.计算,比较日期 */
	//type Duration int64,单位是纳秒

}
