package main

import "unicode"

/*
有⼀个机器⼈，给⼀串指令，L左转 R右转，F前进⼀步，B后退⼀步，问最后机器⼈的坐标，最开始，机器⼈位于
0 0，⽅向为正Y。
可以输⼊重复指令n ： ⽐如 R2(LF) 这个等于指令 RLFLF。
问最后机器⼈的坐标是多少？
*/

const (
	Left = iota
	Forward
	Right
	Back
)

func main() {
	println(move("R2(LF)", 0, 0, Forward))
}

//给出指令,以及起始的坐标,解析并执行后返回最终的坐标
func move(cmd string, x0, y0, z0 int) (x, y, z int) {
	x, y, z = x0, y0, z0
	repeat := 0
	repeatCmd := ""
	for _, s := range cmd {
		switch {
		case unicode.IsNumber(s):
			repeat = repeat*10 + (int(s) - '0')
		case s == ')':
			for i := 0; i < repeat; i++ {
				x, y, z = move(repeatCmd, x, y, z)
			}
			repeat = 0
			repeatCmd = ""
		case repeat > 0 && s != '(' && s != ')':
			repeatCmd = repeatCmd + string(s)
		case s == 'L':
			z = (z + 1) % 4
		case s == 'R':
			z = (z - 1 + 4) % 4
		case s == 'F':
			switch {
			case z == Left || z == Right:
				x = x - z + 1
			case z == Forward || z == Back:
				y = y - z + 2
			}
		case s == 'B':
			switch {
			case z == Left || z == Right:
				x = x + z - 1
			case z == Forward || z == Back:
				y = y + z - 2
			}
		}
	}
	return
}
