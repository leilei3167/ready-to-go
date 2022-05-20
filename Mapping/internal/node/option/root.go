package option

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "执行任务的节点,需指定kafka的服务器集群",
}

var (
	Brokers string //以,分隔
	Topic   string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&Brokers, "b", "", "指定brokers的地址端口,多个地址以,分隔")
	rootCmd.PersistentFlags().StringVar(&Topic, "t", "test_10", "指定topic")
}

func Execut() {
	log.Fatal(rootCmd.Execute())
}
