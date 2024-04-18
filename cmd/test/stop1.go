package main

import (
	"context"
	"fmt"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"os"
)

var Stop1Cmd = &cobra.Command{
	Use:   "stopone",
	Short: "stop a container",
	Long:  `stop a container by its id`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查命令行参数
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go stopone <container-name>")
			return
		}

		// 获取容器名称
		containerName := os.Args[2]

		// 创建Docker客户端
		//ctx := context.Background()
		//cli, err := client.NewClientWithOpts(client.FromEnv)
		//if err != nil {
		//	fmt.Println("Error creating Docker client:", err)
		//	return
		//}
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer cli.Close()

		// 获取所有容器的列表
		containers, err := cli.ContainerList(ctx, containertypes.ListOptions{})
		if err != nil {
			fmt.Println("Error listing containers:", err)
			return
		}

		// 遍历容器列表，找到匹配的容器并关闭
		for _, container := range containers {
			if container.Names[0] == "/"+containerName {
				fmt.Printf("Container %s found, stopping...\n", containerName)
				//timeout := int64(10)
				//options := types.ContainerStopOptions{
				//	Timeout: &timeout,
				//}
				noWaitTimeout := 0
				if err := cli.ContainerStop(ctx, container.ID, containertypes.StopOptions{Timeout: &noWaitTimeout}); err != nil {
					fmt.Printf("Error stopping container %s: %s\n", containerName, err)
				} else {
					fmt.Printf("Container %s stopped successfully.\n", containerName)
				}
				return
			}
		}

		fmt.Printf("Container %s not found.\n", containerName)
	}}

func init() {
	rootCmd.AddCommand(Stop1Cmd)
}
