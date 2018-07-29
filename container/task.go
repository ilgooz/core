package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
)

// ListTasks returns all docker tasks
func (c *Container) ListTasks(namespace []string) (tasks []swarm.Task, err error) {
	return c.client.TaskList(context.Background(), types.TaskListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + Namespace(namespace),
		}),
	})
}

// TasksError returns the error of matching tasks
func (c *Container) TasksError(namespace []string) (errors []string, err error) {
	tasks, err := c.ListTasks(namespace)
	if err != nil {
		return errors, err
	}
	for _, task := range tasks {
		if task.Status.Err != "" {
			errors = append(errors, task.Status.Err)
		}
	}
	return errors, err
}
