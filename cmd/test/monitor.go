package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

// Event 结构体定义
type Event struct {
	ID         string            `json:"id"`
	Status     string            `json:"status"`
	Time       time.Time         `json:"time"`
	From       string            `json:"from"`
	Action     string            `json:"action"`
	Attributes map[string]string `json:"attributes"`
}

// eventsCmd represents the events command
var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Monitor Docker events and store logs",
	Long:  `Monitor Docker events and store logs to specified directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		logDir, err := cmd.Flags().GetString("log-dir")
		if err != nil {
			log.Fatalf("Error getting log directory: %v", err)
		}

		if logDir == "" {
			log.Fatalf("Please provide a log directory with the --log-dir flag.")
		}

		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatalf("Error creating log directory: %v", err)
		}

		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Fatalf("Error creating Docker client: %v", err)
		}

		events, errChan := cli.Events(context.Background(), types.EventsOptions{})
		if err != nil {
			log.Fatalf("Error getting Docker events: %v", err)
		}

		for {
			select {
			case event, ok := <-events:
				if !ok {
					// 通道已关闭，没有更多的事件
					return
				}
				eventJSON, err := json.Marshal(event)
				if err != nil {
					log.Printf("Error marshaling event: %v", err)
					continue
				}

				//eventID := event.ID
				//if eventID == "" {
				//	eventID = "unknown"
				//}
				fmt.Printf(string(eventJSON) + "\n")
				//logFilePath := filepath.Join(logDir, fmt.Sprintf("%s.log", eventID))
				//logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				//if err != nil {
				//	log.Printf("Error opening log file: %v", err)
				//	continue
				//}
				//defer logFile.Close()

				//if _, err := logFile.WriteString(string(eventJSON) + "\n"); err != nil {
				//	log.Printf("Error writing to log file: %v", err)
				//}
			case err := <-errChan:
				// 从错误通道读取到错误
				log.Fatalf("Error reading Docker events: %v", err)
			}
		}
	},
}

func init() {
	eventsCmd.Flags().StringP("log-dir", "l", "", "Path to the directory where logs will be stored")
	rootCmd.AddCommand(eventsCmd)
}

//func main() {
//	if err := rootCmd.Execute(); err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//}
