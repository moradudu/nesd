package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

type Event struct {
	ID        string    // 容器ID
	Action    Action    // 动作
	Timestamp time.Time // 时间戳
}

// Action represents a Docker event action
type Action string

//const (
//	ActionStart = Action("start")
//	ActionStop  = Action("stop")
//	ActionKill  = Action("kill")
//)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Start container monitoring",
	Run: func(cmd *cobra.Command, args []string) {
		startMonitoring(cmd)
	},
}

func startMonitoring(cmd *cobra.Command) {
	logFile, _ := cmd.Flags().GetString("log")

	cli, err := client.NewClientWithOpts()
	if err != nil {
		log.Fatal(err)
	}

	events := make(map[string][]Event)

	ctx := context.Background()

	eventChan, errChan := cli.Events(ctx, types.EventsOptions{})

	fmt.Println("容器监控程序已启动...")

	logFilePtr, err := os.Create(logFile)
	if err != nil {
		log.Fatal(err)
	}
	defer logFilePtr.Close()

	logger := log.New(logFilePtr, "", log.LstdFlags)

	for {
		select {
		case event := <-eventChan:
			e := Event{
				ID:        event.Actor.ID,
				Action:    Action(event.Action),
				Timestamp: time.Now(),
			}

			events[event.Actor.ID] = append(events[event.Actor.ID], e)

			logger.Printf("容器ID：%s，动作：%s，时间：%s\n", e.ID, e.Action, e.Timestamp.String())

			fmt.Printf("接收到事件：%s\n", e.Action)

		case err := <-errChan:
			fmt.Printf("发生错误：%v\n", err)
		}
	}
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
