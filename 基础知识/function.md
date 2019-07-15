在Go语言中，function被视为一种类型，描述了一组相同签名和返回值的函数。
既然是一种类型，那么我们就可以把function类型当作普通的类型来操作。
```
package main

import (
    "fmt"
)

type A func()


func main(){
   var a A = func(){
       fmt.Println("A func type")
   }
   a()
}
```

不支持默认参数，不支持重载，支持不定长参数，多返回值，命名返回值参数，匿名函数，闭包。
```
package main

import "fmt"

func main() {
    a := A(10)
    fmt.Println(a(2))
}

func A(x int) func(y int) int {
    return func (y int) int {
        return x + y
    }
}

```
上面就是一个简单的闭包的例子，我们可以简单分析一下，调用A(10)会返回一个函数，就是A函数中的一个匿名函数，在这个匿名函数中，我们可以使用自由变量x的值，虽然x的声明并不在该函数的作用域范围内。

#1. 定义函数
```
func main() {
    ......
}
```
```
func Add(a, b int) int {
    return a+b
}
```
```
func AddAll(a ...int) int {
    sum := 0
    for _, v := range a {
        sum += v
    }
    return v
}
```
不定长参数只能是参数列表中最后的一个，后面不能再出现其他的参数，比如：
```
func AddAll(a ...int, b, c stiring) int {
    ....
}
```
```
package main

import "fmt"

func main() {
    fmt.Println(add("sum", "all arguments", 1,2,3,4,4,5))
}

func add(b, c string, a ...int) int {
    fmt.Println(b, c)
    sum := 0
    for _, v := range a {
        sum += v
    }
    return sum
}
```
# 2.defer函数
该函数在函数体执行完成后，以逆顺序逐个执行，即便程序发生严重错误时也会执行，支持匿名函数调用，常用于资源清理，文件关闭。GO没有异常机制，但有panic/recover模式来处理错误。panic 可以在任何地方引发，但recover只有在defer调用但函数中有效。
```
package main

import "fmt"

func main() {
    fmt.Println("a")
    defer fmt.Println("b")
    defer fmt.Println("c")
}
```
Output:
```
a
c
b
```
```
package main

import "fmt"


func main() {
    for i:=0;i<3;i++ {
        defer func () {
            fmt.Println(i)
        }()
    }
}
```
Output
```
3
3
3
```
为什么会出现这种情况呢，我们可以发现，在defer调用但函数中，我们并没有显示但定义i变量，那么i变量就会向外层函数偷到i变量，但是结果拿到的是地址。循环遍历结束后，i地址所引用的值为3，所以打印结果也都是3.如果想结果是我们预期的那样，可以用下面的方法实现。
```
package main

import "fmt"


func main() {
    for i:=0;i<3;i++ {
        defer func (a int) {
            fmt.Println(a)
        }(i)
    }
}
```
Output:
```
2
1
0
```
```
package main

import "fmt"

func main() {
    var a, b int = 8 ,0
    fmt.Println(divFunc(a, b))
    c := A(10)
    fmt.Println(c(2))
}

func A(x int) func(y int) int {
    return func (y int) int {
        return x + y
    }
}

func divFunc (a, b int) int {
    defer func () {
        if err := recover(); err != nil{
            fmt.Println("panic in divFunc, the argument is a negative")
        }
    }()
    return a / b
}
```
Output:
```
0
12
```