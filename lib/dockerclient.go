package lib

import (
	"bytes"
	"fmt"
	"os"

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
	endpoint := "tcp://192.168.99.100:2376"
	path := os.Getenv("DOCKER_CERT_PATH")
	ca := fmt.Sprintf("%s/ca.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)
	key := fmt.Sprintf("%s/key.pem", path)
	client, err := docker.NewTLSClient(endpoint, cert, key, ca)
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
	config := docker.Config{
		Image:        imageName,
		Entrypoint:   []string{"ls"},
		Cmd:          []string{"-1", "/usr/local/bin/"},
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true}
	opts := docker.CreateContainerOptions{
		Config: &config,
		Name:   containerName}
	container, err := d.client.CreateContainer(opts)
	if err != nil {
		panic(err)
	}
	d.container = container
	hostConfig := docker.HostConfig{}
	d.client.StartContainer(
		container.ID,
		&hostConfig)

	// var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	d.client.Logs(docker.LogsOptions{
		Container:    container.ID,
		Stdout:       true,
		Stderr:       true,
		Follow:       true,
		Tail:         "all",
		OutputStream: &stdout,
		ErrorStream:  &stderr})
	d.client.RemoveContainer(docker.RemoveContainerOptions{
		ID: container.ID})
	fmt.Println(stdout.String())
	fmt.Println(stderr.String())
	fmt.Println("---")

	config = docker.Config{
		Image:        imageName,
		Entrypoint:   []string{"sh"},
		OpenStdin:    true,
		StdinOnce:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true}
	container, err = d.client.CreateContainer(docker.CreateContainerOptions{
		Config: &config,
		Name:   containerName})
	if err != nil {
		panic(err)
	}

	d.client.StartContainer(
		container.ID,
		&hostConfig)

	exec, err := d.client.CreateExec(docker.CreateExecOptions{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"npm", "update"},
		Container:    container.ID})
	if err != nil {
		panic(err)
	}
	stdout.Reset()
	stderr.Reset()
	d.client.StartExec(exec.ID, docker.StartExecOptions{
		Detach:       false,
		OutputStream: &stdout,
		ErrorStream:  &stderr})
	fmt.Println(stdout.String())
	fmt.Println(stderr.String())
	fmt.Println("---")

	exec, err = d.client.CreateExec(docker.CreateExecOptions{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"npm", "install", "-yg", "mared"},
		Container:    container.ID})
	if err != nil {
		panic(err)
	}
	stdout.Reset()
	stderr.Reset()
	d.client.StartExec(exec.ID, docker.StartExecOptions{
		Detach:       false,
		OutputStream: &stdout,
		ErrorStream:  &stderr})
	fmt.Println(stdout.String())
	fmt.Println("---")
	fmt.Println(stderr.String())
	fmt.Println("---")

	err = d.client.StopContainer(
		container.ID,
		1)
	if err != nil {
		panic(err)
	}

}

// Exec is run
func (d *DockerClient) Exec(command []string) {
	//d.client.StartExec(id string, opts docker.StartExecOptions)
	// opts := docker.CreateExecOptions{
	// 	AttachStdin:  false,
	// 	AttachStdout: true,
	// 	AttachStderr: true,
	// 	Cmd:          command,
	// 	Container:    d.container.ID}
	// exec, err := d.client.CreateExec(opts)
	// if err != nil {
	// 	panic(err)
	// }
	// d.exec = exec
}
