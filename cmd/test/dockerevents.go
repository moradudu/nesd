package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

type DockerEvent struct {
	Type   string `json:"Type"`
	Action string `json:"Action"`
	Actor  struct {
		ID         string `json:"ID"`
		Attributes struct {
			Name string `json:"name"`
		} `json:"Attributes"`
	} `json:"Actor"`
}

var eventCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Listen to Docker events",
	Run:   listenToDockerEvents,
}

func init() {

	rootCmd.AddCommand(eventCmd)

}

func listenToDockerEvents(cmd *cobra.Command, args []string) {
	socketPath := "/var/run/docker.sock" // Docker socket路径

	// 连接Docker socket
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法连接到Docker socket：%v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 发送获取Docker事件的请求
	_, err = conn.Write([]byte("GET /events HTTP/1.1\r\nHost: docker\r\n\r\n"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法发送请求：%v\n", err)
		os.Exit(1)
	}

	// 读取并解析事件数据
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "无法读取数据：%v\n", err)
			os.Exit(1)
		}

		var event DockerEvent
		err = json.Unmarshal(line, &event)
		if err != nil {
			fmt.Fprintf(os.Stderr, "无法解析JSON数据：%v\n", err)
			continue
		}

		// 处理事件
		fmt.Printf("收到事件：类型=%s，动作=%s，容器ID=%s，容器名称=%s\n", event.Type, event.Action, event.Actor.ID, event.Actor.Attributes.Name)
	}
}
