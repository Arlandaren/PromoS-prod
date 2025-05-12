package services

import (
	"context"
	"fmt"
	"reflect"
)

type BindingMap struct {
	factoryFunc reflect.Value
	lifecycle
}

var services = make(map[reflect.Type]BindingMap)

// Add new service to service map
func addService(life lifecycle, factoryFunc interface{}) (err error) {
	factoryFuncType := reflect.TypeOf(factoryFunc)

	// factory function type checking, output should be 1
	if factoryFuncType.Kind() == reflect.Func && factoryFuncType.NumOut() == 1 {
		// okay, it's factory function
		services[factoryFuncType.Out(0)] = BindingMap{
			// service => factory function + lifecycle
			factoryFunc: reflect.ValueOf(factoryFunc),
			lifecycle:   life,
		}
	} else {
		err = fmt.Errorf("type cannot be used as service: %v", factoryFuncType)
	}
	return
}

var contextReference = (*context.Context)(nil)
var contextReferenceType = reflect.TypeOf(contextReference).Elem()

func resolveServiceFromValue(c context.Context, val reflect.Value) (err error) {
	serviceType := val.Elem().Type()

	if serviceType == contextReferenceType {
		val.Elem().Set(reflect.ValueOf(c))
	} else if binding, found := services[serviceType]; found {
		// find service by service type from services map
		if binding.lifecycle == Scoped {
			resolveScopedService(c, val, binding)
		} else {
			val.Elem().Set(invokeFunction(c, binding.factoryFunc)[0])
		}
	} else {
		err = fmt.Errorf("cannot find service %v", serviceType)
	}
	return
}

func resolveScopedService(c context.Context, val reflect.Value, binding BindingMap) (err error) {
	sMap, ok := c.Value(ServiceKey).(serviceMap)
	if ok {
		serviceVal, ok := sMap[val.Type()]
		if !ok {
			serviceVal = invokeFunction(c, binding.factoryFunc)[0]
			sMap[val.Type()] = serviceVal
		}
		val.Elem().Set(serviceVal)
	} else {
		val.Elem().Set(invokeFunction(c, binding.factoryFunc)[0])
	}
	return
}

func resolveFunctionArguments(c context.Context, f reflect.Value, othersArgs ...interface{}) []reflect.Value {
	// make parameter slice from function type
	params := make([]reflect.Value, f.Type().NumIn())

	i := 0
	if othersArgs != nil {
		// map otherArgs array into params (reflect value slice)
		for ; i < len(othersArgs); i++ {
			params[i] = reflect.ValueOf(othersArgs[i])
		}
	}

	// got params => reflect value slice
	for ; i < len(params); i++ {
		pType := f.Type().In(i) // get ith param's type
		pVal := reflect.New(pType)

		// get arguments from services
		err := resolveServiceFromValue(c, pVal)
		if err != nil {
			panic(err)
		}
		params[i] = pVal.Elem()
	}
	return params
}

func invokeFunction(c context.Context, f reflect.Value, otherArgs ...interface{}) []reflect.Value {
	// call function with arguments
	return f.Call(resolveFunctionArguments(c, f, otherArgs...))
}
