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
首先实例化一个`service`,然后创建一个`options`初始化一些`options`的选项，之后将其赋值给`service.opts`字段

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

然后调用服务的初始化方法
```
// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *service) Init(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

	s.once.Do(func() {
		// setup the plugins
		for _, p := range strings.Split(os.Getenv("MICRO_PLUGIN"), ",") {
			if len(p) == 0 {
				continue
			}

			// load the plugin
			c, err := plugin.Load(p)
			if err != nil {
				logger.Fatal(err)
			}

			// initialise the plugin
			if err := plugin.Init(c); err != nil {
				logger.Fatal(err)
			}
		}

		// set cmd name
		if len(s.opts.Cmd.App().Name) == 0 {
			s.opts.Cmd.App().Name = s.Server().Options().Name
		}

		// Initialise the command flags, overriding new service
		if err := s.opts.Cmd.Init(
			cmd.Auth(&s.opts.Auth),
			cmd.Broker(&s.opts.Broker),
			cmd.Registry(&s.opts.Registry),
			cmd.Runtime(&s.opts.Runtime),
			cmd.Transport(&s.opts.Transport),
			cmd.Client(&s.opts.Client),
			cmd.Config(&s.opts.Config),
			cmd.Server(&s.opts.Server),
			cmd.Store(&s.opts.Store),
			cmd.Profile(&s.opts.Profile),
		); err != nil {
			logger.Fatal(err)
		}

		// Explicitly set the table name to the service name
		name := s.opts.Cmd.App().Name
		s.opts.Store.Init(store.Table(name))
	})
}

```

该方法会初始化相关`options`,并解析命令行参数覆盖新的`service`.

现在基本的初始化操作基本完成，调用相关的方法，可以获取需要的对象,比如，需要服务器对象时

```
func (s *service) Server() server.Server {
	return s.opts.Server
}
```

下一步就是注册处理函数了，调用的是有`proto`文件生成的函数，具体内容可以看一下
第二个参数是一个接口类型，需要实现的方法是Hello.
```
func RegisterGreeterHandler(s server.Server, hdlr GreeterHandler, opts ...server.HandlerOption) error {
	type greeter interface {
		Hello(ctx context.Context, in *HelloRequest, out *HelloResponse) error
	}
	type Greeter struct {
		greeter
	}
	h := &greeterHandler{hdlr}
	return s.Handle(s.NewHandler(&Greeter{h}, opts...))
}

```


最后调用`Run`方法启动服务

```
func (s *service) Run() error {
	// register the debug handler
	s.opts.Server.Handle(
		s.opts.Server.NewHandler(
			handler.NewHandler(s.opts.Client),
			server.InternalHandler(true),
		),
	)

	// start the profiler
	if s.opts.Profile != nil {
		// to view mutex contention
		rtime.SetMutexProfileFraction(5)
		// to view blocking profile
		rtime.SetBlockProfileRate(1)

		if err := s.opts.Profile.Start(); err != nil {
			return err
		}
		defer s.opts.Profile.Stop()
	}

	if logger.V(logger.InfoLevel, logger.DefaultLogger) {
		logger.Infof("Starting [service] %s", s.Name())
	}

	if err := s.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, signalutil.Shutdown()...)
	}

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-s.opts.Context.Done():
	}

	return s.Stop()
}
```



未完待续.....
