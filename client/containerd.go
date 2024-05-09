package client

import (
	"errors"
	"fmt"
)

type Containerd struct {
}

func (c Containerd) ContainerList() ([]*Container, error) {
	fmt.Print("this is Containerd ContainerList ")
	return nil, errors.New("not support")
}

func (c Containerd) ImagesList() ([]Images, error) {
	fmt.Print("this is Containerd ImagesList ")
	con := Container{
		ContainerName: []string{container.name},
		ContainerID:   "",
		Status:        "",
		Time:          0,
	}
	return nil, errors.New("not support")
}

func NewContrainerd() Client {
	return &Containerd{}
}
