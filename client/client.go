package client

type Container struct {
	ContainerName []string
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

type event struct {
}

type Client interface {
	ContainerList() ([]*Container, error)
	ImagesList() ([]Images, error)
	Events() (chan event, error)
	Save() error
	InspectContainer(cid string)
	InspectImage(iid string)
	stop(cid string)
	stopAll(cid string)
}

func NewClient(engine string) Client {
	switch engine {
	case "docker":
		return NewDockerClient()
	case "containerd":
		return NewContrainerd()
	case "crio":
	default:
		return nil
	}
	return nil
}
