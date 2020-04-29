map是一堆键值对的未排序集合，类似Python中字典的概念，它的格式为map[keyType]valueType，是一个key-value的hash结构。map的读取和设置也类似slice一样，通过key来操作，只是slice的index只能是int类型，而map多了很多类型，可以是int，可以是string及所有完全定义了==与!=操作的类型。
map线程不安全的数据结构，如果需要线程安全，需要加锁，或者直接只用sync包中的map

# 1.声明
`var 变量名 map[keytype] valuetype`
# 2.初始化
```
var amay map[int] string = map[int] string {1:"zhao",2:"qian",3:"sun",4:"li"}
var bmap map[string] string = make(map[string] string)
var cmap map[string] string = make(map[string] string, 4) //预先设置键值对数量，避免添加键值对时，重新分配内存
```
# 3.元素查找
```
value, ok := myMap["1234"]
if ok{
    //处理找到的value
}
```
# 4.元素删除
```
delete(bmap, "location")
```
# 5.遍历
```
for k, v := range amap {
    .....
}
```
**注：由于v只是一个值的拷贝，所以，如果想对map的修改生效，必须amap[i] = xxx这种形式**
```
for i := range amap {
    amap[i] = xxx
}
```
字典设置键值，修改键值，都通过map[key]=value这种形式
# 6.demo
```
package main

import (
    "fmt"
)

func main(){
    mapa:= make(map[string]int, 1)
    mapa["zhao"] = 1
    mapa["qian"] = 2
    fmt.Println(mapa)
}
```
反汇编后
```
   0x0000000001092fe1 <+33>:    callq  0x100ba00 <runtime.makemap_small>
   0x0000000001092fe6 <+38>:    mov    (%rsp),%rax
   0x0000000001092fea <+42>:    mov    %rax,0x30(%rsp)
   0x0000000001092fef <+47>:    lea    0x1772a(%rip),%rcx        # 0x10aa720 <type.*+95744>
   0x0000000001092ff6 <+54>:    mov    %rcx,(%rsp)
   0x0000000001092ffa <+58>:    mov    %rax,0x8(%rsp)
   0x0000000001092fff <+63>:    lea    0x32f4e(%rip),%rdx        # 0x10c5f54 <go.string.*+500>
   0x0000000001093006 <+70>:    mov    %rdx,0x10(%rsp)
   0x000000000109300b <+75>:    movq   $0x4,0x18(%rsp)
   0x0000000001093014 <+84>:    callq  0x100fa70 <runtime.mapassign_faststr>
   0x0000000001093019 <+89>:    mov    0x20(%rsp),%rax
   0x000000000109301e <+94>:    movq   $0x1,(%rax)
   0x0000000001093025 <+101>:   lea    0x176f4(%rip),%rax        # 0x10aa720 <type.*+95744>
   0x000000000109302c <+108>:   mov    %rax,(%rsp)
   0x0000000001093030 <+112>:   mov    0x30(%rsp),%rcx
   0x0000000001093035 <+117>:   mov    %rcx,0x8(%rsp)
   0x000000000109303a <+122>:   lea    0x32ef3(%rip),%rdx        # 0x10c5f34 <go.string.*+468>
   0x0000000001093041 <+129>:   mov    %rdx,0x10(%rsp)
   0x0000000001093046 <+134>:   movq   $0x4,0x18(%rsp)
   0x000000000109304f <+143>:   callq  0x100fa70 <runtime.mapassign_faststr>
   0x0000000001093054 <+148>:   mov    0x20(%rsp),%rax
   0x0000000001093059 <+153>:   movq   $0x2,(%rax)
   0x0000000001093060 <+160>:   xorps  %xmm0,%xmm0
   0x0000000001093063 <+163>:   movups %xmm0,0x38(%rsp)
   0x0000000001093068 <+168>:   lea    0x176b1(%rip),%rax        # 0x10aa720 <type.*+95744>
   0x000000000109306f <+175>:   mov    %rax,0x38(%rsp)
   0x0000000001093074 <+180>:   mov    0x30(%rsp),%rax
   0x0000000001093079 <+185>:   mov    %rax,0x40(%rsp)
   0x000000000109307e <+190>:   lea    0x38(%rsp),%rax
   0x0000000001093083 <+195>:   mov    %rax,(%rsp)
   0x0000000001093087 <+199>:   movq   $0x1,0x8(%rsp)
   0x0000000001093090 <+208>:   movq   $0x1,0x10(%rsp)
   0x0000000001093099 <+217>:   callq  0x108c9e0 <fmt.Println>
```
创建map执行的是`runtime.makemap_small`,该方法在使用make进行创建map的时候使用。此时map分配在堆上。

```
func makemap_small() *hmap {
    h := new(hmap)
    h.hash0 = fastrand()
    return h
}
```

该函数返回的就是map的指针,其结构体原型为：
```
// A header for a Go map.
type hmap struct {
    count     int // 元素个数，必须位于第一项，调用 len(map) 时，直接返回此值
    flags     uint8
    B         uint8  // B 是 buckets 数组的长度以2为底的对数, (可以容纳 loadFactor * 2^B 个元素)
    noverflow uint16 // 溢出桶的大概数量
    hash0     uint32 // 哈希种子

    buckets    unsafe.Pointer // 数据存储桶
    oldbuckets unsafe.Pointer // 扩容前的桶，只有在扩容的时候非空。
    nevacuate  uintptr        // progress counter for evacuation (buckets less than this have been evacuated)

    extra *mapextra // 可选字段
}

// 桶的数据结构原型结构体.
type bmap struct {
    // tophash 通常包含桶中每个键的hash值的最高字节(8位)
    tophash [bucketCnt]uint8   // 8字节长度的数组
    // 之后是bucketCnt数量的键，然后是bucketCnt数量的值。
    // 后跟一个溢出指针。
}
```
Map就是一个哈希表。 数据被安排在一系列存储桶中。 每个存储桶最多包含8个键/值对。键的哈希值的低位用于选择存储桶。 每个存储桶包含每个键哈希值的高阶位，以区分单个存储桶中的条目。如果有8个以上的键散列到存储桶中，我们将链接
额外的桶。当哈希表增长时，我们将分配一个新的存储桶数组作为两倍大。 将存储桶以增量方式从旧存储桶阵列复制到新存储桶阵列。映射迭代器遍历存储桶数组，并按行走顺序返回键（存储桶编号，然后是溢出链顺序，然后是存储桶索引）。 为了维持迭代语义，我们绝不会在键的存储桶中移动键（如果这样做，键可能会返回0或2次）。 在扩大表时，迭代器将继续在旧表中进行迭代，并且必须检查新表是否将要迭代的存储桶移至新表中。

demo中map的键值插入使用的是`runtime.mapassign_faststr`  
```
func mapassign_faststr(t *maptype, h *hmap, s string) unsafe.Pointer {
    if h == nil {
        panic(plainError("assignment to entry in nil map"))
    }
    if raceenabled {
        callerpc := getcallerpc()
        racewritepc(unsafe.Pointer(h), callerpc, funcPC(mapassign_faststr))
    }
    // 检测该map是否存在其他协程进行修改操作，如果有则直接报错
    if h.flags&hashWriting != 0 {
        throw("concurrent map writes")
    }
    // 对键进行哈希
    key := stringStructOf(&s)
    hash := t.key.alg.hash(noescape(unsafe.Pointer(&s)), uintptr(h.hash0))

    // 哈希后设置hashWriting防止其他协程进行同时修改
    h.flags ^= hashWriting

    // 如果是新的map，buchket数量为0，则先进行资源分配
    if h.buckets == nil {
        h.buckets = newobject(t.bucket) // newarray(t.bucket, 1)
    }

again:
    // 获取桶的编号
    // 例如：hmap的B值为3，bucketMask()之后值为7，二进制值为0111，与hash位与之后，获取的是hash值最后三位
    // 也就是键的hash值的最后B位的值
    bucket := hash & bucketMask(h.B)
    if h.growing() {
        growWork_faststr(t, h, bucket)
    }
    // 定位桶的位置
    b := (*bmap)(unsafe.Pointer(uintptr(h.buckets) + bucket*uintptr(t.bucketsize)))
    // 键的hash值的高字节位，此来获取Key在桶中的位置
    top := tophash(hash)

    var insertb *bmap
    var inserti uintptr
    var insertk unsafe.Pointer

bucketloop:
    for {
        // 在桶中遍历
        for i := uintptr(0); i < bucketCnt; i++ {
            if b.tophash[i] != top {
                if isEmpty(b.tophash[i]) && insertb == nil {
                    insertb = b
                    inserti = i
                }
                if b.tophash[i] == emptyRest {
                    break bucketloop
                }
                continue
            }
            k := (*stringStruct)(add(unsafe.Pointer(b), dataOffset+i*2*sys.PtrSize))
            if k.len != key.len {
                continue
            }
            if k.str != key.str && !memequal(k.str, key.str, uintptr(key.len)) {
                continue
            }
            // 如果已经存在此键，则直接更新它
            inserti = i
            insertb = b
            goto done
        }
        ovf := b.overflow(t)
        if ovf == nil {
            break
        }
        b = ovf
    }

    // 如果没有找到键，则分配新的内存，并添加键进去

    // 如果我们达到了最大的负载因子，或者有太多的溢出桶
    // 此时又不处在扩容的过程中，那么就可以开始新的扩容了
    if !h.growing() && (overLoadFactor(h.count+1, h.B) || tooManyOverflowBuckets(h.noverflow, h.B)) {
        hashGrow(t, h)
        goto again // 扩容后需要重新检测桶的编码和键的位置, so try again
    }

    if insertb == nil {
        // 当前的桶都慢了，分配新的
        insertb = h.newoverflow(t, b)
        inserti = 0 // not necessary, but avoids needlessly spilling inserti
    }
    insertb.tophash[inserti&(bucketCnt-1)] = top // mask inserti to avoid bounds checks

    insertk = add(unsafe.Pointer(insertb), dataOffset+inserti*2*sys.PtrSize)
    // 插入新的key
    *((*stringStruct)(insertk)) = *key
    // 元素数量加一
    h.count++

done:
    // 键值value的地址
    val := add(unsafe.Pointer(insertb), dataOffset+bucketCnt*2*sys.PtrSize+inserti*uintptr(t.valuesize))
    if h.flags&hashWriting == 0 {
        throw("concurrent map writes")
    }
    // 解除写标志
    h.flags &^= hashWriting
    return val
}
```
对比一下反汇编的代码来看一下这个函数的参数：`(t *maptype, h *hmap, s string)`
```
   0x0000000001092fe1 <+33>:    callq  0x100ba00 <runtime.makemap_small>
   0x0000000001092fe6 <+38>:    mov    (%rsp),%rax
   0x0000000001092fea <+42>:    mov    %rax,0x30(%rsp)
   0x0000000001092fef <+47>:    lea    0x1772a(%rip),%rcx        # 0x10aa720 <type.*+95744>
   0x0000000001092ff6 <+54>:    mov    %rcx,(%rsp)
   0x0000000001092ffa <+58>:    mov    %rax,0x8(%rsp)
   0x0000000001092fff <+63>:    lea    0x32f4e(%rip),%rdx        # 0x10c5f54 <go.string.*+500>
   0x0000000001093006 <+70>:    mov    %rdx,0x10(%rsp)
   0x000000000109300b <+75>:    movq   $0x4,0x18(%rsp)
```
在调用`runtime.makemap_small`后返回的hmap指针，这个是作为`runtime.mapassign_fasterstr`第二个参数，
rip寄存器中的指令作为第一个参数，rip寄存器0x32f4e地址处的值作为第三个参数，因为是字符串，所以在放入到rsp寄存器时，把字符串长度也作为字符串的一部分放进去了。

在该函数中，我们最后得到的仅仅是value的地址，并没有实际进行`ma[key] = value`的赋值操作，那么实际意义上的赋值操作在哪里进行的呢？
再看一下反汇编的代码：
```
   0x000000000109304f <+143>:   callq  0x100fa70 <runtime.mapassign_faststr>
   0x0000000001093054 <+148>:   mov    0x20(%rsp),%rax
   0x0000000001093059 <+153>:   movq   $0x2,(%rax)
```
