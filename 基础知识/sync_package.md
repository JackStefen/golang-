# 1.线程安全的map
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

在一个逻辑处理器上并发运行协程
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