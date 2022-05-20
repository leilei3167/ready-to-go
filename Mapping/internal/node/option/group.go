package option

import "github.com/spf13/cobra"

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "消费者组模式运行,需指定Group的名称",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var Group string

func init() {
	rootCmd.AddCommand(groupCmd)
	groupCmd.PersistentFlags().StringVar(&Group, "g", "", "加入或创建一个消费者组")
}
