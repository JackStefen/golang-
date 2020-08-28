本笔记以`go-micro v2.9.1`为基准，其他版本可能会存在差异，可根据具体情况具体查看详情

通过一个简单的示例，看一看整个流程的运作

```
package main

import (
	"context"
	"fmt"

	micro "github.com/micro/go-micro"
	proto "just.for.test/microtest"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("greeter"),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

```

首先，第一步启动了一个服务，这个服务的启动选项是给定了该服务的名称。先来看看这个新创建的服务是个怎样的流程

```
// NewService creates and returns a new Service based on the packages within.
func NewService(opts ...Option) Service {
	return newService(opts...)
}
```

`Service`是一个接口类型，包裹底层的库

```
// Service is an interface that wraps the lower level libraries
// within go-micro. Its a convenience method for building
// and initialising services.
type Service interface {
	// The service name
	Name() string
	// Init initialises options
	Init(...Option)
	// Options returns the current options
	Options() Options
	// Client is used to call services
	Client() client.Client
	// Server is for handling requests and events
	Server() server.Server
	// Run the service
	Run() error
	// The service implementation
	String() string
}
```

`newService`实现一个内部的结构体，该结构体实现了Service接口

```
func newService(opts ...Option) Service {
	service := new(service)
	options := newOptions(opts...)

	// service name
	serviceName := options.Server.Options().Name

	// we pass functions to the wrappers since the values can change during initialisation
	authFn := func() auth.Auth { return options.Server.Options().Auth }
	cacheFn := func() *client.Cache { return options.Client.Options().Cache }

	// wrap client to inject From-Service header on any calls
	options.Client = wrapper.FromService(serviceName, options.Client)
	options.Client = wrapper.TraceCall(serviceName, trace.DefaultTracer, options.Client)
	options.Client = wrapper.CacheClient(cacheFn, options.Client)
	options.Client = wrapper.AuthClient(authFn, options.Client)

	// wrap the server to provide handler stats
	options.Server.Init(
		server.WrapHandler(wrapper.HandlerStats(stats.DefaultStats)),
		server.WrapHandler(wrapper.TraceHandler(trace.DefaultTracer)),
		server.WrapHandler(wrapper.AuthHandler(authFn)),
	)

	// set opts
	service.opts = options

	return service
}
```

未完待续.....
