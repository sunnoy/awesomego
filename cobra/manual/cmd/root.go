/*
 *@Description
 *@author          lirui
 *@create          2020-10-03 21:50
 */
package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var dong string
var rootCmd = &cobra.Command{
	Use: "lr",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println(args)

		fmt.Println(dong)

	},
}

var verCmd = &cobra.Command{
	Use: "ha",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("这是个自命令输出")
	},
}

var sssCmd = &cobra.Command{
	Use: "sss",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("这是个sss输出")
	},
}

func main() {

	rootCmd.PersistentFlags().StringVarP(&dong, "d", "d", "dongdong", "sssss")
	rootCmd.AddCommand(verCmd)
	verCmd.AddCommand(sssCmd)
	// 这个函数调用一定要放在最后
	rootCmd.Execute()
}
