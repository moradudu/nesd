package client

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type DockerClient struct {
}

func (c *DockerClient) ContainerList() ([]*Container, error) {
	fmt.Print("this is DockerClient ContainerList ")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	rawcontainers, err := cli.ContainerList(ctx, container.ListOptions{})

	containers := make([]*Container, 0)
	for _, c := range rawcontainers {
		con := &Container{
			ContainerName: c.Names,
			ContainerID:   c.ID,
			Status:        c.Status,
			Time:          c.Created,
		}
		containers = append(containers, con)
	}

	return containers, err
}

func (c *DockerClient) ImagesList() ([]Images, error) {
	fmt.Print("this is DockerClient ImagesList ")
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Println(image.ID, "|", image.Labels, "|", image.Size)
	}

	return nil, nil
}

var _ Client = &DockerClient{}

func NewDockerClient() Client {
	return &DockerClient{}
}
