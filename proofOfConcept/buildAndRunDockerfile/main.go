package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func main() {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	buildImage(ctx, cli, "test")
	id := runImage(ctx, cli, "test")

	fmt.Println(id)
}

func getContext(filePath string) io.Reader {
	ctx, _ := archive.TarWithOptions(filePath, &archive.TarOptions{})
	return ctx
}

func buildImage(ctx context.Context, cli *client.Client, imageTag string) {
	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		getContext("."),
		types.ImageBuildOptions{
			Tags:       []string{imageTag},
			Dockerfile: "Dockerfile",
			Remove:     true})
	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
}

func runImage(ctx context.Context, cli *client.Client, imageTag string) string {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageTag,
		Tty:   false,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID
}
