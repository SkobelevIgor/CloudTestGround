package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func main() {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	cDel := containersToDelete(ctx, cli, "test")
	deleteContainers(ctx, cli, cDel)
	// iDel := imagesToDelete(ctx, cli, "test")
	// deleteImages(iDel)

	// buildImage(ctx, cli, "test")
	// id := runImage(ctx, cli, "test")

	// fmt.Println(id)
}

func containersToDelete(ctx context.Context, cli *client.Client, imageTag string) []string {
	filter := filters.NewArgs()
	filter.Add("name", imageTag)
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		panic(err)
	}

	var cIds []string
	for _, ctr := range containers {
		cIds = append(cIds, ctr.ID)
	}

	return cIds
}

func imagesToDelete(ctx context.Context, cli *client.Client, imageTag string) []string {

	images, err := cli.ImageList(ctx, types.ImageListOptions{})

	if err != nil {
		panic(err)
	}

	var iIds []string

	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == imageTag+":latest" {
				iIds = append(iIds, img.ID[7:])
			}
		}
	}

	return iIds

}

func deleteContainers(ctx context.Context, cli *client.Client, cIds []string) (st bool) {
	status := true
	for _, cID := range cIds {
		err := cli.ContainerRemove(ctx, cID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			status = false
		}
	}
	return status
}

// func deleteImages(iIds []string) (st bool) {

// }

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
	}, &container.HostConfig{
		AutoRemove: true}, nil, "ctg_"+imageTag)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID
}

func removePreviousContainers(imageTag string) {

}

func removePreviousImages(imageTag string) {

}
