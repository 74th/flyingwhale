package lib

import (
	"bytes"
	"fmt"
	"time"

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

// ExecShortCommand executes docker exec with stdout
func (d *DockerClient) ExecShortCommand(commands []string, allowErrExit bool) (stdout string, stderr string) {

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
	if inspect.ExitCode != 0 {
		fmt.Print("\x1b[31m", stderrStream.String(), "\x1b[0m")
		panic("something wrong")
	}

	return stdoutStream.String(), stderrStream.String()
}

// ExecWithShowingStdout executes docker exec with stdout
func (d *DockerClient) ExecWithShowingStdout(commands []string, allowErrExit bool) {

	exec, err := d.client.CreateExec(docker.CreateExecOptions{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          commands,
		Container:    d.container.ID})
	if err != nil {
		panic(err)
	}

	isRunGuard := true
	var stdoutStream bytes.Buffer
	var stderrStream bytes.Buffer

	go func() {
		// goroutineでexecを実行
		err = d.client.StartExec(exec.ID, docker.StartExecOptions{
			Detach:       false,
			OutputStream: &stdoutStream,
			ErrorStream:  &stderrStream})
		if err != nil {
			panic(err)
		}
		// 完了のフラグ
		isRunGuard = false
	}()

	for true {
		// 終了するまで表示する
		if stdoutStream.Len() > 0 {
			fmt.Print(stdoutStream.String())
			stdoutStream.Reset()
		}
		if stderrStream.Len() > 0 {
			fmt.Print("\x1b[31m", stderrStream.String(), "\x1b[0m")
			stderrStream.Reset()
		}
		if err != nil {
			panic(err)
		}
		time.Sleep(100 * time.Millisecond)
		if !isRunGuard {
			break
		}
	}
	// 最後の出力
	if stdoutStream.Len() > 0 {
		fmt.Print(stdoutStream.String())
		stdoutStream.Reset()
	}
	if stderrStream.Len() > 0 {
		fmt.Print("\x1b[31m", stderrStream.String(), "\x1b[0m")
		stderrStream.Reset()
	}

	inspect, err := d.client.InspectExec(exec.ID)
	if inspect.ExitCode != 0 {
		if !allowErrExit {
			panic("an error ocuured")
		}
	}
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
