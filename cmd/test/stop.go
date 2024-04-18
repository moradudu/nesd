package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var StopCmd = &cobra.Command{
	Use:   "stopall",
	Short: "stop all container",
	Long:  `stop all container by its id`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer cli.Close()

		containers, err := cli.ContainerList(ctx, containertypes.ListOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Print(args)
		for _, container := range containers {
			fmt.Print("Stopping container ", container.ID[:10], "... ")
			noWaitTimeout := 0 // to not wait for the container to exit gracefully
			if err := cli.ContainerStop(ctx, container.ID, containertypes.StopOptions{Timeout: &noWaitTimeout}); err != nil {
				panic(err)
			}
			fmt.Println("Success")
		}
	},
}

func init() {
	rootCmd.AddCommand(StopCmd)
}
