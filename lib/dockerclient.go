package lib

import (
	"bytes"
	"fmt"

	"github.com/fsouza/go-dockerclient"
)

// DockerClient is docker client
type DockerClient struct {
	client        *docker.Client
	auth          docker.AuthConfiguration
	container     *docker.Container
	exec          *docker.Exec
	ContainerName string
}

// CreateDockerClient gets docker client
func CreateDockerClient() *DockerClient {
	d := DockerClient{}
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	d.client = client
	d.auth = docker.AuthConfiguration{}

	return &d
}

// HasImage return name of image is available
func (d *DockerClient) HasImage(name string) bool {
	// TOOD filter
	//filter := map[string][]string{"label": {"RepoTag=debian"}}
	//arg := docker.ListImagesOptions{All: true, Filters: filter}
	arg := docker.ListImagesOptions{All: true}
	list, _ := d.client.ListImages(arg)
	for _, l := range list {
		for _, l2 := range l.RepoTags {
			if l2 == name {
				return true
			}
		}
	}
	return false
}

// PullImage pull image
func (d *DockerClient) PullImage(name string) {
	opts := docker.PullImageOptions{
		Repository: name,
		Tag:        "latest"}
	err := d.client.PullImage(opts, d.auth)
	if err != nil {
		panic(err)
	}
}

// CreateContainer makes a container
func (d *DockerClient) CreateContainer(containerName string, imageName string, entrypoint []string, cmd []string) {
	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        imageName,
			Entrypoint:   []string{"sh"},
			AttachStdin:  true,
			OpenStdin:    true,
			StdinOnce:    true,
			AttachStdout: true,
			AttachStderr: true},
		Name: containerName}
	container, err := d.client.CreateContainer(opts)
	if err != nil {
		panic(err)
	}
	d.container = container
	d.client.StartContainer(
		container.ID,
		&docker.HostConfig{})
}

// StopContainer stops a container
func (d *DockerClient) StopContainer() {
	err := d.client.StopContainer(
		d.container.ID,
		1)
	if err != nil {
		panic(err)
	}
}

// Exec is run
func (d *DockerClient) Exec(commands []string) (stdout string, stderr string) {

	exec, err := d.client.CreateExec(docker.CreateExecOptions{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          commands,
		Container:    d.container.ID})
	if err != nil {
		panic(err)
	}

	var stdoutStream bytes.Buffer
	var stderrStream bytes.Buffer
	err = d.client.StartExec(exec.ID, docker.StartExecOptions{
		Detach:       false,
		OutputStream: &stdoutStream,
		ErrorStream:  &stderrStream})
	if err != nil {
		panic(err)
	}
	inspect, err := d.client.InspectExec(exec.ID)
	if err != nil {
		panic(err)
	}
	if inspect.ExitCode != 0 {
		fmt.Println(stdoutStream.String())
		fmt.Println(stderrStream.String())
		panic("an error ocuured")
	}

	return stdoutStream.String(), stderrStream.String()
}

// RemoveContainer remove a container
func (d *DockerClient) RemoveContainer() {
	err := d.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:    d.container.ID,
		Force: true})
	if err != nil {
		panic(err)
	}
}

// CommitContainer commit container
func (d *DockerClient) CommitContainer(imageName string) {
	_, err := d.client.CommitContainer(docker.CommitContainerOptions{
		Container:  d.container.ID,
		Repository: imageName,
		Tag:        "latest",
		Message:    "created by flying whale"})
	if err != nil {
		panic(err)
	}
}
