```
package main

import (
    "fmt"
)

func main(){
    ch := make(chan int)
    ch <- 1
    fmt.Println(<-ch)
}

```
有人问我，上面的程序能不能正常运行？

我当时就生气了，感觉自己的智商受到了极大的侮辱。他让我冷静冷静，分析分析为啥？

我话不多说，上面先来一顿操作
```
	0x0024 00036 (demo1.go:8)	PCDATA	$0, $0
	0x0024 00036 (demo1.go:8)	LEAQ	type.chan int(SB), AX
	0x002b 00043 (demo1.go:8)	PCDATA	$2, $0
	0x002b 00043 (demo1.go:8)	MOVQ	AX, (SP)
	0x002f 00047 (demo1.go:8)	MOVQ	$0, 8(SP)
	0x0038 00056 (demo1.go:8)	CALL	runtime.makechan(SB)
	0x003d 00061 (demo1.go:8)	PCDATA	$2, $1
	0x003d 00061 (demo1.go:8)	MOVQ	16(SP), AX
	0x0042 00066 (demo1.go:8)	PCDATA	$0, $1
	0x0042 00066 (demo1.go:8)	MOVQ	AX, "".ch+56(SP)
```
我们发现,创建chan实际上调用的是`runtime.makechan`。该函数提供了两个参数，第一个参数是`type.chan int`类型， 第二个参数是0，也就是说我们在make的时候并未提供第二个参数。返回值是一个`hchan`类型指针。之后对ch进行了赋值操作。
先看一眼`hchan`是个啥，混个眼熟
```
type hchan struct {
	qcount   uint           // 队列中的已存数据长度
	dataqsiz uint           // 循环队列的大小
	buf      unsafe.Pointer // 指向循环队列数组
	elemsize uint16         // 元素大小
	closed   uint32         // 当前channel是否关闭
	elemtype *_type // 元素类型
	sendx    uint   // 发送索引
	recvx    uint   // 接受索引
	recvq    waitq  // 等待接受的接受者的链表
	sendq    waitq  // 等待发送的发送者链表

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}
```
看一下该函数的原型
```
func makechan(t *chantype, size int) *hchan {
	elem := t.elem

	// compiler checks this but be safe.
	if elem.size >= 1<<16 {
		throw("makechan: invalid channel element type")
	}
	if hchanSize%maxAlign != 0 || elem.align > maxAlign {
		throw("makechan: bad alignment")
	}

	mem, overflow := math.MulUintptr(elem.size, uintptr(size))
	if overflow || mem > maxAlloc-hchanSize || size < 0 {
		panic(plainError("makechan: size out of range"))
	}

	// Hchan does not contain pointers interesting for GC when elements stored in buf do not contain pointers.
	// buf points into the same allocation, elemtype is persistent.
	// SudoG's are referenced from their owning thread so they can't be collected.
	// TODO(dvyukov,rlh): Rethink when collector can move allocated objects.
	var c *hchan
	switch {
	case mem == 0:
		// Queue or element size is zero.
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		// Race detector uses this location for synchronization.
		c.buf = c.raceaddr()
	case elem.kind&kindNoPointers != 0:
		// Elements do not contain pointers.
		// Allocate hchan and buf in one call.
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
		// Elements contain pointers.
		c = new(hchan)
		c.buf = mallocgc(mem, elem, true)
	}

	c.elemsize = uint16(elem.size)
	c.elemtype = elem
	c.dataqsiz = uint(size)

	if debugChan {
		print("makechan: chan=", c, "; elemsize=", elem.size, "; elemalg=", elem.alg, "; dataqsiz=", size, "\n")
	}
	return c
}
```
因为我们`make`时，提供的size大小为0，所以`mem==0`,`hchan`的`dataqsiz`也是0.在对ch变量赋值时，赋得是hchan的指针。

`ch <- 1`
这一步是向chan中写数据，该动作调用的是`runtime.chansend1`
```
0x0047 00071 (demo1.go:11)	PCDATA	$2, $0
	0x0047 00071 (demo1.go:11)	MOVQ	AX, (SP)
	0x004b 00075 (demo1.go:11)	PCDATA	$2, $1
	0x004b 00075 (demo1.go:11)	LEAQ	"".statictmp_0(SB), AX
	0x0052 00082 (demo1.go:11)	PCDATA	$2, $0
	0x0052 00082 (demo1.go:11)	MOVQ	AX, 8(SP)
	0x0057 00087 (demo1.go:11)	CALL	runtime.chansend1(SB)
```
看一下`"".statictmp_0`这个东西是什么，
```
"".statictmp_0 SRODATA size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
```
如果我们往chan中写的是11，这个东西的值又是什么
```
"".statictmp_0 SRODATA size=8
	0x0000 0b 00 00 00 00 00 00 00                          ........
```
就是一个十六进制表示的数字，继续看
```
func chansend1(c *hchan, elem unsafe.Pointer) {
	chansend(c, elem, true, getcallerpc())
}

/*
 * generic single channel send/recv
 * If block is not nil,
 * then the protocol will not
 * sleep but return if it could
 * not complete.
 *
 * sleep can wake up with g.param == nil
 * when a channel involved in the sleep has
 * been closed.  it is easiest to loop and re-run
 * the operation; we'll see that it's now closed.
 */
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	if c == nil {
		if !block {
			return false
		}
		gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}

	if debugChan {
		print("chansend: chan=", c, "\n")
	}

	if raceenabled {
		racereadpc(c.raceaddr(), callerpc, funcPC(chansend))
	}

	// Fast path: check for failed non-blocking operation without acquiring the lock.
	//
	// After observing that the channel is not closed, we observe that the channel is
	// not ready for sending. Each of these observations is a single word-sized read
	// (first c.closed and second c.recvq.first or c.qcount depending on kind of channel).
	// Because a closed channel cannot transition from 'ready for sending' to
	// 'not ready for sending', even if the channel is closed between the two observations,
	// they imply a moment between the two when the channel was both not yet closed
	// and not ready for sending. We behave as if we observed the channel at that moment,
	// and report that the send cannot proceed.
	//
	// It is okay if the reads are reordered here: if we observe that the channel is not
	// ready for sending and then observe that it is not closed, that implies that the
	// channel wasn't closed during the first observation.
	if !block && c.closed == 0 && ((c.dataqsiz == 0 && c.recvq.first == nil) ||
		(c.dataqsiz > 0 && c.qcount == c.dataqsiz)) {
		return false
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	lock(&c.lock)

	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
	}

	if sg := c.recvq.dequeue(); sg != nil {
		// Found a waiting receiver. We pass the value we want to send
		// directly to the receiver, bypassing the channel buffer (if any).
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
	}

	if c.qcount < c.dataqsiz {
		// Space is available in the channel buffer. Enqueue the element to send.
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			raceacquire(qp)
			racerelease(qp)
		}
		typedmemmove(c.elemtype, qp, ep)
		c.sendx++
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		c.qcount++
		unlock(&c.lock)
		return true
	}

	if !block {
		unlock(&c.lock)
		return false
	}

	// Block on the channel. Some receiver will complete our operation for us.
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	// No stack splits between assigning elem and enqueuing mysg
	// on gp.waiting where copystack can find it.
	mysg.elem = ep
	mysg.waitlink = nil
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.waiting = mysg
	gp.param = nil
	c.sendq.enqueue(mysg)
	goparkunlock(&c.lock, waitReasonChanSend, traceEvGoBlockSend, 3)
	// Ensure the value being sent is kept alive until the
	// receiver copies it out. The sudog has a pointer to the
	// stack object, but sudogs aren't considered as roots of the
	// stack tracer.
	KeepAlive(ep)

	// someone woke us up.
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	if gp.param == nil {
		if c.closed == 0 {
			throw("chansend: spurious wakeup")
		}
		panic(plainError("send on closed channel"))
	}
	gp.param = nil
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	mysg.c = nil
	releaseSudog(mysg)
	return true
}
```
第一个参数是上面创建的`*hchan`，第二个参数是`"".statictmp_0(SB)`地址。因为我们创建的是不带缓存的chan,其qcount和dataqsiz都是0，在对chan加锁后`lock(&c.lock)`,阻塞在channel上，等待接受者可以帮我们完成任务，因为当前代码无接受者goroutine,导致死锁问题。

再看来一个示例
```
package main

import (
    "fmt"
)

func main(){
    var ch chan int
    ch <- 11
    fmt.Println(<-ch)
}

```
看一下汇编后的情况
```
	0x0024 00036 (demo1.go:10)	PCDATA	$0, $1
	0x0024 00036 (demo1.go:10)	MOVQ	$0, "".ch+56(SP)
	0x002d 00045 (demo1.go:12)	MOVQ	$0, (SP)
	0x0035 00053 (demo1.go:12)	PCDATA	$2, $1
	0x0035 00053 (demo1.go:12)	LEAQ	"".statictmp_0(SB), AX
	0x003c 00060 (demo1.go:12)	PCDATA	$2, $0
	0x003c 00060 (demo1.go:12)	MOVQ	AX, 8(SP)
	0x0041 00065 (demo1.go:12)	CALL	`runtime.chansend1`(SB)
```

`runtime.chansend1`函数的第一个参数赋个0.

该实例运行也是异常的，异常的原因都是死锁，只不过
```
➜  channeltest go run demo1.go
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan send (nil chan)]:
main.main()
	demo1.go:12 +0x3a
exit status 2
```
这次死锁的原因是，向nil channel发送数据，也就是说，var声明的ch是一个nil chanl,在`chansend`方法中
```
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	if c == nil {
		if !block {
			return false
		}
		gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}
	...
}
```
检查到`c==nil`,之后，会将当前goroutine,gopark掉。park掉的原因是`waitReasonChanSendNilChan`
```
const (
    waitReasonChanSendNilChan                         // "chan send (nil chan)"
)
```

要想让上面的示例，成功运行，该怎么办呢？
```
func main(){
    // ch := make(chan int)
    ch:= make(chan int, 1)
    // var ch chan int
    ch <- 11
    fmt.Println(<-ch)  //11
}
```
加上chan的缓存。上面的示例就能正常运行了，很稳，为啥呢？仔细看看`chansend`
```
	if c.qcount < c.dataqsiz {
		// Space is available in the channel buffer. Enqueue the element to send.
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			raceacquire(qp)
			racerelease(qp)
		}
	`	typedmemmove(c.elemtype, qp, ep)`
		c.sendx++
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		c.qcount++
		unlock(&c.lock)
		return true
	}
```

上述代码节选自`runtime.chansend`函数，当我们首次将元素放到chan里面的时候，因为c.dataqsiz=1,c.qcount=0.
`	qp := chanbuf(c, c.sendx)`指向c.buf的第c.sendx个槽。`	typedmemmove(c.elemtype, qp, ep)`该函数将ep的值复制到qp.之后c.sendx++将发送索引加一。如果发送索引达到了循环队列的最大长度，则从头开始。

c的数据量自增1.解锁在`chansend`中加的锁，并返回true.
继续往下走
```
	0x0065 00101 (demo1.go:13)	PCDATA	$0, $0
	0x0065 00101 (demo1.go:13)	MOVQ	"".ch+56(SP), AX
	0x006a 00106 (demo1.go:13)	PCDATA	$2, $0
	0x006a 00106 (demo1.go:13)	MOVQ	AX, (SP)
	0x006e 00110 (demo1.go:13)	PCDATA	$2, $1
	0x006e 00110 (demo1.go:13)	LEAQ	""..autotmp_2+48(SP), AX
	0x0073 00115 (demo1.go:13)	PCDATA	$2, $0
	0x0073 00115 (demo1.go:13)	MOVQ	AX, 8(SP)
	0x0078 00120 (demo1.go:13)	CALL	runtime.chanrecv1(SB)
	0x007d 00125 (demo1.go:13)	MOVQ	""..autotmp_2+48(SP), AX
```
下面的接受数据`<-ch`,涉及的函数调用是`runtime.chanrecv1`,其参数为第一个参数是ch,第二个参数是需要存放读取出来的数据的寄存器地址`""..autotmp_2+48(SP)`。没有返回值。因为我们并没有申明变量来接受数据。

```
// entry points for <- c from compiled code
//go:nosplit
func chanrecv1(c *hchan, elem unsafe.Pointer) {
	chanrecv(c, elem, true)
}

// chanrecv 从channel中接受数据并将接受的数据写入到ep.
// ep可能为nil,这种情况下接受的数据将被忽略。
// 如果block == false并且没有元素可以使用，则返回(false, false)
// 否则，如果c已经关闭，零值的*ep并返回(true, false)
// 否则，使用元素填充ep并返回(true, true)
// 非空的ep必须指向堆或者调用者栈
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
	// raceenabled: don't need to check ep, as it is always on the stack
	// or is new memory allocated by reflect.

	if debugChan {
		print("chanrecv: chan=", c, "\n")
	}

	if c == nil {
		if !block {
			return
		}
		gopark(nil, nil, waitReasonChanReceiveNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}

	// Fast path: check for failed non-blocking operation without acquiring the lock.
	//
	// After observing that the channel is not ready for receiving, we observe that the
	// channel is not closed. Each of these observations is a single word-sized read
	// (first c.sendq.first or c.qcount, and second c.closed).
	// Because a channel cannot be reopened, the later observation of the channel
	// being not closed implies that it was also not closed at the moment of the
	// first observation. We behave as if we observed the channel at that moment
	// and report that the receive cannot proceed.
	//
	// The order of operations is important here: reversing the operations can lead to
	// incorrect behavior when racing with a close.
	if !block && (c.dataqsiz == 0 && c.sendq.first == nil ||
		c.dataqsiz > 0 && atomic.Loaduint(&c.qcount) == 0) &&
		atomic.Load(&c.closed) == 0 {
		return
	}

	var t0 int64
	if blockprofilerate > 0 {
		t0 = cputicks()
	}

	lock(&c.lock)

	if c.closed != 0 && c.qcount == 0 {
		if raceenabled {
			raceacquire(c.raceaddr())
		}
		unlock(&c.lock)
		if ep != nil {
			typedmemclr(c.elemtype, ep)
		}
		return true, false
	}

	if sg := c.sendq.dequeue(); sg != nil {
		// Found a waiting sender. If buffer is size 0, receive value
		// directly from sender. Otherwise, receive from head of queue
		// and add sender's value to the tail of the queue (both map to
		// the same buffer slot because the queue is full).
		recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true, true
	}

	if c.qcount > 0 {
		// Receive directly from queue
		qp := chanbuf(c, c.recvx)
		if raceenabled {
			raceacquire(qp)
			racerelease(qp)
		}
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}
		typedmemclr(c.elemtype, qp)
		c.recvx++
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
		c.qcount--
		unlock(&c.lock)
		return true, true
	}

	if !block {
		unlock(&c.lock)
		return false, false
	}

	// no sender available: block on this channel.
	gp := getg()
	mysg := acquireSudog()
	mysg.releasetime = 0
	if t0 != 0 {
		mysg.releasetime = -1
	}
	// No stack splits between assigning elem and enqueuing mysg
	// on gp.waiting where copystack can find it.
	mysg.elem = ep
	mysg.waitlink = nil
	gp.waiting = mysg
	mysg.g = gp
	mysg.isSelect = false
	mysg.c = c
	gp.param = nil
	c.recvq.enqueue(mysg)
	goparkunlock(&c.lock, waitReasonChanReceive, traceEvGoBlockRecv, 3)

	// someone woke us up
	if mysg != gp.waiting {
		throw("G waiting list is corrupted")
	}
	gp.waiting = nil
	if mysg.releasetime > 0 {
		blockevent(mysg.releasetime-t0, 2)
	}
	closed := gp.param == nil
	gp.param = nil
	mysg.c = nil
	releaseSudog(mysg)
	return true, !closed
}
```

因为c.qcount>0,直接从队列中接受值，整个流程结束，其实上述的示例是比较简单的。更复杂一些的应用场景有哪些呢？
- 多个goroutine的同步控制
- 结合select使用
- 任务队列


关于第一个应用场景示例
```
import (
	"time"
)
func worker(done chan bool) {
	time.Sleep(time.Second)
	done <- true
}
func main() {
	done := make(chan bool, 1)
	go worker(done)
	<-done
}
```
只有在工作goroutine完成了任务处理之后，整个程序才能顺利结束.


第二个应用场景示例

```
package main

import (
	"fmt"
)
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("fibonacci quit")
			return
		}
	}
}
func main() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}
```
第三个应用场景示例
```
package main

import (
	"fmt"
)
func main() {
	done := make(chan int, 10) // 带 10 个缓存

	// 开N个后台打印线程
	for i := 0; i < cap(done); i++ {
		go func(){
			fmt.Println("你好, 世界")
			done <- 1
		}()
	}

	// 等待N个后台线程完成
	for i := 0; i < cap(done); i++ {
		<-done
	}
}
```
当然了，上述的示例，有更方便的替换方案实现：
```
package main

import (
	"fmt"
	"sync"
)
func main() {
	var wg sync.WaitGroup

	// 开N个后台打印线程
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			fmt.Println("你好, 世界")
			wg.Done()
		}()
	}

	// 等待N个后台线程完成
	wg.Wait()
}
```

有关关闭channel的原则：
- 一个 sender，多个 receiver，由 sender 来关闭 channel，通知数据已发送完毕。

- 一旦 sender 有多个，可能就无法判断数据是否完毕了。这时候可以借助外部额外 channel 来做信号广播。这种做法类似于 done channel，或者 stop channel。

- 如果确定不会有 goroutine 在通信过程中被阻塞，也可以不关闭 channel，等待 GC 对其进行回收。


之所以出现类限制，主要是因为channel类型没办法多次关闭，在关闭的chan上调用close将panic,向close的chan发送数据也将panic.从close的chan上读取数据，将获取chan元素类型的零值，此种场景下要做第二个返回值的判断，以防止无法确切的知道到底是chan关闭了，还是从chan中读取的数就是零值。


至于有关chan作为参数时，是值传递还是引用传递，我们可以看一个小示例：
```
package main

import (
	"fmt"
)

func worker(stop chan bool) {
	
	for i:=0;i<10;i++ {
		fmt.Println("干活....")
	}
	stop <- true
}


func main() {	
	stop := make(chan bool)
	go worker(stop)
	<- stop
}
```

检点看一下汇编
```
	0x001d 00029 (main.go:17)	PCDATA	$0, $0
	0x001d 00029 (main.go:17)	LEAQ	type.chan bool(SB), AX
	0x0024 00036 (main.go:17)	PCDATA	$2, $0
	0x0024 00036 (main.go:17)	MOVQ	AX, (SP)
	0x0028 00040 (main.go:17)	MOVQ	$0, 8(SP)
	0x0031 00049 (main.go:17)	CALL	runtime.makechan(SB)
	0x0036 00054 (main.go:17)	PCDATA	$2, $1
	0x0036 00054 (main.go:17)	MOVQ	16(SP), AX
	0x003b 00059 (main.go:17)	PCDATA	$0, $1
	0x003b 00059 (main.go:17)	MOVQ	AX, "".stop+24(SP)
	0x0040 00064 (main.go:18)	MOVL	$8, (SP)
	0x0047 00071 (main.go:18)	PCDATA	$2, $2
	0x0047 00071 (main.go:18)	LEAQ	"".worker·f(SB), CX
	0x004e 00078 (main.go:18)	PCDATA	$2, $1
	0x004e 00078 (main.go:18)	MOVQ	CX, 8(SP)
	0x0053 00083 (main.go:18)	PCDATA	$2, $0
	0x0053 00083 (main.go:18)	MOVQ	AX, 16(SP)
	0x0058 00088 (main.go:18)	CALL	runtime.newproc(SB)
	0x005d 00093 (main.go:19)	PCDATA	$2, $1
	0x005d 00093 (main.go:19)	PCDATA	$0, $0
	0x005d 00093 (main.go:19)	MOVQ	"".stop+24(SP), AX
	0x0062 00098 (main.go:19)	PCDATA	$2, $0
	0x0062 00098 (main.go:19)	MOVQ	AX, (SP)
	0x0066 00102 (main.go:19)	MOVQ	$0, 8(SP)
	0x006f 00111 (main.go:19)	CALL	runtime.chanrecv1(SB)
	0x0074 00116 (main.go:20)	MOVQ	32(SP), BP
```
我们在新起一个goroutine进行执行业务处理是，使用`runtime.newproc`该函数有两个参数，第一个参数是siz=8,第二个参数是一个函数地址
```
// 创建一个新的g运行fn，带有size字节的参数。
// 将它放在g的等待队列中。
// 编译器将go语句转化成对该方法的调用。
// 
// Cannot split the stack because it assumes that the arguments
// are available sequentially after &fn; they would not be
// copied if a stack split occurred.
//go:nosplit
func newproc(siz int32, fn *funcval) {
	argp := add(unsafe.Pointer(&fn), sys.PtrSize)
	gp := getg()
	pc := getcallerpc()
	systemstack(func() {
		newproc1(fn, (*uint8)(argp), siz, gp, pc)
	})
}

```

本例中第三个参数是函数的参数，也就是我们`makechan`出来的`*hchan`,这也就是为何，我们可以在业务函数中可以修改chan的原因。


本系列文章：
- [我可能并不会使用golang之slice](https://juejin.im/post/5ec2030ee51d454de777380d)
- [我可能并不会使用golang之map](https://juejin.im/post/5ec3473be51d454d952bd7f0)
- [我可能并不会使用golang之channel](https://juejin.im/post/5ec4ccbce51d4578a2553e4c)
- [我可能并不会使用golang之goroutines](https://juejin.im/post/5ec72e0951882542f346e672)

有任何问题，欢迎留言
