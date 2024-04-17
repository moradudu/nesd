package main

import (
	"github.com/nesd/client"
)

func main() {
	c := client.NewClient()
	c.ContainerList()
}
