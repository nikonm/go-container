package go_container

import (
	"fmt"
	"github.com/nikonm/go-container/errors"
	"reflect"
	"strings"
)

// If service implements it, then method called when container StopAll called
type Stopper interface {
	Stop() error
}

type Container struct {
	basePkg     string
	rawServices map[string]interface{}
	services    map[string]interface{}
	stoppes     []func() error
}

// Retrieve service by keyName
func (c *Container) GetService(key string) (interface{}, error) {
	if c.basePkg != "" && strings.Index(key, ".") == -1 {
		key = c.basePkg + "." + key
	}

	s, ok := c.services[key]
	if !ok {
		return nil, c.error(key, errors.ServiceNotFound)
	}
	return s, nil
}

// Stopping all services which implements Stopper
func (c *Container) StopAll() error {
	for _, s := range c.stoppes {
		if err := s(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Container) init() (*Container, error) {
	for name, fn := range c.rawServices {
		if err := c.initService(name, fn, nil); err != nil {
			return c, err
		}
	}
	return c, nil
}

func (c *Container) initService(name string, fn interface{}, chain map[string]bool) error {
	if _, ok := c.services[name]; ok { // Already initialized
		return nil
	}
	v := reflect.TypeOf(fn)

	if v.Kind() == reflect.Func {
		args := make([]reflect.Value, 0)
		if v.NumOut() != 2 {
			return c.error(name, errors.InvalidOutArgumentsCount)
		}

		v3 := v.Out(1)

		if v3.String() != "error" {
			return c.error(name, errors.InvalidOutSign)
		}

		if v.NumIn() > 0 {
			for i := 0; i < v.NumIn(); i++ {
				val := v.In(0)
				var cKey string

				switch val.Kind() {
				case reflect.Ptr:
					cKey = val.Elem().String()
				default:
					cKey = val.String()
				}

				rawArg, ok := c.rawServices[cKey]
				if !ok {
					return c.error(name, errors.NotResolveDep)
				}
				arg, found := c.services[cKey]
				if !found {
					if chain == nil {
						chain = make(map[string]bool)
					}
					if _, ok := chain[name]; ok {
						return c.error(name, errors.CircleLoop)
					}
					chain[name] = true
					if err := c.initService(cKey, rawArg, chain); err != nil {
						return err
					}
					arg = c.services[cKey]
				}
				args = append(args, reflect.ValueOf(arg))
			}
		}
		r := reflect.ValueOf(fn).Call(args)
		c.services[name] = r[0].Interface()
		err := r[1].Interface()
		if err != nil {
			return c.error(name, errors.ServiceInitError(err.(error)))
		}
		if s, ok := c.services[name].(Stopper); ok {
			c.stoppes = append(c.stoppes, s.Stop)
		}

		return nil
	}
	return c.error(name, errors.InvalidInitSign)
}

func (c *Container) error(name string, err errors.Error) error {
	fn, ok := c.rawServices[name]
	if ok {
		err.SetData(fmt.Sprintf("(container key = %s func signature = %s)", name, reflect.TypeOf(fn).String()))
	}
	err.SetData(fmt.Sprintf("(container key = %s)", name))
	return err
}

func New(opt *Options) (*Container, error) {
	rs, err := opt.getServices()
	if rs == nil {
		return nil, err
	}
	c := &Container{
		basePkg:     opt.BasePkg,
		rawServices: rs,
		services:    map[string]interface{}{},
	}

	return c.init()
}

type Options struct {
	// Base package, if set then all services use it
	BasePkg string
	// List of service constructor function, ex. func(s *OtherService) (*Service, error)
	Services []interface{}
}

func (opt *Options) getServices() (map[string]interface{}, errors.Error) {
	services := map[string]interface{}{}
	for _, sFn := range opt.Services {
		vt := reflect.TypeOf(sFn)

		if vt.NumOut() < 1 {
			return nil, errors.InvalidInitSign
		}
		val := vt.Out(0)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		services[val.String()] = sFn
	}
	return services, errors.Error{}
}
