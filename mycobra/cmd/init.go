package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "add",
	Short: "short init",
	Long:  "Long init",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run init cmd begin")
		fmt.Println(
			cmd.Flags().Lookup("viper").Value,
			cmd.Flags().Lookup("author").Value,
			cmd.Flags().Lookup("config").Value,
			viper.GetString("author"),
			cmd.Flags().Lookup("license").Value,
			cmd.Parent().Flags().Lookup("source").Value,
		)
		fmt.Println("run init cmd end")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
