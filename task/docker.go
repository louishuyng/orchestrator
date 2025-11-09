package task

import (
	"context"
	"io"
	"log"
	"math"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

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

func NewDocker(c *Config) *Docker {
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	return &Docker{
		Client: dc,
		Config: *c,
	}
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()

	reader, error := d.Client.ImagePull(ctx, d.Config.Image, types.ImagePullOptions{})

	if error != nil {
		log.Printf("Error pulling image %s: %v", d.Config.Image, error)
		return DockerResult{Error: error}
	}

	io.Copy(os.Stdout, reader)
	rp := container.RestartPolicy{Name: d.Config.RestartPolicy}

	r := container.Resources{
		Memory:   d.Config.Memory,
		NanoCPUs: int64(d.Config.Cpu * math.Pow10(9)),
	}

	cc := container.Config{
		Image:        d.Config.Image,
		Tty:          false,
		Env:          d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}

	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container from image %s: %v", d.Config.Image, err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Printf("Error starting container %s: %v", resp.ID, err)
		return DockerResult{Error: err}
	}

	out, err := d.Client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Printf("Error getting logs for container %s: %v", resp.ID, err)
		return DockerResult{Error: err}
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{
		Action:      "start",
		ContainerID: resp.ID,
		Result:      "success",
	}
}

func (d *Docker) Stop(containerID string) DockerResult {
	log.Printf("Stopping container %s", containerID)
	ctx := context.Background()

	err := d.Client.ContainerStop(ctx, containerID, nil)
	if err != nil {
		log.Printf("Error stopping container %s: %v", containerID, err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	})
	if err != nil {
		log.Printf("Error removing container %s: %v", containerID, err)
		return DockerResult{Error: err}
	}

	return DockerResult{
		Action:      "stop",
		ContainerID: containerID,
		Result:      "success",
	}
}

func (d *Docker) Inspect(containerID string) DockerInspectResponse {
	ctx := context.Background()
	resp, err := d.Client.ContainerInspect(ctx, containerID)
	if err != nil {
		log.Printf("Error inspecting container: %s\n", err)
		return DockerInspectResponse{Error: err}
	}

	return DockerInspectResponse{Container: &resp}
}
