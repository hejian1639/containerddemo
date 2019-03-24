package main

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"log"
	"syscall"
)

func main() {
	if err := stopContainer(); err != nil {
		log.Fatal(err)
	}
}

func stopContainer() error {
	// create a new client connected to the default socket path for containerd
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	// create a new context with an "demo" namespace
	ctx := namespaces.WithNamespace(context.Background(), "demo")


	// create a container
	container, err := client.LoadContainer(
		ctx,
		"redis-server")
	if err != nil {
		return err
	}

	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	// create a task from the container
	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}
	defer task.Delete(ctx)
	//return nil

	// make sure we wait before calling start
	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// kill the process and get the exit status
	if err := task.Kill(ctx, syscall.SIGTERM); err != nil {
		return err
	}

	// wait for the process to fully exit and print out the exit status

	status := <-exitStatusC
	code, _, err := status.Result()
	if err != nil {
		return err
	}
	fmt.Printf("redis-server exited with status: %d\n", code)

	return nil
}
