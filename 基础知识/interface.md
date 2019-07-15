# 1.概述
C++,Java 使用"侵入式"接口，主要表现在实现类需要明确声明自己实现了某个接口。这种强制性的接口继承方式是面向对象编程思想发展过程中一个遭受相当多质疑的特性。

Go语言采用的是“非侵入式接口",Go语言的接口有其独到之处：只要类型T的公开方法完全满足接口I的要求，就可以把类型T的对象用在需要接口I的地方，所谓类型T的公开方法完全满足接口I的要求，也即是类型T实现了接口I所规定的一组成员。这种做法的学名叫做Structural Typing，有人也把它看作是一种静态的Duck Typing。
# 2.接口说明
接口类型，表达了固定的一个方法集合。一个接口变量可以存储任意实际值（非接口），只要这个值实现了接口的方法。

# 3.接口几点注意
- 接口只有方法声明，没有实现
- 接口可以嵌入到其他接口，或者结构体中
- 将对象赋值给接口时，会发生拷贝，而接口内部存储的是指这个复制品的指针，既无法修改复制品的状态，也无法获得指针。
- 只有当接口存储的类型和对象都为nil 时，接口才等于nil
- 接口调用不会做receiver的自动转换
- 接口同样支持匿名字段方法
- 接口可 实现类似面向对象中的多态
- 空接口：如果一个接口里面没有定义任何方法，那么它就是空接口，任意结构体都隐式地实现了空接口。
- 接口变量只包含两个指针字段，那么它的内存占用应该是 2 个机器字

```
package main

import (
    "fmt"
)

type Men interface {
    sayHi() string
    eatFood()
}

type Student struct {
    name string
    age int
}
func (st Student) sayHi() string {
    return fmt.Sprint(st.name, " say Hi!")

}
func (st Student) eatFood() {
    fmt.Println(st.name, "eat food.")
}

func shutdown(men Men) {
    if st, ok := men.(Student); ok {
        fmt.Println("student  " + st.name +" shutdown!")
        return
    }
    fmt.Println("Unknow")
}

func main() {
    var iface interface {}
    var a int = 21
    iface = a  //由于任何值都有零个或者多个方法，所以任何值都可以满足它。
    fmt.Println(iface)
    var men Men
    var st = Student{
        name: "zhaojunwei",
        age: 12,
    }
    men = st  //因为student 结构体定义了 Men接口的所有方法，所以该接口变量可以接受该实际值
    fmt.Println(men.sayHi())
    men.eatFood()
    shutdown(st)
}
```
类型判断更广泛的用法是switch type
```
func shutdown(men Men) {
    switch v := men.(type) {
        case Student:
            fmt.Println("student  " + v.name +" shutdown!")
    default:
            fmt.Println("Unknow")
    }
}
```
# 4.类型转换
可以把拥有超集的接口转化为子集的接口。
```
// intf
// intf
package main

import (
	"fmt"
)

type Men interface {
	eatFood()
	student
}

type student interface {
	learn()
}

type Student struct {
	name string
	age  int
}

func (st Student) learn() {
	fmt.Println("learn function")
}

func (st Student) eatFood() {
	fmt.Println(st.name, "eat food.")
}
func main() {
	var men = Student{
		name: "zhaojun",
		age:  14,
	} //该赋值的含义是将实现了接口功能的结构体赋值给接口变量
	fmt.Println(men)
	var st student
	st = men //子集接口被赋值超集接口
	st.learn()
}
```
# 5.组合
```
package main

import (
    "fmt"
)

type User struct {
    Name string
    Email string
}

func (u *User) Notify() error {
    fmt.Printf("User: Sending User Email to %s<%s>\n", u.Name, u.Email)
    return nil
}

type Notifier interface {
    Notify() error
}

func SendNotification(notify Notifier) error {
    return notify.Notify()
}

type Admin struct {
    User
    Level string
}

func (a *Admin) Notify() error {
    fmt.Printf("Admin: Sending Admin Email To %s<%s>\n", a.Name, a.Email)
    return nil
}

func main() {
    user := User{
        Name: "jane",
        Email: "jane@email.com",
    }
    admin := &Admin{
        User: user,
        Level: "super",
    }
    //SendNotification(&user)
    SendNotification(admin)
}
```
Output:
```
Admin: Sending Admin Email To jane<jane@email.com>
```
