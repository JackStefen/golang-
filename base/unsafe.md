Pointer代表一个指向任意类型的指针。有四种特殊的操作，可以在Pointer类型上进行，但对其他类型并不能使用：
- 任意类型的指针都可以转化为一个Pointer
- 一个Pointer可以转化为任意类型的指针
- 一个uintptr可以转化为一个Pointer
- 一个Pointer可以转化为一个uintptr

因此，指针允许程序越过类型系统并读写任意内存。 使用时应格外小心。

以下涉及Pointer的模式是有效的

不使用这些模式的代码今天可能无效，或者将来变得无效。即使以下有效模式也带有重要的警告。

## `1.Conversion of a *T1 to Pointer to *T2.`

T2不大于T1，并且两个共享相同的内存布局，这种转换允许将一种类型的数据重新解释为另一种类型的数据.例如`math.Float64bits`的实现：
```
func Float64bits(f float64) uint64 {
		return *(*uint64)(unsafe.Pointer(&f))
}
```

## `2. Conversion of a Pointer to a uintptr (but not back to Pointer).`

将`Pointer`转换为`uintptr`会生成所指向的值的内存地址（整数）。 这种`uintptr`的通常用法是打印它。

**将uintptr转换回Pointer通常是无效的。**

**`uintptr`是整数，不是引用。**
将`Pointer`转换为`uintptr`会创建一个没有指针语义的整数值。 **即使uintptr保留了某个对象的地址，垃圾回收器也不会在对象移动时更新该uintptr的值，该uintptr也不会使该对象被回收**。

其余模式枚举了从`uintptr`到`Pointer`的唯一有效转换。

## `3.Conversion of a Pointer to a uintptr and back, with arithmetic.`

如果p指向已分配的对象，则可以通过转换为`uintptr`，添加偏移量并将其转换回`Pointer`的方式进入对象。

```
p = unsafe.Pointer(uintptr(p) + offset)
```
此模式最常见的用法是访问结构或数组元素中的字段 
```
// 二者是等效的 f := unsafe.Pointer(&s.f)
//	f := unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Offsetof(s.f))
```

```
// 二者是等效的 e := unsafe.Pointer(&x[i])
//	e := unsafe.Pointer(uintptr(unsafe.Pointer(&x[0])) + i*unsafe.Sizeof(x[0]))
```

以这种方式从指针添加和减去偏移量都是有效的。在所有情况下，结果都必须继续指向原始分配的对象。

与C语言不同，将指针移到其原始分配的末尾是无效的：
```
	// INVALID: end points outside allocated space.
	var s thing
	end = unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Sizeof(s))
```

```
	// INVALID: end points outside allocated space.
	b := make([]byte, n)
	end = unsafe.Pointer(uintptr(unsafe.Pointer(&b[0])) + uintptr(n))
```

请注意，两个转换必须出现在相同的表达式中，并且它们之间只有中间的算术：

```
	// INVALID: uintptr cannot be stored in variable
	// before conversion back to Pointer.
	u := uintptr(p)
	p = unsafe.Pointer(u + offset)
```

请注意，指针必须指向已分配的对象 ，因此不可能会为`nil`。

```
	// INVALID: conversion of nil pointer
	u := unsafe.Pointer(nil)
	p := unsafe.Pointer(uintptr(u) + offset)
```

## `4.Conversion of a Pointer to a uintptr when calling syscall.Syscall.`

软件包`syscall`中的`Syscall`函数将其`uintptr`参数直接传递给操作系统，然后，操作系统可以根据调用的详细信息将其中一些参数重新解释为指针。 也就是说，系统调用实现正在将某些参数从`uintptr`隐式转换回指针。

如果必须将指针参数转换为uintptr以用作参数，则该转换必须出现在调用表达式本身中：
```
syscall.Syscall(SYS_READ, uintptr(fd), uintptr(unsafe.Pointer(p)), uintptr(n))
```
编译器处理汇编函数调用的参数列表中`Pointer`到`uintptr`转换，方法是安排保留引用的分配对象（如果有），直到调用完成后才移动，即使从类型本身来看，调用期间对象不再需要。


为了使编译器能够识别这种模式，转换必须出现在参数列表中：

```
	// INVALID: uintptr cannot be stored in variable
	// before implicit conversion back to Pointer during system call.
	u := uintptr(unsafe.Pointer(p))
	syscall.Syscall(SYS_READ, uintptr(fd), u, uintptr(n))
```
## `5.Conversion of the result of reflect.Value.Pointer or reflect.Value.UnsafeAddr from uintptr to Pointer.`

包`reflect`的名为`Pointer`和`UnsafeAddr`的`Value`方法返回`uintptr`而不是`unsafe.Pointer`类型，以防止调用者将结果更改为任意类型，而无需首先导入`unsafe`。 但是，这意味着结果很脆弱，必须在调用后立即使用相同的表达式将其转换为`Pointer`：

```
p := (*int)(unsafe.Pointer(reflect.ValueOf(new(int)).Pointer()))
```

与上述情况一样，在转换之前存储结果无效：
```
	// INVALID: uintptr cannot be stored in variable
	// before conversion back to Pointer.
	u := reflect.ValueOf(new(int)).Pointer()
	p := (*int)(unsafe.Pointer(u))
```

## `6.Conversion of a reflect.SliceHeader or reflect.StringHeader Data field to or from Pointer.`

与前面的情况一样，反射数据结构`SliceHeader`和`StringHeader`将字段`Data`声明为`uintptr`，以防止调用者在不首先导入`unsafe`的情况下将结果更改为任意类型。 但是，这意味着`SliceHeader`和`StringHeader`仅在解释实际切片或字符串值的内容时才有效。

```
	var s string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s)) // case 1
	hdr.Data = uintptr(unsafe.Pointer(p))              // case 6 (this case)
	hdr.Len = n
```
在这种用法中，`hdr.Data`实际上是引用字符串头部的指针的替代方法，而不是`uintptr`变量本身 。

通常，`reflect.SliceHeader`和`reflect.StringHeader`只能用作指向实际切片或字符串的`* reflect.SliceHeader`和`* reflect.StringHeader`，而不能用作纯结构。

程序不应声明或分配这些结构类型的变量。
```
	// INVALID: a directly-declared header will not hold Data as a reference.
	var hdr reflect.StringHeader
	hdr.Data = uintptr(unsafe.Pointer(p))
	hdr.Len = n
	s := *(*string)(unsafe.Pointer(&hdr)) // p possibly already lost
```

