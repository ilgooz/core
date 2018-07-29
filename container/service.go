package container

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
)

// ListServices returns existing docker services matching a specific label name
func (c *Container) ListServices(label string) (services []swarm.Service, err error) {
	return c.client.ServiceList(context.Background(), types.ServiceListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: label,
		}),
	})
}

// FindService returns the Docker Service. Return error if not found.
func (c *Container) FindService(namespace []string) (service swarm.Service, err error) {
	service, _, err = c.client.ServiceInspectWithRaw(
		context.Background(),
		Namespace(namespace),
		types.ServiceInspectOptions{},
	)
	return service, err
}

// StartService starts a docker service
func (c *Container) StartService(options ServiceOptions) (serviceID string, err error) {
	service := options.toSwarmServiceSpec()
	response, err := c.client.ServiceCreate(context.Background(), service, types.ServiceCreateOptions{})
	if err != nil {
		return "", err
	}
	return response.ID, c.waitForStatus(options.Namespace, RUNNING)
}

// StopService stops a docker service
func (c *Container) StopService(namespace []string) (err error) {
	status, err := c.ServiceStatus(namespace)
	if err != nil || status == STOPPED {
		return
	}
	if err := c.client.ServiceRemove(context.Background(), Namespace(namespace)); err != nil {
		return err
	}
	return c.waitForStatus(namespace, STOPPED)
}

// ServiceStatus return the status of the Docker Swarm Servicer
func (c *Container) ServiceStatus(namespace []string) (status StatusType, err error) {
	_, err = c.FindService(namespace)
	if docker.IsErrNotFound(err) {
		return STOPPED, nil
	}
	return RUNNING, err
}

// ServiceLogs returns the logs of a service
func (c *Container) ServiceLogs(namespace []string) (reader io.ReadCloser, err error) {
	return c.client.ServiceLogs(
		context.Background(),
		Namespace(namespace),
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: false,
			Follow:     true,
		},
	)
}
