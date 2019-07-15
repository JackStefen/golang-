# 1.协程（goroutiue)
在go语言中，并发的实现是通过协程来实现的。语法上使用go关键字加函数即可让函数以并发的方式执行。
注意， go语言是以通信的方式来共享内存的，而不是通过共享内存来通信的。
协程的异常保护，子协程的异常如果自己不处理，会向上抛出给父协程。直至影响主协程挂掉，程序停止。举个简单的例子:
```
package main

import (
    "fmt"
)

func main() {
    c := make(chan bool)
    go func () {
        fmt.Println("go func....")
        c <- true
        close(c)
        c <- true
    }()
    for v:= range c {
        fmt.Println(v)
    }
    fmt.Println("main done...")
}
//==================
go func....
panic: send on closed channel

goroutine 17 [running]:
main.main.func1(0xc000084000)
	/Users/xx/workspace/src/just.for.test/goroutinetest/demo.go:13 +0xbc
created by main.main
	/Users/xx/workspace/src/just.for.test/goroutinetest/demo.go:9 +0x5c
exit status 2
```
使用defer func()调用recover方法处理子协程中的异常，不让其抛出给父协程
```
package main

import (
    "fmt"
)

func main() {
    c := make(chan bool)
    go func () {
        defer func() {
            if err:= recover(); err != nil {
                  fmt.Println(err)
            }
        }()
        fmt.Println("go func....")
        c <- true
        close(c)
        c <- true
    }()
    for v:= range c {
        fmt.Println(v)
    }
    fmt.Println("main done...")
}
//===================
go func....
send on closed channel
true
main done...
```

# 2.管道（channel）
上面的例子中，我们实际上已经用到了channel他是各协程通信的方式，包含有缓冲和无缓存两种，无缓存的话是阻塞方式 的，就是说，如果读空管道会阻塞直到有数据写入，如果写入非空管道，就会阻塞，直到有数据被读出。有缓冲就像异步的了，只要管道不满就可以一直写，只要管道不为空就可以随时读。
- 创建channel使用make方法，是否有第二个参数标识channel是否带缓存
- 读channel使用<- channel语法，写channel使用channel<-语法，
- 向已满的channel写数据将阻塞，向已空的channel读数据将阻塞
- 向以关闭的channel读数据，将得到元素类型的零值，向已关闭的channel写数据将抛出异常
- 在单读单写，单写多读的场景下，最好只由写channel任务关闭channel
- 在多写场景下，需要使用第三方的协程来管理，这个协程等其他所有写协程都完成后，再关闭channel。
结合select，可实现循环读取写入管道
```
package main

import (
    "fmt"
    "time"
    "math/rand"
)

func main() {
    channel := make(chan string)
    rand.Seed(time.Now().Unix())
    go func() {
        cnt := rand.Intn(10)
        fmt.Println("message cnt: ", cnt)
        for i:=0;i<cnt;i++ {
            channel <- fmt.Sprintf("message-%2d", i)
        }
        close(channel)
    }()
    var more bool = true
    var msg string
    for more {
        select{
            case msg, more = <- channel:
                if more {
                    fmt.Println(msg)
                } else {
                    fmt.Println("channel closed")
                }
        }
    }
}
```
Output:
```
message cnt:  8
message- 0
message- 1
message- 2
message- 3
message- 4
message- 5
message- 6
message- 7
channel closed
```
下面看一个问题，下面的例子之所以会出现死锁的情况，是因为range结束的条件是遍历完成channel关闭，如果channel未关闭，其将一直尝试从channel中读取数据，其效果如同向一个空的channel读数据一样
```
package main

import (
    "fmt"
)

func putchan(ch chan int) {
        ch <- 2
}

func main(){
    ch := make(chan int,12)
    for i:=0;i<5;i++ {
        go putchan(ch)
    }
    for i:= range ch {
        fmt.Println(i)
    }
}
```
结果：
```
2
2
2
2
2
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan receive]:
main.main()
	/Users/wscn/Desktop/CodeView/chanwirtedo.go:22 +0x10b
exit status 2
wscndeMac-mini:CodeView wscn$ vi chanwirtedo.go
```
[解决方法](https://stackoverflow.com/questions/34572122/fatal-error-all-goroutines-are-asleep-deadlock)
```
package main

import (
    "fmt"
)

func putchan(ch chan int) {
        ch <- 2
}

func main(){
    ch := make(chan int,12)
    for i:=0;i<5;i++ {
        go putchan(ch)
    }
    for i:=0;i<5;i++ {
        fmt.Println(<-ch)
    }
}
```
上面这种，有多少数据就读多少数据的做法，就不会产生异常。
# 3.通过不同方式来处理多个channel 通信

- 带缓存的channel

```
package main
import (
        "fmt"
        "runtime"
)
// 从 1 至 1 亿循环叠加，并打印结果。
func print(c chan bool, n int) {
        x := 0
        for i := 1; i <= 100000000; i++ {
                x += i
        }
        fmt.Println(n, x)
        c <- false
}

func main() {
    // 使用多核运行程序
        runtime.GOMAXPROCS(runtime.NumCPU())
        //c := make(chan bool, 10)
        c := make(chan bool)
        for i := 0; i < 10; i++ {
                go print(c, i)
        }
        for i := 0; i < 10; i++ {
            fmt.Println(<-c)
        }
        //<-c
        fmt.Println("DONE.")
}
// =================
4 5000000050000000
false
0 5000000050000000
false
9 5000000050000000
false
1 5000000050000000
false
7 5000000050000000
false
5 5000000050000000
false
2 5000000050000000
false
8 5000000050000000
false
6 5000000050000000
false
3 5000000050000000
false
DONE.
```
- 通过waitgroup

```
package main

import (
    "fmt"
    "runtime"
    "sync"
)

func print(wg *sync.WaitGroup, n int) {
    defer wg.Done()
    x := 0
    for i:=1;i<=100000;i++ {
        x += i
    }
    fmt.Println(n,  x)
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    //wg := sync.WaitGroup{}
    var wg sync.WaitGroup
    wg.Add(10)
    for i:=0;i<10;i++ {
        go print(&wg, i)
    }
    wg.Wait()
    fmt.Println("Done.")
}
//= ================
0 5000050000
1 5000050000
9 5000050000
5 5000050000
7 5000050000
3 5000050000
8 5000050000
4 5000050000
6 5000050000
2 5000050000
Done.
```
# 4.多channel复用select
多个channel写，读的时候效果如同，多个channel汇聚成一个channel，单从这个channel读数据即可。
```
package main

import (
    "fmt"
    "time"
    "math/rand"
)

func main() {
    channel := make(chan string)
    rand.Seed(time.Now().Unix())
    go func() {
        cnt := rand.Intn(10)
        fmt.Println("message cnt: ", cnt)
        for i:=0;i<cnt;i++ {
            channel <- fmt.Sprintf("message-%2d", i)
        }
        close(channel)
    }()
    channel2 := make(chan string)
    rand.Seed(time.Now().Unix())
    go func() {
        cnt := rand.Intn(10)
        fmt.Println("message cnt2: ", cnt)
        for i:=0;i<cnt;i++ {
            channel2 <- fmt.Sprintf("message2-%2d", i)
        }
        close(channel2)
    }()
    for {
        select{
            case msg, more := <- channel:
                if more {
                    fmt.Println(msg)
                } else {
                    fmt.Println("channel closed")
                    return
                }
            case msg2,more2 := <- channel2:
               if more2 {
                   fmt.Println(msg2)
               } else {
                   fmt.Println("channel2 colsed")
                   //break
                   // break在这里没有用
                   return
               }
            default:
                fmt.Println("read from default")
        }
    }
}
//=======================
message cnt:  6
message cnt2:  8
read from default
message- 0
message2- 0
message- 1
message2- 1
message- 2
message2- 2
read from default
message2- 3
read from default
message2- 4
read from default
message- 3
message2- 5
message- 4
message2- 6
message2- 7
message- 5
channel closed
```
