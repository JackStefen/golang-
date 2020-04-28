接口类型，表达了固定的一个方法集合。一个接口变量可以存储任意实际值，只要这个值实现了接口的方法。
# 接口声明
```
type Men interface {
    sayHi() string
    eatFood()
}
var men Men
```
只要在接口中定义方法的声明，无需给出方法实现
# 接口特征：
- 接口只有方法声明，没有实现
- 接口可以嵌入到其他接口，或者结构体中
- 将对象赋值给接口时，会发生拷贝，而接口内部存储的是指这个复制品的指针，既无法修改复制品的状态，也无法获得指针。
- 只有当接口存储的类型和对象都为nil时，接口才等于nil
- 接口调用不会做receiver的自动转换
- 接口同样支持匿名字段方法
- 接口可实现类似面向对象中的多态
- 空接口：如果一个接口里面没有定义任何方法，那么它就是空接口，任意结构体都隐式地实现了空接口。
- 接口变量只包含两个指针字段，那么它的内存占用应该是2个指针字节长度
- 可以把拥有超集的接口转化为子集的接口。

# 类型判断
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
# demo

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

func main() {
    user := User{
        Name: "jane",
        Email: "jane@email.com",
    }
    admin := &Admin{
        User: user,
        Level: "super",
    }
    SendNotification(admin) //User: Sending User Email To jane<jane@email.com>
}
```
# 反射

每个interface变量都有一个对应pair，pair中记录了实际变量的值和类型:(value, type)
反射主要与Golang的interface类型相关（它的type是concrete type），只有interface类型才有反射一说。
`reflect`包实现了运行时反射，从而允许程序处理任意类型的对象。 典型的用法是使用静态类型`interface {}`获取值，并通过调用`TypeOf`来提取其动态类型信息，该类型将返回`Type`。调用`ValueOf()`返回一个代表运行时数据的`Value`。 零值采用一个类型，并返回一个表示该类型的零值的值。
看一下`reflect`包中具体的几个需要了解的结构体，首先，`Type`类型，该类型是一个**接口**类型，是Go类型的表示。 并非所有方法都适用于所有类型。 在每种方法的文档中都注明了限制（如果有）。 使用Kind方法先找出类型。调用特定于种类的方法。 调用不适合该类型的方法会导致运行时恐慌。 类型值是可比较的，例如==运算符，因此它们可用作字典的键。 如果两个Type值表示相同的类型，则它们相等。
```
type Type interface {

    // Aligin返回内存分配时该类型以字节为单位的对齐方式，方式适用于所有类型
    Align() int

    // FielAlign 返回该类型作为结构体字段时以字节为单位的对齐方式
    FieldAlign() int

    // Method返回该type方法集合中第i个方法，如果i的值不在[0,NumMethod())区间中则panic
    // For a non-interface type T or *T, the returned Method's Type and Func
    // fields describe a function whose first argument is the receiver.
    //
    // For an interface type, the returned Method's Type field gives the
    // method signature, without a receiver, and the Func field is nil.
    Method(int) Method

    // MethodByName returns the method with that name in the type's
    // method set and a boolean indicating if the method was found.
    //
    // For a non-interface type T or *T, the returned Method's Type and Func
    // fields describe a function whose first argument is the receiver.
    //
    // For an interface type, the returned Method's Type field gives the
    // method signature, without a receiver, and the Func field is nil.
    MethodByName(string) (Method, bool)

    // NumMethod返回类型的方法集合中导出的方法数量
    NumMethod() int

    // Name 返回其包中已定义类型的类型名称。
    // 对于其他（未定义）类型，它返回空字符串。
    Name() string

    // PkgPath returns a defined type's package path, that is, the import path
    // that uniquely identifies the package, such as "encoding/base64".
    // If the type was predeclared (string, error) or not defined (*T, struct{},
    // []int, or A where A is an alias for a non-defined type), the package path
    // will be the empty string.
    PkgPath() string

    // Size returns the number of bytes needed to store
    // a value of the given type; it is analogous to unsafe.Sizeof.
    Size() uintptr

    // String returns a string representation of the type.
    // The string representation may use shortened package names
    // (e.g., base64 instead of "encoding/base64") and is not
    // guaranteed to be unique among types. To test for type identity,
    // compare the Types directly.
    String() string

    // Kind 返回类型特定的kind类型.
    Kind() Kind

    // Implements reports whether the type implements the interface type u.
    Implements(u Type) bool

    // AssignableTo reports whether a value of the type is assignable to type u.
    AssignableTo(u Type) bool

    // ConvertibleTo reports whether a value of the type is convertible to type u.
    ConvertibleTo(u Type) bool

    // Comparable 标识这种类型的值是否是可进行比较的
    Comparable() bool

    // Methods applicable only to some types, depending on Kind.
    // The methods allowed for each kind are:
    //
    //  Int*, Uint*, Float*, Complex*: Bits
    //  Array: Elem, Len
    //  Chan: ChanDir, Elem
    //  Func: In, NumIn, Out, NumOut, IsVariadic.
    //  Map: Key, Elem
    //  Ptr: Elem
    //  Slice: Elem
    //  Struct: Field, FieldByIndex, FieldByName, FieldByNameFunc, NumField

    // Bits returns the size of the type in bits.
    // It panics if the type's Kind is not one of the
    // sized or unsized Int, Uint, Float, or Complex kinds.
    Bits() int

    // ChanDir returns a channel type's direction.
    // It panics if the type's Kind is not Chan.
    ChanDir() ChanDir

    // IsVariadic reports whether a function type's final input parameter
    // is a "..." parameter. If so, t.In(t.NumIn() - 1) returns the parameter's
    // implicit actual type []T.
    //
    // For concreteness, if t represents func(x int, y ... float64), then
    //
    //  t.NumIn() == 2
    //  t.In(0) is the reflect.Type for "int"
    //  t.In(1) is the reflect.Type for "[]float64"
    //  t.IsVariadic() == true
    //
    // IsVariadic panics if the type's Kind is not Func.
    IsVariadic() bool

    // Elem返回一个type的元素类型，如果type的Kind类型不是Array, Chan, Map, Ptr, or Slice则panic
    Elem() Type

    // Field返回一个结构体类型第i个字段
    // 如果type的Kind类型不是Struct则panic
    // 如果i不在[0, NumField())区间内则panic
    Field(i int) StructField

    // FieldByIndex返回与索引序列相对应的嵌套字段。
    // 等效于为每个索引i依次调用Field。
    // 如果类型的Kind不是Struct则panic
    FieldByIndex(index []int) StructField

    // FieldByName returns the struct field with the given name
    // and a boolean indicating if the field was found.
    FieldByName(name string) (StructField, bool)

    // FieldByNameFunc returns the struct field with a name
    // that satisfies the match function and a boolean indicating if
    // the field was found.
    //
    // FieldByNameFunc considers the fields in the struct itself
    // and then the fields in any embedded structs, in breadth first order,
    // stopping at the shallowest nesting depth containing one or more
    // fields satisfying the match function. If multiple fields at that depth
    // satisfy the match function, they cancel each other
    // and FieldByNameFunc returns no match.
    // This behavior mirrors Go's handling of name lookup in
    // structs containing embedded fields.
    FieldByNameFunc(match func(string) bool) (StructField, bool)

    // In返回函数类型第i个参数类型
    // 如果类型的Kind不是Func则panic
    // 如果i不在[0,NumIn())区间内则panic
    In(i int) Type

    // Key返回一个字典类型key的类型
    // 如果该type的Kind类型不是一个Map则panic
    Key() Type

    // It panics if the type's Kind is not Array.
    // Len返回一个数组类型的长度
    // 如果类型的Kind不是一个数组，则panic，
    Len() int

    // NumField返回一个结构体类型字段的数量
    // 如果该类型的Kind不是一个Struct则panic
    NumField() int

    // NumIn返回一个函数类型的参数数量
    // 如果类型的Kind不是一个Func则panic
    NumIn() int

    // NumOut返回一个函数类型的返回值数量
    // 如果类型的Kind不是一个Func则panic
    NumOut() int

    // Out返回函数类型第i个返回值的类型
    // 如果类型的Kind不是一个Func则panic
    // 如果i不在[0, NumOut())区间内则panic
    Out(i int) Type

    common() *rtype
    uncommon() *uncommonType
}
```
关于Kind,其类型为非负整数类型的别名，代表Type表示的特定类型，零值表示非法的Kind类型。
```
type Kind uint

const (
    Invalid Kind = iota
    Bool
    Int
    Int8
    Int16
    Int32
    Int64
    Uint
    Uint8
    Uint16
    Uint32
    Uint64
    Uintptr
    Float32
    Float64
    Complex64
    Complex128
    Array
    Chan
    Func
    Interface
    Map
    Ptr
    Slice
    String
    Struct
    UnsafePointer
)
```
另一个重要的数据结构为Value,Value是Go值的反射接口。 并非所有方法都适用于所有类型的值。 在每种方法的文档中都注明了限制（如果有）。 在调用特定于种类的方法之前，请使用Kind方法找出值的种类。 调用方法不适合该类型的类型会导致运行时出现恐慌。零值表示无值。其IsValid方法返回false，其Kind方法返回`Invalid`，其String方法返回` <invalid Value>`，而所有其他方法均会panic。大多数函数和方法从不返回无效值。
Value类型不同于Type,为一个结构体,这也决定了其和Type类型的使用上的不同。
- 调用它的`Interface()`方法会得到接口变量的真实内容，然后可以通过类型判断进行转换，转换为原有真实类型。
- 调用它的`Type()`方法等同于调用`TypeOf()`:`reflect.TypeOf(a) == reflect.ValueOf(a).Type()`
```
type Value struct {
    // typ持有Value类型值的类型
    typ *rtype

    // Pointer-valued data or, if flagIndir is set, pointer to data.
    // Valid when either flagIndir is set or typ.pointers() is true.
    ptr unsafe.Pointer

    // flag holds metadata about the value.
    // The lowest bits are flag bits:
    //  - flagStickyRO: obtained via unexported not embedded field, so read-only
    //  - flagEmbedRO: obtained via unexported embedded field, so read-only
    //  - flagIndir: val holds a pointer to the data
    //  - flagAddr: v.CanAddr is true (implies flagIndir)
    //  - flagMethod: v is a method value.
    // The next five bits give the Kind of the value.
    // This repeats typ.Kind() except for method values.
    // The remaining 23+ bits give a method number for method values.
    // If flag.kind() != Func, code can assume that flagMethod is unset.
    // If ifaceIndir(typ), code can assume that flagIndir is set.
    flag

    // A method value represents a curried method invocation
    // like r.Read for some receiver r. The typ+val+flag bits describe
    // the receiver r, but the flag's Kind bits say Func (methods are
    // functions), and the top bits of the flag give the method number
    // in r's type's method table.
}

```
reflect 包提供了一些基础反射方法，分别是 TypeOf() 和 ValueOf() 方法，分别用于获取变量的类型和值，定义如下：
```
// TypeOf返回i的反射类型
// 如果i为一个nil接口值，TypeOf返回nil
func TypeOf(i interface{}) Type {
    eface := *(*emptyInterface)(unsafe.Pointer(&i))
    return toType(eface.typ)
}

// ValueOf返回一个新的Value，初始化为存储在接口i中的具体值。ValueOf(nil)返回Value零值
func ValueOf(i interface{}) Value {
    if i == nil {
        return Value{}
    }

    // TODO: Maybe allow contents of a Value to live on the stack.
    // For now we make the contents always escape to the heap. It
    // makes life easier in a few places (see chanrecv/mapassign
    // comment below).
    escapes(i)

    return unpackEface(i)
}

// Indirect 返回参数v指向的value
// 如果v是一个nil指针，返回一个Value空值
// 如果v不是一个指针，原样返回v
func Indirect(v Value) Value {
    if v.Kind() != Ptr {
        return v
    }
    return v.Elem()
}
```
`net/rpc`中反射使用的例子
```
func (server *Server) register(rcvr interface{}, name string, useName bool) error {
    s := new(service)
    s.typ = reflect.TypeOf(rcvr)
    s.rcvr = reflect.ValueOf(rcvr)
    sname := reflect.Indirect(s.rcvr).Type().Name()
    if useName {
        sname = name
    }
    if sname == "" {
        s := "rpc.Register: no service name for type " + s.typ.String()
        log.Print(s)
        return errors.New(s)
    }
    if !isExported(sname) && !useName {
        s := "rpc.Register: type " + sname + " is not exported"
        log.Print(s)
        return errors.New(s)
    }
    s.name = sname

    // Install the methods
    s.method = suitableMethods(s.typ, true)

    if len(s.method) == 0 {
        str := ""

        // To help the user, see if a pointer receiver would work.
        method := suitableMethods(reflect.PtrTo(s.typ), false)
        if len(method) != 0 {
            str = "rpc.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
        } else {
            str = "rpc.Register: type " + sname + " has no exported methods of suitable type"
        }
        log.Print(str)
        return errors.New(str)
    }

    if _, dup := server.serviceMap.LoadOrStore(sname, s); dup {
        return errors.New("rpc: service already defined: " + sname)
    }
    return nil
}
```
反射的应用：
- 通过反射动态的调用方法
- 反射中匿名结构体的获取
- 反射中的字段和方法遍历
- 通过reflec.Value修改实际变量的值
```
package main

import (
    "fmt"
    "reflect"
)

type User struct {
    Id int
    Name string
    Age int
}

type Manager struct {
    User
    salary float64
}

func (u User) Hello(name string) {
    fmt.Println("Hello ", name, ", my name is ", u.Name)
}

func main() {
    u := User{1, "yangmi", 12}
    v := reflect.ValueOf(u)
    mv := v.MethodByName("Hello")
    args := []reflect.Value{reflect.ValueOf("joe")}
    //1.通过反射动态的调用方法
    mv.Call(args) // Hello  joe , my name is  yangmi
    
    var mg = Manager{User: u, salary: 100000.00}
    t := reflect.TypeOf(mg)
    tv := reflect.TypeOf(mg)
    // 2.反射中匿名结构体的获取
    // reflect.StructField{Name:"User", PkgPath:"", Type:(*reflect.rtype)(0x10dd900), Tag:"", Offset:0x0, Index:[]int{0}, Anonymous:true}
    fmt.Printf("%#v\n", t.FieldByIndex([]int{0}))
    //reflect.StructField{Name:"salary", PkgPath:"main", Type:(*reflect.rtype)(0x10c5aa0), Tag:"", Offset:0x20, Index:[]int{1}, Anonymous:false}
    fmt.Printf("%#v\n", t.FieldByIndex([]int{1}))
    // 3.反射中的字段和方法遍历
    for i := 0; i < t.NumField(); i++ {
        f := tu.Field(i)
        val := tv.Field(i)
        fmt.Printf("%6s: %v = %v\n", f.Name, f.Type, val)
    } // User: main.User = {1 yangmi 12}
      // salary: float64 = 100000
    
    for i := 0; i < t.NumMethod(); i++ {
        m := t.Method(i)
        fmt.Printf("%s - %v\n", m.Name, m.Type)
    } // Hello - func(main.Manager, string)
    var aint = 21
    avo := reflect.ValueOf(&aint).Elem()
    fmt.Println(avo.Type()) // int
    fmt.Println(avo.CanSet()) // true
    // 4.通过reflec.Value修改实际变量的值
    avo.SetInt(12)
    fmt.Println(avo) // 12
}
```
**当我们想通过反射来修改变量的值时，需要传入变量的指针**
# 指针
我们一般使用*T作为一个指针类型，标识一个指向类型为T变量的指针，为了安全考虑，
- 两个不同的指针类型不能转换，比如`*int`和`*int64`,
- 声明什么类型的指针，就赋值指向什么类型的。
- go中的指针不同于c语言，是不能进行数学运算的
- 不同类型的指针不能使用==或!=比较，但是可以与nil 作比较。
- 不同类型的指针变量不能相互赋值
否则，会报诸如此类的错误：
```
cannot use &b (type *int64) as type *int in assignment
```
# 使用指针的注意事项

- 一般情况下，不要通过指针分享内建类型的值.
- 通常，使用指针分享结构体类型的值，除非那个结构体类型实现的是代表私有类型的值。
- 引用类型像数组切片，字典，管道，接口，和函数值，我们很少使用指针来分享这些值。
- 通常，不要使用指针分享一个引用类型的值，除非你实现unMarshal类型的功能。

# unsafe包中的黑科技

unsafe 包用于 Go 编译器，在编译阶段使用。从名字就可以看出来，它是不安全的，官方并不建议使用。

```
type ArbitraryType int
type Pointer *ArbitraryType
```
该类型类似于c语言中的void *, 可以指向任意的类型。
unsafe 包还有其他三个函数：
```
// Sizeof takes an expression x of any type and returns the size in bytes
// of a hypothetical variable v as if v was declared via var v = x.
// The size does not include any memory possibly referenced by x.
// For instance, if x is a slice, Sizeof returns the size of the slice
// descriptor, not the size of the memory referenced by the slice.
// The return value of Sizeof is a Go constant.
func Sizeof(x ArbitraryType) uintptr

// Offsetof returns the offset within the struct of the field represented by x,
// which must be of the form structValue.field. In other words, it returns the
// number of bytes between the start of the struct and the start of the field.
// The return value of Offsetof is a Go constant.
func Offsetof(x ArbitraryType) uintptr

// Alignof takes an expression x of any type and returns the required alignment
// of a hypothetical variable v as if v was declared via var v = x.
// It is the largest value m such that the address of v is always zero mod m.
// It is the same as the value returned by reflect.TypeOf(x).Align().
// As a special case, if a variable s is of struct type and f is a field
// within that struct, then Alignof(s.f) will return the required alignment
// of a field of that type within a struct. This case is the same as the
// value returned by reflect.TypeOf(s.f).FieldAlign().
// The return value of Alignof is a Go constant.
func Alignof(x ArbitraryType) uintptr
```
Sizeof 返回类型 x 所占据的字节数，但不包含 x 所指向的内容的大小。
Offsetof 返回结构体成员在内存中的位置离结构体起始处的字节数，所传参数必须是结构体的成员。

- 同类型的指针之间不能相互转化，但确实需要转化的时候，还是可以做到的，使用unsafe.Pointer
```
package main

import (
    "fmt"
    "unsafe"
    "reflect"
)


func main(){
    var b int64 = 1
    bpointer := unsafe.Pointer(&b)
    fmt.Println(reflect.TypeOf(bpointer))  // unsafe.Pointer
    var c = (*int)(bpointer)
    fmt.Println(reflect.TypeOf(c)) // *int
    fmt.Println(*c) // 1
}
```
- 将unsafe.Pointer类型转化为uintptr类型，可进行指针的数学运算


```
    us := &User{Name: "zhangsan", Mobile:"18012345678"}
    usPointer:= unsafe.Pointer(us)
    fmt.Println(*(*string)(unsafe.Pointer(uintptr(usPointer) + unsafe.Offsetof(us.Mobile)))) // 18012345678
```
