package container

import (
	"context"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

// CreateNetwork creates a Docker Network with a namespace
func (c *Container) CreateNetwork(namespace []string) (networkID string, err error) {
	network, err := c.FindNetwork(namespace)
	if docker.IsErrNotFound(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if network.ID != "" {
		networkID = network.ID
		return
	}
	namespaceFlat := Namespace(namespace)
	response, err := c.client.NetworkCreate(context.Background(), namespaceFlat, types.NetworkCreate{
		CheckDuplicate: true, // Cannot have 2 network with the same name
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": namespaceFlat,
		},
	})
	if err != nil {
		return
	}
	networkID = response.ID
	return
}

// DeleteNetwork deletes a Docker Network associated with a namespace
func (c *Container) DeleteNetwork(namespace []string) (err error) {
	network, err := c.FindNetwork(namespace)
	if docker.IsErrNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return c.client.NetworkRemove(context.Background(), network.ID)

}

// FindNetwork finds a Docker Network by a namespace. If no network if found, an error is returned.
func (c *Container) FindNetwork(namespace []string) (network types.NetworkResource, err error) {
	return c.client.NetworkInspect(context.Background(), Namespace(namespace), types.NetworkInspectOptions{})
}
