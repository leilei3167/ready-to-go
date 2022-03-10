package main

import "fmt"

func hasCycle(head *ListNode) bool {
	//创建一个map
	key := make(map[*ListNode]int)
	//一直循环到nil
	for head != nil {
		_, ok := key[head]
		if ok { //知识点,判断map中是否存在某个key,如果存在ok会是true,并且_会被对应value赋值
			return true
		} else { //ok为false说明map没有对应的key
			key[head] = 1 //每遍历完一个节点可以向里面放一个值
			head = head.Next
		}
	}

	return false

}

type ListNode struct {
	Val  int
	Next *ListNode
}

func main() {

	node := new(ListNode)
	node.Val = 1

	node1 := new(ListNode)
	node1.Val = 2
	node1.Next = node

	node2 := new(ListNode)
	node2.Val = 3
	node2.Next = node1

	node3 := new(ListNode)
	node3.Val = 4
	node3.Next = node2

	node4 := new(ListNode)
	node4.Val = 5
	node4.Next = node4

	fmt.Println(hasCycle(node))

}
