package tasks

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func RemoveContainer(client client.APIClient, containerName string) (err error) {
	err = client.ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{
		RemoveVolumes: false,
		RemoveLinks:   false,
		Force:         true,
	})
	if err != nil {
		logrus.Errorf("%v", err)
	}
	return
}

func ImageBuild(client client.APIClient, imageName string) {
	// create temporary dir which will contain all files needed by image
	tmpDirPath, err := ioutil.TempDir("", "docker-build-"+imageName)
	// remove dir at the end of building
	defer os.RemoveAll(tmpDirPath)

	DockerfileName := "Dockerfile"
	DockerfileBody := "FROM openjdk:8-buster\n" // FIXME different images for GS
	DockerfilePath := filepath.Join(tmpDirPath, DockerfileName)
	err = ioutil.WriteFile(DockerfilePath, []byte(DockerfileBody), 0666)
	//logrus.Info("dockerfile: " + DockerfilePath)
	//logrus.Infof("tmpdir: %v", tmpDirPath)
	ctx, _ := archive.TarWithOptions(tmpDirPath, &archive.TarOptions{})
	res, err := client.ImageBuild(context.Background(), ctx, types.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: DockerfileName, // FIXME
	})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, res.Body)
}

func StartWrapperInDocker(gsId string) {

	containerName := gsId
	imageName := gsId

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	//containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	//if err != nil {
	//	panic(err)
	//}
	//for _, container := range containers {
	//	fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	//}

	//reader, err := cli.ImagePull(context.TODO(), "alpine", types.ImagePullOptions{})
	//if err != nil {
	//	panic(err)
	//}
	//io.Copy(os.Stdout, reader)

	ImageBuild(cli, containerName)
	RemoveContainer(cli, containerName)

	var mounts []mount.Mount
	mounts = append(mounts, mount.Mount{
		Type:          "bind",
		Source:        os.Args[0],
		Target:        "/opt/lvlup/exp/slave", //FIXME
		ReadOnly:      true,
		Consistency:   "",
		BindOptions:   nil,
		VolumeOptions: nil,
		TmpfsOptions:  nil,
	})
	mounts = append(mounts, mount.Mount{
		Type:   "bind",
		Source: config.REDIS_HOST,
		Target: "/opt/lvlup/exp/redis.sock", //FIXME
		//ReadOnly:      true,
		ReadOnly:      false,
		Consistency:   "",
		BindOptions:   nil,
		VolumeOptions: nil,
		TmpfsOptions:  nil,
	})
	mounts = append(mounts, mount.Mount{
		Type:   "bind",
		Source: config.GS_DIR + "/" + gsId + "/",
		Target: "/opt/lvlup/exp/gs/" + gsId + "/", // FIXME
		//ReadOnly:      true,
		ReadOnly:      false,
		Consistency:   "",
		BindOptions:   nil,
		VolumeOptions: nil,
		TmpfsOptions:  nil,
	})
	containerCreateRes, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Hostname:     "",
			Domainname:   "",
			User:         "",
			AttachStdin:  false,
			AttachStdout: false,
			AttachStderr: false,
			ExposedPorts: nil,
			Tty:          false,
			OpenStdin:    false,
			StdinOnce:    false,
			Env: []string{
				"REDIS_HOST=/opt/lvlup/exp/redis.sock", //FIXME needs variable
				"GS_DIR=/opt/lvlup/exp/gs",             //FIXME needs variable
			},
			Healthcheck: nil,
			ArgsEscaped: false,
			Image:       imageName,
			//Volumes: map[string]struct{}{
			//	"/tmp/:/zzzzz/": {},
			//},
			WorkingDir:      "",
			Entrypoint:      strslice.StrSlice{"/opt/lvlup/exp/slave"}, //FIXME
			Cmd:             strslice.StrSlice{"wrapper", gsId},        //FIXME
			NetworkDisabled: false,
			MacAddress:      "",
			OnBuild:         nil,
			Labels:          nil,
			StopSignal:      "",
			StopTimeout:     nil,
			Shell:           nil,
		},
		&container.HostConfig{
			Binds:           nil,
			ContainerIDFile: "",
			LogConfig:       container.LogConfig{},
			NetworkMode:     "",
			PortBindings:    nil,
			RestartPolicy:   container.RestartPolicy{},
			AutoRemove:      false,
			VolumeDriver:    "",
			VolumesFrom:     nil,
			CapAdd:          nil,
			CapDrop:         nil,
			Capabilities:    nil,
			DNS:             nil,
			DNSOptions:      nil,
			DNSSearch:       nil,
			ExtraHosts:      nil,
			GroupAdd:        nil,
			IpcMode:         "",
			Cgroup:          "",
			Links:           nil,
			OomScoreAdj:     0,
			PidMode:         "",
			Privileged:      false,
			PublishAllPorts: false,
			ReadonlyRootfs:  false,
			SecurityOpt:     nil,
			StorageOpt:      nil,
			Tmpfs:           nil,
			UTSMode:         "",
			UsernsMode:      "",
			ShmSize:         0,
			Sysctls:         nil,
			Runtime:         "",
			ConsoleSize:     [2]uint{},
			Isolation:       "",
			Resources:       container.Resources{},
			Mounts:          mounts,
			MaskedPaths:     nil,
			ReadonlyPaths:   nil,
			Init:            nil,
		},
		&network.NetworkingConfig{},
		containerName,
	)
	if err != nil {
		logrus.Errorf("%v", err)
	}
	logrus.Infof("%v", containerCreateRes)

	err = cli.ContainerStart(context.Background(), containerName, types.ContainerStartOptions{})
	if err != nil {
		logrus.Errorf("%v", err)
	}
}
