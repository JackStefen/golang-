struct(结构体)也是一种聚合的数据类型，struct可以包含多个任意类型的值，这些值被称为struct的字段。
# 1.声明一个结构体
```
type Person struct {
    Name string
    Age int
    Birthday time.Time
}
```
结构体字段的首字母标识字段的权限访问，大写的包外可访问，小写的属于保内私有变量，仅在该结构体内容使用
# 2. 初始化
```
    ts := "2000-01-01 12:00:32"
    timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
    loc, _ := time.LoadLocation("Local")                            //重要：获取时区
    theTime, _ := time.ParseInLocation(timeLayout, ts, loc) //使用模板在对应时区转化为time.time类型
    fmt.Println(theTime)
    p := Person{Name:"zhanglinpeng",Age: 18, Birthday: theTime,} //最后的,不能少
    p2 := new(Persion) //new() 函数来创建一个「零值」结构体，所有的字段都被初始化为相应类型的零值。返回的是结构体指针
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
注意顺序不能变

# 5嵌套
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
# 6. 值传递
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
# 11.接受者的类型决定了能否修改绑定的结构体
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
    fmt.Println(a.Name)  //""
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
    a.print()
    fmt.Println(a.Name)  //"FuncA"
    (*A).print(&a)
}
```
# 15.方法权限
最后说下访问权限，因为Go是以大小写来区分是公有还是私有，**但都是针对包级别的**，
所以在包内所有的都能访问，而方法绑定本身只能绑定包内的类型，所以方法可以访问接收者所有成员。

# 16.一个有趣的例子
```
package main

import (
    "fmt"
)

type User struct {
    subject [100]byte
}

func main() {
    var a [3]int
    fmt.Println(a)
    user := new(User)
    fmt.Println(user.subject)
    fmt.Println(len(user.subject))
    fmt.Println(reflect.TypeOf(user))
}
```
Output:
```
[0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
100
*main.User
```
可以看出，作为结构体的一个字段，new结构体的时候，会同时初始化其字段为相应类型的零值，然后返回其结构体指针

# json tag
```
package main

import (
     "fmt"
     "encoding/json"
)

type UserInfo struct {
	UsrId int64      `json:"user_id,omitempty"`
	NickName string  `json:"nickname"`
	Address string   `json:"-" `
}

func main () {
    var u UserInfo = UserInfo{
    	NickName: "zhanglinpeng",
    	Address: "",
    }
    rl, err:= json.Marshal(u)
    if err != nil {
    	fmt.Println("json marshal error: ", err)
    }
    // os.Stdout.Write(rl)
    fmt.Println(string(rl))
    var vlr UserInfo
    err1 := json.Unmarshal(rl, &vlr)
    if err != nil {
    	fmt.Println("json unmarshal error: ", err1)
    }
    fmt.Printf("%v\n", vlr)
}
```
Output:
```
{"nickname":"zhanglinpeng"}
{0 zhanglinpeng }
```
注意事项，**在标签中`json:"nickname"` 外层符号为键盘Tab健上方的键。json冒号和后面的字符串之间不能有空格, omitempty和逗号之间也不能有空格，总之在标签中能不用空格就不用空格。**
- UsrId字段没有显示的原因是，我们在实例化结构体的时候未实例化UsrId字段，那么json结果输出中就没有这个字段，但是实际上，他的默认值是0，是存在的，这也是为什么在最后的转化回去的输出结果中该字段的值为0.
- Address字段无论设置与否，结构都不会显示，因为标签设置为了`-`。

我们在`-`字符的后面加了一个逗号。结果就会输出Address字段，字段名为"-"

```
package main

import (
     "fmt"
     "encoding/json"
)

type UserInfo struct {
	UsrId int64      `json:",omitempty"`
	NickName string  `json:"nickname"`
	Address string   `json:"-," `
}

func main () {
    var u UserInfo = UserInfo{
    	NickName: "zhanglinpeng",
    	Address: "",
    }
    rl, err:= json.Marshal(u)
    if err != nil {
    	fmt.Println("json marshal error: ", err)
    }
    // os.Stdout.Write(rl)
    fmt.Println(string(rl))
    var vlr UserInfo
    err1 := json.Unmarshal(rl, &vlr)
    if err != nil {
    	fmt.Println("json unmarshal error: ", err1)
    }
    fmt.Printf("%v\n", vlr)
}
```
Output:
```
{"nickname":"zhanglinpeng","-":"shanghai"}
{0 zhanglinpeng shanghai}
```
如果在结构体实例化中没有实例化字段就会跳过，注意omitempty前面有逗号。