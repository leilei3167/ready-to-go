package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

type ErrorHandling int

const (
	ContinueOnParseError  ErrorHandling = 1 // 解析错误尝试继续处理
	ExitOnParseError      ErrorHandling = 2 // 解析错误程序停止
	PanicOnParseError     ErrorHandling = 3 // 解析错误 panic
	ReturnOnDividedByZero ErrorHandling = 4 // 除0返回
	PanicOnDividedByZero  ErrorHandling = 5 // 除0 painc
)

type OpType int

const (
	ADD      OpType = 1
	MINUS    OpType = 2
	MULTIPLY OpType = 3
	DIVIDE   OpType = 4
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "一个简单的计算器",
}

//设置flag选项
var (
	parseHandling int
)

func init() {
	rootCmd.PersistentFlags().IntVar(&parseHandling, "parse_err", int(ContinueOnParseError),
		"do what when parse arg error")
}
func Execute() {
	log.Fatal(rootCmd.Execute())
}
