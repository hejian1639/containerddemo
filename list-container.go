package main

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"log"
)

func main() {
	if err := listContainer(); err != nil {
		log.Fatal(err)
	}
}

func listContainer() error {
	// create a new client connected to the default socket path for containerd
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	// create a new context with an "demo" namespace
	ctx := namespaces.WithNamespace(context.Background(), "demo")

	containers, err := client.Containers(
		ctx,
	)
	if err != nil {
		return err
	}

	fmt.Printf("container number: %d\n", len(containers))

	for _, container := range containers {
		image, _ := container.Image(ctx)

		fmt.Printf("id: %s\timage: %s\n", container.ID(), image.Name())
	}

	return nil
}
