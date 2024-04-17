package main

import (
	"encoding/json"
	"fmt"
	"github.com/nesd/client"
	"github.com/spf13/cobra"
)

var ImageListCmd = &cobra.Command{
	Use:   "pa",
	Short: "list image from unix socket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("homohomo")
		imageClinet := client.NewClient()
		list, error := imageClinet.ImagesList()
		if error != nil {
			fmt.Println(error)
		}
		buf, _ := json.Marshal(list)
		fmt.Print(string(buf))

	},
}

func init() {
	rootCmd.AddCommand(ImageListCmd)
}
