# 微服务
微服务是一种软件架构模式，用于将大型单片应用程序分解为较小的可管理独立服务，这些服务通过协议进行通信，并且每个服务都专注于做好一件事。微服务的概念并不新鲜，这是面向服务架构的重新构想，但其方法更整体地与unix进程和管道保持一致。
随着组织扩展技术和人数，管理单一代码库变得更加困难。我们通过第一手经验看到，微服务系统可以实现更快的开发周期，更高的生产力和卓越的可扩展系统。
优点：
- 易于规模化开发
- 更容易理解
- 更容易频繁地部署新版本的服务
- 提高容错能力和隔离
- 提高执行速度
- 可重用服务和快速原型


Micro解决了构建微服务系统的关键要求。 它采用微服务架构模式并将其转换为一组工具，充当可扩展平台的构建块。 Micro隐藏了分布式系统的复杂性，并为开发人员提供了很好理解的概念。Micro是一个提供查询和访问微服务的工具集。
运行时有以下功能点组成:
- API网关：使用服务发现进行动态请求路由的单个入口点。 API网关允许您在后端构建可扩展的微服务架构，并在前端整合服务公共API。 micro api通过发现和可插拔处理程序提供强大的路由，以提供http，grpc，websockets，发布事件等
- 集成CLI：用于直接描述，查询和与终端的平台和服务进行交互的CLI。 CLI为您提供了您希望了解微服务发生情况的所有命令。 它还包括一个交互模式。
- 服务代理：基于Go Micro构建的透明代理。 将服务发现，负载平衡，容错，消息编码，中间件，监控等卸载到单个位置。 独立运行或与服务一起运行。
- 模板生成器：创建新的服务模板以快速入门。 Micro提供用于编写微服务的预定义模板。 始终以相同的方式开始，构建相同的服务以提高工作效率。
- SlackOps Bot：一个在您的平台上运行的机器人，允许您从Slack本身管理您的应用程序。 微型机器人支持ChatOps，使您能够通过消息传递与团队一起完成所有工作。 它还包括创建松弛命令作为动态发现的服务的能力
- web：通过Web仪表板，您可以浏览服务，描述其端点，请求和响应格式，甚至可以直接查询它们。 仪表板还为想要即时进入终端的开发人员提供了内置的CLI体验。
- go框架：利用功能强大的Go Micro框架轻松快速地开发微服务。 Go Micro抽象出分布式系统的复杂性，并提供更简单的抽象来构建高度可扩展的微服务。

Go Config
管理复杂的配置
- 动态的-按需动态加载配置
- 可插拔-选择哪个源需要加载，file，环境变量，consul
- Mergeable-覆盖多配置源
- 回滚-键值不存在时指定回滚值
-可观察的- 检测配置的改变
Go Plugins 
- 为go-micro/micro提供组件
- 包含很多受欢迎的后端技术
- grpc,k8s,etcd,kafka等等
- 生产测试



go Mrcro是一个用于微服务开发的框架，提供了分布式系统开发的核心要求，包括RPC和事件驱动的通信
主要功能：
- 服务发现
- 负载均衡
- 消息编码
- 请求、响应
- 异步消息
- 可插拔接口

go micro由一系列包组成：
- transport:同步消息
- broker:异步消息
- codec：消息编码
- registry：服务发现
注册中心提供了一个服务发现机制来将名称解析为地址，用于将名称解析为地址。 它可以由consul，etcd，zookeeper，dns，gossip等支持。服务应该在启动时使用注册表注册，并在关闭时注销。 服务可以选择提供到期TTL并在一定时间间隔内重新注册以确保活跃，并且如果服务终止则清理服务。
- selector:负载均衡
selector是基于registry建立的负载均衡抽象。允许服务使用过滤函数过滤，或使用一个算法被选择，比如，随机等，在发出请求时，客户端会使用选择器。 客户端将使用选择器而不是注册表，因为它提供了内置的负载平衡机制。
- client：发起请求
- server：处理请求


# 1.安装Protobuf

protobuf在代码生成中是必须的，在上一节中，我们已经安装了`protoc`和*   [protoc-gen-go](https://github.com/golang/protobuf)
本节，我们需要安装的是`protoc-gen-micro`
```
go get github.com/micro/protoc-gen-micro
```
# 2. 服务发现

服务发现用于解析服务名到相关地址，默认的发现系统为广播DNS，需要zeroconf.另外实际开发中，可以使用consul.因为发现是可插拔的，也可以使用etcd,k8s,zookeeper等其他组件。
```
MICRO_REGISTRY=consul go run main.go
```
# 3.实现一个服务

```
syntax = "proto3";

service Greeter {
	rpc Hello(HelloRequest) returns (HelloResponse) {}
}

message HelloRequest {
	string name = 1;
}

message HelloResponse {
	string greeting = 2;
}

```
然后编译他，产生相关代码
```
protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. path/to/greeter.proto
```
实现服务

- 实现服务定义的接口
- 初始化一个micro.Service
- 注册handler
- 运行服务service
- 
```
package main

import (
	"context"
	"fmt"

	micro "github.com/micro/go-micro"
	proto "github.com/micro/examples/service/proto"
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
运行服务
```
go run examples/service/main.go
```
定义一个客户端

```
package main

import (
	"context"
	"fmt"

	micro "github.com/micro/go-micro"
	proto "github.com/micro/examples/service/proto"
)


func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(micro.Name("greeter.client"))
	service.Init()

	// Create new greeter client
	greeter := proto.NewGreeterService("greeter", service.Client())

	// Call the greeter
	rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "John"})
	if err != nil {
		fmt.Println(err)
	}

	// Print response
	fmt.Println(rsp.Greeting)
}
```
运行客户端
```
go run client.go
```
Output

```
Hello John
```

# 订阅和发布

Go-micro具有用于事件驱动架构的内置消息代理接口。

Publish

创建一个带名称的发布器

```
p := micro.NewPublisher("events", service.Client())
```
发布一个proto消息

```
p.Publish(context.TODO(), &proto.Event{Name: "event"})
```

Subscribe
创建一个消息处理器，其签名需是`func(context.Context, v interface{}) error`.

```
func ProcessEvent(ctx context.Context, event *proto.Event) error {
	fmt.Printf("Got event %+v\n", event)
	return nil
}
```
注册到一个主题上
```
micro.RegisterSubscriber("events", ProcessEvent)
```
# Plugins

默认情况下，go-micro仅在核心处提供了每个接口的一些实现，但它是完全可插拔的。

Build with plugins

如果需要整合插件，只需将它们链接到一个单独的文件中并重建即可。
创建一个plugins.go
```
import (
        // etcd v3 registry
        _ "github.com/micro/go-plugins/registry/etcdv3"
        // nats transport
        _ "github.com/micro/go-plugins/transport/nats"
        // kafka broker
        _ "github.com/micro/go-plugins/broker/kafka"
)
```
编译
```
// For local use
go build -i -o service ./main.go ./plugins.go
```

插件使用
```
service --registry=etcdv3 --transport=nats --broker=kafka

```
Plugin as option 

```
import (
        "github.com/micro/go-micro" 
        // etcd v3 registry
        "github.com/micro/go-plugins/registry/etcdv3"
        // nats transport
        "github.com/micro/go-plugins/transport/nats"
        // kafka broker
        "github.com/micro/go-plugins/broker/kafka"
)

func main() {
	registry := etcdv3.NewRegistry()
	broker := kafka.NewBroker()
	transport := nats.NewTransport()

        service := micro.NewService(
                micro.Name("greeter"),
                micro.Registry(registry),
                micro.Broker(broker),
                micro.Transport(transport),
        )

	service.Init()
	service.Run()
}
```
另外，可以将插件作为服务的可选项。

# Wrappers
Go-micro包含中间件作为包装器的概念。 客户端或处理程序可以使用装饰器模式进行包装。

Handler 
下面是一个示例，请求日志包装
```
// implements the server.HandlerWrapper
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		fmt.Printf("[%v] server request: %s", time.Now(), req.Endpoint())
		return fn(ctx, req, rsp)
	}
}
```
可以在创建服务的时候初始化：
```
service := micro.NewService(
	micro.Name("greeter"),
	// wrap the handler
	micro.WrapHandler(logWrapper),
)
```
下面是一个客户端收集请求日志的示例：
```
type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	fmt.Printf("[wrapper] client request to service: %s endpoint: %s\n", req.Service(), req.Endpoint())
	return l.Client.Call(ctx, req, rsp)
}

// implements client.Wrapper as logWrapper
func logWrap(c client.Client) client.Client {
	return &logWrapper{c}
}
```
创建服务时，可以如下初始化：
```
service := micro.NewService(
	micro.Name("greeter"),
	// wrap the client
	micro.WrapClient(logWrap),
)

```
