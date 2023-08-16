package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/diezfx/split-app-backend/pkg/logger"
)

func main() {
	ctx := context.Background()
	logger.SetConsolLogger()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	client.Pipeline("test")
	defer client.Close()
	// initialize Dagger client
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		logger.Fatal(ctx, errors.New("not enough parameters")).Msg("check arguments")
	}
	command := os.Args[1]

	if command == "lint" {
		result, err := lint(ctx, client)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)

	}

	contextDir := client.Host().Directory(".")
	//_ := client.SetSecret("ghApiToken", ghToken)

	ref, err := contextDir.
		DockerBuild(dagger.DirectoryDockerBuildOpts{
			Dockerfile: "deployment/Dockerfile",
			BuildArgs:  []dagger.BuildArg{{Name: "APP_NAME", Value: "split-app-backend"}}}).
		Publish(ctx, "ghcr.io/diezfx/split-app-backend:latest")
	if err != nil {
		panic(err)
	}
	fmt.Printf("publish image %s", ref)
}

func lint(ctx context.Context, client *dagger.Client) (string, error) {
	c := client.Pipeline("lint")

	contextDir := client.Host().Directory(".")
	out, err := c.Container().
		From("golangci/golangci-lint:v1.54").
		WithMountedDirectory("/app", contextDir).
		WithWorkdir("/app").
		WithExec([]string{"/usr/bin/golangci-lint", "run", "-v", "--timeout", "5m"}).
		Stderr(ctx)
	if err != nil {
		return "", err
	}
	return out, nil
}
