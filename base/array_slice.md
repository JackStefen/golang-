# 1. golang数组

数组是指一系列同一种类型数据的集合，是一个由固定长度的特定类型元素组成的序列，一个数组可以由零个或多个元素组成。
**特别地，数组的长度是数组类型的组成部分。**不同长度或不同类型的数据组成的数组都是不同的类型
数组长度定义好以后，是不能修改的。
- 内置函数len可以用于计算数组的长度
- cap函数可以用于计算数组的容量。

不过对于数组类型来说，len和cap函数返回的结果始终是一样的，都是对应数组类型的长度。

## 1.1 声明方式：
```
[12]int
[2*N] struct {x,y int32}
[10]float64
[2][4]int // 二维数组
[2][2][2]float32
```
## 1.2 遍历

```
for i:=0;i<len(arr);i++ {
    fmt.Println("element", i, "of arr is", arr[i])
}
for i,v := range arr {
    fmt.Println("element", i ,"of arr is ", v)
}
```
用`for range`方式迭代的性能可能会更好一些，因为这种迭代可以保证不会出现数组越界的情形，每轮迭代对数组元素的访问时可以省去对下标越界的判断。

数组的缺点：

- 数组是定长的，一但定义好长度就无法修改
- 相同元素类型，不同长度的数组属于不同的类型，之间是不能相互赋值的
- 数组是一个值类型，一个数组变量即表示整个数组，它并不是隐式的指向第一个元素的指针（比如C语言的数组），而是一个完整的值。所有的值类型变量在赋值和作为参数传递时都将产生一次复制动作

数组虽然有诸多缺点，平时业务上使用较少，但是确实非常重要的数据结构，数组类型是切片和字符串等结构的基础。
```
package main

import (
    "fmt"
)

func modify(in [5]int) {
    in[0] = 19
    fmt.Println(in)
}

func main(){
     s := [4]int{1,2,4}
     fmt.Println(s[1: len(s)-1])
     modify(s)
}
```
运行该示例会报错：
```
➜  slicetest go run arraytest.go
# command-line-arguments
./arraytest.go:15:12: cannot use s (type [4]int) as type [5]int in argument to modify
```
# 2. 切片
切片分为切片头和内部数组两部分，数组切片的数据结构可以抽象出三个变量

- 一个指向原生数组的指针
- 切片中元素的个数
- 已分配的存储空间

具体可以看`runtime`包中的`slice`结构体
```
type slice struct {
    array unsafe.Pointer 
    len   int
    cap   int
}

```

## 2.1 切片的创建方式
- 基于数组，切割后产生的就是一个切片类型的数据
- 直接创建，使用make方法，例如：make([]int, 10)
- 基于切片，切割后产生的也是一个切片类型的数据
- 内容复制，copy方法,例如:copy(des, org)

注：切割语法和python中的切割语法非常相似，唯一不同的是golang中切割不支持负数的索引

## 2.2 操作
- 切片长度可通过len()函数获得
- 存储空间值可以通过cap()函数获得。

- 通过append()函数可以，追加元素到切片中，当切片超过存储空间的时候，会重新分配内存，append函数写法支持像js中的rest逆运算三点...。切片追加后形成了一个新的切片变量,而老的切片变量的三个域其实并不会改变，改变的只是底层的数组。
- 切片元素遍历同数组：for索引遍历和range遍历
- 修改传入参数的元素就可以通过切片来完成

使用数组切片时的注意事项：
```
package main

import (
    "fmt"
)

func main() {
    var a []int = []int{1,2,3,4,5}
    b := a[2:5]
    c := a[1:3]
    fmt.Println(b, c)
    b[0] = 9
    fmt.Println(b, c) // [9 4 5] [2 9]
}
```

可以看出对数组切片元素的修改，影响到所有产生于原切片中共同位置等元素的值。但是，如果因为超出切片容量，而重新分配了内存大小，那么，重新分配后的遍历是不会受到影响的。
```
package main

import (
    "fmt"
)

func main() {
    var a []int = []int{1,2,3,4,5}
    b := a[2:5]
    c := a[1:3]
    fmt.Println(b, c)
    append(c, 1,2,3,4,5,5)
    b[0] = 9
    fmt.Println(b, c) //[9 4 5] [2 3 1 2 3 4 5 5]
}
```
**当使用slice作为参数进行传递时，发生了一次值copy,在函数中，对传入的参数的修改，都不会影响到调用者函数中的源slice中的值**
请看一下下面的例子
```
package main


import (
	"fmt"
)

func main() {


	var a = []int{1,2,3,4}
	funcA(a)
	c := make([]int, 0)
	funcA(c)
	d := make([]int, 0, 5)
	funcA(d)
	fmt.Println(a)
	fmt.Println(c)
	fmt.Println(d)
}


func funcA(b []int) {
	b = append(b, 4, 5, 6, 7, 8, 9, 10)
}

```
**但是map作为参数进行传递的时候，就不一样了，函数中对map的修改，都会影响到，传入的参数的原始map**
```
package main


import (
	"fmt"
)

func main() {
	var a map[int]int = map[int]int{
		1: 1,
		2:2,
	}

	funcA(a)
	fmt.Println(a)
}

func funcA(b map[int]int) {
	b[1] = 190
}
```

# 题外话
- make用于内建类型（map、slice 和channel）的内存分配。
- new用于各种类型的内存分配。

new(T)分配了零值填充的T类型的内存空间，并且返回其地址，即一个*T类型的值

内建函数make(T, args)与new(T)有着不同的功能，make只能创建slice、map和channel，并且返回一个有初始值(非零)的T类型，而不是*T。

本质来讲，导致这三个类型有所不同的原因是指向数据结构的引用在使用前必须被初始化.

# 3.切片的源码分析

```
package main

import (
        "fmt"
)

func main() {
        a := make([]int, 10, 10)
        fmt.Println(a)
        a = append(a, 1)
        fmt.Println(a)
}
```

在main函数上加断点，使用反汇编查看程序的运行情况
```
(gdb) disass
Dump of assembler code for function main.main:
=> 0x0000000001092fc0 <+0>: mov    %gs:0x30,%rcx
   0x0000000001092fc9 <+9>: cmp    0x10(%rcx),%rsp
   0x0000000001092fcd <+13>:    jbe    0x1093111 <main.main+337>
   0x0000000001092fd3 <+19>:    sub    $0x60,%rsp
   0x0000000001092fd7 <+23>:    mov    %rbp,0x58(%rsp)
   0x0000000001092fdc <+28>:    lea    0x58(%rsp),%rbp
   0x0000000001092fe1 <+33>:    lea    0x10418(%rip),%rax        # 0x10a3400 <type.*+66176>
   0x0000000001092fe8 <+40>:    mov    %rax,(%rsp)
   0x0000000001092fec <+44>:    movq   $0xa,0x8(%rsp)
   0x0000000001092ff5 <+53>:    movq   $0xa,0x10(%rsp)
   0x0000000001092ffe <+62>:    callq  0x103b480 <runtime.makeslice>
   0x0000000001093003 <+67>:    mov    0x18(%rsp),%rax
   0x0000000001093008 <+72>:    mov    %rax,0x40(%rsp)
   0x000000000109300d <+77>:    mov    %rax,(%rsp)
   0x0000000001093011 <+81>:    movq   $0xa,0x8(%rsp)
   0x000000000109301a <+90>:    movq   $0xa,0x10(%rsp)
   0x0000000001093023 <+99>:    callq  0x1008780 <runtime.convTslice>
   0x0000000001093028 <+104>:   mov    0x18(%rsp),%rax
   0x000000000109302d <+109>:   xorps  %xmm0,%xmm0
   0x0000000001093030 <+112>:   movups %xmm0,0x48(%rsp)
   0x0000000001093035 <+117>:   lea    0xeb24(%rip),%rcx        # 0x10a1b60 <type.*+59872>
   0x000000000109303c <+124>:   mov    %rcx,0x48(%rsp)
   0x0000000001093041 <+129>:   mov    %rax,0x50(%rsp)
   0x0000000001093046 <+134>:   lea    0x48(%rsp),%rax
   0x000000000109304b <+139>:   mov    %rax,(%rsp)
   0x000000000109304f <+143>:   movq   $0x1,0x8(%rsp)
   0x0000000001093058 <+152>:   movq   $0x1,0x10(%rsp)
   0x0000000001093061 <+161>:   callq  0x108c9e0 <fmt.Println>
   0x0000000001093066 <+166>:   lea    0x10393(%rip),%rax        # 0x10a3400 <type.*+66176>
   0x000000000109306d <+173>:   mov    %rax,(%rsp)
   0x0000000001093071 <+177>:   mov    0x40(%rsp),%rax
   0x0000000001093076 <+182>:   mov    %rax,0x8(%rsp)
   0x000000000109307b <+187>:   movq   $0xa,0x10(%rsp)
   0x0000000001093084 <+196>:   movq   $0xa,0x18(%rsp)
   0x000000000109308d <+205>:   movq   $0xb,0x20(%rsp)
   0x0000000001093096 <+214>:   callq  0x103b580 <runtime.growslice>
   0x000000000109309b <+219>:   mov    0x28(%rsp),%rax
   0x00000000010930a0 <+224>:   mov    0x30(%rsp),%rcx
   0x00000000010930a5 <+229>:   mov    0x38(%rsp),%rdx
   0x00000000010930aa <+234>:   movq   $0x1,0x50(%rax)
   0x00000000010930b2 <+242>:   mov    %rax,(%rsp)
   0x00000000010930b6 <+246>:   lea    0x1(%rcx),%rax
   0x00000000010930ba <+250>:   mov    %rax,0x8(%rsp)
   0x00000000010930bf <+255>:   mov    %rdx,0x10(%rsp)
   0x00000000010930c4 <+260>:   callq  0x1008780 <runtime.convTslice>
   0x00000000010930c9 <+265>:   mov    0x18(%rsp),%rax
   0x00000000010930ce <+270>:   xorps  %xmm0,%xmm0
   0x00000000010930d1 <+273>:   movups %xmm0,0x48(%rsp)
   0x00000000010930d6 <+278>:   lea    0xea83(%rip),%rcx        # 0x10a1b60 <type.*+59872>
   0x00000000010930dd <+285>:   mov    %rcx,0x48(%rsp)
   0x00000000010930e2 <+290>:   mov    %rax,0x50(%rsp)
   0x00000000010930e7 <+295>:   lea    0x48(%rsp),%rax
   0x00000000010930ec <+300>:   mov    %rax,(%rsp)
--Type <RET> for more, q to quit, c to continue without paging--c
   0x00000000010930f0 <+304>:   movq   $0x1,0x8(%rsp)
   0x00000000010930f9 <+313>:   movq   $0x1,0x10(%rsp)
   0x0000000001093102 <+322>:   callq  0x108c9e0 <fmt.Println>
   0x0000000001093107 <+327>:   mov    0x58(%rsp),%rbp
   0x000000000109310c <+332>:   add    $0x60,%rsp
   0x0000000001093110 <+336>:   retq
   0x0000000001093111 <+337>:   callq  0x104f200 <runtime.morestack_noctxt>
   0x0000000001093116 <+342>:   jmpq   0x1092fc0 <main.main>
   0x000000000109311b <+347>:   int3
   0x000000000109311c <+348>:   int3
   0x000000000109311d <+349>:   int3
   0x000000000109311e <+350>:   int3
   0x000000000109311f <+351>:   int3
End of assembler dump.
```

可以看到当我们在创建一个slice时，底层调用的是`runtime.makeslice`,看一下这个函数
```
func makeslice(et *_type, len, cap int) unsafe.Pointer {
    mem, overflow := math.MulUintptr(et.size, uintptr(cap))
    if overflow || mem > maxAlloc || len < 0 || len > cap {
        // NOTE: Produce a 'len out of range' error instead of a
        // 'cap out of range' error when someone does make([]T, bignumber).
        // 'cap out of range' is true too, but since the cap is only being
        // supplied implicitly, saying len is clearer.
        // See golang.org/issue/4085.
        mem, overflow := math.MulUintptr(et.size, uintptr(len))
        if overflow || mem > maxAlloc || len < 0 {
            panicmakeslicelen()
        }
        panicmakeslicecap()
    }

    return mallocgc(mem, et, true)
}
```
实际调用的是`mallocgc(mem, et, true)`来进行资源分配,当分配的资源使用完之后，如果再进行append操作的话，比如
当append(a, 1)一个新值到切片中时，将会首先进行一次`runtime.growslice`的动作，进行扩容，然后再进行对新的slice的append操作

```
// growslice 处理slice append操作时的扩容。
// 该函数传入的参数为元素类型，老的slice，和渴望的新的容量最小值
// 返回一个新的slice，容量大小至少为传入的cap大小，老的slice数据将copy到新的slice中。
// 新的slice长度设置为老的slice长度，而不是新设的容量大小值
// 这是为了方便代码生成。老的slice长度用于计算在append操作写入新的值时的位置 
func growslice(et *_type, old slice, cap int) slice {
    if raceenabled {
        callerpc := getcallerpc()
        racereadrangepc(old.array, uintptr(old.len*int(et.size)), callerpc, funcPC(growslice))
    }
    if msanenabled {
        msanread(old.array, uintptr(old.len*int(et.size)))
    }

    if cap < old.cap {
        panic(errorString("growslice: cap out of range"))
    }

    if et.size == 0 {
        // append should not create a slice with nil pointer but non-zero len.
        // We assume that append doesn't need to preserve old.array in this case.
        return slice{unsafe.Pointer(&zerobase), old.len, cap}
    }

    newcap := old.cap
    doublecap := newcap + newcap
    if cap > doublecap {
        newcap = cap
    } else {
        if old.len < 1024 {
            newcap = doublecap
        } else {
            // Check 0 < newcap to detect overflow
            // and prevent an infinite loop.
            for 0 < newcap && newcap < cap {
                newcap += newcap / 4
            }
            // Set newcap to the requested cap when
            // the newcap calculation overflowed.
            if newcap <= 0 {
                newcap = cap
            }
        }
    }

    var overflow bool
    var lenmem, newlenmem, capmem uintptr
    // Specialize for common values of et.size.
    // For 1 we don't need any division/multiplication.
    // For sys.PtrSize, compiler will optimize division/multiplication into a shift by a constant.
    // For powers of 2, use a variable shift.
    switch {
    case et.size == 1:
        lenmem = uintptr(old.len)
        newlenmem = uintptr(cap)
        capmem = roundupsize(uintptr(newcap))
        overflow = uintptr(newcap) > maxAlloc
        newcap = int(capmem)
    case et.size == sys.PtrSize:
        lenmem = uintptr(old.len) * sys.PtrSize
        newlenmem = uintptr(cap) * sys.PtrSize
        capmem = roundupsize(uintptr(newcap) * sys.PtrSize)
        overflow = uintptr(newcap) > maxAlloc/sys.PtrSize
        newcap = int(capmem / sys.PtrSize)
    case isPowerOfTwo(et.size):
        var shift uintptr
        if sys.PtrSize == 8 {
            // Mask shift for better code generation.
            shift = uintptr(sys.Ctz64(uint64(et.size))) & 63
        } else {
            shift = uintptr(sys.Ctz32(uint32(et.size))) & 31
        }
        lenmem = uintptr(old.len) << shift
        newlenmem = uintptr(cap) << shift
        capmem = roundupsize(uintptr(newcap) << shift)
        overflow = uintptr(newcap) > (maxAlloc >> shift)
        newcap = int(capmem >> shift)
    default:
        lenmem = uintptr(old.len) * et.size
        newlenmem = uintptr(cap) * et.size
        capmem, overflow = math.MulUintptr(et.size, uintptr(newcap))
        capmem = roundupsize(capmem)
        newcap = int(capmem / et.size)
    }

    // The check of overflow in addition to capmem > maxAlloc is needed
    // to prevent an overflow which can be used to trigger a segfault
    // on 32bit architectures with this example program:
    //
    // type T [1<<27 + 1]int64
    //
    // var d T
    // var s []T
    //
    // func main() {
    //   s = append(s, d, d, d, d)
    //   print(len(s), "\n")
    // }
    if overflow || capmem > maxAlloc {
        panic(errorString("growslice: cap out of range"))
    }

    var p unsafe.Pointer
    if et.kind&kindNoPointers != 0 {
        p = mallocgc(capmem, nil, false)
        // The append() that calls growslice is going to overwrite from old.len to cap (which will be the new length).
        // Only clear the part that will not be overwritten.
        memclrNoHeapPointers(add(p, newlenmem), capmem-newlenmem)
    } else {
        // Note: can't use rawmem (which avoids zeroing of memory), because then GC can scan uninitialized memory.
        p = mallocgc(capmem, et, true)
        if writeBarrier.enabled {
            // Only shade the pointers in old.array since we know the destination slice p
            // only contains nil pointers because it has been cleared during alloc.
            bulkBarrierPreWriteSrcOnly(uintptr(p), uintptr(old.array), lenmem)
        }
    }
    memmove(p, old.array, lenmem)

    return slice{p, old.len, newcap}
}
```
growslice需要将旧切片也作为参数传入，在生成新的切片时，需要将旧切片的内容复制进来，第三个参数为新的切片最小的容量，新的切片也是通过`mallocgc(capmem, et, true)`调用生成的。最后返回的是新`slice`.而slice的array就是新生成的底层数组指针。老切片的长度作为新切片的长度，新的最小容量作为当前的切片容量。
