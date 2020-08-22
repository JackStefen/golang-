# 1.函数概述
在Go语言中，function被视为一种类型,既然是一种类型，那么我们就可以把function类型当作普通的类型来操作。
函数特征：
- 不支持默认参数
- 不支持重载
- 支持不定长参数, 不定长参数只能是参数列表中最后的一个，后面不能再出现其他的参数
- 多返回值
- 命名返回值参数
- 匿名函数
- 闭包。
```
package main

func main() {
    a := 1
    b := 1
    _ = sum(a, b)
}

func sum(a int, b int) (c int) {
    c =  a+ b
    return c
}
```
从汇编层面理解上面的函数逻辑`go tool compile  -S -N -l fundemo1.go`(不优化,不内联)：
```
    "".main STEXT size=87 args=0x0 locals=0x30
    0x0000 00000 (fundemo1.go:5)    TEXT    "".main(SB), ABIInternal, $48-0 // 函数栈空间为48字节，参数和返回值大小为0，栈不分裂
    0x0000 00000 (fundemo1.go:5)    MOVQ    (TLS), CX
    0x0009 00009 (fundemo1.go:5)    CMPQ    SP, 16(CX)   // 判断栈空间是否需要分配更多的栈
    0x000d 00013 (fundemo1.go:5)    JLS 80
    0x000f 00015 (fundemo1.go:5)    SUBQ    $48, SP    // 分配栈空间，48个字节
    0x0013 00019 (fundemo1.go:5)    MOVQ    BP, 40(SP) // 将基址指针存储到栈上(40)SP，一个字节
    0x0018 00024 (fundemo1.go:5)    LEAQ    40(SP), BP // 把 40(SP) 的地址放到 BP 里。
    0x001d 00029 (fundemo1.go:5)    FUNCDATA    $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x001d 00029 (fundemo1.go:5)    FUNCDATA    $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x001d 00029 (fundemo1.go:5)    FUNCDATA    $3, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x001d 00029 (fundemo1.go:6)    PCDATA  $2, $0
    0x001d 00029 (fundemo1.go:6)    PCDATA  $0, $0
    0x001d 00029 (fundemo1.go:6)    MOVQ    $1, "".a+32(SP)  // 定义局部变量a
    0x0026 00038 (fundemo1.go:7)    MOVQ    $1, "".b+24(SP)  // 定义局部变量b
    0x002f 00047 (fundemo1.go:8)    MOVQ    "".a+32(SP), AX
    0x0034 00052 (fundemo1.go:8)    MOVQ    AX, (SP)      // 将变量a的值放在sp开启的一个字节里面
    0x0038 00056 (fundemo1.go:8)    MOVQ    $1, 8(SP)     // 第二个参数放在第二个字节里面
    0x0041 00065 (fundemo1.go:8)    CALL    "".sum(SB)   // 调用sum函数
    0x0046 00070 (fundemo1.go:9)    MOVQ    40(SP), BP  // 来恢复栈基址指针
    0x004b 00075 (fundemo1.go:9)    ADDQ    $48, SP // 销毁已经失去作用的48字节空间。
    0x004f 00079 (fundemo1.go:9)    RET
    0x0050 00080 (fundemo1.go:9)    NOP
    0x0050 00080 (fundemo1.go:5)    PCDATA  $0, $-1
    0x0050 00080 (fundemo1.go:5)    PCDATA  $2, $-1
    0x0050 00080 (fundemo1.go:5)    CALL    runtime.morestack_noctxt(SB)
    0x0055 00085 (fundemo1.go:5)    JMP 0
    0x0000 65 48 8b 0c 25 00 00 00 00 48 3b 61 10 76 41 48  eH..%....H;a.vAH
    0x0010 83 ec 30 48 89 6c 24 28 48 8d 6c 24 28 48 c7 44  ..0H.l$(H.l$(H.D
    0x0020 24 20 01 00 00 00 48 c7 44 24 18 01 00 00 00 48  $ ....H.D$.....H
    0x0030 8b 44 24 20 48 89 04 24 48 c7 44 24 08 01 00 00  .D$ H..$H.D$....
    0x0040 00 e8 00 00 00 00 48 8b 6c 24 28 48 83 c4 30 c3  ......H.l$(H..0.
    0x0050 e8 00 00 00 00 eb a9                             .......
    rel 5+4 t=16 TLS+0
    rel 66+4 t=8 "".sum+0
    rel 81+4 t=8 runtime.morestack_noctxt+0
"".sum STEXT nosplit size=25 args=0x18 locals=0x0
    0x0000 00000 (fundemo1.go:11)   TEXT    "".sum(SB), NOSPLIT|ABIInternal, $0-24 //参数大小为24字节
    0x0000 00000 (fundemo1.go:11)   FUNCDATA    $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x0000 00000 (fundemo1.go:11)   FUNCDATA    $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x0000 00000 (fundemo1.go:11)   FUNCDATA    $3, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x0000 00000 (fundemo1.go:11)   PCDATA  $2, $0
    0x0000 00000 (fundemo1.go:11)   PCDATA  $0, $0
    0x0000 00000 (fundemo1.go:11)   MOVQ    $0, "".c+24(SP)   // 定义变量c
    0x0009 00009 (fundemo1.go:12)   MOVQ    "".a+8(SP), AX    // 将a放在寄存器中
    0x000e 00014 (fundemo1.go:12)   ADDQ    "".b+16(SP), AX   // 执行加操作
    0x0013 00019 (fundemo1.go:12)   MOVQ    AX, "".c+24(SP)
    0x0018 00024 (fundemo1.go:13)   RET
    0x0000 48 c7 44 24 18 00 00 00 00 48 8b 44 24 08 48 03  H.D$.....H.D$.H.
    0x0010 44 24 10 48 89 44 24 18 c3                       D$.H.D$..
```

![](https://user-gold-cdn.xitu.io/2020/1/20/16fc176069535e82?w=1046&h=774&f=png&s=280194)

```
package main

//import "fmt"

func main() {
    a := A(10)
    //fmt.Println(a(2))
    _ = a(2)
}

func A(x int) func(y int) int {
    return func (y int) int {
        return x + y
    }
}
```

上面就是一个简单的闭包的例子，我们可以简单分析一下，调用A(10)会返回一个函数，就是A函数中的一个匿名函数，在这个匿名函数中，我们可以使用自由变量x的值，虽然x的声明并不在该函数的作用域范围内。

```
"".main STEXT size=82 args=0x0 locals=0x20
    0x0000 00000 (demo5.go:5)   TEXT    "".main(SB), ABIInternal, $32-0
    0x0000 00000 (demo5.go:5)   MOVQ    (TLS), CX
    0x0009 00009 (demo5.go:5)   CMPQ    SP, 16(CX)
    0x000d 00013 (demo5.go:5)   JLS 75
    0x000f 00015 (demo5.go:5)   SUBQ    $32, SP
    0x0013 00019 (demo5.go:5)   MOVQ    BP, 24(SP)
    0x0018 00024 (demo5.go:5)   LEAQ    24(SP), BP
    0x001d 00029 (demo5.go:5)   FUNCDATA    $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x001d 00029 (demo5.go:5)   FUNCDATA    $1, gclocals·2a5305abe05176240e61b8620e19a815(SB)
    0x001d 00029 (demo5.go:5)   FUNCDATA    $3, gclocals·ebb0e8ce1793da18f0378b883cb3e122(SB)
    0x001d 00029 (demo5.go:6)   PCDATA  $2, $0
    0x001d 00029 (demo5.go:6)   PCDATA  $0, $0
    0x001d 00029 (demo5.go:6)   MOVQ    $10, (SP)         // 调用A函数的参数
    0x0025 00037 (demo5.go:6)   CALL    "".A(SB)          // 调用A函数
    0x002a 00042 (demo5.go:6)   PCDATA  $2, $1
    0x002a 00042 (demo5.go:6)   MOVQ    8(SP), DX         // A函数的返回值
    0x002f 00047 (demo5.go:6)   MOVQ    DX, "".a+16(SP)   // 变量a的定义，值为A(10)
    0x0034 00052 (demo5.go:8)   MOVQ    $2, (SP)          // a(2)函数的参数
    0x003c 00060 (demo5.go:8)   MOVQ    (DX), AX
    0x003f 00063 (demo5.go:8)   PCDATA  $2, $0 
    0x003f 00063 (demo5.go:8)   CALL    AX               //调用a函数a(2)
    0x0041 00065 (demo5.go:9)   MOVQ    24(SP), BP
    0x0046 00070 (demo5.go:9)   ADDQ    $32, SP
    0x004a 00074 (demo5.go:9)   RET
    0x004b 00075 (demo5.go:9)   NOP
    0x004b 00075 (demo5.go:5)   PCDATA  $0, $-1
    0x004b 00075 (demo5.go:5)   PCDATA  $2, $-1
    0x004b 00075 (demo5.go:5)   CALL    runtime.morestack_noctxt(SB)
    0x0050 00080 (demo5.go:5)   JMP 0
    0x0000 65 48 8b 0c 25 00 00 00 00 48 3b 61 10 76 3c 48  eH..%....H;a.v<H
    0x0010 83 ec 20 48 89 6c 24 18 48 8d 6c 24 18 48 c7 04  .. H.l$.H.l$.H..
    0x0020 24 0a 00 00 00 e8 00 00 00 00 48 8b 54 24 08 48  $.........H.T$.H
    0x0030 89 54 24 10 48 c7 04 24 02 00 00 00 48 8b 02 ff  .T$.H..$....H...
    0x0040 d0 48 8b 6c 24 18 48 83 c4 20 c3 e8 00 00 00 00  .H.l$.H.. ......
    0x0050 eb ae                                            ..
    rel 5+4 t=16 TLS+0
    rel 38+4 t=8 "".A+0
    rel 63+0 t=11 +0
    rel 76+4 t=8 runtime.morestack_noctxt+0
"".A STEXT size=117 args=0x10 locals=0x20
    0x0000 00000 (demo5.go:11)  TEXT    "".A(SB), ABIInternal, $32-16
    0x0000 00000 (demo5.go:11)  MOVQ    (TLS), CX
    0x0009 00009 (demo5.go:11)  CMPQ    SP, 16(CX)
    0x000d 00013 (demo5.go:11)  JLS 110
    0x000f 00015 (demo5.go:11)  SUBQ    $32, SP
    0x0013 00019 (demo5.go:11)  MOVQ    BP, 24(SP)
    0x0018 00024 (demo5.go:11)  LEAQ    24(SP), BP
    0x001d 00029 (demo5.go:11)  FUNCDATA    $0, gclocals·ffd148479e14c29ee3c68361945c5d25(SB)
    0x001d 00029 (demo5.go:11)  FUNCDATA    $1, gclocals·663f8c6bfa83aa777198789ce63d9ab4(SB)
    0x001d 00029 (demo5.go:11)  FUNCDATA    $3, gclocals·9fb7f0986f647f17cb53dda1484e0f7a(SB)
    0x001d 00029 (demo5.go:11)  PCDATA  $2, $0
    0x001d 00029 (demo5.go:11)  PCDATA  $0, $0
    0x001d 00029 (demo5.go:11)  MOVQ    $0, "".~r1+48(SP)       // 定义返回值变量
    0x0026 00038 (demo5.go:12)  PCDATA  $2, $1
    0x0026 00038 (demo5.go:12)  LEAQ    type.noalg.struct { F uintptr; "".x int }(SB), AX
    0x002d 00045 (demo5.go:12)  PCDATA  $2, $0
    0x002d 00045 (demo5.go:12)  MOVQ    AX, (SP)
    0x0031 00049 (demo5.go:12)  CALL    runtime.newobject(SB)    // 创建一个结构体对象
    0x0036 00054 (demo5.go:12)  PCDATA  $2, $1
    0x0036 00054 (demo5.go:12)  MOVQ    8(SP), AX                // 上面的函数返回值放到AX
    0x003b 00059 (demo5.go:12)  PCDATA  $0, $1
    0x003b 00059 (demo5.go:12)  MOVQ    AX, ""..autotmp_3+16(SP)
    0x0040 00064 (demo5.go:12)  LEAQ    "".A.func1(SB), CX       // 函数地址放到CX
    0x0047 00071 (demo5.go:12)  PCDATA  $2, $0
    0x0047 00071 (demo5.go:12)  MOVQ    CX, (AX)                // 将匿名函数地址放入(AX)指向的地址，为F赋值
    0x004a 00074 (demo5.go:12)  PCDATA  $2, $1
    0x004a 00074 (demo5.go:12)  MOVQ    ""..autotmp_3+16(SP), AX   // AX现在里面放的是struct地址
    0x004f 00079 (demo5.go:12)  TESTB   AL, (AX)
    0x0051 00081 (demo5.go:12)  MOVQ    "".x+40(SP), CX
    0x0056 00086 (demo5.go:12)  PCDATA  $2, $0
    0x0056 00086 (demo5.go:12)  MOVQ    CX, 8(AX)              // 将x值赋值给struct的x
    0x005a 00090 (demo5.go:12)  PCDATA  $2, $1
    0x005a 00090 (demo5.go:12)  PCDATA  $0, $0
    0x005a 00090 (demo5.go:12)  MOVQ    ""..autotmp_3+16(SP), AX
    0x005f 00095 (demo5.go:12)  PCDATA  $2, $0
    0x005f 00095 (demo5.go:12)  PCDATA  $0, $2
    0x005f 00095 (demo5.go:12)  MOVQ    AX, "".~r1+48(SP) // 返回返回值变量，
    0x0064 00100 (demo5.go:12)  MOVQ    24(SP), BP
    0x0069 00105 (demo5.go:12)  ADDQ    $32, SP
    0x006d 00109 (demo5.go:12)  RET
    0x006e 00110 (demo5.go:12)  NOP
    0x006e 00110 (demo5.go:11)  PCDATA  $0, $-1
    0x006e 00110 (demo5.go:11)  PCDATA  $2, $-1
    0x006e 00110 (demo5.go:11)  CALL    runtime.morestack_noctxt(SB)
    0x0073 00115 (demo5.go:11)  JMP 0
    0x0000 65 48 8b 0c 25 00 00 00 00 48 3b 61 10 76 5f 48  eH..%....H;a.v_H
    0x0010 83 ec 20 48 89 6c 24 18 48 8d 6c 24 18 48 c7 44  .. H.l$.H.l$.H.D
    0x0020 24 30 00 00 00 00 48 8d 05 00 00 00 00 48 89 04  $0....H......H..
    0x0030 24 e8 00 00 00 00 48 8b 44 24 08 48 89 44 24 10  $.....H.D$.H.D$.
    0x0040 48 8d 0d 00 00 00 00 48 89 08 48 8b 44 24 10 84  H......H..H.D$..
    0x0050 00 48 8b 4c 24 28 48 89 48 08 48 8b 44 24 10 48  .H.L$(H.H.H.D$.H
    0x0060 89 44 24 30 48 8b 6c 24 18 48 83 c4 20 c3 e8 00  .D$0H.l$.H.. ...
    0x0070 00 00 00 eb 8b                                   .....
    rel 5+4 t=16 TLS+0
    rel 41+4 t=15 type.noalg.struct { F uintptr; "".x int }+0
    rel 50+4 t=8 runtime.newobject+0
    rel 67+4 t=15 "".A.func1+0
    rel 111+4 t=8 runtime.morestack_noctxt+0
"".A.func1 STEXT nosplit size=55 args=0x10 locals=0x10
    0x0000 00000 (demo5.go:12)  TEXT    "".A.func1(SB), NOSPLIT|NEEDCTXT|ABIInternal, $16-16
    0x0000 00000 (demo5.go:12)  SUBQ    $16, SP
    0x0004 00004 (demo5.go:12)  MOVQ    BP, 8(SP)
    0x0009 00009 (demo5.go:12)  LEAQ    8(SP), BP
    0x000e 00014 (demo5.go:12)  FUNCDATA    $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x000e 00014 (demo5.go:12)  FUNCDATA    $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x000e 00014 (demo5.go:12)  FUNCDATA    $3, gclocals·ebb0e8ce1793da18f0378b883cb3e122(SB)
    0x000e 00014 (demo5.go:12)  PCDATA  $2, $0
    0x000e 00014 (demo5.go:12)  PCDATA  $0, $0
    0x000e 00014 (demo5.go:12)  MOVQ    8(DX), AX          // //DX是闭包对象的地址，(取内容)，偏移量8上放的刚好是上下文的int
    0x0012 00018 (demo5.go:12)  MOVQ    AX, "".x(SP)       // 此时将上下文的int放入x中
    0x0016 00022 (demo5.go:12)  MOVQ    $0, "".~r1+32(SP)  // 定义返回值
    0x001f 00031 (demo5.go:13)  MOVQ    "".x(SP), AX       // 将x值放入AX
    0x0023 00035 (demo5.go:13)  ADDQ    "".y+24(SP), AX    // 执行加法运算，并放入AX
    0x0028 00040 (demo5.go:13)  MOVQ    AX, "".~r1+32(SP)  // 返回结果
    0x002d 00045 (demo5.go:13)  MOVQ    8(SP), BP
    0x0032 00050 (demo5.go:13)  ADDQ    $16, SP
    0x0036 00054 (demo5.go:13)  RET
    0x0000 48 83 ec 10 48 89 6c 24 08 48 8d 6c 24 08 48 8b  H...H.l$.H.l$.H.
    0x0010 42 08 48 89 04 24 48 c7 44 24 20 00 00 00 00 48  B.H..$H.D$ ....H
    0x0020 8b 04 24 48 03 44 24 18 48 89 44 24 20 48 8b 6c  ..$H.D$.H.D$ H.l
    0x0030 24 08 48 83 c4 10 c3                             $.H....
```

可以看到，闭包通过
```
type.noalg.struct {
    F uintptr
    "".x int 
} 
```
实现了持有上下文的变量x的值。

作为一种特殊的函数：方法函数，是在某个类型上，定义的函数。其实他们和全局函数是一样的,只不过函数名被修饰为特定的名称,例如下面的方法函数和全局函数

```
type A struct {
    Name string
}

func (a *A) print() {
    a.Name = "FuncA"
    fmt.Println("function A")
}

func print() {
    fmt.Println("function A")
}

```

接受者是指针类型
```
TEXT    "".(*A).print(SB), ABIInternal, $88-8
```
如果是非指针类型的，则
```
TEXT    "".A.print(SB), ABIInternal, $88-16
```
全局函数：
```
TEXT    "".print(SB), ABIInternal, $88-0
```


# 2.defer函数
该函数在函数体执行完成后，以逆顺序逐个执行，即便程序发生严重错误时也会执行，支持匿名函数调用，常用于资源清理，文件关闭。GO没有异常机制，但有panic/recover模式来处理错误。panic 可以在任何地方引发，但recover只有在defer调用但函数中有效。
**注：defer关键词后面必须跟函数调用，而不是仅仅给出函数原型**

```
package main

import "fmt"

func main() {
    var a, b int = 8 ,0
    fmt.Println(divFunc(a, b))
    for i:=0;i<3;i++ {
        defer func () {
            fmt.Println(i)
        }()
    } // 3 3 3
}

func divFunc (a, b int) int {
    defer func () {
        if err := recover(); err != nil{
            fmt.Println("panic in divFunc, the argument is a invalid")
        }
    }()
    return a / b
}
```
main函数中的for循环中的输出结果是`3\n3\n3`,稍微调整一下defer函数内容
```
for i:=0;i<3;i++ {
        defer func (a int) {
            fmt.Println(a)
        }(i)
    } // 2 1 0
```

`go tool compile`(不带参数的defer函数)看一下：
```
 "".main STEXT size=191 args=0x0 locals=0x30
    0x0000 00000 (demo1.go:5)   TEXT    "".main(SB), ABIInternal, $48-0
    0x0000 00000 (demo1.go:5)   MOVQ    (TLS), CX
    0x0009 00009 (demo1.go:5)   CMPQ    SP, 16(CX)
    0x000d 00013 (demo1.go:5)   JLS 181
    0x0013 00019 (demo1.go:5)   SUBQ    $48, SP
    0x0017 00023 (demo1.go:5)   MOVQ    BP, 40(SP)
    0x001c 00028 (demo1.go:5)   LEAQ    40(SP), BP
    0x0021 00033 (demo1.go:5)   FUNCDATA    $0, gclocals·69c1753bd5f81501d95132d08af04464(SB)
    0x0021 00033 (demo1.go:5)   FUNCDATA    $1, gclocals·568470801006e5c0dc3947ea998fe279(SB)
    0x0021 00033 (demo1.go:5)   FUNCDATA    $3, gclocals·6e8d7ea4abad763909b26991048ee1fe(SB)
    0x0021 00033 (demo1.go:8)   PCDATA  $2, $1
    0x0021 00033 (demo1.go:8)   PCDATA  $0, $0
    0x0021 00033 (demo1.go:8)   LEAQ    type.int(SB), AX
    0x0028 00040 (demo1.go:8)   PCDATA  $2, $0
    0x0028 00040 (demo1.go:8)   MOVQ    AX, (SP)
    0x002c 00044 (demo1.go:8)   CALL    runtime.newobject(SB)   // 创建一个int类型的变量
    0x0031 00049 (demo1.go:8)   PCDATA  $2, $1
    0x0031 00049 (demo1.go:8)   MOVQ    8(SP), AX               
    0x0036 00054 (demo1.go:8)   PCDATA  $0, $1
    0x0036 00054 (demo1.go:8)   MOVQ    AX, "".&i+32(SP)        // 该变量放在i里面32（SP）
    0x003b 00059 (demo1.go:8)   PCDATA  $2, $0
    0x003b 00059 (demo1.go:8)   MOVQ    $0, (AX)                // i = 0
    0x0042 00066 (demo1.go:8)   JMP 68
    0x0044 00068 (demo1.go:8)   PCDATA  $2, $1
    0x0044 00068 (demo1.go:8)   MOVQ    "".&i+32(SP), AX        // 将i值放到AX
    0x0049 00073 (demo1.go:8)   PCDATA  $2, $0
    0x0049 00073 (demo1.go:8)   CMPQ    (AX), $3               // i与3进行比较
    0x004d 00077 (demo1.go:8)   JLT 81                         //小于3跳转81
    0x004f 00079 (demo1.go:8)   JMP 165                        //大于等于3跳转165
    0x0051 00081 (demo1.go:9)   PCDATA  $2, $1
    0x0051 00081 (demo1.go:9)   MOVQ    "".&i+32(SP), AX         // 将i值放到AX
    0x0056 00086 (demo1.go:11)  MOVQ    AX, ""..autotmp_3+24(SP)
    0x005b 00091 (demo1.go:9)   MOVL    $8, (SP)                // 8放到(SP)处
    0x0062 00098 (demo1.go:9)   PCDATA  $2, $2
    0x0062 00098 (demo1.go:9)   LEAQ    "".main.func1·f(SB), CX
    0x0069 00105 (demo1.go:9)   PCDATA  $2, $1
    0x0069 00105 (demo1.go:9)   MOVQ    CX, 8(SP)              // defer函数放到8(SP)处
    0x006e 00110 (demo1.go:9)   PCDATA  $2, $0
    0x006e 00110 (demo1.go:9)   MOVQ    AX, 16(SP)             // i放到16(SP)处
    0x0073 00115 (demo1.go:9)   CALL    runtime.deferproc(SB)  //调用derferproc函数
    0x0078 00120 (demo1.go:9)   TESTL   AX, AX
    0x007a 00122 (demo1.go:9)   JNE 149                        // deferproc调用不正常
    0x007c 00124 (demo1.go:9)   JMP 126                        // 正常
    0x007e 00126 (demo1.go:8)   PCDATA  $2, $-2
    0x007e 00126 (demo1.go:8)   PCDATA  $0, $-2
    0x007e 00126 (demo1.go:8)   JMP 128
    0x0080 00128 (demo1.go:8)   PCDATA  $2, $1
    0x0080 00128 (demo1.go:8)   PCDATA  $0, $1
    0x0080 00128 (demo1.go:8)   MOVQ    "".&i+32(SP), AX       // 将i值放到AX
    0x0085 00133 (demo1.go:8)   PCDATA  $2, $0
    0x0085 00133 (demo1.go:8)   MOVQ    (AX), AX
    0x0088 00136 (demo1.go:8)   PCDATA  $2, $3
    0x0088 00136 (demo1.go:8)   MOVQ    "".&i+32(SP), CX
    0x008d 00141 (demo1.go:8)   INCQ    AX                     // 自增i
    0x0090 00144 (demo1.go:8)   PCDATA  $2, $0
    0x0090 00144 (demo1.go:8)   MOVQ    AX, (CX)
    0x0093 00147 (demo1.go:8)   JMP 68
    0x0095 00149 (demo1.go:9)   PCDATA  $0, $0
    0x0095 00149 (demo1.go:9)   XCHGL   AX, AX
    0x0096 00150 (demo1.go:9)   CALL    runtime.deferreturn(SB)
    0x009b 00155 (demo1.go:9)   MOVQ    40(SP), BP
    0x00a0 00160 (demo1.go:9)   ADDQ    $48, SP
    0x00a4 00164 (demo1.go:9)   RET
    0x00a5 00165 (demo1.go:13)  XCHGL   AX, AX
    0x00a6 00166 (demo1.go:13)  CALL    runtime.deferreturn(SB)
    0x00ab 00171 (demo1.go:13)  MOVQ    40(SP), BP
    0x00b0 00176 (demo1.go:13)  ADDQ    $48, SP
    0x00b4 00180 (demo1.go:13)  RET
    0x00b5 00181 (demo1.go:13)  NOP
    0x00b5 00181 (demo1.go:5)   PCDATA  $0, $-1
    0x00b5 00181 (demo1.go:5)   PCDATA  $2, $-1
    0x00b5 00181 (demo1.go:5)   CALL    runtime.morestack_noctxt(SB)
    0x00ba 00186 (demo1.go:5)   JMP 0
    0x0000 65 48 8b 0c 25 00 00 00 00 48 3b 61 10 0f 86 a2  eH..%....H;a....
    0x0010 00 00 00 48 83 ec 30 48 89 6c 24 28 48 8d 6c 24  ...H..0H.l$(H.l$
    0x0020 28 48 8d 05 00 00 00 00 48 89 04 24 e8 00 00 00  (H......H..$....
    0x0030 00 48 8b 44 24 08 48 89 44 24 20 48 c7 00 00 00  .H.D$.H.D$ H....
    0x0040 00 00 eb 00 48 8b 44 24 20 48 83 38 03 7c 02 eb  ....H.D$ H.8.|..
    0x0050 54 48 8b 44 24 20 48 89 44 24 18 c7 04 24 08 00  TH.D$ H.D$...$..
    0x0060 00 00 48 8d 0d 00 00 00 00 48 89 4c 24 08 48 89  ..H......H.L$.H.
    0x0070 44 24 10 e8 00 00 00 00 85 c0 75 19 eb 00 eb 00  D$........u.....
    0x0080 48 8b 44 24 20 48 8b 00 48 8b 4c 24 20 48 ff c0  H.D$ H..H.L$ H..
    0x0090 48 89 01 eb af 90 e8 00 00 00 00 48 8b 6c 24 28  H..........H.l$(
    0x00a0 48 83 c4 30 c3 90 e8 00 00 00 00 48 8b 6c 24 28  H..0.......H.l$(
    0x00b0 48 83 c4 30 c3 e8 00 00 00 00 e9 41 ff ff ff     H..0.......A...
    rel 5+4 t=16 TLS+0
    rel 36+4 t=15 type.int+0
    rel 45+4 t=8 runtime.newobject+0
    rel 101+4 t=15 "".main.func1·f+0
    rel 116+4 t=8 runtime.deferproc+0
    rel 151+4 t=8 runtime.deferreturn+0
    rel 167+4 t=8 runtime.deferreturn+0
    rel 182+4 t=8 runtime.morestack_noctxt+0
"".main.func1 STEXT size=184 args=0x8 locals=0x78
    0x0000 00000 (demo1.go:9)   TEXT    "".main.func1(SB), ABIInternal, $120-8
    0x0000 00000 (demo1.go:9)   MOVQ    (TLS), CX
    0x0009 00009 (demo1.go:9)   CMPQ    SP, 16(CX)
    0x000d 00013 (demo1.go:9)   JLS 174
    0x0013 00019 (demo1.go:9)   SUBQ    $120, SP
    0x0017 00023 (demo1.go:9)   MOVQ    BP, 112(SP)
    0x001c 00028 (demo1.go:9)   LEAQ    112(SP), BP
    0x0021 00033 (demo1.go:9)   FUNCDATA    $0, gclocals·533adcd55fa5ed3e2fd959716125aef9(SB)
    0x0021 00033 (demo1.go:9)   FUNCDATA    $1, gclocals·439b0b339525dcecc112fff85820bb4d(SB)
    0x0021 00033 (demo1.go:9)   FUNCDATA    $3, gclocals·f6aec3988379d2bd21c69c093370a150(SB)
    0x0021 00033 (demo1.go:9)   FUNCDATA    $4, "".main.func1.stkobj(SB)
    0x0021 00033 (demo1.go:10)  PCDATA  $2, $1
    0x0021 00033 (demo1.go:10)  PCDATA  $0, $1
    0x0021 00033 (demo1.go:10)  MOVQ    "".&i+128(SP), AX       // 取i值(i地址+128(SP)),放到AX
    0x0029 00041 (demo1.go:10)  PCDATA  $2, $0
    0x0029 00041 (demo1.go:10)  MOVQ    (AX), AX
    0x002c 00044 (demo1.go:10)  MOVQ    AX, ""..autotmp_2+48(SP)
    0x0031 00049 (demo1.go:10)  MOVQ    AX, (SP)
    0x0035 00053 (demo1.go:10)  CALL    runtime.convT64(SB)
    0x003a 00058 (demo1.go:10)  PCDATA  $2, $1
    0x003a 00058 (demo1.go:10)  MOVQ    8(SP), AX
    0x003f 00063 (demo1.go:10)  PCDATA  $2, $0
    0x003f 00063 (demo1.go:10)  PCDATA  $0, $2
    0x003f 00063 (demo1.go:10)  MOVQ    AX, ""..autotmp_3+64(SP)
    0x0044 00068 (demo1.go:10)  PCDATA  $0, $3
    0x0044 00068 (demo1.go:10)  XORPS   X0, X0
    0x0047 00071 (demo1.go:10)  MOVUPS  X0, ""..autotmp_1+72(SP)
    0x004c 00076 (demo1.go:10)  PCDATA  $2, $1
    0x004c 00076 (demo1.go:10)  PCDATA  $0, $2
    0x004c 00076 (demo1.go:10)  LEAQ    ""..autotmp_1+72(SP), AX
    0x0051 00081 (demo1.go:10)  MOVQ    AX, ""..autotmp_5+56(SP)
    0x0056 00086 (demo1.go:10)  TESTB   AL, (AX)
    0x0058 00088 (demo1.go:10)  PCDATA  $2, $2
    0x0058 00088 (demo1.go:10)  PCDATA  $0, $1
    0x0058 00088 (demo1.go:10)  MOVQ    ""..autotmp_3+64(SP), CX
    0x005d 00093 (demo1.go:10)  PCDATA  $2, $3
    0x005d 00093 (demo1.go:10)  LEAQ    type.int(SB), DX
    0x0064 00100 (demo1.go:10)  PCDATA  $2, $2
    0x0064 00100 (demo1.go:10)  MOVQ    DX, ""..autotmp_1+72(SP)
    0x0069 00105 (demo1.go:10)  PCDATA  $2, $1
    0x0069 00105 (demo1.go:10)  MOVQ    CX, ""..autotmp_1+80(SP)
    0x006e 00110 (demo1.go:10)  TESTB   AL, (AX)
    0x0070 00112 (demo1.go:10)  JMP 114
    0x0072 00114 (demo1.go:10)  MOVQ    AX, ""..autotmp_4+88(SP)
    0x0077 00119 (demo1.go:10)  MOVQ    $1, ""..autotmp_4+96(SP)
    0x0080 00128 (demo1.go:10)  MOVQ    $1, ""..autotmp_4+104(SP)
    0x0089 00137 (demo1.go:10)  PCDATA  $2, $0
    0x0089 00137 (demo1.go:10)  MOVQ    AX, (SP)
    0x008d 00141 (demo1.go:10)  MOVQ    $1, 8(SP)
    0x0096 00150 (demo1.go:10)  MOVQ    $1, 16(SP)
    0x009f 00159 (demo1.go:10)  CALL    fmt.Println(SB)
    0x00a4 00164 (demo1.go:11)  MOVQ    112(SP), BP
    0x00a9 00169 (demo1.go:11)  ADDQ    $120, SP
    0x00ad 00173 (demo1.go:11)  RET
    0x00ae 00174 (demo1.go:11)  NOP
    0x00ae 00174 (demo1.go:9)   PCDATA  $0, $-1
    0x00ae 00174 (demo1.go:9)   PCDATA  $2, $-1
    0x00ae 00174 (demo1.go:9)   CALL    runtime.morestack_noctxt(SB)
    0x00b3 00179 (demo1.go:9)   JMP 0
    0x0000 65 48 8b 0c 25 00 00 00 00 48 3b 61 10 0f 86 9b  eH..%....H;a....
    0x0010 00 00 00 48 83 ec 78 48 89 6c 24 70 48 8d 6c 24  ...H..xH.l$pH.l$
    0x0020 70 48 8b 84 24 80 00 00 00 48 8b 00 48 89 44 24  pH..$....H..H.D$
    0x0030 30 48 89 04 24 e8 00 00 00 00 48 8b 44 24 08 48  0H..$.....H.D$.H
    0x0040 89 44 24 40 0f 57 c0 0f 11 44 24 48 48 8d 44 24  .D$@.W...D$HH.D$
    0x0050 48 48 89 44 24 38 84 00 48 8b 4c 24 40 48 8d 15  HH.D$8..H.L$@H..
    0x0060 00 00 00 00 48 89 54 24 48 48 89 4c 24 50 84 00  ....H.T$HH.L$P..
    0x0070 eb 00 48 89 44 24 58 48 c7 44 24 60 01 00 00 00  ..H.D$XH.D$`....
    0x0080 48 c7 44 24 68 01 00 00 00 48 89 04 24 48 c7 44  H.D$h....H..$H.D
    0x0090 24 08 01 00 00 00 48 c7 44 24 10 01 00 00 00 e8  $.....H.D$......
    0x00a0 00 00 00 00 48 8b 6c 24 70 48 83 c4 78 c3 e8 00  ....H.l$pH..x...
    0x00b0 00 00 00 e9 48 ff ff ff                          ....H...
    rel 5+4 t=16 TLS+0
    rel 54+4 t=8 runtime.convT64+0
    rel 96+4 t=15 type.int+0
    rel 160+4 t=8 fmt.Println+0
    rel 175+4 t=8 runtime.morestack_noctxt+0
```
defer函数主要涉及两个比较重要的函数：

- `runtime.deferproc`

```
// 创建一个新的siz字节参数的defered函数fn,编辑器将一个defer语句转化为对该函数的调用
func deferproc(siz int32, fn *funcval) { // arguments of fn follow fn
    if getg().m.curg != getg() {
        // go code on the system stack can't defer
        throw("defer on system stack")
    }

    // the arguments of fn are in a perilous state. The stack map
    // for deferproc does not describe them. So we can't let garbage
    // collection or stack copying trigger until we've copied them out
    // to somewhere safe. The memmove below does that.
    // Until the copy completes, we can only call nosplit routines.
    sp := getcallersp()
    argp := uintptr(unsafe.Pointer(&fn)) + unsafe.Sizeof(fn)
    callerpc := getcallerpc()

    d := newdefer(siz)
    if d._panic != nil {
        throw("deferproc: d.panic != nil after newdefer")
    }
    d.fn = fn
    d.pc = callerpc
    d.sp = sp
    switch siz {
    case 0:
        // Do nothing.
    case sys.PtrSize:
        *(*uintptr)(deferArgs(d)) = *(*uintptr)(unsafe.Pointer(argp))
    default:
        memmove(deferArgs(d), unsafe.Pointer(argp), uintptr(siz))
    }

    // deferproc 正常情况下返回0
    // 恐慌的defer返回1
    // 如果deferproc返回！= 0，则编译器生成的代码将始终检查返回值并跳转到函数的末尾。
    return0()
    // 没有代码可以到这里-C返回寄存器已设置且不能破坏。
}

//将与延迟调用关联的参数存储在内存 _defer头之后。
//go:nosplit
func deferArgs(d *_defer) unsafe.Pointer {
    if d.siz == 0 {
        // Avoid pointer past the defer allocation.
        return nil
    }
    return add(unsafe.Pointer(d), unsafe.Sizeof(*d))
}
```
- `runtime.deferreturn`
```
// 如果有，请运行一个延迟函数。 编译器在调用defer的任何函数的末尾插入对此的调用。
// 如果有被 defer 的函数的话，这里会调用 runtime·jmpdefer 跳到对应的位置
// 实际效果是会一遍遍地调用 deferreturn 直到 _defer 链表被清空
// 无法拆分堆栈，因为我们重用了调用方的框架来调用延迟的函数。

// 单个参数实际上并没有使用-它只是获取了地址，因此可以与待处理的延迟匹配。
//go:nosplit
func deferreturn(arg0 uintptr) {
    gp := getg()
    d := gp._defer
    if d == nil {
        return
    }
    sp := getcallersp()
    if d.sp != sp {
        return
    }

    // Moving arguments around.
    //
    // Everything called after this point must be recursively
    // nosplit because the garbage collector won't know the form
    // of the arguments until the jmpdefer can flip the PC over to
    // fn.
    switch d.siz {
    case 0:
        // Do nothing.
    case sys.PtrSize: //64位机上该值为8
        *(*uintptr)(unsafe.Pointer(&arg0)) = *(*uintptr)(deferArgs(d))
    default:
        memmove(unsafe.Pointer(&arg0), deferArgs(d), uintptr(d.siz))
    }
    fn := d.fn
    d.fn = nil
    gp._defer = d.link
    freedefer(d)
    jmpdefer(fn, uintptr(unsafe.Pointer(&arg0)))
}
```


TEST 指令设置零标志位, ZF, 当两个操作数进行And操作的结果值是0的时候，如果两个操作数相同，他们按位与结果为0，当他们都是0的时候。如果结果为0标识`deferproc`函数调用正常，跳转到0x1092ff9地址处开始执行，后面进行正常的自增变量，判断判断条件等等， 否则跳转到0x109302a继续执行，该指令后面就是处理`deferreturn`
```
   0x0000000001093024 <+100>:   test   %eax,%eax
   0x0000000001093026 <+102>:   jne    0x109302a <main.main+106>
   0x0000000001093028 <+104>:   jmp    0x1092ff9 <main.main+57>
```

自增i,当i大于3的时候，调到0x109303a继续执行。后面就是执行`runtime.deferreturn`
否则，
```
   0x0000000001092ffe <+62>:    incq   (%rax)
   0x0000000001093001 <+65>:    cmpq   $0x3,(%rax)
   0x0000000001093005 <+69>:    jge    0x109303a <main.main+122>
```
通过上面的分析，我们可看到，在没有传参的情况下defer函数，得到的是i的地址引用，当真正执行到defer函数时，此时的i值已经为`3`了，地址引用的值也是`3`.

到底是先执行return还是先执行defer函数，我们来看一下：
```
package main

import "fmt"

func main() {
     fmt.Println(deferReturn()) // 101
}

func deferReturn()  (ret int) {
      ret = 1
      defer func() {
           ret++
      }()
      return 100
}
```
go tool 看一下(为了方便看具体逻辑，去除了fmt.Prinltn调用)
```
"".main STEXT size=51 args=0x0 locals=0x10
    0x0000 00000 (demo4.go:7)   TEXT    "".main(SB), ABIInternal, $16-0
    0x0000 00000 (demo4.go:7)   MOVQ    (TLS), CX
    0x0009 00009 (demo4.go:7)   CMPQ    SP, 16(CX)
    0x000d 00013 (demo4.go:7)   JLS 44
    0x000f 00015 (demo4.go:7)   SUBQ    $16, SP
    0x0013 00019 (demo4.go:7)   MOVQ    BP, 8(SP)
    0x0018 00024 (demo4.go:7)   LEAQ    8(SP), BP
    0x001d 00029 (demo4.go:7)   FUNCDATA    $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x001d 00029 (demo4.go:7)   FUNCDATA    $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x001d 00029 (demo4.go:7)   FUNCDATA    $3, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x001d 00029 (demo4.go:10)  PCDATA  $2, $0
    0x001d 00029 (demo4.go:10)  PCDATA  $0, $0
    0x001d 00029 (demo4.go:10)  CALL    "".deferReturn(SB)   // deferReturn函数调用
    0x0022 00034 (demo4.go:11)  MOVQ    8(SP), BP
    0x0027 00039 (demo4.go:11)  ADDQ    $16, SP
    0x002b 00043 (demo4.go:11)  RET
    0x002c 00044 (demo4.go:11)  NOP
    0x002c 00044 (demo4.go:7)   PCDATA  $0, $-1
    0x002c 00044 (demo4.go:7)   PCDATA  $2, $-1
    0x002c 00044 (demo4.go:7)   CALL    runtime.morestack_noctxt(SB)
    0x0031 00049 (demo4.go:7)   JMP 0
    0x0000 65 48 8b 0c 25 00 00 00 00 48 3b 61 10 76 1d 48  eH..%....H;a.v.H
    0x0010 83 ec 10 48 89 6c 24 08 48 8d 6c 24 08 e8 00 00  ...H.l$.H.l$....
    0x0020 00 00 48 8b 6c 24 08 48 83 c4 10 c3 e8 00 00 00  ..H.l$.H........
    0x0030 00 eb cd                                         ...
    rel 5+4 t=16 TLS+0
    rel 30+4 t=8 "".deferReturn+0
    rel 45+4 t=8 runtime.morestack_noctxt+0
"".deferReturn STEXT size=138 args=0x8 locals=0x20
    0x0000 00000 (demo4.go:13)  TEXT    "".deferReturn(SB), ABIInternal, $32-8
    0x0000 00000 (demo4.go:13)  MOVQ    (TLS), CX
    0x0009 00009 (demo4.go:13)  CMPQ    SP, 16(CX)
    0x000d 00013 (demo4.go:13)  JLS 128
    0x000f 00015 (demo4.go:13)  SUBQ    $32, SP
    0x0013 00019 (demo4.go:13)  MOVQ    BP, 24(SP)
    0x0018 00024 (demo4.go:13)  LEAQ    24(SP), BP
    0x001d 00029 (demo4.go:13)  FUNCDATA    $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x001d 00029 (demo4.go:13)  FUNCDATA    $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
    0x001d 00029 (demo4.go:13)  FUNCDATA    $3, gclocals·9fb7f0986f647f17cb53dda1484e0f7a(SB)
    0x001d 00029 (demo4.go:13)  PCDATA  $2, $0
    0x001d 00029 (demo4.go:13)  PCDATA  $0, $0
    0x001d 00029 (demo4.go:13)  MOVQ    $0, "".ret+40(SP)   // 定义返回值变量
    0x0026 00038 (demo4.go:14)  MOVQ    $1, "".ret+40(SP)   // 返回值变量赋值1
    0x002f 00047 (demo4.go:15)  MOVL    $8, (SP)            // SP第一个字节是常数8(1000)
    0x0036 00054 (demo4.go:15)  PCDATA  $2, $1
    0x0036 00054 (demo4.go:15)  LEAQ    "".deferReturn.func1·f(SB), AX    // f函数地址
    0x003d 00061 (demo4.go:15)  PCDATA  $2, $0
    0x003d 00061 (demo4.go:15)  MOVQ    AX, 8(SP)                         // SP第二个字节
    0x0042 00066 (demo4.go:15)  PCDATA  $2, $1
    0x0042 00066 (demo4.go:15)  LEAQ    "".ret+40(SP), AX
    0x0047 00071 (demo4.go:15)  PCDATA  $2, $0
    0x0047 00071 (demo4.go:15)  MOVQ    AX, 16(SP)                       // SP第三个自己，放返回值
    0x004c 00076 (demo4.go:15)  CALL    runtime.deferproc(SB)
    0x0051 00081 (demo4.go:15)  TESTL   AX, AX
    0x0053 00083 (demo4.go:15)  JNE 112                        // deferproc函数调用异常
    0x0055 00085 (demo4.go:15)  JMP 87                         // deferproc函数调用正常
    0x0057 00087 (demo4.go:18)  MOVQ    $100, "".ret+40(SP)    // 返回值赋值为100
    0x0060 00096 (demo4.go:18)  XCHGL   AX, AX
    0x0061 00097 (demo4.go:18)  CALL    runtime.deferreturn(SB)  
    0x0066 00102 (demo4.go:18)  MOVQ    24(SP), BP
    0x006b 00107 (demo4.go:18)  ADDQ    $32, SP
    0x006f 00111 (demo4.go:18)  RET
    0x0070 00112 (demo4.go:15)  XCHGL   AX, AX
    0x0071 00113 (demo4.go:15)  CALL    runtime.deferreturn(SB)
    0x0076 00118 (demo4.go:15)  MOVQ    24(SP), BP
    0x007b 00123 (demo4.go:15)  ADDQ    $32, SP
    0x007f 00127 (demo4.go:15)  RET
    0x0080 00128 (demo4.go:15)  NOP
    0x0080 00128 (demo4.go:13)  PCDATA  $0, $-1
    0x0080 00128 (demo4.go:13)  PCDATA  $2, $-1
    0x0080 00128 (demo4.go:13)  CALL    runtime.morestack_noctxt(SB)
    0x0085 00133 (demo4.go:13)  JMP 0
    0x0000 65 48 8b 0c 25 00 00 00 00 48 3b 61 10 76 71 48  eH..%....H;a.vqH
    0x0010 83 ec 20 48 89 6c 24 18 48 8d 6c 24 18 48 c7 44  .. H.l$.H.l$.H.D
    0x0020 24 28 00 00 00 00 48 c7 44 24 28 01 00 00 00 c7  $(....H.D$(.....
    0x0030 04 24 08 00 00 00 48 8d 05 00 00 00 00 48 89 44  .$....H......H.D
    0x0040 24 08 48 8d 44 24 28 48 89 44 24 10 e8 00 00 00  $.H.D$(H.D$.....
    0x0050 00 85 c0 75 1b eb 00 48 c7 44 24 28 64 00 00 00  ...u...H.D$(d...
    0x0060 90 e8 00 00 00 00 48 8b 6c 24 18 48 83 c4 20 c3  ......H.l$.H.. .
    0x0070 90 e8 00 00 00 00 48 8b 6c 24 18 48 83 c4 20 c3  ......H.l$.H.. .
    0x0080 e8 00 00 00 00 e9 76 ff ff ff                    ......v...
    rel 5+4 t=16 TLS+0
    rel 57+4 t=15 "".deferReturn.func1·f+0
    rel 77+4 t=8 runtime.deferproc+0
    rel 98+4 t=8 runtime.deferreturn+0
    rel 114+4 t=8 runtime.deferreturn+0
    rel 129+4 t=8 runtime.morestack_noctxt+0
"".deferReturn.func1 STEXT nosplit size=20 args=0x8 locals=0x0
    0x0000 00000 (demo4.go:15)  TEXT    "".deferReturn.func1(SB), NOSPLIT|ABIInternal, $0-8
    0x0000 00000 (demo4.go:15)  FUNCDATA    $0, gclocals·1a65e721a2ccc325b382662e7ffee780(SB)
    0x0000 00000 (demo4.go:15)  FUNCDATA    $1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
    0x0000 00000 (demo4.go:15)  FUNCDATA    $3, gclocals·1cf923758aae2e428391d1783fe59973(SB)
    0x0000 00000 (demo4.go:16)  PCDATA  $2, $1
    0x0000 00000 (demo4.go:16)  PCDATA  $0, $0
    0x0000 00000 (demo4.go:16)  MOVQ    "".&ret+8(SP), AX
    0x0005 00005 (demo4.go:16)  PCDATA  $2, $0
    0x0005 00005 (demo4.go:16)  MOVQ    (AX), AX
    0x0008 00008 (demo4.go:16)  PCDATA  $2, $2
    0x0008 00008 (demo4.go:16)  PCDATA  $0, $1
    0x0008 00008 (demo4.go:16)  MOVQ    "".&ret+8(SP), CX
    0x000d 00013 (demo4.go:16)  INCQ    AX
    0x0010 00016 (demo4.go:16)  PCDATA  $2, $0
    0x0010 00016 (demo4.go:16)  MOVQ    AX, (CX)
    0x0013 00019 (demo4.go:17)  RET
    0x0000 48 8b 44 24 08 48 8b 00 48 8b 4c 24 08 48 ff c0  H.D$.H..H.L$.H..
    0x0010 48 89 01 c3                                      H...
```
所以，严格意义上说，defer函数是在调用defer函数的函数返回后执行的。