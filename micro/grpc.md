# 1.简介

 在gRPC中，客户端应用程序可以直接调用不同计算机上的服务器应用程序上的方法，就像它是本地对象一样，使您可以更轻松地创建分布式应用程序和服务。与许多RPC系统一样，gRPC基于定义服务的思想，指定可以使用其参数和返回类型远程调用的方法。在服务器端，服务器实现此接口并运行gRPC服务器来处理客户端调用。在客户端，客户端有一个存根（在某些语言中称为客户端），它提供与服务器相同的方法。gRPC可以使用protocol buffers作为其接口定义语言（IDL）和其基础消息交换格式,来序列化结构化数据,关于详细的Proto语法介绍，可以看一下另一篇文章[https://www.jianshu.com/p/434ac0fbcf59](https://www.jianshu.com/p/434ac0fbcf59)


![图片来自gRPC doc.png](https://upload-images.jianshu.io/upload_images/3004516-e6b52b9c7f39d69d.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)


# 2.基本概念

## 2.1.服务定义
与许多RPC系统一样，gRPC基于定义服务的思想，指定可以使用其参数和返回类型远程调用的方法。 默认情况下，gRPC使用protocol buffers作为接口定义语言（IDL）来描述服务接口和有效负载消息的结构。 如果需要，可以使用其他替代方案。
```
service HelloService {
  rpc SayHello (HelloRequest) returns (HelloResponse);
}

message HelloRequest {
  string greeting = 1;
}

message HelloResponse {
  string reply = 1;
}
```

gRPC允许您定义四种服务方法：
- Unary RPCs：客户端发送一个请求到服务端，并从服务端得到一个响应，如同常规的函数调用。
`rpc SayHello(HelloRequest) returns (HelloResponse){
}`
- Server streaming RPCs:客户机向服务器发送一个请求，并获取一个流来读取返回的消息序列。客户端从返回的流中读取，直到没有更多消息。 gRPC保证单个RPC调用中的消息排序。
`rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse){
}`
- Client streaming RPCs:客户端再次使用提供的流写入一系列消息并将它们发送到服务器。 一旦客户端写完消息，它就等待服务器读取它们并返回它的响应。 gRPC再次保证在单个RPC调用中的消息排序。
`rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse) {
}`
- Bidirectional streaming RPCs:双方使用读写流发送一系列消息。 这两个流独立运行，因此客户端和服务器可以按照自己喜欢的顺序进行读写：例如，服务器可以在写入响应之前等待接收所有客户端消息，或者它可以交替地读取消息然后写入消息， 或者其他一些读写组合。 保留每个流中的消息顺序。
`rpc BidiHello(stream HelloRequest) returns (stream HelloResponse){
}`
## 2.2API使用
从.proto文件中的服务定义开始，gRPC提供协议缓冲区编译器插件，用于生成客户机和服务器端代码。gRPC用户通常在客户端调用这些API，并在服务器端实现相应的API。
- 在服务器端，服务器实现服务声明的方法，并运行gRPC服务器来处理客户端调用。 gRPC基础结构解码传入请求，执行服务方法并对服务响应进行编码。
- 在客户端，客户端有一个称为存根的本地对象（对于某些语言，首选术语是客户端），它实现与服务相同的方法。 然后，客户端可以在本地对象上调用这些方法，将调用的参数包装在适当的协议缓冲区消息类型中 - gRPC在将请求发送到服务器并返回服务器的protocol buffers的响应。
## 2.3同步和异步模式
同步RPC调用会阻塞直到服务端的响应到达，是最接近RPC所期望的过程调用的抽象。另一方面，网络本质上是异步的，在许多情况下，能够在不阻塞当前线程的情况下启动rpc是很有用的。
# 3.RPC的生命周期
现在让我们仔细看看当gRPC客户端调用gRPC服务器方法时会发生什么。
## 3.1Unary RPC
首先让我们看一下最简单的RPC类型，客户端发送单个请求，得到单个响应。
- 一但客户端调用了stub/clientd对象上的方法，服务端就会得到通知，RPC被调用了，携带客户端关于本次调用的元数据，方法名，和指定的截止时间（如果提供了的话）。
- 然后，服务器可以立即发送回自己的初始元数据（必须在任何响应之前发送），或者等待客户端的请求消息 - 首先发生的是特定于应用程序的消息。
- 一旦服务器具有客户端的请求消息，它就会执行创建和填充其响应所需的任何工作。 然后将响应与状态详细信息（状态代码和可选状态消息）以及可选的尾随元数据一起返回（如果成功）到客户端。
- 如果status 是OK,客户端得到响应，就结束了整个调用。

## 3.2 Server streaming RPC
Server streaming RPC，在得到客户端的请求信息后，期待服务端发送响应的流，发送完所有的响应之后，服务端状态细节和可选的尾元数据也会被服务端发送来结束调用。一旦客户端拥有所有服务器的响应，客户端就会完成。
## 3.3 Client streaming RPC
客户端发送一个请求的流而不是单个请求，服务器发送回单个响应，通常但不一定在收到所有客户端请求后，以及其状态详细信息和可选的尾随元数据。

## 3.4Bidirectional streaming RPC
在双向流式RPC中，调用再次由调用方法的客户端和接收客户端元数据，方法名称和截止时间的服务器启动。 服务器再次可以选择发回其初始元数据或等待客户端开始发送请求。接下来会发生什么取决于应用程序，因为客户端和服务器可以按任何顺序读写 - 流完全独立地运行。 因此，例如，服务器可以等到它收到所有客户端的消息之后再写入其响应，或者服务器和客户端可以“乒乓”：服务器获取请求，然后发回响应，然后客户端发送 另一个基于响应的请求，等等。

## 3.5Deadlines/Timeouts
gRPC允许客户端指定它愿意等待多久待RPC调用完成，直到RPC被中断，并带有DEADLINE_EXCEEDED错误。服务端，可以查询一个特定的RPC是否已经超时，或者还有多久待调用完成。如果指定deadline或者timeout不同语言，方式可能不同。

## 3.6RPC termination
客户端和服务器都对调用的成功做出独立的和本地的决定，并且他们的结论可能不同，这就意味着，你可能在服务端收到（“我已经发送完所有的响应”），但是客户端缺失败了（“响应超时”），服务端也可能在客户端发送完所有请求之前决定完成。
## 3.7Cancelling RPCs
客户端和服务端在任何时候都可以取消RPC调用，取消立即终止RPC，以便不再进行进一步的工作。 它不是“撤消”：取消之前所做的更改将不会被回滚。
## 3.8 Metadata
元数据是以键值对列表形式的特定RPC调用（例如身份验证详细信息）的信息，其中键是字符串，值通常是字符串（但可以是二进制数据）。 元数据对gRPC本身是不透明的 - 它允许客户端提供与服务器调用相关的信息，反之亦然。
## 3.9Channels
gRPC通道提供与指定主机和端口上的gRPC服务器的连接，并在创建客户端存根（或某些语言中的“客户端”）时使用。 客户端可以指定通道参数来修改gRPC的默认行为，例如打开和关闭消息压缩。 通道具有状态，包括已连接和空闲。
# 4.安装
## 4.1. Install gRPC
`go get -u google.golang.org/grpc`
## 4.2. Install Protocol Buffers v3
安装protoc编译器，用于产生gRPC服务代码，下载地址：
[https://github.com/google/protobuf/releases](https://github.com/google/protobuf/releases)
- 解压下载的文件
- 更新PATH环境变量,将protoc二进制可执行文件路径加到环境变量中。
## 4.3 install protoc plugin for golang
`go get -u github.com/golang/protobuf/protoc-gen-go`
# 5.编译示例
示例代码在grpc项目下的examples目录下
```
cd $GOPATH/src/google.golang.org/grpc/examples/helloworld
```
gRPC服务定义在`.proto`文件中，该文件被用于编译产生相关的`.pb.go`文件，`.pb.go`文件是使用protocol编译器`protoc`编译`.proto`文件产生的。示例代码中该文件已经产生，内容涵盖一下两点：
- 产生客户端和服务端代码
- 用于填充，序列化和检索HelloRequest和HelloReply消息类型的代码。
测试：
- `go run greeter_server/main.go`启动服务端运行
- `go run greeter_client/main.go`在的终端里，启动客户端运行

如果在运行上面命令的时候，出现依赖包问题，比如：
```
➜  helloworld git:(master) go run greeter_server/main.go
../../status/status.go:37:2: cannot find package "google.golang.org/genproto/googleapis/rpc/status" in any of:
	/usr/local/go/src/google.golang.org/genproto/googleapis/rpc/status (from $GOROOT)
	/Users/xxx/workspace/src/google.golang.org/genproto/googleapis/rpc/status (from $GOPATH)
```
安装 google.golang.org/genproto:
```
$ wget https://github.com/google/go-genproto/archive/master.tar.gz -O ${GOPATH}/src/google.golang.org/genproto.tar.gz
$ cd ${GOPATH}/src/google.golang.org && tar zxvf genproto.tar.gz && mv go-genproto-master genproto
```
如果顺利，将会看到客户端标准输出：
```
➜  helloworld git:(master) go run greeter_client/main.go
2019/07/12 17:21:47 Greeting: Hello world
```
# 6.更新服务
## 6.1定义新服务
上面已经成功运行了我们的gRPC示例代码，现在当我们需要新增服务需求时，在`.proto`文件中定义相关服务,比如，下面我们新增一个SayHelloAgain方法，方法的参数和返回值和之前的保持不变
```
// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
  // Sends another greeting
  rpc SayHelloAgain (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```
## 6.2.编译proto文件产生新的服务代码
此时，需要使用protoc编译器重新编译一下我们修改后的文件
```
$ protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld
```
执行此命令后，新的helloworld.pb.go文件有了新的变化。
## 6.3更新我们应用程序，重新运行
修改`greeter_server/main.go`文件：
```
/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```
修改`greeter_client/main.go`文件：
```
/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
	r, err = c.SayHelloAgain(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}

```
## 6.4 运行
- `go run greeter_server/main.go` 运行服务端
- `go run greeter_client/main.go` 运行客户端
顺利的话，从控制台打印:
```
➜  helloworld git:(master) go run greeter_client/main.go
2019/07/12 17:40:44 Greeting: Hello world
2019/07/12 17:40:44 Greeting: Hello again world
```
# 7.总结
## 7.1 定义服务
定义一个服务，需要在`.proto`文件中指定一个service :
```
service RouteGuide {
   ...
}
```
然后在service中定义rpc方法，指定请求参数类型，和返回值类型。gRPC允许我们定义四种类型的服务方法，所有这些类型的方法都在RouteGuide服务中。
- 第一种：最简单的，客服端通过stub发送一个请求到服务端，等待一个响应返回，就像普通的方法调用。
`rpc GetFeature(Point) returns (Feature) {}`
- 第二种：服务端流RPC,客户端发送一个请求到服务端，并获取一个返回流用来读取信息序列，客户端读取流知道无更多消息。就像在示例中看到的，将stream关键字放在返回值类型前面就可以定义一个服务端流方法。
`rpc ListFeatures(Rectangle) returns (stream Feature) {}`
- 客户端流RPC,客户端写一个消息序列并将它们发送到服务端，客户端完成写消息后，等待服务端读取它们并返回响应，你可以在请求类型前面加stream关键字来定义一个客户端流RPC.
`rpc RecordRoute(stream Point) returns (RouteSummary) {}`
- 双向RPC,两端都是用消息系列来进行读写，两个流操作是独立的，所以客户端和服务端可以以任何顺序来进行读写，例如，服务端可以等到客户端所有的信息到达后再返回响应。或者读一个写一个，或者其他方式。保留每个流中的消息顺序，你可以在请求参数和响应参数前面加上stream来定义这类方法。
`rpc RouteChat(stream RouteNote) returns (stream RouteNote) {}`
## 7.2.定义方法的请求参数类型和响应参数类型
我们`.proto`文件同样包含protocol buffer 请求和响应的消息类型在方法定义中，如下：
```
message Point {
  int32 latitude = 1;
  int32 longitude = 2;
}
```
## 7.3.编译产生客户端和服务端代码
接下来，通过我们`.proto`文件中的服务定义，产生gRPC客户端和服务端接口，使用protocol buffer的编译器 `protoc`带有gRPC的go语言插件。
```
 protoc -I routeguide/ routeguide/route_guide.proto --go_out=plugins=grpc:routeguide
```
运行上面的命令，可以产生我们需要的.pb.go文件
## 7.4创建服务端
首先，我们来看看如果创建一个RouteGuide服务端
- 实现服务定义产生的服务接口，以此来做实际的工作。
- 运行一个gRPC服务来监听客户端请求并分发他们到正确的服务上
实现RouteGuide：
```
type routeGuideServer struct {
        ...
}
...

func (s *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
    for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	// No feature was found, return an unnamed feature
	return &pb.Feature{"", point}, nil
}
...

func (s *routeGuideServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
        for _, feature := range s.savedFeatures {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}
...

func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
       var pointCount, featureCount, distance int32
	var lastPoint *pb.Point
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointCount++
		for _, feature := range s.savedFeatures {
			if proto.Equal(feature.Location, point) {
				featureCount++
			}
		}
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}
		lastPoint = point
	}
}
...

func (s *routeGuideServer) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
   for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := serialize(in.Location)
                ... // look for notes to be sent to client
		for _, note := range s.routeNotes[key] {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}
...
```
一但我们实现了所有的方法，我们还需要开启一个gRPC服务，客户端才能实际使用我们的服务，如下：
```
flag.Parse()
lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
if err != nil {
        log.Fatalf("failed to listen: %v", err)
}
grpcServer := grpc.NewServer()
pb.RegisterRouteGuideServer(grpcServer, &routeGuideServer{})
... // determine whether to use TLS
grpcServer.Serve(lis)
```
步骤如下：
- 指定客户端请求的端口  err := net.Listen("tcp", fmt.Sprintf(":%d", *port)).
- 创建一个gRPC服务实例  grpc.NewServer().
- 将服务注册到gRPC服务端上 
- 调用Serve()方法来阻塞监听，直到进程被杀或者Stop()被调用。
## 7.5创建客户端
创建一个client stub：
为了调用服务方法，我们需要创建一个gRPC管道来与服务端通信，我们通过传入服务端地址和端口到grpc.Dial()方法来实现：
```
conn, err := grpc.Dial(*serverAddr)
if err != nil {
    ...
}
defer conn.Close()
```
在grpc.Dial方法中可以通过DialOptions来设置权限验证，我们的例子中，目前不需要这样做。
一旦gRPC的管道建立，我们需要一个客户端stub来进行RPC交互，我们可以通过pb包中的NewRouteGuideClient 方法来实现，
```
client := pb.NewRouteGuideClient(conn)
```
调用服务方法：
在gRPC-go中，RPC操作都是同步阻塞模式，这意味着，RPC调用要等待服务端响应。
简单的RPC调用，就像调用本地的方法:
```
feature, err := client.GetFeature(context.Background(), &pb.Point{409146138, -746188906})
if err != nil {
        ...
}
```
如你所见，我们可以在我们建立的stub上进行方法调用，在方法调用的参数上，提供了一个请求的protocol buffer类型的值，并传入了一个context.Context对象，该对象可以在需要的时候改变RPC调用的行为，例如，超时取消，如果调用未返回一个错误，我们就可以读取返回值信息从第一个返回值中。
```
log.Println(feature)
```
服务端流RPC
```
rect := &pb.Rectangle{ ... }  // initialize a pb.Rectangle
stream, err := client.ListFeatures(context.Background(), rect)
if err != nil {
    ...
}
for {
    feature, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
    }
    log.Println(feature)
}
```
和简单的RPC类似，我们传给方法一个context参数，一个请求protocol buffer参数，然而，在获取相应的时候，我们的得到的是一个RouteGuide_ListFeaturesClient实例，客户端可以使用该stream来读取服务端响应。
我们使用RouteGuide_ListFeaturesClient的Recv()方法来重复读取服务端响应到一个protocol buffer对象中（示例中为Feature）直到没有更多的信息。客户端在每次调用Recv()方法后都需要检查异常，如果err 是nil,表示该stream还正常，可以继续读取，如果err == io.EOF表示消息已经读取完了，否则就是一个RPC错误。

客户端流RPC:
```
// Create a random number of random points
r := rand.New(rand.NewSource(time.Now().UnixNano()))
pointCount := int(r.Int31n(100)) + 2 // Traverse at least two points
var points []*pb.Point
for i := 0; i < pointCount; i++ {
	points = append(points, randomPoint(r))
}
log.Printf("Traversing %d points.", len(points))
stream, err := client.RecordRoute(context.Background())
if err != nil {
	log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
}
for _, point := range points {
	if err := stream.Send(point); err != nil {
		if err == io.EOF {
			break
		}
		log.Fatalf("%v.Send(%v) = %v", stream, point, err)
	}
}
reply, err := stream.CloseAndRecv()
if err != nil {
	log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
}
log.Printf("Route summary: %v", reply)
```
RouteGuide_RecordRouteClient 有一个Send方法，我们可以使用它向服务端发送请求，一旦我们结束写入客户端请求到stream中，我们需要调用stream上的CloseAndRecv()方法来告知gRPC我们已经完成写入请求，等待服务端响应。我们通过CloseAndRecv()方法的返回值err可以得到RPC的状态，如果err 是nil 表示该方法的第一个返回值是一个合法的服务端响应。
双端的streaming RPC:
```
stream, err := client.RouteChat(context.Background())
waitc := make(chan struct{})
go func() {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			close(waitc)
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
		}
		log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
	}
}()
for _, note := range notes {
	if err := stream.Send(note); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}
}
stream.CloseSend()
<-waitc
```
语法和客户端stream方法类似，除了我们在结束我们的调用时，需要使用stream山的CloseSend()方法。由于每个端获取双发的消息的顺序都是双发写入消息的顺序，所以客户端和服务端可以任意顺序的读取和写入消息，双端的stream操作时独立的。
## 7.6最后就可以进行相应的测试了
