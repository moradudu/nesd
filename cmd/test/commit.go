package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var CommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit a container to create an image from its contents:",
	Long:  `Commit a container to create an image from its contents:`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer cli.Close()

		createResp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: "c9fe3bce8a6d",
			Cmd:   []string{"touch", "/helloworld"},
		}, nil, nil, nil, "")
		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(ctx, createResp.ID, container.StartOptions{}); err != nil {
			panic(err)
		}

		statusCh, errCh := cli.ContainerWait(ctx, createResp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}

		commitResp, err := cli.ContainerCommit(ctx, createResp.ID, container.CommitOptions{Reference: "helloworld"})
		if err != nil {
			panic(err)
		}

		fmt.Println(commitResp.ID)
	},
}

func init() {
	rootCmd.AddCommand(CommitCmd)
}
