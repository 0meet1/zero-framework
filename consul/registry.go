package consul

import (
	"strings"

	"github.com/0meet1/zero-framework/global"
	"github.com/hashicorp/consul/api"
)

const (
	REGISTRY_TRUNK = "zero.registry.trunk"
)

type ZeroServeRegistryTrunk interface {
	Registry(*api.AgentServiceRegistration) error
	Deregister(*api.AgentServiceRegistration) error
	FindAll() (map[string]*api.AgentService, error)
}

type xZeroServeRegistryTrunk struct {
	apiConfig *api.Config
	apiClient *api.Client
}

func (registry *xZeroServeRegistryTrunk) Registry(serv *api.AgentServiceRegistration) error {
	return registry.apiClient.Agent().ServiceRegister(serv)
}

func (registry *xZeroServeRegistryTrunk) Deregister(serv *api.AgentServiceRegistration) error {
	return registry.apiClient.Agent().ServiceDeregister(serv.ID)
}

func (registry *xZeroServeRegistryTrunk) FindAll() (map[string]*api.AgentService, error) {
	return registry.apiClient.Agent().Services()
}

func (registry *xZeroServeRegistryTrunk) runServeRegistry() error {
	registry.apiConfig = api.DefaultConfig()

	if strings.TrimSpace(global.StringValue("zero.consul.serverAddr")) != "" {
		registry.apiConfig.Address = global.StringValue("zero.consul.serverAddr")
	}

	client, err := api.NewClient(registry.apiConfig)
	if err != nil {
		return err
	}
	registry.apiClient = client
	return nil
}
