## 1.Hystrix

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
GoC函数首先实例化cmd，该结构是在熔断器上的运行使用的。常用于描述在熔断器上的run、fallback函数。

errChan用于记录函数或者熔断器的错误，是一个带有缓存的chan
finished 标识是否结束的chan,带缓存

然后GetCircuit根据参数name 从circuitBreakers中获取circuit.如果获取不到newCircuitBreaker
创建一个然后放入到circuitBreakers中，因为这个操作需要线程安全的，所以加了两道加锁机制

如果报错则将错误写入到errChan,然后return

获取到circuit之后赋值给cmd的circuit属性。 

接下来，根据cmd创建了一个Cond,用于返还ticket,在returnTicket函数临时变量中，
```
returnTicket := func() {
		cmd.Lock()
		// Avoid releasing before a ticket is acquired.
		for !ticketChecked {
			ticketCond.Wait()
		}
		cmd.circuit.executorPool.Return(cmd.ticket)
		cmd.Unlock()
	}
```
这也是sync中Cond的最佳实践：加锁，检测条件，执行逻辑，解锁。

然后创建一个sync.Once，用于接下来两个goroutine.确保只有最快执行的goroutine才会运行errWithFallback()和reportAllEvent函数

reportAllEvent函数用于上报事件。


接下来就是两个goroutine用于检测各种状态，并做相应状态下的动作。

第一个goroutine,在结束的时候，会defer写入finished管道。

然后判断熔断器是否打开，如果最新的执行有较高的错误率，将拒绝新的请求，来使后端进行恢复，直到感觉状态正常一些后，才会允许新的流量

如果当前不允许访问，则ticketChecked=true,Cond Singal()执行returnOnce.之后直接return.本次该goroutine结束。

如果当前环境宽裕，运行流量进来，那么要么可以从executorPool中获取到ticket，要么因为获取不到ticket导致执行select中default中的
的逻辑，该逻辑其实和上面熔断器打开有点像，只不过errWithFallback中的错误信息为ErrMaxConcurrency.标识有太多的并行逻辑执行。


如果获取到ticket之后。就可以运行我们的run函数啦，根据起止时间，可以在returnOnce中计算run函数的运行时间runDuration

并根据run函数运行的返回值，来判断是上报success还是上报错误信息。

以上就是第一个goroutine的全部逻辑了。在所有的内容执行完成后，将cmd.finished写入值标识熔断器包裹的逻辑执行完成啦.

下面的goroutine主要是用来判断上下文的控制和超时控制，还有根据cmd.finished直接结束

这个goroutine中，根据我们设置的setting获取Timeout超时设置，来定义一个时间计时器

然后select判断哪个case会成功，来执行相应的逻辑，除了cmd.finished分支，其他的分支都是需要上报的。


以上就是GoC中的全部逻辑了。

我们发现在每个上报逻辑中，首先归还ticket,然后根据不同的状态执行errorWithFallback,最后reportAllEvent

```
func (c *command) errorWithFallback(ctx context.Context, err error) {
	eventType := "failure"
	if err == ErrCircuitOpen {
		eventType = "short-circuit"
	} else if err == ErrMaxConcurrency {
		eventType = "rejected"
	} else if err == ErrTimeout {
		eventType = "timeout"
	} else if err == context.Canceled {
		eventType = "context_canceled"
	} else if err == context.DeadlineExceeded {
		eventType = "context_deadline_exceeded"
	}

	c.reportEvent(eventType)
	fallbackErr := c.tryFallback(ctx, err)
	if fallbackErr != nil {
		c.errChan <- fallbackErr
	}
}
```

errorWithFallback会根据我们的err参数来决定时间类型


```
func (c *command) tryFallback(ctx context.Context, err error) error {
	if c.fallback == nil {
		// If we don't have a fallback return the original error.
		return err
	}

	fallbackErr := c.fallback(ctx, err)
	if fallbackErr != nil {
		c.reportEvent("fallback-failure")
		return fmt.Errorf("fallback failed with '%v'. run error was '%v'", fallbackErr, err)
	}

	c.reportEvent("fallback-success")

	return nil
}
```

如果参数c.fallback为nil, 直接返回err,如果fallback不为空，调用fallback判断返回值是否为nil, 不为空则上报。为空则上报fallback成功

```
func (c *command) reportEvent(eventType string) {
	c.Lock()
	defer c.Unlock()

	c.events = append(c.events, eventType)
}

```

上报事件，就是让cmd中的events列表中增加标识事件类型的字符串

```
// ReportEvent records command metrics for tracking recent error rates and exposing data to the dashboard.
func (circuit *CircuitBreaker) ReportEvent(eventTypes []string, start time.Time, runDuration time.Duration) error {
	if len(eventTypes) == 0 {
		return fmt.Errorf("no event types sent for metrics")
	}

	circuit.mutex.RLock()
	o := circuit.open
	circuit.mutex.RUnlock()
	if eventTypes[0] == "success" && o {
		circuit.setClose()
	}

	var concurrencyInUse float64
	if circuit.executorPool.Max > 0 {
		concurrencyInUse = float64(circuit.executorPool.ActiveCount()) / float64(circuit.executorPool.Max)
	}

	select {
	case circuit.metrics.Updates <- &commandExecution{
		Types:            eventTypes,
		Start:            start,
		RunDuration:      runDuration,
		ConcurrencyInUse: concurrencyInUse,
	}:
	default:
		return CircuitError{Message: fmt.Sprintf("metrics channel (%v) is at capacity", circuit.Name)}
	}

	return nil
}
```
发现么，reportAllEvent函数在上报事件的时候，就是将cmd.events中的事件作为commanExecution一部分上传到circuit.metrics.Updates中
否则报 metrics chann 满负荷的错误。


到现在为止，一切都还未结束，我们先插入一部分配置的内容，然后再继续往下面看


一般情况下，在应用程序启动的时候，可以调用`hystrix.ConfigureCommand()`方法来调整每个命令的设置。

```
hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
	Timeout:               1000,
	MaxConcurrentRequests: 100,
	ErrorPercentThreshold: 25,
})
```
也可以使用`hystrix.Configure()` 接受一个map[string]CommandConfig

我们首先看看settings.go文件中的init()方法

```
var circuitSettings map[string]*Settings
var settingsMutex *sync.RWMutex
var log logger

func init() {
	circuitSettings = make(map[string]*Settings)
	settingsMutex = &sync.RWMutex{}
	log = DefaultLogger
}

```

全局变量的初始化。circuitSettings 是熔断器的设置，字段的键是各个cmd的名称,值为`Settings`,具体为下：

```
type Settings struct {
	Timeout                time.Duration
	MaxConcurrentRequests  int
	RequestVolumeThreshold uint64
	SleepWindow            time.Duration
	ErrorPercentThreshold  int
}
```
上面的几个配置项，都有默认值
```
var (
	// DefaultTimeout is how long to wait for command to complete, in milliseconds
	DefaultTimeout = 1000
	// DefaultMaxConcurrent is how many commands of the same type can run at the same time
	DefaultMaxConcurrent = 10
	// DefaultVolumeThreshold is the minimum number of requests needed before a circuit can be tripped due to health
	DefaultVolumeThreshold = 20
	// DefaultSleepWindow is how long, in milliseconds, to wait after a circuit opens before testing for recovery
	DefaultSleepWindow = 5000
	// DefaultErrorPercentThreshold causes circuits to open once the rolling measure of errors exceeds this percent of requests
	DefaultErrorPercentThreshold = 50
	// DefaultLogger is the default logger that will be used in the Hystrix package. By default prints nothing.
	DefaultLogger = NoopLogger{}
)
```
- Timeout: 执行cmd的超时时间，默认值为1000
- MaxConcurrentRequests: 最大的并发请求数，默认为10
- RequestVolumeThreshold: 因为健康状况导致的熔断器打开的最小请求数据，默认为20
- SleepWindow: 当熔断器打开之后需要等多久之后测试是否恢复的时间窗口，以毫秒为单位，默认为5s
- ErrorPercentThreshold： 一旦错误的度量数超过指定的百分比，将导致熔断器打开,默认值为50%

上述这些配置项，如果不进行配置，将使用默认的配置项。


说完配置之后，再来看看我们上报的信息，被用来怎么了，上面我们说了，上报的信息都放到了circuit.metrics.Updates中。那谁来处理呢？

要明白这一点，需要知道metrics是干啥的。首先每个circuit创建用于每个ExecutorPool去跟踪请求是否应该被允许，或者拒绝，如果熔断器健康状态比较低

```
type CircuitBreaker struct {
	Name                   string
	open                   bool
	forceOpen              bool
	mutex                  *sync.RWMutex
	openedOrLastTestedTime int64

	executorPool *executorPool
	metrics      *metricExchange
}
```

open 属性由于标识，熔断器是否处于打卡状态，
在上面我们判断是否允许新的请求进来的时候，用到了cmd.AllowRequst()
我们来看一下是什么原理

```
// AllowRequest is checked before a command executes, ensuring that circuit state and metric health allow it.
// When the circuit is open, this call will occasionally return true to measure whether the external service
// has recovered.
func (circuit *CircuitBreaker) AllowRequest() bool {
	return !circuit.IsOpen() || circuit.allowSingleTest()
}

func (circuit *CircuitBreaker) allowSingleTest() bool {
	circuit.mutex.RLock()
	defer circuit.mutex.RUnlock()

	now := time.Now().UnixNano()
	openedOrLastTestedTime := atomic.LoadInt64(&circuit.openedOrLastTestedTime)
	if circuit.open && now > openedOrLastTestedTime+getSettings(circuit.Name).SleepWindow.Nanoseconds() {
		swapped := atomic.CompareAndSwapInt64(&circuit.openedOrLastTestedTime, openedOrLastTestedTime, now)
		if swapped {
			log.Printf("hystrix-go: allowing single test to possibly close circuit %v", circuit.Name)
		}
		return swapped
	}

	return false
}
```

可以发现，其不仅会判断Isopen()，还会进行`allowSingleTest`判断
而后者会在熔断器在打开时，根据时间窗口来适当的放出一部分测试的请求，来测试是否需要关闭熔断器


metrics属性是一个metricExchange结构体
```
type metricExchange struct {
	Name    string
	Updates chan *commandExecution
	Mutex   *sync.RWMutex

	metricCollectors []metricCollector.MetricCollector
}
```

Updates 就是用来记录上报数据的。metricCollectors是一个MetricCollector列表
MetricCollector是一个接口类型，代表所有收集器必须遵守的合同，以收集熔断路统计的数据。
只要未在hystrix context上下文之外进行修改，此接口的实现就不必维护其数据存储周围的锁。核心方法为
```
type MetricCollector interface {
	Update(MetricResult)
	Reset()
}

```
Update 从一个命令执行中接受一组metrics
Reset 重启内部的计数器和计时器。


我们先来看看是如何实例化该结构体的

```
func newMetricExchange(name string) *metricExchange {
	m := &metricExchange{}
	m.Name = name

	m.Updates = make(chan *commandExecution, 2000)
	m.Mutex = &sync.RWMutex{}
	m.metricCollectors = metricCollector.Registry.InitializeMetricCollectors(name)
	m.Reset()

	go m.Monitor()

	return m
}
```
在初始化m.metricCollectors属性时，调用了Registry的InitializeMetricCollectors方法
Registry是一个默认的metricCollectorRegistry，熔断器用来收集关于熔断器的健康情况数据
而InitializeMetricCollectors方法运行已经注册的MetricCollector初始化操作，创建一个MetricCollectors数组

```
var Registry = metricCollectorRegistry{
	lock: &sync.RWMutex{},
	registry: []func(name string) MetricCollector{
		newDefaultMetricCollector,
	},
}

type metricCollectorRegistry struct {
	lock     *sync.RWMutex
	registry []func(name string) MetricCollector
}

func (m *metricCollectorRegistry) InitializeMetricCollectors(name string) []MetricCollector {
	m.lock.RLock()
	defer m.lock.RUnlock()

	metrics := make([]MetricCollector, len(m.registry))
	for i, metricCollectorInitializer := range m.registry {
		metrics[i] = metricCollectorInitializer(name)
	}
	return metrics
}
```
看到newDefaultMetricCollector了吧，这个就是我们默认的收集器

```
func newDefaultMetricCollector(name string) MetricCollector {
	m := &DefaultMetricCollector{}
	m.mutex = &sync.RWMutex{}
	m.Reset()
	return m
}

type DefaultMetricCollector struct {
	mutex *sync.RWMutex

	numRequests *rolling.Number
	errors      *rolling.Number

	successes               *rolling.Number
	failures                *rolling.Number
	rejects                 *rolling.Number
	shortCircuits           *rolling.Number
	timeouts                *rolling.Number
	contextCanceled         *rolling.Number
	contextDeadlineExceeded *rolling.Number

	fallbackSuccesses *rolling.Number
	fallbackFailures  *rolling.Number
	totalDuration     *rolling.Timing
	runDuration       *rolling.Timing
}
```
`DefaultMetricCollector`默认的收集器持有熔断器的状态信息，看一下完整的`DefaultMetricCollector`,可以发现，有两种数据类型`rolling.Number`
和`rolling.Timing`,这两个数据结构才是真实存储数据的地方。先看一下`rolling.Number`
```

type Number struct {
	Buckets map[int64]*numberBucket
	Mutex   *sync.RWMutex
}

type numberBucket struct {
	Value float64
}

```
Number在一定数量的时间段内跟踪numberBucket。当前时间段长一秒，仅保留最后10秒。也就是说每个numberBucket记录当前秒数的记录
Buckets的key是当前的秒数


我们先看metricExchange是如何监控的，然后再看如何在`rolling.Number`中增删改查数据的。在newMetricExchange最后的时候开启了一个monitor监控，
```
func (m *metricExchange) Monitor() {
	for update := range m.Updates {
		// we only grab a read lock to make sure Reset() isn't changing the numbers.
		m.Mutex.RLock()

		totalDuration := time.Since(update.Start)
		wg := &sync.WaitGroup{}
		for _, collector := range m.metricCollectors {
			wg.Add(1)
			go m.IncrementMetrics(wg, collector, update, totalDuration)
		}
		wg.Wait()

		m.Mutex.RUnlock()
	}
}

```
在监控中，会遍历收集器列表，然后异步执行IncrementMetrics来更新相关数据项存储

```

func (m *metricExchange) IncrementMetrics(wg *sync.WaitGroup, collector metricCollector.MetricCollector, update *commandExecution, totalDuration time.Duration) {
	// granular metrics
	r := metricCollector.MetricResult{
		Attempts:         1,
		TotalDuration:    totalDuration,
		RunDuration:      update.RunDuration,
		ConcurrencyInUse: update.ConcurrencyInUse,
	}

	switch update.Types[0] {
	case "success":
		r.Successes = 1
	case "failure":
		r.Failures = 1
		r.Errors = 1
	case "rejected":
		r.Rejects = 1
		r.Errors = 1
	case "short-circuit":
		r.ShortCircuits = 1
		r.Errors = 1
	case "timeout":
		r.Timeouts = 1
		r.Errors = 1
	case "context_canceled":
		r.ContextCanceled = 1
	case "context_deadline_exceeded":
		r.ContextDeadlineExceeded = 1
	}

	if len(update.Types) > 1 {
		// fallback metrics
		if update.Types[1] == "fallback-success" {
			r.FallbackSuccesses = 1
		}
		if update.Types[1] == "fallback-failure" {
			r.FallbackFailures = 1
		}
	}

	collector.Update(r)

	wg.Done()
}

```
执行默认收集器的更新操作,增加存储项的值

```
func (d *DefaultMetricCollector) Update(r MetricResult) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	d.numRequests.Increment(r.Attempts)
	d.errors.Increment(r.Errors)
	d.successes.Increment(r.Successes)
	d.failures.Increment(r.Failures)
	d.rejects.Increment(r.Rejects)
	d.shortCircuits.Increment(r.ShortCircuits)
	d.timeouts.Increment(r.Timeouts)
	d.fallbackSuccesses.Increment(r.FallbackSuccesses)
	d.fallbackFailures.Increment(r.FallbackFailures)
	d.contextCanceled.Increment(r.ContextCanceled)
	d.contextDeadlineExceeded.Increment(r.ContextDeadlineExceeded)

	d.totalDuration.Add(r.TotalDuration)
	d.runDuration.Add(r.RunDuration)
}
```


`Update`操作各项`rollings.Number`值时，执行的`Increment`操作，该操作首先获取当前时间的Bucket。如果没有，就新建，如果已经有了，会和参数相加后付给当前时间，
然后删除掉历史数据（10s）之前的老数据，整个过程当然都是需要持有锁的。

```
func (r *Number) getCurrentBucket() *numberBucket {
	now := time.Now().Unix()
	var bucket *numberBucket
	var ok bool

	if bucket, ok = r.Buckets[now]; !ok {
		bucket = &numberBucket{}
		r.Buckets[now] = bucket
	}

	return bucket
}

func (r *Number) Increment(i float64) {
	if i == 0 {
		return
	}

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	b := r.getCurrentBucket()
	b.Value += i
	r.removeOldBuckets()
}

func (r *Number) removeOldBuckets() {
	now := time.Now().Unix() - 10

	for timestamp := range r.Buckets {
		// TODO: configurable rolling window
		if timestamp <= now {
			delete(r.Buckets, timestamp)
		}
	}
}

```

在记录单个执行记录的运行时长，和总时长时，使用的记录结构是`rolling.Timing`

```
type Timing struct {
	Buckets map[int64]*timingBucket
	Mutex   *sync.RWMutex

	CachedSortedDurations []time.Duration
	LastCachedTime        int64
}

type timingBucket struct {
	Durations []time.Duration
}
```
这个和`rolling.Number`逻辑其实是非常像的，定时维护每个时间段的持续时间（从Start到上报数据时过了多久，可以在Monitor逻辑中看一下）。
持续时间保存在CachedSortedDurations中，这个数组是排好序的，1分钟内的记录，以允许各种要根据源数据计算的统计信息。


最后，还要提一下的就是这么多参数，如何更加直观的展示出来呢，就像`pprof`有一个web页面可以查看当前程序运行指标一样，`Hystrix`也有一个指标查看的`dashboard`


只要在我们项目中增加
```
hystrixStreamHandler := hystrix.NewStreamHandler()
hystrixStreamHandler.Start()
go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)
```

这个其实和`pprof`很像，在新的goroutine中开启一个web服务。

这块逻辑其实主要是一个loop

```
func (sh *StreamHandler) loop() {
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			circuitBreakersMutex.RLock()
			for _, cb := range circuitBreakers {
				sh.publishMetrics(cb)
				sh.publishThreadPools(cb.executorPool)
			}
			circuitBreakersMutex.RUnlock()
		case <-sh.done:
			return
		}
	}
}
```

每一秒钟进行数据收集一次,然后进行汇总计算

```
func (sh *StreamHandler) publishMetrics(cb *CircuitBreaker) error {
	now := time.Now()
	reqCount := cb.metrics.Requests().Sum(now)
	errCount := cb.metrics.DefaultCollector().Errors().Sum(now)
	errPct := cb.metrics.ErrorPercent(now)

	eventBytes, err := json.Marshal(&streamCmdMetric{
		Type:           "HystrixCommand",
		Name:           cb.Name,
		Group:          cb.Name,
		Time:           currentTime(),
		ReportingHosts: 1,

		RequestCount:       uint32(reqCount),
		ErrorCount:         uint32(errCount),
		ErrorPct:           uint32(errPct),
		CircuitBreakerOpen: cb.IsOpen(),

		RollingCountSuccess:            uint32(cb.metrics.DefaultCollector().Successes().Sum(now)),
		RollingCountFailure:            uint32(cb.metrics.DefaultCollector().Failures().Sum(now)),
		RollingCountThreadPoolRejected: uint32(cb.metrics.DefaultCollector().Rejects().Sum(now)),
		RollingCountShortCircuited:     uint32(cb.metrics.DefaultCollector().ShortCircuits().Sum(now)),
		RollingCountTimeout:            uint32(cb.metrics.DefaultCollector().Timeouts().Sum(now)),
		RollingCountFallbackSuccess:    uint32(cb.metrics.DefaultCollector().FallbackSuccesses().Sum(now)),
		RollingCountFallbackFailure:    uint32(cb.metrics.DefaultCollector().FallbackFailures().Sum(now)),

		LatencyTotal:       generateLatencyTimings(cb.metrics.DefaultCollector().TotalDuration()),
		LatencyTotalMean:   cb.metrics.DefaultCollector().TotalDuration().Mean(),
		LatencyExecute:     generateLatencyTimings(cb.metrics.DefaultCollector().RunDuration()),
		LatencyExecuteMean: cb.metrics.DefaultCollector().RunDuration().Mean(),

		// TODO: all hard-coded values should become configurable settings, per circuit

		RollingStatsWindow:         10000,
		ExecutionIsolationStrategy: "THREAD",

		CircuitBreakerEnabled:                true,
		CircuitBreakerForceClosed:            false,
		CircuitBreakerForceOpen:              cb.forceOpen,
		CircuitBreakerErrorThresholdPercent:  uint32(getSettings(cb.Name).ErrorPercentThreshold),
		CircuitBreakerSleepWindow:            uint32(getSettings(cb.Name).SleepWindow.Seconds() * 1000),
		CircuitBreakerRequestVolumeThreshold: uint32(getSettings(cb.Name).RequestVolumeThreshold),
	})
	if err != nil {
		return err
	}
	err = sh.writeToRequests(eventBytes)
	if err != nil {
		return err
	}

	return nil
}

```

将汇总后的数据写到`requests`里面，这个`requests`是一个map.键是`http.Request`,值就是我们的汇总数据

更多`dashboard`内容，可以看看[hystrix-dashboard](`https://github.com/Netflix-Skunkworks/hystrix-dashboard`)


有关流量控制的部分，`hystrix`中使用的是令牌方式的流控。

```
type executorPool struct {
	Name    string
	Metrics *poolMetrics
	Max     int
	Tickets chan *struct{}
}

func newExecutorPool(name string) *executorPool {
	p := &executorPool{}
	p.Name = name
	p.Metrics = newPoolMetrics(name)
	p.Max = getSettings(name).MaxConcurrentRequests

	p.Tickets = make(chan *struct{}, p.Max)
	for i := 0; i < p.Max; i++ {
		p.Tickets <- &struct{}{}
	}

	return p
}

func (p *executorPool) Return(ticket *struct{}) {
	if ticket == nil {
		return
	}

	p.Metrics.Updates <- poolMetricsUpdate{
		activeCount: p.ActiveCount(),
	}
	p.Tickets <- ticket
}

```

在小于最大并发请求的时候，可以直接从`Tickets`中获取`ticket`,用完之后返还。在`GoC`中我们也留意到，
如果没有更多的`ticket`可以获取到的话，返还`ticket`,上报`ErrMaxConcurrency`的错误。