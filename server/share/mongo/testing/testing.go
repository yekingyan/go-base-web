package mgotesting

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Clear is a func to clean.
type Clear func()

// RunDockerMongo return mongo url in docker.
func RunDockerMongo() (string, Clear) {
	ctx := context.Background()
	// 连接到docker
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	// 创建容器
	res, err := c.ContainerCreate(
		ctx,
		&container.Config{
			Image: "mongo",
			ExposedPorts: nat.PortSet{
				"27017/tcp": {},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"27017/tcp": []nat.PortBinding{{
					HostIP:   "127.0.0.1",
					HostPort: "0",
				}},
			},
		},
		nil, nil, "")
	if err != nil {
		panic(err)
	}
	containerID := res.ID
	// 启动容器
	c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	// 删除容器
	clear := func() {
		c.ContainerRemove(ctx, containerID,
			types.ContainerRemoveOptions{Force: true})
	}

	// 获取容器状态
	cInfo, _ := c.ContainerInspect(ctx, containerID)
	ports := cInfo.NetworkSettings.Ports["27017/tcp"][0]
	fmt.Println(ports.HostIP, ports.HostPort)
	mongoURI := fmt.Sprintf("mongodb://%s:%s", ports.HostIP, ports.HostPort)
	return mongoURI, clear
}
