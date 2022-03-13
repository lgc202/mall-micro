package registry

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type Protocol string

const (
	HTTP Protocol = "HTTP"
	GRPC Protocol = "GRPC"
)

type Registry interface {
	Register(address, serviceID, name string, tags []string, port int, protocol Protocol) error
	DeRegister(serviceId string) error
}

type RegistryClient struct {
	Host   string
	Port   int
	Client *api.Client
}

func NewRegisry(host string, port int) (Registry, error) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", host, port)
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &RegistryClient{
		Host:   host,
		Port:   port,
		Client: client,
	}, nil
}

func (r RegistryClient) Register(address, serviceID, name string, tags []string, port int, protocol Protocol) error {
	// 生成对应的检查对象
	check := &api.AgentServiceCheck{
		Timeout:                        "5s",
		Interval:                       "5s",  // 5s做一次检查
		DeregisterCriticalServiceAfter: "15s", // 检查到不健康15s后注销
	}

	switch protocol {
	case HTTP:
		check.HTTP = fmt.Sprintf("http://%s:%d/health", address, port)
	case GRPC:
		check.GRPC = fmt.Sprintf("%s:%d", address, port)
	default:
		return fmt.Errorf("unsupport protocol: %s", protocol)
	}

	// 生成注册对象
	registration := &api.AgentServiceRegistration{
		Name:    name,
		ID:      serviceID,
		Address: address,
		Port:    port,
		Tags:    tags,
		Check:   check,
	}

	return r.Client.Agent().ServiceRegister(registration)
}

func (r RegistryClient) DeRegister(serviceId string) error {
	return r.Client.Agent().ServiceDeregister(serviceId)
}
