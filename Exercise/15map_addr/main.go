package main

type Student struct {
	name string
}

func main() {
	//map的value是不可寻址的,如要修改 需要将value修改为指针
	m := map[string]Student{"people": {"zhoujielun"}}
	m["people"].name = "wuyanzu"
}
