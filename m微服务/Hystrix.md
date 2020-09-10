Hystrix

在分布式环境中，许多服务依赖关系中的一些将不可避免地失败。Hystrix是一个库，它通过添加延迟容忍和容错逻辑来帮助您控制这些分布式服务之间的交互。
Hystrix通过隔离服务之间的访问点、停止它们之间的级联故障并提供回调选项来实现这一点，所有这些都提高了系统的整体弹性。


## 2.Hystrix目标

Hystrix被设计用来做以下工作
- 保护和控制通过第三方客户的库访问依赖项(通常通过网络)的延迟和故障。
- 在复杂的分布式系统中停止级联故障
- 快速失败和快速恢复。
- 在可能的情况下进行回退并优雅地降级。
- 使接近实时监测，警报，和操作控制。

## 3.Hystrix解决了什么问题
在复杂的分布式架构中，应用程序有几十个依赖项，每个依赖项都会在某一时刻不可避免地失败。如果主机应用程序不能与这些外部故障隔离，它就有被关闭的风险。

例如，对于依赖30个服务的应用程序，其中每个服务的正常运行时间为99.99%，您可以期望如下结果：

```
99.9930 = 99.7% uptime
0.3% of 1 billion requests = 3,000,000 failures
即便是所有的依赖项都良好的运行，那么每个月也都将有2个小时以上的时间处于不可用状态.
```

现实通常更糟。

即使所有依赖项都运行良好，如果不设计整个系统以实现弹性，那么即使是0.01%的宕机时间对数十个服务的总体影响也相当于每月宕机时间的潜在小时数。


当许多后端系统中的一个成为潜在系统时，它可以阻止整个用户请求.随着流量的增加，一个潜在的后端依赖关系会导致所有服务器上的所有资源在几秒钟内饱和。

应用程序中通过网络或进入客户库(可能导致网络请求)的每个点都是潜在故障的来源。
比故障更糟糕的是，这些应用程序还可能导致服务之间的延迟增加，从而导致排队、线程和其他系统资源，甚至导致系统中更多的级联故障。


当通过第三方客户端执行网络访问时，这些问题就会加剧。
第三方客户端是一个黑箱，其中隐藏了实现细节，并且可能随时更改，每个客户端库的网络或资源配置都不同，并且通常难以监视和更改。

更糟糕的是传递依赖关系，它们执行可能昂贵或容易出错的网络调用，而应用程序没有显式地调用它们。

网络连接失败或降级。服务和服务器出现故障或变慢。新的库或服务部署会改变行为或性能特征。客户端库有缺陷。

所有这些表示失败和延迟的都需要隔离和管理，以便单个依赖项失败不会导致整个应用程序或系统崩溃。


### 设计原则
- 防止任何单一依赖关系耗尽所有容器(如Tomcat)用户线程。
- 减少负载和快速失败而不是排队。
- 在任何可行的地方提供回调，以保护用户避免失败。
- 使用隔离技术来限制任何一个依赖关系的影响。
- 通过接近实时的度量、监视和警报优化发现时间
- 在Hystrix的大多数方面，通过低延迟传播配置更改和支持动态属性更改来优化恢复时间，这允许您使用低延迟反馈循环进行实时操作修改。
- 防止依赖关系客户端在整个执行过程中出现故障，而不仅仅是在网络流量中。

## Hystrix如何实现其目标

Hystrix是这样做的
- 将所有对外部系统(或依赖项)的调用包装在HystrixCommand或HystrixObservableCommand对象中，通常在单独的线程中执行
- 时间长于您定义的阈值的超时调用。有一个默认值，但是对于大多数依赖项，您可以通过属性自定义设置这些超时，使它们略高于每个依赖项性能99.5的百分比。
- 为每个依赖项维护一个小的线程池(或信号量)。如果该依赖项已满，针对该依赖项的请求将立即被拒绝，而不是排队。
- 测定成功、失败(客户端抛出的异常)、超时和线程拒绝的数量。
- 熔断器以在一段时间内手动或自动地停止对特定服务的所有请求，如果该服务的错误百分比超过阈值。
- 当请求失败、被拒绝、超时或短路时执行回退逻辑
- 近实时地监视指标和配置更改。

当您使用Hystrix来包装每个底层依赖项时，上面图中所示的体系结构会发生如下图所示的变化。
每个依赖项相互隔离，在发生延迟时限制资源的饱和，并在回退逻辑中覆盖，回退逻辑决定在依赖项中发生任何类型的故障时应该做出什么响应

hystrix是Netflix开源的一个JAVA项目，hystrix-go是golang的实现版本

作为Hystrix命令执行代码

定义依赖于外部系统的应用程序逻辑，将函数传递给hystrix.Go。当系统正常时，这是唯一执行的事情。

```
hystrix.Go("my_command", func() error {
	// talk to other services
	return nil
}, nil)
```
Go函数的原型为，主要用于运行你指定的函数跟踪函数的健康状况。
如果函数开始变得缓慢或者多次失败，我们将阻塞新的调用，以给服务时间来修复。
在中断期间如果需要执行一些其他代码时，可以指定一个自定义的回调函数。
```
func Go(name string, run runFunc, fallback fallbackFunc) chan error {
	runC := func(ctx context.Context) error {
		return run()
	}
	var fallbackC fallbackFuncC
	if fallback != nil {
		fallbackC = func(ctx context.Context, err error) error {
			return fallback(err)
		}
	}
	return GoC(context.Background(), name, runC, fallbackC)
}
```
上述函数时异步执行函数时使用的，
如果需要同步执行，直接调用Do方法。Do以同步模式的方式运行你指定的函数，阻塞直到函数执行成功，或者返回错误，包括hystrix的断路器错误。

```
func Do(name string, run runFunc, fallback fallbackFunc) error {
	runC := func(ctx context.Context) error {
		return run()
	}
	var fallbackC fallbackFuncC
	if fallback != nil {
		fallbackC = func(ctx context.Context, err error) error {
			return fallback(err)
		}
	}
	return DoC(context.Background(), name, runC, fallbackC)
}
```

从更符合我们思维习惯的角度，先看同步模式的，这样更加容易理解一些
Do方法调用的是Doc函数，该函数是加上Cotext上下文版本的Do函数。

```
func DoC(ctx context.Context, name string, run runFuncC, fallback fallbackFuncC) error {
	done := make(chan struct{}, 1)

	r := func(ctx context.Context) error {
		err := run(ctx)
		if err != nil {
			return err
		}

		done <- struct{}{}
		return nil
	}

	f := func(ctx context.Context, e error) error {
		err := fallback(ctx, e)
		if err != nil {
			return err
		}

		done <- struct{}{}
		return nil
	}

	var errChan chan error
	if fallback == nil {
		errChan = GoC(ctx, name, r, nil)
	} else {
		errChan = GoC(ctx, name, r, f)
	}

	select {
	case <-done:
		return nil
	case err := <-errChan:
		return err
	}
}
```
这个函数主要分为几部分，
重新定义指定的run函数
重新定义回调函数
上述逻辑主要增加了在成功调用后写done chan方便DoC进行检测是否执行结束

接下来的注意内容就是调用GoC函数，这个函数也是Go函数的主要调用逻辑。所以这个才是重中之重


在DoC的最后一部分，就是检测done chan是有数据，还是errChan有数据，根据判断指定的函数逻辑是成功执行了，还是有错误发生。


DoC就这么简单，也非常符合我们的思维方式。

那么重点来了，GoC是Go增加了Context上下文的Go函数，具体逻辑需要详细看一下

```
func GoC(ctx context.Context, name string, run runFuncC, fallback fallbackFuncC) chan error {
	cmd := &command{
		run:      run,
		fallback: fallback,
		start:    time.Now(),
		errChan:  make(chan error, 1),
		finished: make(chan bool, 1),
	}

	// dont have methods with explicit params and returns
	// let data come in and out naturally, like with any closure
	// explicit error return to give place for us to kill switch the operation (fallback)

	circuit, _, err := GetCircuit(name)
	if err != nil {
		cmd.errChan <- err
		return cmd.errChan
	}
	cmd.circuit = circuit
	ticketCond := sync.NewCond(cmd)
	ticketChecked := false
	// When the caller extracts error from returned errChan, it's assumed that
	// the ticket's been returned to executorPool. Therefore, returnTicket() can
	// not run after cmd.errorWithFallback().
	returnTicket := func() {
		cmd.Lock()
		// Avoid releasing before a ticket is acquired.
		for !ticketChecked {
			ticketCond.Wait()
		}
		cmd.circuit.executorPool.Return(cmd.ticket)
		cmd.Unlock()
	}
	// Shared by the following two goroutines. It ensures only the faster
	// goroutine runs errWithFallback() and reportAllEvent().
	returnOnce := &sync.Once{}
	reportAllEvent := func() {
		err := cmd.circuit.ReportEvent(cmd.events, cmd.start, cmd.runDuration)
		if err != nil {
			log.Printf(err.Error())
		}
	}

	go func() {
		defer func() { cmd.finished <- true }()

		// Circuits get opened when recent executions have shown to have a high error rate.
		// Rejecting new executions allows backends to recover, and the circuit will allow
		// new traffic when it feels a healthly state has returned.
		if !cmd.circuit.AllowRequest() {
			cmd.Lock()
			// It's safe for another goroutine to go ahead releasing a nil ticket.
			ticketChecked = true
			ticketCond.Signal()
			cmd.Unlock()
			returnOnce.Do(func() {
				returnTicket()
				cmd.errorWithFallback(ctx, ErrCircuitOpen)
				reportAllEvent()
			})
			return
		}

		// As backends falter, requests take longer but don't always fail.
		//
		// When requests slow down but the incoming rate of requests stays the same, you have to
		// run more at a time to keep up. By controlling concurrency during these situations, you can
		// shed load which accumulates due to the increasing ratio of active commands to incoming requests.
		cmd.Lock()
		select {
		case cmd.ticket = <-circuit.executorPool.Tickets:
			ticketChecked = true
			ticketCond.Signal()
			cmd.Unlock()
		default:
			ticketChecked = true
			ticketCond.Signal()
			cmd.Unlock()
			returnOnce.Do(func() {
				returnTicket()
				cmd.errorWithFallback(ctx, ErrMaxConcurrency)
				reportAllEvent()
			})
			return
		}

		runStart := time.Now()
		runErr := run(ctx)
		returnOnce.Do(func() {
			defer reportAllEvent()
			cmd.runDuration = time.Since(runStart)
			returnTicket()
			if runErr != nil {
				cmd.errorWithFallback(ctx, runErr)
				return
			}
			cmd.reportEvent("success")
		})
	}()

	go func() {
		timer := time.NewTimer(getSettings(name).Timeout)
		defer timer.Stop()

		select {
		case <-cmd.finished:
			// returnOnce has been executed in another goroutine
		case <-ctx.Done():
			returnOnce.Do(func() {
				returnTicket()
				cmd.errorWithFallback(ctx, ctx.Err())
				reportAllEvent()
			})
			return
		case <-timer.C:
			returnOnce.Do(func() {
				returnTicket()
				cmd.errorWithFallback(ctx, ErrTimeout)
				reportAllEvent()
			})
			return
		}
	}()

	return cmd.errChan
}
```


