  
[![codecov](https://codecov.io/gh/nikonm/go-container/branch/master/graph/badge.svg)](https://codecov.io/gh/nikonm/go-container)
[![Build Status](https://travis-ci.org/nikonm/go-container.svg?branch=master)](https://travis-ci.org/nikonm/go-container)

### DI Container
 - Container uses reflect on initialization
 - Uses auto-substitution in the service initialization function, return error if circle dependency loop found
 - Supports stopping services, can be used for graceful exit. Use the Stopper interface for it
 - There is a utility function "GetPkgPath" to get the package name + structure
 
#### Example
 ```go
type Service1 struct {
    Service2 *Service2 
}
func NewService1(s *Service2) (*Service1, error) {
	return &Service1{Service2: s}, nil
}

type Service2 struct {}
func (s2 *Service2) Stop() error {
    return nil
}

func NewService2() (*Service2, error) {
	return &Service2{}, nil
}

func main() {
    container, err := New(&Options{
		BasePkg:  "go_container", // you can use it for shorten name resolving
		Services: []interface{}{
			NewService2,
			NewService1,
		},
	})
    service2, err := container.GetService("Service2").(*Service2) // Returning Service2 instance

    err := container.StopAll() // Call all Stop functions in each service(Graceful exit) 
}
```