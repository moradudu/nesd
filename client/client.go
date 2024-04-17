package client

import (
	"github.com/docker/docker/api/types"
)

type Containers struct {
	ContainerName string
	ContainerID   string
	Status        string
	Time          int64
}
type Images struct {
	ImageName string
	ImageID   string
	Status    string
	Time      int64
}

type Client interface {
	ContainerList() ([]types.Container, error)
	ImagesList() ([]Images, error)
}

func NewClient() Client {
	Enginer := "docker"
	switch Enginer {
	case "docker":
		return NewDockerClient()
	case "containerd":
	case "crio":
	default:
		return nil
	}
	return nil
}
