package container

import (
	"context"
	"testing"

	docker "github.com/docker/docker/client"
	"github.com/stvp/assert"
)

func removeSharedNetworkIfExist(c *Container) error {
	_, err := c.sharedNetwork()
	if docker.IsErrNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return c.client.NetworkRemove(context.Background(), Namespace(sharedNetworkNamespace))
}

func TestCreateSharedNetworkIfNeeded(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	err = removeSharedNetworkIfExist(c)
	assert.Nil(t, err)
	err = c.createSharedNetworkIfNeeded()
	assert.Nil(t, err)
}

func TestSharedNetwork(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	network, err := c.sharedNetwork()
	assert.Nil(t, err)
	assert.NotEqual(t, "", network.ID)
}

func TestSharedNetworkID(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	networkID, err := c.SharedNetworkID()
	assert.Nil(t, err)
	assert.NotEqual(t, "", networkID)
}
