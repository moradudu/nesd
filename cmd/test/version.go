package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of nesd",
	Long:  `All software has versions. This is nesd's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nesd version 0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(VersionCmd)
}
