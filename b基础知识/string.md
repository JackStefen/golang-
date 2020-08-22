# 1.字符串
一个字符串是一个不可改变的字节序列。和数组不同的是，字符串的元素不可修改，是一个只读的字节数组。
每个字符串的长度虽然也是固定的，但是字符串的长度并不是字符串类型的一部分。
StringHeader是字符串的运行时表示形式。它不能安全或便携地使用，其表示可能在以后的版本中更改。
而且，数据字段不足以保证数据它的引用不会被垃圾收集，所以程序必须保留单独，正确的类型指针指向底层数据。
```
// StringHeader is the runtime representation of a string.
// It cannot be used safely or portably and its representation may
// change in a later release.
// Moreover, the Data field is not sufficient to guarantee the data
// it references will not be garbage collected, so programs must keep
// a separate, correctly typed pointer to the underlying data.
type StringHeader struct {
	Data uintptr
	Len  int
}
```
由此可知，字符串有两部分信息组成，其一为字符串指向的底层字节数组，其二是字符串的字节的长度。
```
package main

import (
   "fmt"
   "reflect"
   "unsafe"
)

func main(){
    s := "this is a test string"
    fmt.Println("len of string: ", (*reflect.StringHeader)(unsafe.Pointer(&s)).Len)
}
```
我们通过一个反汇编示例来看看字符串
```
package main

import (
   "fmt"
)

func main() {
    a := "liudehua "
    b := a + " i love you"
    fmt.Println(b)
}
```
反汇编查看一下
```
   0x0000000001092fe9 <+41>:	lea    0x33794(%rip),%rax        # 0x10c6784 <go.string.*+2628>
   0x0000000001092ff0 <+48>:	mov    %rax,0x8(%rsp)
   0x0000000001092ff5 <+53>:	movq   $0x9,0x10(%rsp)
   0x0000000001092ffe <+62>:	lea    0x339dc(%rip),%rax        # 0x10c69e1 <go.string.*+3233>
   0x0000000001093005 <+69>:	mov    %rax,0x18(%rsp)
   0x000000000109300a <+74>:	movq   $0xb,0x20(%rsp)
```
可以看到，字符串在汇编层面上就是地址（mov    %rax,0x8(%rsp)）和和字符串长度（movq   $0x9,0x10(%rsp)）
```
   0x0000000001093013 <+83>:	callq  0x103f1e0 <runtime.concatstring2>
   0x0000000001093018 <+88>:	mov    0x30(%rsp),%rax
   0x000000000109301d <+93>:	mov    0x28(%rsp),%rcx
```
后面执行字符串拼接操作的是`runtime.concatstring2`.
单步调试可以看到，接下来的函数调用链
```
(gdb) n
Single stepping until exit from function main.main,
which has no line number information.
0x000000000103f1e0 in runtime.concatstring2 ()
(gdb) n
Single stepping until exit from function runtime.concatstring2,
which has no line number information.
0x000000000103ef00 in runtime.concatstrings ()
(gdb) n
Single stepping until exit from function runtime.concatstrings,
which has no line number information.
0x000000000103f4b0 in runtime.rawstringtmp ()
(gdb) n
Single stepping until exit from function runtime.rawstringtmp,
which has no line number information.
0x000000000103fb10 in runtime.rawstring ()
(gdb) n
Single stepping until exit from function runtime.rawstring,
which has no line number information.
0x000000000100a2f0 in runtime.mallocgc ()
(gdb) n
Single stepping until exit from function runtime.mallocgc,
which has no line number information.
0x0000000001050810 in runtime.publicationBarrier ()
(gdb) n
Single stepping until exit from function runtime.publicationBarrier,
which has no line number information.
0x000000000100a5d9 in runtime.mallocgc ()
(gdb) n
Single stepping until exit from function runtime.mallocgc,
which has no line number information.
0x000000000103fb5f in runtime.rawstring ()
(gdb) n
Single stepping until exit from function runtime.rawstring,
which has no line number information.
0x000000000103f522 in runtime.rawstringtmp ()
(gdb) n
Single stepping until exit from function runtime.rawstringtmp,
which has no line number information.
0x000000000103efae in runtime.concatstrings ()
(gdb) n
Single stepping until exit from function runtime.concatstrings,
which has no line number information.
0x000000000103f227 in runtime.concatstring2 ()
(gdb) n
Single stepping until exit from function runtime.concatstring2,
which has no line number information.
0x0000000001093018 in main.main ()
```

`runtime.concatstring2`调用的是`concatstrings`,该函数实现了字符串拼接操作x+y+x...，操作数传入到参数a中。
```
// concatstrings implements a Go string concatenation x+y+z+...
// The operands are passed in the slice a.
// If buf != nil, the compiler has determined that the result does not
// escape the calling function, so the string data can be stored in buf
// if small enough.
func concatstrings(buf *tmpBuf, a []string) string {
	idx := 0
	l := 0
	count := 0
	for i, x := range a {
		n := len(x)
		if n == 0 {
			continue
		}
		if l+n < l {
			throw("string concatenation too long")
		}
		l += n
		count++
		idx = i
	}
	if count == 0 {
		return ""
	}

	// If there is just one string and either it is not on the stack
	// or our result does not escape the calling frame (buf != nil),
	// then we can return that string directly.
	if count == 1 && (buf != nil || !stringDataOnStack(a[idx])) {
		return a[idx]
	}
	s, b := rawstringtmp(buf, l)
	for _, x := range a {
		copy(b, x)
		b = b[len(x):]
	}
	return s
}

func concatstring2(buf *tmpBuf, a [2]string) string {
	return concatstrings(buf, a[:])
}
```

```
func rawstringtmp(buf *tmpBuf, l int) (s string, b []byte) {
	if buf != nil && l <= len(buf) {
		b = buf[:l]
		s = slicebytetostringtmp(b)
	} else {
		s, b = rawstring(l)
	}
	return
}

// rawstring allocates storage for a new string. The returned
// string and byte slice both refer to the same storage.
// The storage is not zeroed. Callers should use
// b to set the string contents and then drop b.
func rawstring(size int) (s string, b []byte) {
	p := mallocgc(uintptr(size), nil, false)

	stringStructOf(&s).str = p
	stringStructOf(&s).len = size

	*(*slice)(unsafe.Pointer(&b)) = slice{p, size, size}

	return
}
```
字符串的结构体表示
```
type stringStruct struct {
	str unsafe.Pointer
	len int
}

func stringStructOf(sp *string) *stringStruct {
	return (*stringStruct)(unsafe.Pointer(sp))
}
```

验证空字符的字节长度
```
package main


import (
    "fmt"
    "unsafe"
)

func main() {
    //var a uint8 //sizeof = 1
    // var b string = "zhao" //sizeof = 16
    // var b *string //sizeof = 8
    var b unsafe.Pointer // sizeof = 8
    fmt.Println(unsafe.Sizeof(b))
}
```

Go语言除了for range语法对UTF8字符串提供了特殊支持外，还对字符串和[]rune类型的相互转换提供了特殊的支持。
- `byte`:int8的别名
- `rune`:int32的别名
```
package main

import (
    "fmt"
)

func main(){
    str := "Hello, 钢铁侠"
    fmt.Println(str)
    for i:=0;i<len(str);i++ {
        fmt.Println(str[i])
    }
    fmt.Println([]byte(str))
    for _, s := range str {
        fmt.Println(s)
    }
    fmt.Println([]rune(str))
    fmt.Println([]int32(str))
}
// ===========
Hello, 钢铁侠
72
101
108
108
111
44
32
233
146
162
233
147
129
228
190
160
[72 101 108 108 111 44 32 233 146 162 233 147 129 228 190 160]
72
101
108
108
111
44
32
38050
38081
20384
[72 101 108 108 111 44 32 38050 38081 20384]
[72 101 108 108 111 44 32 38050 38081 20384]
```
可以看到，通过len()遍历字符串，和`[]byte`的返回结果是一样的，也就是说len获取到的是字符串的字节长度,通过range遍历获取的是rune


