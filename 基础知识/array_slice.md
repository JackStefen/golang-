> 数组部分

数组是指一系列同一种类型数据的集合，是一个由固定长度的特定类型元素组成的序列，一个数组可以由零个或多个元素组成。特别地，数组的长度是数组类型的组成部分。不同长度或不同类型的数据组成的数组都是不同的类型,因此在Go语言中很少直接使用数组.Golang中声明数组的方法有：
```
[12]int
[2*N] struct {x,y int32}
[10]float64
[2][4]int
[2][2][2]float32
```
后面两个为多维数组。数组长度定义好以后，是不能修改的。内置函数len可以用于计算数组的长度，cap函数可以用于计算数组的容量。不过对于数组类型来说，len和cap函数返回的结果始终是一样的，都是对应数组类型的长度。
# 1. 元素访问有两种方式
- `for i:=0;i<len(arr);i++{fmt.Println("element", i, "of arr is", arr[i])}`
- `for i,v := range arr {fmt.Println("element", i ,"of arr is ", v)}`
用for range方式迭代的性能可能会更好一些，因为这种迭代可以保证不会出现数组越界的情形，每轮迭代对数组元素的访问时可以省去对下标越界的判断。
示例：
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
     modify(s) // cannot use s (type [4]int) as type [5]int in argument to modify
     s1 := [5]int{1,2,4}  
     modify(s1) // [19 2 4 0 0]
}

```

数组的缺点：
- 数组是定长的，一但定义好长度就无法修改
- 相同元素类型，不同长度的数组属于不同的类型，之间是不能相互赋值的
- 数组是一个值类型，一个数组变量即表示整个数组，它并不是隐式的指向第一个元素的指针（比如C语言的数组），而是一个完整的值。所有的值类型变量在赋值和作为参数传递时都将产生一次复制动作

数组虽然有诸多缺点，平时业务上使用较少，但是确实非常重要的数据结构，数组类型是切片和字符串等结构的基础。

> 切片部分

切片分为切片头和内部数组两部分，数组切片的数据结构可以抽象出三个变量
- 一个指向原生数组的指针
- 切片中元素的个数
- 已分配的存储空间

# 1.切片的创建有多种方式
- 1.1基于数组，切割后产生的就是一个切片类型的数据
- 1.2直接创建，使用make方法
- 1.3基于切片，切割后产生的也是一个切片类型的数据
- 1.4内容复制，copy方法
**注：切割语法和python中的切割语法非常相似，唯一不同的是golang中切割不支持负数的索引**
```
package main

import (
   "fmt"
   "reflect"
)
func main() {
    // 1.基于数组切割创建
    var arr [10]int = [10]int{1,2,3,4,5,6,7,8,9,10} //声明一个数组并赋值
    var arrSlice []int = arr[:5]  //切割操作产生的切片类型
    fmt.Println(arrSlice) //[1 2 3 4 5]
    fmt.Println(reflect.TypeOf(arrSlice)) //[]int
    // 2.通过直接的make方法创建
    // make函数可以指定三个参数，也可以只指定前两个参数
    // 第二个参数时指定切片的个数，第三个参数时指定分配的空间，切片元素的默认值为0
    var sli1 []int = make([]int, 0)
    var sli2 []string = make([]string, 2, 5)
    sli1 = []int{1,2}
    sli2 = []string{"zhao", "qian"}
    fmt.Println(sli1, sli2) // [1 2] [zhao qian]
    // 3.基于切片
    sli3 := sli2[:1]
    fmt.Println(sli3) //[zhao]
    //4.copy方法
    var sli4 []string = []string{"sun","li"}
    copy(sli4, sli3)
    fmt.Println(sli4) // [zhao li]
}
```
# 2.切片的追加操作和遍历操作
- 切片长度可通过len()函数获得，存储空间值可以通过cap()函数获得。
- 通过append()函数可以，追加元素到切片中，当切片超过存储空间的时候，会重新分配内存，append函数写法支持像js中的rest逆运算三点`...`。切片追加后形成了一个新的切片变量,而老的切片变量的三个域其实并不会改变，改变的只是底层的数组。
- 切片元素遍历同数组：for索引遍历和range遍历
- 修改传入参数的元素就可以通过数组切片来完成
```
package main

import (
    "fmt"
)

type ByteSlice []byte

func (slice ByteSlice) append(data []byte) []byte {
    fmt.Println("append args: ", data, slice)
    slice = append(slice, data...)
    return slice
}

func (p *ByteSlice) Append(data []byte) {
    fmt.Println("Append args: ", p, data)
    slice := *p
    *p = append(slice, data...)
}

func modify(in []int) {
    in[0] = 10
}

func main(){
    var mySlice []int

    mySlice = append(mySlice, 1)
    mySlice = append(mySlice, 5, 6, 7)
    mySlice = append(mySlice, []int{8, 9}...)
    fmt.Println(mySlice) // [1 5 6 7 8 9]
    fmt.Println(len(mySlice), cap(mySlice))//6 8
    
    n:= ByteSlice{1,2,3}
    // 通过实现使用接受者是指针类型来修改调用者
    n.Append([]byte{4,5})
    fmt.Println(n) //[1 2 3 4 5]
    // 而不是修改后返回修改后的值然后重新赋值，
    n = n.append([]byte{6,7})
    fmt.Println(n)//[1 2 3 4 5 6 7]
    modify(mySlice)
    fmt.Println(mySlice) //[10 5 6 7 8 9]
    // range 遍历
    for i, v := range mySlice {
        fmt.Println(i, v)
    }
    // =============
    0 10
    1 5
    2 6
    3 7
    4 8
    5 9
}
```
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
    fmt.Println(b, c)
}
```
结果输出为：
```
[3 4 5] [2 3]
[9 4 5] [2 9]
```
可以看出对数组切片元素等修改，影响到所有产生于原切片中共同位置等元素的值。但是，如果因为超出切片容量，而重新分配了内存大小，那么，重新分配后的遍历是不会受到影响的。
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
    fmt.Println(b, c)
}
```
输出结果：
```
[3 4 5] [2 3]
[9 4 5] [2 3 1 2 3 4 5 5]
```
# 小贴士：make / new
make用于内建类型（map、slice 和channel）的内存分配。new用于各种类型的内存分配。
new(T)分配了零值填充的T类型的内存空间，并且返回其地址，即一个`*T`类型的值
内建函数make(T, args)与new(T)有着不同的功能，make只能创建slice、map和channel，并且返回一个有初始值(非零)的T类型，而不是`*T`。本质来讲，导致这三个类型有所不同的原因是指向数据结构的引用在使用前必须被初始化.
# 延伸，切片的内部结构
上面我们提到，切片的数据结构可以抽象出三个变量
- 一个指向原生数组的指针
- 切片中元素的个数
- 已分配的存储空间
我们可以通过unsafe包里面的工具来了解一下
```
package main

import (
    "fmt"
    "unsafe"
    "reflect"
)

func main(){
    sli := make([]int, 10)
    fmt.Println(sli) //[0 0 0 0 0 0 0 0 0 0],零值
    sli = []int{1,2,3,4}
    slilen := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&sli)) + uintptr(8)))
    slicap := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&sli)) + uintptr(16)))
    fmt.Println(*slilen, *slicap)  //4 4
    sliArr := (**[8]int)(unsafe.Pointer(&sli))
    fmt.Println(**sliArr) //[1 2 3 4 0 0 0 0]
    fmt.Println(reflect.TypeOf(**sliArr)) //[8]int 数组类型
}
```
可以看到，切片的底层还是一个数组结构，切片头部还有长度和容量信息，具体结构可以使用下面的图解释
![切片结构.png](https://upload-images.jianshu.io/upload_images/3004516-baacd0b71b27e79b.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)
