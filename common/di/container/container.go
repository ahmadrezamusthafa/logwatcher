package container

import (
	"github.com/ahmadrezamusthafa/logwatcher/common/logger"
	"github.com/facebookgo/inject"
)

type Service interface {
	StartUp()
	Shutdown()
}

type ContainerAware interface {
	SetContainer(c ServiceContainer)
}

type ServiceContainer interface {
	Ready() error
	GetService(id string) (interface{}, bool)
	RegisterService(id string, svc interface{})
	RegisterServices(services map[string]interface{})
	Shutdown()
}

type ServiceRegistry struct {
	graph    inject.Graph
	services map[string]interface{}
	order    map[int]string
	ready    bool
}

func (reg *ServiceRegistry) GetService(id string) (interface{}, bool) {
	svc, ok := reg.services[id]
	return svc, ok
}

func (reg *ServiceRegistry) RegisterService(id string, svc interface{}) {
	err := reg.graph.Provide(&inject.Object{Name: id, Value: svc, Complete: false})
	if err != nil {
		panic(err.Error())
	}
	reg.order[len(reg.order)] = id
	reg.services[id] = svc
}

func (reg *ServiceRegistry) RegisterServices(services map[string]interface{}) {
	for id, svc := range reg.services {
		reg.RegisterService(id, svc)
	}
}

func (reg *ServiceRegistry) Ready() error {
	if reg.ready {
		return nil
	}
	err := reg.graph.Populate()
	if err != nil {
		return err
	}
	for i := 0; i < len(reg.order); i++ {
		k := reg.order[i]
		obj := reg.services[k]
		if s, success := obj.(Service); success {
			containerAware, ok := s.(ContainerAware)
			if ok {
				containerAware.SetContainer(reg)
			}
			s.StartUp()
		}
	}
	reg.ready = true
	return err
}

func (reg *ServiceRegistry) Shutdown() {
	for _, name := range reg.order {
		if service, ok := reg.services[name]; ok {
			if s, success := service.(Service); success {
				s.Shutdown()
			}
		}
	}
}

func NewContainer() *ServiceRegistry {
	logger.Info("Init service container...")
	return &ServiceRegistry{services: make(map[string]interface{}), order: make(map[int]string), ready: false}
}
