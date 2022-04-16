package main

import "fmt"

/* slice
底层数组分配在堆上


var a []int //声明一个切片,值为零值 nil
b:=[]int{} //声明一个空切片,值不为nil,而是为空,其指针指向空结构体(空结构体不会分配内存)

*/

func main() {
	//1.以空切片开始,append将复制a(24字节),并向其追加值,若容量不够则发生扩容(分配新内存),旧底层数组的值将被复制到新的数组,完成后旧数组无人引用将被GC
	var a []int
	for i := 0; i < 10; i++ {
		a = append(a, i+1)
		fmt.Printf("len %v cap %v value:%v\n", len(a), cap(a), a[i])
	}
	//2.如果提前直到预计的容量一定要提前指定长度,可减少扩容导致的内存分配和GC
	b := make([]int, 10) //100个零值的数组
	fmt.Printf("%#v\n", b)
	for i := 0; i < 10; i++ {
		b[i] = i + 1 //append会产生一次slice的复制
		fmt.Printf("len %v cap %v value:%v\n", len(a), cap(a), a[i])
	}
	//3.使用[:]切分切片要注意引用同一个底层数组的问题,使用Copy可以彻底复制一个切片到一个新切片,不用担心引用的问题

	//4.使用append比较容易造成副作用,尤其是其发生扩容时,如果在扩容前底层数组被某个地方引用,append发生扩容后,旧底层数组的引用不会随新底层数组的创建而发生迁移,将导致
	//内存泄漏

	//5.内存泄漏:保持对堆中的某个值的引用,并且永远不释放(因为有引用GC也不会回收),导致内存泄漏,要排查泄漏唯一的方法是GC trace(每次GC时内存是否会增加)
	//经典的出现内存泄漏的场景
	//一:该关闭的协程没有关闭,协程持有的东西没有被释放;
	//二:map,map只放入却没有删除,map将会无限增长,占用内存;(严格不算内存泄漏)
	//三:slice的append操作,
	//四:该close时忘了close
	s := "世界 is world"
	for i, v := range s {
		fmt.Printf("index:%v value:%v\n", i, v)
		//索引的值并非是从0连续的,而是0-3-6-7-8...因为一个汉字需要三个字节来编码

	}
	//range遍历的是目标的副本,因特别注意当使用指针语义时
	f := []string{"leilei", "yangzhen", "haoyun"}
	for i := range f { //i为0-2
		f = f[:1]                   //"leilei","yangzhen"
		fmt.Printf("v[%s]\n", f[i]) //因为f被分割只有2个长度,而i是原切片的拷贝,始终为3

	}
	//一定要理解使用值语义和指针语义的区别,分别的代价是什么?将可能造成什么结果?

}
