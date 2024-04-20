package main

import (
	"context"
	"fmt"
	"log"

	//"github.com/containerd/nerdctl/pkg/api/types"

	//"github.com/docker/cli/cli/command/container"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var ContainerRunCmd = &cobra.Command{
	Use:   "containerrun",
	Short: "list image from unix socket",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取镜像名称和命令作为命令行参数
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run main.go <image-name> <command>")
			return
		}
		imageName := os.Args[2]
		command := os.Args[3:]

		// 创建 Docker 客户端
		///shishikan
		//1
		//2
		//3

		fmt.Printf("")
		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			fmt.Println("Error creating Docker client:", err)
			return
		}
		defer cli.Close()

		// 创建容器
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: imageName,
			Cmd:   command,
			Tty:   true,
		}, &container.HostConfig{}, &network.NetworkingConfig{}, nil, "")
		if err != nil {
			fmt.Println("Error creating container:", err)
			return
		}

		// 启动容器
		if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
			fmt.Println("Error starting container:", err)
			return
		}

		// 等待容器退出

		statusCh, errCh := cli.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				log.Fatal(err)
			}
		case <-statusCh:
		}

		// 获取容器日志输出
		logs, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
		if err != nil {
			fmt.Println("Error retrieving container logs:", err)
			return
		}
		defer logs.Close()

		// 将容器日志输出打印到标准输出
		_, err = io.Copy(os.Stdout, logs)
		if err != nil {
			fmt.Println("Error printing container logs:", err)
			return
		}

		// 删除容器
		if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
			fmt.Println("Error removing container:", err)
			return
		}

		fmt.Println("Container execution completed.")
	},
}

func init() {
	rootCmd.AddCommand(ContainerRunCmd)
}
