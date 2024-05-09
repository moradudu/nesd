package main

import (
	"github.com/spf13/cobra"
)

var engine = "docker"
var endpoint = "unix:///var/log/docker.sock"
var root = ""

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nerd",
	Short: "",
	Long:  ``,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&engine, "engine", "e", "", "容器引擎类型")
	rootCmd.PersistentFlags().StringVarP(&endpoint, "socket", "s", "", "容器引擎socket")
	rootCmd.MarkFlagRequired("engine")
}
