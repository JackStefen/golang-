struct(结构体)也是一种聚合的数据类型，struct可以包含多个任意类型的值，这些值被称为struct的字段。
# 1.声明一个结构体
```
type Person struct {
    Name string
    Age int
    Birthday time.Time
}
```
结构体字段的首字母标识字段的权限访问，大写的包外可访问，小写的属于包内私有变量，仅在该结构体内部使用
# 2. 初始化
```
    ts := "2000-01-01 12:00:32"
    timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
    loc, _ := time.LoadLocation("Local")                            //重要：获取时区
    theTime, _ := time.ParseInLocation(timeLayout, ts, loc) //使用模板在对应时区转化为time.time类型
    fmt.Println(theTime)
    p := Person{Name:"zhanglinpeng",Age: 18, Birthday: theTime,} //最后的,不能少
    p2 := new(Persion) //new()函数来创建一个「零值」结构体，所有的字段都被初始化为相应类型的零值。返回的是结构体指针
    zerostruct := struct{}{} //空结构体,字节长度为0
```
# 3.匿名结构体
```
var apr = struct {
        Name string
        Age  int
    }{
        Name: "zhanglinpeng",
        Age: 13,
    }
    fmt.Println(apr)
```
# 4.匿名字段
```
var Student struct {
    string
    int
}
a := Student{"zhaoll", 19}
```
注意初始化时的顺序不能变

# 5.嵌套
```
type Person struct {
    Name string
    Age     int
    Contact  struct {
        Phone, Email, QQ string
    }
}
```
内嵌匿名结构体的初始化只能通过以下方式实现：
```
var apr = struct {
        Name string
        Age  int
        Contact struct {
            Phone, Email string
        }
    }{
        Name: "zhanglinpeng",
        Age: 13,
    }
    apr.Contact.Phone = "110"
    apr.Contact.Email = "110@qq.com"
    fmt.Println(apr)
```
# 6.值传递
结构体作为参数传递给函数时，也是值传递，如果想修改原始结构体，可以使用指针
```

func changeName(pr *Person) {
    pr.name = "zhanglinpeng"
}
...
changeName(&person)
fmt.Println(person)
```
# 7.属性值获取
可以使用点操作符获取属性值，点操作符还可以应用在struct指针上。
```
person := &person{Name:"weishihao",Age:14,...}
person.Name = "zhaoyuhao"
```
# 8.结构体比较
```
var Person struct {
    Name string
    Age int
} 
p1 := Person{Name:"zhaojj", Age: 14}
p2 := Person{Name:"zhaojj", Age: 14}
p3 := Person{Name:"zhaojj", Age: 15}
fmt.Println(p1 == p2) // true
fmt.Println(p1 == p3) // false
```
# 9.嵌入（“继承”）
```
type Person struct {
    Gender int
}

type teacher struct {
    Person
    Name string
    Age int
}

type student struct {
    Person
    Name string
    Age int
}

t1 := teacher{Name:"mayun", Age: 44, Person: Person{Gender: 0}}
s1 := student{Name: "zhaojj", Age: 12, Person: Person{Gender: 1}}
t1.Name = "yangmi"
s1.Gender = 0
```
我们发现，修改“继承”来的属性，可以直接点操作，而不用s1.Person.Gender = 0, 虽然这样做也是可行的。
如果在嵌入的结构体中存在同名的属性字段，那么在访问不同结构体中的属性字段时，需要指明，比如上述的那种访问方式。
如果同级别的嵌入结构体存在同名属性字段，就会报错。

# 10.方法
我们只需要在普通函数前面加个接受者（receiver，写在函数名前面的括号里面），这样编译器就知道这个函数（方法）属于哪个struct了。 需要注意的是，因为Go不支持函数重载，所以某个接收者（receiver）的某个方法只能对应一个函数，比如下面的就属于方法重复。
```
type A struct {
    Name string
}
type B struct {
    Name string
}

func (a A) print() {
    fmt.Println("function A")
}

func (b B) print() {
    fmt.Println("function B")
}

func (b B) print(i int) {
    fmt.Println("function B with argument")
}
```
针对A, B不同结构体print是不同的方法，所以可以和平相处，但是针对B结构体，存在两个同名的print方法，那么就会报错。
# 11.结构体指针
```
package main

import (
    "fmt"
)

type User struct {
    Id   int
    Name string
}

func (u User) displayId() {
    fmt.Println(u.Id)
}

func (u *User) displayName() {
    fmt.Println(u.Name)
}

func main() {
    us := User{Id: 1, Name: "zhao"}
    us.displayId() // 1
    us.displayName() // zhao
    us2 := &User{Id: 2, Name: "qian"}
    us2.displayId() // 2
    us2.displayName() // qian
}
```
可以看出，无论是结构体变量还是结构体指针变量，都是可以调用接受者不管是结构体还是结构体指针的方法。但是，传递给接口的时候会有所不同
```
package main

import (
    "fmt"
)

type DisplayInfo interface {
    displayId()
    displayName()
}

type User struct {
    Id   int
    Name string
}

func (u User) displayId() {
    fmt.Println(u.Id)
}

func (u *User) displayName() {
    fmt.Println(u.Name)
}


func DisplayUserInfo(ds DisplayInfo) {
    ds.displayId()
    ds.displayName()
}

func main() {
    us := User{Id: 1, Name: "zhao"}
    us.displayId()
    us.displayName()
    us2 := &User{Id: 2, Name: "qian"}
    us2.displayId()
    us2.displayName()

    us3 :=User{Id:3,Name:"sun"} // 如果这里使用&User{Id:3,Name:"sun"}是可以运行的
    DisplayUserInfo(us3) // cannot use us3 (type User) as type DisplayInfo in argument to DisplayUserInfo
    // User does not implement DisplayInfo (displayName method has pointer receiver)

}
```
错误信息中说，User类型没有实现DisplayInfo接口原因是displayName方法接受者是指针。但是为什么`us3=&User{Id:3,Name:"sun"}`可以呢？这是因为**接受者是指针类型的时候，说明指针指向的结构体实现了接口 接受者是值类型的时候，说明的是结构体本身实现了接口**.接受者是T的属于一个方法集，接受者是\*T的是另一个方法集，该方法及包含接受者是*T和T的。
# 12.接受者的类型决定了能否修改绑定的结构体
```
package main

import "fmt"

type A struct {
    Name string
}
type B struct {
    Name string
}

func (a A) print() {
    a.Name = "FuncA"
    fmt.Println("function A")
}

func (b *B) print() {
    b.Name = "FuncB"
    fmt.Println("function B")
}
func main() {
    a := A{}
    a.print()
    fmt.Println(a.Name)  // ""
    b := B{}
    b.print()
    fmt.Println(b.Name) // "FuncB"
}
```
# 13.方法绑定本身只能绑定包内的类型
如果是包外的类型，我们是无法绑定方法的。这就是为什么类型别名，无法将类型上的方法带到当前的包中的原因，比如，我在当前包中定义了一个int 类型，那么int类型上的方法，只有我们自己去实现
# 14. method value VS method expression
其实有两种调用方式，上面讲的那种官方管它叫method value，还有另一种调用方式，叫method expression
```
package main

import "fmt"

type A struct {
    Name string
}

func (a *A) print() {
    a.Name = "FuncA"
    fmt.Println("function A")
}

func main() {
    a := A{}
    a.print() // method value
    (*A).print(&a)  // method expression
    (&a).print()
}
```
# 15.方法权限
最后说下访问权限，因为Go是以大小写来区分是公有还是私有，**但都是针对包级别的**，
所以在包内所有的都能访问，而方法绑定本身只能绑定包内的类型，所以方法可以访问接收者所有成员。

# 16.demo
```
package main

import (
    "fmt"
    "unsafe"
)

type User struct {
    subject [10]byte
}

func main() {
    user := new(User)
    fmt.Println(user.subject) // [0 0 0 0 0 0 0 0 0 0]
    fmt.Println(len(user.subject)) // 10
    fmt.Println(reflect.TypeOf(user)) // *main.User
    fmt.Println(unsafe.Sizeof(struct{}{})) // 0
}
```
