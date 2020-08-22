# 1.Pool
首先看看Pool的应用和分析,在`net/http`包中，有一下Pool的使用案例
```
var (
	bufioReaderPool   sync.Pool
	bufioWriter2kPool sync.Pool
	bufioWriter4kPool sync.Pool
)

var copyBufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 32*1024)
		return &b
	},
}

func bufioWriterPool(size int) *sync.Pool {
	switch size {
	case 2 << 10:
		return &bufioWriter2kPool
	case 4 << 10:
		return &bufioWriter4kPool
	}
	return nil
}

func newBufioReader(r io.Reader) *bufio.Reader {
	if v := bufioReaderPool.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	// Note: if this reader size is ever changed, update
	// TestHandlerBodyClose's assumptions.
	return bufio.NewReader(r)
}

func putBufioReader(br *bufio.Reader) {
	br.Reset(nil)
	bufioReaderPool.Put(br)
}

func newBufioWriterSize(w io.Writer, size int) *bufio.Writer {
	pool := bufioWriterPool(size)
	if pool != nil {
		if v := pool.Get(); v != nil {
			bw := v.(*bufio.Writer)
			bw.Reset(w)
			return bw
		}
	}
	return bufio.NewWriterSize(w, size)
}

func putBufioWriter(bw *bufio.Writer) {
	bw.Reset(nil)
	if pool := bufioWriterPool(bw.Available()); pool != nil {
		pool.Put(bw)
	}
}
```

上面的例子，大家可能不太明白为啥这么用，咱们先来了解一下Pool的原始面貌：
```
type Pool struct {
	noCopy noCopy

	local     unsafe.Pointer // local fixed-size per-P pool, actual type is [P]poolLocal
	localSize uintptr        // size of the local array

	// New optionally specifies a function to generate
	// a value when Get would otherwise return nil.
	// It may not be changed concurrently with calls to Get.
    // New可选的指定一个函数，用于产生一个值，当调用Get方法时，否则返回nil.
    // 它可能不会与Get调用同时更改。
    // Pool的New函数通常只应返回指针类型，因为无需分配即可将指针放入返回接口值
	New func() interface{}
}
```
池是一组可以单独保存和检索的临时对象。存储在池中的任何项都可以在不通知的情况下随时自动删除。如果发生这种情况时，池持有惟一的引用，则项可能被释放。
一个池对于多个goroutine同时使用是安全的。
Pool的目的是缓存已分配但未使用的项目以供以后重用，从而减轻了垃圾收集器的压力。也就是说，它使构建高效，线程安全的空闲列表变得容易。 但是，它并不适合所有空闲列表。
池的适当用法是管理一组临时项目，这些临时项目在程序包的并发独立客户端之间静默共享并有可能被重用。Pool提供了一种跨许多客户端分摊分配开销的方法。

良好使用Pool的一个示例是fmt软件包，该软件包维护着动态大小的临时输出缓冲区存储。store在负载时滚动扩展，在空闲时收缩。
另一方面，作为短期对象的一部分维护的空闲列表不适用于Pool，因为在这种情况下开销无法很好地摊销。让这样的对象实现自己的空闲列表会更有效。
首次使用后不得复制池。

```
// Get从池中选择任意项，将其从池中移除，并将其返回给调用者。Get可以选择忽略池并将其视为空。
// 调用者不应假定传递给Put的值和Get返回的值之间有任何关系。
// 如果 p.New非空，Get将获取New的返回值，否则它将返回nil
func (p *Pool) Get() interface{}   
// Put adds x to the pool.
func (p *Pool) Put(x interface{})  
```
# 2. Once
再来看看`net/http`包中的`sync.Once`的应用。
```
// onceCloseListener wraps a net.Listener, protecting it from
// multiple Close calls.
type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (oc *onceCloseListener) Close() error {
	oc.once.Do(oc.close)
	return oc.closeErr
}

func (oc *onceCloseListener) close() { oc.closeErr = oc.Listener.Close() }
```

具体看看Once的样子
```
// Once is an object that will perform exactly one action.
type Once struct {
	m    Mutex
	done uint32
}

func (o *Once) Do(f func())
```
事实上，sync.Once是一个只执行一个操作的对象。多次调用其上的Do()方法，仅首次调用起效。 看一下Do方法的原型func (o *Once) Do(f func()) 因为对Do的调用只有在对f的调用返回时才会返回，如果f导致Do被调用，那么它就会死锁。 如果f函数异常，Do将任务其已经结束返回，后续对Do的调用将不再调用f。
看一下具体的例子
```
package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	onceBody := func() {
		fmt.Println("Only once")
	}
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}

```
onceBody只执行一次。

# 3.使用Mutex互斥锁实现的线程安全的map
```
package main

import (
        "fmt"
        "sync"
)

type SafeMap struct {
        m map[int]int
        sync.Mutex
}

func newSafeMap() *SafeMap {
        newMap := new(SafeMap)
        newMap.m = make(map[int]int)
        return newMap
}

func (s *SafeMap) Put(key int, value int) {
        s.Lock()
        s.m[key] = value
        s.Unlock()
}

func (s *SafeMap) Get(key int) (value int) {
        s.Lock()
        value = s.m[key]
        s.Unlock()
        return
}
func main() {
        var m sync.Map
        // var mn = make(map[int]int)
        ms := newSafeMap()
        go func(ms *SafeMap) {
                // mn[1] = 1
                m.LoadOrStore(1, 3)
                ms.Put(1, 5)
        }(ms)

        go func(ms *SafeMap) {
                // mn[1] = 2
                m.LoadOrStore(1, 4)
                ms.Put(1, 6)
        }(ms)

        go func(ms *SafeMap) {
                fmt.Println()
                fmt.Println("goroutine 1 slock: ", ms.Get(1)) //goroutine 1 slock:  0
                v, _ := m.Load(1)
                fmt.Println("goroutine 1: ", v) //goroutine 1:  4
        }(ms)
        go func(ms *SafeMap) {
                fmt.Println()
                fmt.Println("goroutine 2 slock: ", ms.Get(1))  //goroutine 2 slock:  0
                v, _ := m.Load(1)
                fmt.Println("goroutine 2: ", v) //goroutine 2:  4
        }(ms)
        v, _ := m.Load(1)
        fmt.Println("main goroutine: ", v) //main goroutine:  <nil>
        // fmt.Println(mn[1])
        fmt.Println("main goroutine slock: ", ms.Get(1)) //main goroutine slock:  5
}
```


卖票
```
package main

import (
    "fmt"
    "runtime"
    "time"
    "math/rand"
    "sync"
)
var total_tickets int32 = 10
var mutex = &sync.Mutex{}
func f(i int) {
    for {
        mutex.Lock()
        if total_tickets > 0 {
            time.Sleep(time.Duration(rand.Intn(5))* time.Millisecond)
            total_tickets--
            fmt.Println("id ", i, " tickets: ", total_tickets)
        } else {
            break
        }
        mutex.Unlock()
    }
}

func main() {
    runtime.GOMAXPROCS(4)
    rand.Seed(time.Now().Unix())
    for i:=0; i<5;i++ {
        go f(i)
    }
    var input string
    fmt.Scanln(&input)
    fmt.Println(total_tickets, "done")
}
```

```
package main
import (
    "fmt"
    "runtime"
    "time"
    "math/rand"
    "sync"
    "sync/atomic"
)
var total_tickets int32 = 10
var mutex = &sync.Mutex{}
func f(i int) {
    for {
        mutex.Lock()
        if total_tickets > 0 {
            time.Sleep(time.Duration(rand.Intn(5))* time.Millisecond)
            total_tickets--
            fmt.Println("id ", i, " tickets: ", total_tickets)
        } else {
            break
        }
        mutex.Unlock()
    }
}

func syncatomic() {
    var cnt uint32 = 0
    for i:=0;i<10;i++ {
        go func() {
            for i:=0;i<20;i++{
                time.Sleep(time.Millisecond)
                atomic.AddUint32(&cnt, 1)
            }
        }()
    }
    time.Sleep(time.Second)
    cntFinal := atomic.LoadUint32(&cnt)
    fmt.Println("cnt: ", cntFinal)
}
func main() {
    runtime.GOMAXPROCS(4)
    rand.Seed(time.Now().Unix())
    for i:=0; i<5;i++ {
        go f(i)
    }
    syncatomic()
    var input string
    fmt.Scanln(&input)
    fmt.Println(total_tickets, "done")
}
```
# 4.RWMutex读写互斥锁。
RWMutex是读取器/写入器互斥锁。 锁可以由任意数量的读取器或单个写入器持有。 RWMutex的零值是未锁定的互斥量。
第一次使用后，不得复制RWMutex。

如果goroutine拥有RWMutex进行读取，而另一个goroutine可能会调用Lock，则在释放初始读取锁之前，任何goroutine都不应期望能够获取读取锁。

特别是，这禁止了递归读取锁定。 这是为了确保锁最终可用。 被锁定的锁定调用将使新读者无法获得锁定。
# 5. WaitGroup在多协程任务管理上的使用
WaitGroup等待goroutine的集合完成。主goroutine调用Add来设置要等待的goroutine的数量。然后每一个goroutines运行并在完成时调用Done。
同时，可以使用Wait来阻塞，直到所有的goroutines完成。

首次使用后，不得复制WaitGroup。
```
// syncdo
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(1)

	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Starting Go Routines")
	go func() {
		defer wg.Done()

		for char := 'a'; char < 'a'+26; char++ {
			fmt.Printf("%c ", char)
		}
	}()

	go func() {
		defer wg.Done()

		for number := 1; number < 27; number++ {
			fmt.Printf("%d ", number)
		}
	}()

	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("\nTerminating Program")
}
```
Output:
```
Starting Go Routines
Waiting To Finish
a b c d e f g h i j k l m n o p q r s t u v w x y z 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26
Terminating Program
```
将第一个协程慢一拍执行
```
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

func main() {
    runtime.GOMAXPROCS(1)

    var wg sync.WaitGroup
    wg.Add(2)

    fmt.Println("Starting Go Routines")
    go func() {
        defer wg.Done()

        time.Sleep(1 * time.Microsecond)
        for char := ‘a’; char < ‘a’+26; char++ {
            fmt.Printf("%c ", char)
        }
    }()

    go func() {
        defer wg.Done()

        for number := 1; number < 27; number++ {
            fmt.Printf("%d ", number)
        }
    }()

    fmt.Println("Waiting To Finish")
    wg.Wait()

    fmt.Println("\nTerminating Program")
}
```
Output:
```
Starting Go Routines
Waiting To Finish
1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 a b c d e f g h i j k l m n o p q r s t u v w x y z 
Terminating Program
```
并行运行
```
package main

import (
    "fmt"
    "runtime"
    "sync"
)

func main() {
    runtime.GOMAXPROCS(2)

    var wg sync.WaitGroup
    wg.Add(2)

    fmt.Println("Starting Go Routines")
    go func() {
        defer wg.Done()

        for char := ‘a’; char < ‘a’+26; char++ {
            fmt.Printf("%c ", char)
        }
    }()

    go func() {
        defer wg.Done()

        for number := 1; number < 27; number++ {
            fmt.Printf("%d ", number)
        }
    }()

    fmt.Println("Waiting To Finish")
    wg.Wait()

    fmt.Println("\nTerminating Program")
}
```
Output:
```
Starting Go Routines
Waiting To Finish
1 2 a b c d e f g h i j k l m n o p q r s t u v 3 4 5 6 7 8 9 w x y z 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 
Terminating Program
```

# 6.条件变量
Cond实现了一个条件变量，它是goroutines等待或宣布事件发生的集合点。每个Cond都有一个关联的Locker L（通常是* Mutex或* RWMutex），在更改条件和调用Wait方法时必须将其保留。

第一次使用后，不得复制Cond。

```
type Cond struct {
	noCopy noCopy

	// L is held while observing or changing the condition
	L Locker

	notify  notifyList
	checker copyChecker
}

// NewCond returns a new Cond with Locker l.
func NewCond(l Locker) *Cond {
	return &Cond{L: l}
}

// 广播唤醒所有等待c的goroutine。
// 允许但不要求调用者在调用过程中持有c.L。
func (c *Cond) Broadcast() 
// 信号唤醒一个等待在c上的goroutine，如果有的话
func (c *Cond) Signal()

// Wait原子地解锁c.L并中止调用goroutine的执行。
// 稍后恢复执行后，等待锁定c.L，然后再返回。 与其他系统不同，等待不会返回，除非被广播或信号唤醒。
func (c *Cond) Wait()
// 因为在等待第一次恢复时c.L未被锁定，所以调用者通常无法假定等待返回时条件为真。
```