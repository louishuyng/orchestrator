package task

import "github.com/docker/docker/client"

type Docker struct {
	Client *client.Client
	Config Config
}

type DockerResult struct {
	Error       error
	Action      string
	ContainerID string
	Result      string
}
