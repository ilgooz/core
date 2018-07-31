// +build integration

package container

import (
	"testing"

	"github.com/stvp/assert"
)

func startTestService(name []string) (serviceID string, err error) {
	c, err := New()
	if err != nil {
		return "", err
	}
	return c.StartService(ServiceOptions{
		Image:     "nginx",
		Namespace: name,
	})
}

func TestIntegrationStartService(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestStartService"}
	serviceID, err := startTestService(namespace)
	defer c.StopService(namespace)
	assert.Nil(t, err)
	assert.NotEqual(t, "", serviceID)
}

func TestIntegrationStartService2Times(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestStartService2Times"}
	startTestService(namespace)
	defer c.StopService(namespace)
	serviceID, err := startTestService(namespace)
	assert.NotNil(t, err)
	assert.Equal(t, "", serviceID)
}

func TestIntegrationStopService(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestStopService"}
	startTestService(namespace)
	err = c.StopService(namespace)
	assert.Nil(t, err)
}

func TestIntegrationStopNotExistingService(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestStopNotExistingService"}
	err = c.StopService(namespace)
	assert.Nil(t, err)
}

func TestIntegrationServiceStatusNeverStarted(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestServiceStatusNeverStarted"}
	status, err := c.ServiceStatus(namespace)
	assert.Nil(t, err)
	assert.NotEqual(t, RUNNING, status)
	assert.Equal(t, STOPPED, status)
}

func TestIntegrationServiceStatusRunning(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestServiceStatusRunning"}
	startTestService(namespace)
	defer c.StopService(namespace)
	status, err := c.ServiceStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, RUNNING)
	assert.NotEqual(t, status, STOPPED)
}

func TestIntegrationServiceStatusStopped(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestServiceStatusStopped"}
	startTestService(namespace)
	c.StopService(namespace)
	status, err := c.ServiceStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
	assert.NotEqual(t, status, RUNNING)
}

func TestIntegrationFindServiceNotExisting(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	_, err = c.FindService([]string{"TestFindServiceNotExisting"})
	assert.NotNil(t, err)
}

func TestIntegrationFindService(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestFindService"}
	startTestService(namespace)
	defer c.StopService(namespace)
	service, err := c.FindService(namespace)
	assert.Nil(t, err)
	assert.NotEqual(t, "", service.ID)
}

func TestIntegrationFindServiceCloseName(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestFindServiceCloseName", "name"}
	namespace1 := []string{"TestFindServiceCloseName", "name2"}
	startTestService(namespace)
	defer c.StopService(namespace)
	startTestService(namespace1)
	defer c.StopService(namespace1)
	service, err := c.FindService(namespace)
	assert.Nil(t, err)
	assert.NotEqual(t, "", service.ID)
}

func TestIntegrationFindServiceStopped(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestFindServiceStopped"}
	startTestService(namespace)
	c.StopService(namespace)
	_, err = c.FindService(namespace)
	assert.NotNil(t, err)
}

func TestIntegrationListServices(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	c.StartService(ServiceOptions{
		Image:     "nginx",
		Namespace: []string{"TestListServices"},
		Labels: map[string]string{
			"label_name": "value_1",
		},
	})
	c.StartService(ServiceOptions{
		Image:     "nginx",
		Namespace: []string{"TestListServiceswithValue2"},
		Labels: map[string]string{
			"label_name_2": "value_2",
		},
	})
	defer c.StopService([]string{"TestListServices"})
	defer c.StopService([]string{"TestListServiceswithValue2"})
	services, err := c.ListServices("label_name")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
	assert.Equal(t, Namespace([]string{"TestListServices"}), services[0].Spec.Name)
}

func TestIntegrationServiceLogs(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestServiceLogs"}
	startTestService(namespace)
	defer c.StopService(namespace)
	reader, err := c.ServiceLogs(namespace)
	assert.Nil(t, err)
	assert.NotNil(t, reader)
}