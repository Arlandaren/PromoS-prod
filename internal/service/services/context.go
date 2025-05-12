package services

import (
	"context"
	"reflect"
)

type ServiceKeyType string

const ServiceKey ServiceKeyType = "services"

type serviceMap map[reflect.Type]reflect.Value

func NewServiceContext(c context.Context) context.Context {
	if c.Value(ServiceKey) == nil {
		return context.WithValue(c, ServiceKey, make(serviceMap))
	} else {
		return c
	}
}
