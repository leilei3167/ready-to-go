package slice

/*
这些加减乘除的过程,可以抽象为: 定义一个初始值,遍历切片,对每一个元素执行某种操作,返回结果

配合泛型和高阶函数,可以在Go中实现Reduce!

// Sum calculates the total from a slice of numbers.
func Sum(numbers []int) int {
	var sum int
	for _, number := range numbers {
		sum += number
	}
	return sum
}

// SumAllTails calculates the sums of all but the first number given a collection of slices.
func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tail := numbers[1:]
			sums = append(sums, Sum(tail))
		}
	}

	return sums
} */

//在拥有高阶函数的编程语言中,经常会提到 Fold(折叠) 这个概念
/*
折叠 可以被看做是用函数来取代数据结构的构成成分

在Reduce中，在集合上运行计算的行为被抽象出来了
必须要将一个默认值作为参数,避免使用默认的Go的默认值,并且可以根据不同的用途自行输入默认值,如加法默认值从0开始
乘法默认值从1开始
*/

//额外增加一种类型约束B,使得对于函数来说可以约束多种类型,Reduce实际是对函数行为的抽象
func Reduce[A, B any](collection []A, accumulator func(B, A) B, initialValue B) B {
	var result = initialValue
	for _, x := range collection {
		result = accumulator(result, x)
	}

	return result
}

//使用Reduce来改写
func Sum(num []int) int {
	add := func(acc, x int) int {
		return acc + x
	}
	return Reduce(num, add, 0)
}

func SumAllTail(num ...[]int) []int {
	sumTail := func(acc, x []int) []int {
		if len(x) == 0 {
			return append(acc, 0)
		} else {
			tail := x[1:]
			return append(acc, Sum(tail))
		}
	}
	return Reduce(num, sumTail, []int{})
}

//v1
/* func BalanceFor(transactions []Transaction, name string) float64 {
	var balance float64
	for _, t := range transactions {
		if t.From == name {
			balance -= t.Sum
		}
		if t.To == name {
			balance += t.Sum
		}
	}
	return balance
}
*/
//v2

type Transaction struct {
	From string
	To   string
	Sum  float64
}

func NewTransaction(from, to Account, sum float64) Transaction {
	return Transaction{From: from.Name, To: to.Name, Sum: sum}
}

type Account struct {
	Name    string
	Balance float64
}

func NewBalanceFor(account Account, transactions []Transaction) Account {
	return Reduce(
		transactions,
		applyTransaction,
		account,
	)
}

func applyTransaction(a Account, transaction Transaction) Account {
	if transaction.From == a.Name {
		a.Balance -= transaction.Sum
	}
	if transaction.To == a.Name {
		a.Balance += transaction.Sum
	}
	return a
}

func Find[A any](items []A, pre func(A) bool) (value A, found bool) {
	for _, item := range items {
		if pre(item) {
			return item, true
		}
	}
	return
}
