# 1.demo：
注：本示例仅仅语法堆积，没有实际意义
- 创建一个带缓存的channel
- 向channel中放数据
- 遍历channel
- 关闭channel

```
package main 
import (
    "fmt"
)
func putchan(ch chan int) {
    ch <- 2
}

func main(){
    ch := make(chan int, 12)
    for i:=0;i<5;i++ {
        go putchan(ch)
        if i == 4 {
            close(ch)
        }
    }
    for i := range ch {
        fmt.Println(i)
    }
}
```

`go tool compile -S -l -N demo.go`(不内联，不优化)
- putchan函数的汇编
```
"".putchan STEXT size=72 args=0x8 locals=0x18
    0x0000 00000 (demo.go:5)    TEXT    "".putchan(SB), ABIInternal, $24-8
    0x0000 00000 (demo.go:5)    MOVQ    (TLS), CX
    0x0009 00009 (demo.go:5)    CMPQ    SP, 16(CX)
    0x000d 00013 (demo.go:5)    JLS 65
    0x000f 00015 (demo.go:5)    SUBQ    $24, SP      // 栈空间分配
    0x0013 00019 (demo.go:5)    MOVQ    BP, 16(SP)   // 栈基址保存
    0x0018 00024 (demo.go:5)    LEAQ    16(SP), BP   // 栈基址移动
    0x001d 00029 (demo.go:5)    FUNCDATA    $0, gclocals·1a65e721a2ccc325b382662e7ffee780(SB)
    0x001d 00029 (demo.go:5)    FUNCDATA    $1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
    0x001d 00029 (demo.go:5)    FUNCDATA    $2, gclocals·9fb7f0986f647f17cb53dda1484e0f7a(SB)
    0x001d 00029 (demo.go:6)    PCDATA  $0, $1
    0x001d 00029 (demo.go:6)    PCDATA  $1, $1
    0x001d 00029 (demo.go:6)    MOVQ    "".ch+32(SP), AX   // ch参数
    0x0022 00034 (demo.go:6)    PCDATA  $0, $0
    0x0022 00034 (demo.go:6)    MOVQ    AX, (SP)           // ch参数放到SP
    0x0026 00038 (demo.go:6)    PCDATA  $0, $1
    0x0026 00038 (demo.go:6)    LEAQ    ""..stmp_0(SB), AX  // 常数2 
    0x002d 00045 (demo.go:6)    PCDATA  $0, $0
    0x002d 00045 (demo.go:6)    MOVQ    AX, 8(SP)           // 2放到8(SP)
    0x0032 00050 (demo.go:6)    CALL    runtime.chansend1(SB)  // 调用发送函数
    0x0037 00055 (demo.go:7)    MOVQ    16(SP), BP         // 栈基址回放
    0x003c 00060 (demo.go:7)    ADDQ    $24, SP            // 回收栈空间
    0x0040 00064 (demo.go:7)    RET
    0x0041 00065 (demo.go:7)    NOP
    0x0041 00065 (demo.go:5)    PCDATA  $1, $-1
    0x0041 00065 (demo.go:5)    PCDATA  $0, $-1
    0x0041 00065 (demo.go:5)    CALL    runtime.morestack_noctxt(SB)
    0x0046 00070 (demo.go:5)    JMP 0
    0x0000 65 48 8b 0c 25 00 00 00 00 48 3b 61 10 76 32 48  eH..%....H;a.v2H
    0x0010 83 ec 18 48 89 6c 24 10 48 8d 6c 24 10 48 8b 44  ...H.l$.H.l$.H.D
    0x0020 24 20 48 89 04 24 48 8d 05 00 00 00 00 48 89 44  $ H..$H......H.D
    0x0030 24 08 e8 00 00 00 00 48 8b 6c 24 10 48 83 c4 18  $......H.l$.H...
    0x0040 c3 e8 00 00 00 00 eb b8                          ........
    rel 5+4 t=16 TLS+0
    rel 41+4 t=15 ""..stmp_0+0
    rel 51+4 t=8 runtime.chansend1+0
    rel 66+4 t=8 runtime.morestack_noctxt+0
""..stmp_0 SRODATA size=8
    0x0000 02 00 00 00 00 00 00 00                          ........
```
- main函数的汇编

```
"".main STEXT size=417 args=0x0 locals=0xa0
    0x0000 00000 (demo.go:9)    TEXT    "".main(SB), ABIInternal, $160-0
    0x0000 00000 (demo.go:9)    MOVQ    (TLS), CX
    0x0009 00009 (demo.go:9)    LEAQ    -32(SP), AX
    0x000e 00014 (demo.go:9)    CMPQ    AX, 16(CX)
    0x0012 00018 (demo.go:9)    JLS 407
    0x0018 00024 (demo.go:9)    SUBQ    $160, SP
    0x001f 00031 (demo.go:9)    MOVQ    BP, 152(SP)
    0x0027 00039 (demo.go:9)    LEAQ    152(SP), BP
    0x002f 00047 (demo.go:9)    FUNCDATA    $0, gclocals·3e27b3aa6b89137cce48b3379a2a6610(SB)
    0x002f 00047 (demo.go:9)    FUNCDATA    $1, gclocals·2819894ff67942ee916593664c8da755(SB)
    0x002f 00047 (demo.go:9)    FUNCDATA    $2, gclocals·3639c5e889acb2527c3873192ba4ec61(SB)
    0x002f 00047 (demo.go:9)    FUNCDATA    $3, "".main.stkobj(SB)
    0x002f 00047 (demo.go:10)   PCDATA  $0, $1
    0x002f 00047 (demo.go:10)   PCDATA  $1, $0
    0x002f 00047 (demo.go:10)   LEAQ    type.chan int(SB), AX
    0x0036 00054 (demo.go:10)   PCDATA  $0, $0
    0x0036 00054 (demo.go:10)   MOVQ    AX, (SP)
    0x003a 00058 (demo.go:10)   MOVQ    $12, 8(SP)
    0x0043 00067 (demo.go:10)   CALL    runtime.makechan(SB)   //创建channel
    0x0048 00072 (demo.go:10)   PCDATA  $0, $1
    0x0048 00072 (demo.go:10)   MOVQ    16(SP), AX    // 函数的返回值
    0x004d 00077 (demo.go:10)   PCDATA  $0, $0
    0x004d 00077 (demo.go:10)   PCDATA  $1, $1
    0x004d 00077 (demo.go:10)   MOVQ    AX, "".ch+80(SP)     //返回值保存到80(SP)
    0x0052 00082 (demo.go:11)   MOVQ    $0, "".i+64(SP)  //定义变量i
    0x005b 00091 (demo.go:11)   JMP 93
    0x005d 00093 (demo.go:11)   CMPQ    "".i+64(SP), $5    //i和5对比
    0x0063 00099 (demo.go:11)   JLT 103                    // i小于5跳到103
    0x0065 00101 (demo.go:11)   JMP 182                    // i大于5跳到182
    0x0067 00103 (demo.go:12)   MOVL    $8, (SP)
    0x006e 00110 (demo.go:12)   PCDATA  $0, $1
    0x006e 00110 (demo.go:12)   LEAQ    "".putchan·f(SB), AX
    0x0075 00117 (demo.go:12)   PCDATA  $0, $0
    0x0075 00117 (demo.go:12)   MOVQ    AX, 8(SP)
    0x007a 00122 (demo.go:12)   PCDATA  $0, $2
    0x007a 00122 (demo.go:12)   MOVQ    "".ch+80(SP), CX
    0x007f 00127 (demo.go:12)   PCDATA  $0, $0
    0x007f 00127 (demo.go:12)   MOVQ    CX, 16(SP)
    0x0084 00132 (demo.go:12)   CALL    runtime.newproc(SB)
    0x0089 00137 (demo.go:13)   CMPQ    "".i+64(SP), $4     // i和4比较
    0x008f 00143 (demo.go:13)   JEQ 147                     // 等于4跳到147
    0x0091 00145 (demo.go:13)   JMP 180
    0x0093 00147 (demo.go:14)   PCDATA  $0, $1
    0x0093 00147 (demo.go:14)   MOVQ    "".ch+80(SP), AX
    0x0098 00152 (demo.go:14)   PCDATA  $0, $0
    0x0098 00152 (demo.go:14)   MOVQ    AX, (SP)
    0x009c 00156 (demo.go:14)   CALL    runtime.closechan(SB)    // 关闭ch
    0x00a1 00161 (demo.go:14)   JMP 163
    0x00a3 00163 (demo.go:11)   PCDATA  $0, $-2
    0x00a3 00163 (demo.go:11)   PCDATA  $1, $-2
    0x00a3 00163 (demo.go:11)   JMP 165
    0x00a5 00165 (demo.go:11)   PCDATA  $0, $0
    0x00a5 00165 (demo.go:11)   PCDATA  $1, $1
    0x00a5 00165 (demo.go:11)   MOVQ    "".i+64(SP), AX
    0x00aa 00170 (demo.go:11)   INCQ    AX                    // i自增
    0x00ad 00173 (demo.go:11)   MOVQ    AX, "".i+64(SP)
    0x00b2 00178 (demo.go:11)   JMP 93
    0x00b4 00180 (demo.go:13)   PCDATA  $0, $-2
    0x00b4 00180 (demo.go:13)   PCDATA  $1, $-2
    0x00b4 00180 (demo.go:13)   JMP 163
    0x00b6 00182 (demo.go:17)   PCDATA  $0, $1
    0x00b6 00182 (demo.go:17)   PCDATA  $1, $0
    0x00b6 00182 (demo.go:17)   MOVQ    "".ch+80(SP), AX
    0x00bb 00187 (demo.go:17)   PCDATA  $0, $0
    0x00bb 00187 (demo.go:17)   PCDATA  $1, $2
    0x00bb 00187 (demo.go:17)   MOVQ    AX, ""..autotmp_3+104(SP)
    0x00c0 00192 (demo.go:17)   JMP 194
    0x00c2 00194 (demo.go:17)   PCDATA  $0, $1
    0x00c2 00194 (demo.go:17)   MOVQ    ""..autotmp_3+104(SP), AX
    0x00c7 00199 (demo.go:17)   PCDATA  $0, $0
    0x00c7 00199 (demo.go:17)   MOVQ    AX, (SP)
    0x00cb 00203 (demo.go:17)   PCDATA  $0, $1
    0x00cb 00203 (demo.go:17)   LEAQ    ""..autotmp_5+72(SP), AX
    0x00d0 00208 (demo.go:17)   PCDATA  $0, $0
    0x00d0 00208 (demo.go:17)   MOVQ    AX, 8(SP)
    0x00d5 00213 (demo.go:17)   CALL    runtime.chanrecv2(SB)   //读取ch
    0x00da 00218 (demo.go:17)   MOVBLZX 16(SP), AX              // 返回值bool
    0x00df 00223 (demo.go:17)   MOVB    AL, ""..autotmp_6+55(SP)
    0x00e3 00227 (demo.go:17)   TESTB   AL, AL
    0x00e5 00229 (demo.go:17)   JNE 236
    0x00e7 00231 (demo.go:17)   JMP 391
    0x00ec 00236 (demo.go:17)   MOVQ    ""..autotmp_5+72(SP), AX
    0x00f1 00241 (demo.go:17)   MOVQ    AX, "".i+56(SP)
    0x00f6 00246 (demo.go:17)   MOVQ    $0, ""..autotmp_5+72(SP)
    0x00ff 00255 (demo.go:18)   MOVQ    "".i+56(SP), AX
    0x0104 00260 (demo.go:18)   MOVQ    AX, (SP)
    0x0108 00264 (demo.go:18)   CALL    runtime.convT64(SB)
    0x010d 00269 (demo.go:18)   PCDATA  $0, $1
    0x010d 00269 (demo.go:18)   MOVQ    8(SP), AX
    0x0112 00274 (demo.go:18)   PCDATA  $0, $0
    0x0112 00274 (demo.go:18)   PCDATA  $1, $3
    0x0112 00274 (demo.go:18)   MOVQ    AX, ""..autotmp_7+96(SP)
    0x0117 00279 (demo.go:18)   PCDATA  $1, $4
    0x0117 00279 (demo.go:18)   XORPS   X0, X0
    0x011a 00282 (demo.go:18)   MOVUPS  X0, ""..autotmp_4+112(SP)
    0x011f 00287 (demo.go:18)   PCDATA  $0, $1
    0x011f 00287 (demo.go:18)   PCDATA  $1, $3
    0x011f 00287 (demo.go:18)   LEAQ    ""..autotmp_4+112(SP), AX
    0x0124 00292 (demo.go:18)   MOVQ    AX, ""..autotmp_9+88(SP)
    0x0129 00297 (demo.go:18)   TESTB   AL, (AX)
    0x012b 00299 (demo.go:18)   PCDATA  $0, $3
    0x012b 00299 (demo.go:18)   PCDATA  $1, $2
    0x012b 00299 (demo.go:18)   MOVQ    ""..autotmp_7+96(SP), CX
    0x0130 00304 (demo.go:18)   PCDATA  $0, $4
    0x0130 00304 (demo.go:18)   LEAQ    type.int(SB), DX
    0x0137 00311 (demo.go:18)   PCDATA  $0, $3
    0x0137 00311 (demo.go:18)   MOVQ    DX, ""..autotmp_4+112(SP)
    0x013c 00316 (demo.go:18)   PCDATA  $0, $1
    0x013c 00316 (demo.go:18)   MOVQ    CX, ""..autotmp_4+120(SP)
    0x0141 00321 (demo.go:18)   TESTB   AL, (AX)
    0x0143 00323 (demo.go:18)   JMP 325
    0x0145 00325 (demo.go:18)   MOVQ    AX, ""..autotmp_8+128(SP)
    0x014d 00333 (demo.go:18)   MOVQ    $1, ""..autotmp_8+136(SP)
    0x0159 00345 (demo.go:18)   MOVQ    $1, ""..autotmp_8+144(SP)
    0x0165 00357 (demo.go:18)   PCDATA  $0, $0
    0x0165 00357 (demo.go:18)   MOVQ    AX, (SP)
    0x0169 00361 (demo.go:18)   MOVQ    $1, 8(SP)
    0x0172 00370 (demo.go:18)   MOVQ    $1, 16(SP)
    0x017b 00379 (demo.go:18)   CALL    fmt.Println(SB)
    0x0180 00384 (demo.go:18)   JMP 386
    0x0182 00386 (demo.go:17)   PCDATA  $0, $-2
    0x0182 00386 (demo.go:17)   PCDATA  $1, $-2
    0x0182 00386 (demo.go:17)   JMP 194
    0x0187 00391 (<unknown line number>)    PCDATA  $0, $0
    0x0187 00391 (<unknown line number>)    PCDATA  $1, $0
    0x0187 00391 (<unknown line number>)    MOVQ    152(SP), BP
    0x018f 00399 (<unknown line number>)    ADDQ    $160, SP
    0x0196 00406 (<unknown line number>)    RET
    0x0197 00407 (<unknown line number>)    NOP
    0x0197 00407 (demo.go:9)    PCDATA  $1, $-1
    0x0197 00407 (demo.go:9)    PCDATA  $0, $-1
    0x0197 00407 (demo.go:9)    CALL    runtime.morestack_noctxt(SB)
    0x019c 00412 (demo.go:9)    JMP 0
    0x0000 65 48 8b 0c 25 00 00 00 00 48 8d 44 24 e0 48 3b  eH..%....H.D$.H;
    0x0010 41 10 0f 86 7f 01 00 00 48 81 ec a0 00 00 00 48  A.......H......H
    0x0020 89 ac 24 98 00 00 00 48 8d ac 24 98 00 00 00 48  ..$....H..$....H
    0x0030 8d 05 00 00 00 00 48 89 04 24 48 c7 44 24 08 0c  ......H..$H.D$..
    0x0040 00 00 00 e8 00 00 00 00 48 8b 44 24 10 48 89 44  ........H.D$.H.D
    0x0050 24 50 48 c7 44 24 40 00 00 00 00 eb 00 48 83 7c  $PH.D$@......H.|
    0x0060 24 40 05 7c 02 eb 4f c7 04 24 08 00 00 00 48 8d  $@.|..O..$....H.
    0x0070 05 00 00 00 00 48 89 44 24 08 48 8b 4c 24 50 48  .....H.D$.H.L$PH
    0x0080 89 4c 24 10 e8 00 00 00 00 48 83 7c 24 40 04 74  .L$......H.|$@.t
    0x0090 02 eb 21 48 8b 44 24 50 48 89 04 24 e8 00 00 00  ..!H.D$PH..$....
    0x00a0 00 eb 00 eb 00 48 8b 44 24 40 48 ff c0 48 89 44  .....H.D$@H..H.D
    0x00b0 24 40 eb a9 eb ed 48 8b 44 24 50 48 89 44 24 68  $@....H.D$PH.D$h
    0x00c0 eb 00 48 8b 44 24 68 48 89 04 24 48 8d 44 24 48  ..H.D$hH..$H.D$H
    0x00d0 48 89 44 24 08 e8 00 00 00 00 0f b6 44 24 10 88  H.D$........D$..
    0x00e0 44 24 37 84 c0 75 05 e9 9b 00 00 00 48 8b 44 24  D$7..u......H.D$
    0x00f0 48 48 89 44 24 38 48 c7 44 24 48 00 00 00 00 48  HH.D$8H.D$H....H
    0x0100 8b 44 24 38 48 89 04 24 e8 00 00 00 00 48 8b 44  .D$8H..$.....H.D
    0x0110 24 08 48 89 44 24 60 0f 57 c0 0f 11 44 24 70 48  $.H.D$`.W...D$pH
    0x0120 8d 44 24 70 48 89 44 24 58 84 00 48 8b 4c 24 60  .D$pH.D$X..H.L$`
    0x0130 48 8d 15 00 00 00 00 48 89 54 24 70 48 89 4c 24  H......H.T$pH.L$
    0x0140 78 84 00 eb 00 48 89 84 24 80 00 00 00 48 c7 84  x....H..$....H..
    0x0150 24 88 00 00 00 01 00 00 00 48 c7 84 24 90 00 00  $........H..$...
    0x0160 00 01 00 00 00 48 89 04 24 48 c7 44 24 08 01 00  .....H..$H.D$...
    0x0170 00 00 48 c7 44 24 10 01 00 00 00 e8 00 00 00 00  ..H.D$..........
    0x0180 eb 00 e9 3b ff ff ff 48 8b ac 24 98 00 00 00 48  ...;...H..$....H
    0x0190 81 c4 a0 00 00 00 c3 e8 00 00 00 00 e9 5f fe ff  ............._..
    0x01a0 ff                                               .
    rel 5+4 t=16 TLS+0
    rel 50+4 t=15 type.chan int+0
    rel 68+4 t=8 runtime.makechan+0
    rel 113+4 t=15 "".putchan·f+0
    rel 133+4 t=8 runtime.newproc+0
    rel 157+4 t=8 runtime.closechan+0
    rel 214+4 t=8 runtime.chanrecv2+0
    rel 265+4 t=8 runtime.convT64+0
    rel 307+4 t=15 type.int+0
    rel 380+4 t=8 fmt.Println+0
    rel 408+4 t=8 runtime.morestack_noctxt+0

```




# 2.gdb调试

我们可以在`-gcflags=-l`来禁用内联
`go build -gcflags='-l' demo.go`编译后使用gdb来进行调试看看每一步的运行情况
```
➜  channeltest gdb ./demo
GNU gdb (GDB) 8.3
Copyright (C) 2019 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
Type "show copying" and "show warranty" for details.
This GDB was configured as "x86_64-apple-darwin16.7.0".
Type "show configuration" for configuration details.
For bug reporting instructions, please see:
<http://www.gnu.org/software/gdb/bugs/>.
Find the GDB manual and other documentation resources online at:
    <http://www.gnu.org/software/gdb/documentation/>.

For help, type "help".
Type "apropos word" to search for commands related to "word"...
Reading symbols from ./demo...
(No debugging symbols found in ./demo)
Loading Go Runtime support.
(gdb)
```
在gdb中，可以使用`b <content>`命令来增加断点，用于断点分析
- `b main.main`在main函数调用上增加断点
- 也可以使用`(gdb) info add main.putchan`查看函数的地址，然后在指定的地址上设断点：`b *0x1093000(函数地址)`


```
(gdb) b main.main
Breakpoint 1 at 0x1093050
```
运行
```
(gdb) r
Starting program: /Users/zhaojunwei/workspace/src/just.for.test/channeltest/demo
[New Thread 0x1103 of process 10987]
[New Thread 0x1403 of process 10987]
[New Thread 0x1503 of process 10987]
[New Thread 0x1117 of process 10987]
[New Thread 0x1603 of process 10987]
[New Thread 0x1703 of process 10987]
[New Thread 0x1803 of process 10987]

Thread 3 hit Breakpoint 1, 0x0000000001093050 in main.main ()
(gdb)
```
到断点后，就可以分析反汇编程序的运行情况
```
(gdb) disass
Dump of assembler code for function main.main:
=> 0x0000000001093050 <+0>: mov    %gs:0x30,%rcx
   0x0000000001093059 <+9>: cmp    0x10(%rcx),%rsp
   0x000000000109305d <+13>:    jbe    0x109316b <main.main+283>
   0x0000000001093063 <+19>:    sub    $0x60,%rsp
   0x0000000001093067 <+23>:    mov    %rbp,0x58(%rsp)
   0x000000000109306c <+28>:    lea    0x58(%rsp),%rbp
   0x0000000001093071 <+33>:    lea    0xfe68(%rip),%rax        # 0x10a2ee0 <type.*+64768>
   0x0000000001093078 <+40>:    mov    %rax,(%rsp)
   0x000000000109307c <+44>:    movq   $0xc,0x8(%rsp)
   0x0000000001093085 <+53>:    callq  0x1004170 <runtime.makechan>
   0x000000000109308a <+58>:    mov    0x10(%rsp),%rax
   0x000000000109308f <+63>:    mov    %rax,0x40(%rsp)
   0x0000000001093094 <+68>:    xor    %ecx,%ecx
   0x0000000001093096 <+70>:    jmp    0x10930a1 <main.main+81>
   0x0000000001093098 <+72>:    lea    0x1(%rax),%rcx
   0x000000000109309c <+76>:    mov    0x40(%rsp),%rax
   0x00000000010930a1 <+81>:    cmp    $0x5,%rcx
   0x00000000010930a5 <+85>:    jge    0x1093147 <main.main+247>
   0x00000000010930ab <+91>:    mov    %rcx,0x30(%rsp)
   0x00000000010930b0 <+96>:    movl   $0x8,(%rsp)
   0x00000000010930b7 <+103>:   lea    0x39db2(%rip),%rcx        # 0x10cce70 <go.func.*+122>
   0x00000000010930be <+110>:   mov    %rcx,0x8(%rsp)
   0x00000000010930c3 <+115>:   mov    %rax,0x10(%rsp)
   0x00000000010930c8 <+120>:   callq  0x1030f90 <runtime.newproc>
   0x00000000010930cd <+125>:   mov    0x30(%rsp),%rax
   0x00000000010930d2 <+130>:   cmp    $0x4,%rax
   0x00000000010930d6 <+134>:   jne    0x1093098 <main.main+72>
   0x00000000010930d8 <+136>:   mov    0x40(%rsp),%rax
   0x00000000010930dd <+141>:   mov    %rax,(%rsp)
   0x00000000010930e1 <+145>:   callq  0x1004c60 <runtime.closechan>
   0x00000000010930e6 <+150>:   mov    0x30(%rsp),%rax
   0x00000000010930eb <+155>:   jmp    0x1093098 <main.main+72>
   0x00000000010930ed <+157>:   mov    0x38(%rsp),%rax
   0x00000000010930f2 <+162>:   movq   $0x0,0x38(%rsp)
   0x00000000010930fb <+171>:   mov    %rax,(%rsp)
   0x00000000010930ff <+175>:   callq  0x10086a0 <runtime.convT64>
   0x0000000001093104 <+180>:   mov    0x8(%rsp),%rax
   0x0000000001093109 <+185>:   xorps  %xmm0,%xmm0
   0x000000000109310c <+188>:   movups %xmm0,0x48(%rsp)
   0x0000000001093111 <+193>:   lea    0x10348(%rip),%rcx        # 0x10a3460 <type.*+66176>
   0x0000000001093118 <+200>:   mov    %rcx,0x48(%rsp)
   0x000000000109311d <+205>:   mov    %rax,0x50(%rsp)
   0x0000000001093122 <+210>:   lea    0x48(%rsp),%rax
   0x0000000001093127 <+215>:   mov    %rax,(%rsp)
   0x000000000109312b <+219>:   movq   $0x1,0x8(%rsp)
   0x0000000001093134 <+228>:   movq   $0x1,0x10(%rsp)
   0x000000000109313d <+237>:   callq  0x108ca20 <fmt.Println>
   0x0000000001093142 <+242>:   mov    0x40(%rsp),%rax
   0x0000000001093147 <+247>:   mov    %rax,(%rsp)
   0x000000000109314b <+251>:   lea    0x38(%rsp),%rcx
   0x0000000001093150 <+256>:   mov    %rcx,0x8(%rsp)
   0x0000000001093155 <+261>:   callq  0x1004f20 <runtime.chanrecv2>
   0x000000000109315a <+266>:   cmpb   $0x0,0x10(%rsp)
--Type <RET> for more, q to quit, c to continue without paging--
```

继续在putchan函数地址上加断点，然后继续单步调试
```
(gdb) info add main.putchan
Symbol "main.putchan" is at 0x1093000 in a file compiled without debugging.
(gdb) b *0x1093000
Breakpoint 1 at 0x1093000
(gdb) n
Thread 3 hit Breakpoint 1, 0x0000000001093050 in main.main ()
(gdb) n
Single stepping until exit from function main.main,
which has no line number information.
0x0000000001004170 in runtime.makechan ()
(gdb) n
Single stepping until exit from function runtime.makechan,
which has no line number information.
0x000000000100a330 in runtime.mallocgc ()
(gdb) n
Single stepping until exit from function runtime.mallocgc,
which has no line number information.
0x000000000100a100 in runtime.(*mcache).nextFree ()
(gdb) n
Single stepping until exit from function runtime.(*mcache).nextFree,
which has no line number information.
0x00000000010109b0 in runtime.(*mspan).nextFreeIndex ()
(gdb) n
Single stepping until exit from function runtime.(*mspan).nextFreeIndex,
which has no line number information.
0x000000000100a14d in runtime.(*mcache).nextFree ()
(gdb) n
Single stepping until exit from function runtime.(*mcache).nextFree,
which has no line number information.
0x000000000100aa9e in runtime.mallocgc ()
(gdb) n
Single stepping until exit from function runtime.mallocgc,
which has no line number information.
0x0000000001004218 in runtime.makechan ()
(gdb) n
Single stepping until exit from function runtime.makechan,
which has no line number information.
0x000000000109308a in main.main ()
(gdb) n
Single stepping until exit from function main.main,
which has no line number information.
[Switching to Thread 0x1703 of process 11969]

Thread 6 hit Breakpoint 2, 0x0000000001093000 in main.putchan ()
(gdb) disass
Dump of assembler code for function main.putchan:
=> 0x0000000001093000 <+0>: mov    %gs:0x30,%rcx
   0x0000000001093009 <+9>: cmp    0x10(%rcx),%rsp
   0x000000000109300d <+13>:    jbe    0x1093041 <main.putchan+65>
   0x000000000109300f <+15>:    sub    $0x18,%rsp
   0x0000000001093013 <+19>:    mov    %rbp,0x10(%rsp)
   0x0000000001093018 <+24>:    lea    0x10(%rsp),%rbp
   0x000000000109301d <+29>:    mov    0x20(%rsp),%rax
   0x0000000001093022 <+34>:    mov    %rax,(%rsp)
   0x0000000001093026 <+38>:    lea    0x4ad23(%rip),%rax        # 0x10ddd50 <main.statictmp_0>
   0x000000000109302d <+45>:    mov    %rax,0x8(%rsp)
   0x0000000001093032 <+50>:    callq  0x10043a0 <runtime.chansend1>
   0x0000000001093037 <+55>:    mov    0x10(%rsp),%rbp
   0x000000000109303c <+60>:    add    $0x18,%rsp
   0x0000000001093040 <+64>:    retq
   0x0000000001093041 <+65>:    callq  0x104f240 <runtime.morestack_noctxt>
   0x0000000001093046 <+70>:    jmp    0x1093000 <main.putchan>
   0x0000000001093048 <+72>:    int3
   0x0000000001093049 <+73>:    int3
   0x000000000109304a <+74>:    int3
   0x000000000109304b <+75>:    int3
   0x000000000109304c <+76>:    int3
   0x000000000109304d <+77>:    int3
   0x000000000109304e <+78>:    int3
   0x000000000109304f <+79>:    int3
End of assembler dump.
(gdb)
```
# 3.函数调用

## 3.1 makechan
使用n进行单步调试，首先调用`runtime.makechan()`,该函数
```
// 缓冲channel --> ch := make(ch int, 10)
func makechan(t *chantype, size int) *hchan {
    elem := t.elem

    // compiler checks this but be safe.
    if elem.size >= 1<<16 {
        throw("makechan: invalid channel element type")
    }
    if hchanSize%maxAlign != 0 || elem.align > maxAlign {
        throw("makechan: bad alignment")
    }
    //  MulUintptr返回a * b以及乘法是否溢出。在受支持的平台上，这是编译器固有的功能。
    mem, overflow := math.MulUintptr(elem.size, uintptr(size))
    if overflow || mem > maxAlloc-hchanSize || size < 0 {
        panic(plainError("makechan: size out of range"))
    }

    // Hchan does not contain pointers interesting for GC when elements stored in buf do not contain pointers.
    // buf points into the same allocation, elemtype is persistent.
    // SudoG's are referenced from their owning thread so they can't be collected.
    // TODO(dvyukov,rlh): Rethink when collector can move allocated objects.
    var c *hchan
    switch {
    case mem == 0:
        // Queue or element size is zero.
        c = (*hchan)(mallocgc(hchanSize, nil, true))
        // Race detector uses this location for synchronization.
        c.buf = c.raceaddr()
    case elem.kind&kindNoPointers != 0:
        // 元素不包含指针
        // 在一个调用中分配hchan和buf.
        c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
        c.buf = add(unsafe.Pointer(c), hchanSize)
    default:
        // 元素包含指针.
        c = new(hchan)
        c.buf = mallocgc(mem, elem, true)
    }

    c.elemsize = uint16(elem.size)
    c.elemtype = elem
    c.dataqsiz = uint(size)

    if debugChan {
        print("makechan: chan=", c, "; elemsize=", elem.size, "; elemalg=", elem.alg, "; dataqsiz=", size, "\n")
    }
    return c
}
```
makechan函数用来创建channel的，每个channel的结构体原型为：
```
type hchan struct {
    qcount   uint           // total data in the queue
    dataqsiz uint           // 循环队列大小
    buf      unsafe.Pointer // points to an array of dataqsiz elements
    elemsize uint16 // channel的元素类型大小
    closed   uint32 // channel是否关闭的标志位
    elemtype *_type // channel的元素类型
    sendx    uint   // send index
    recvx    uint   // receive index
    recvq    waitq  // list of recv waiters
    sendq    waitq  // list of send waiters

    // lock protects all fields in hchan, as well as several
    // fields in sudogs blocked on this channel.
    //
    // Do not change another G's status while holding this lock
    // (in particular, do not ready a G), as this can deadlock
    // with stack shrinking.
    lock mutex
}
```
makechan第一个参数为chantype类型的指针
```
type chantype struct {
    typ  _type
    elem *_type
    dir  uintptr
}
```
在我们示例中就是，channel类型是int,那么elem.size大小就是32或者64位，
makechan调用mallocgc(hchanSize+mem, nil, true)方法进行资源分配，分配一个size字节大小的对象，小于32kb的小对象直接在P 缓存的空闲列表中分配，大对象则直接在堆上分配。hchanSize为常量，mem值为`8*12=96`
```
const (
    maxAlign  = 8
    hchanSize = unsafe.Sizeof(hchan{}) + uintptr(-int(unsafe.Sizeof(hchan{}))&(maxAlign-1))
)
```

计算后得到hchanSize=96(结构体内的指针类型占用4个字节)
在垃圾回收标记终止阶段（MarkTermination）执行mallocgc()会直接报错，因为该阶段属于STW。
具体看一下该mallocgc函数：
```
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
    if gcphase == _GCmarktermination {
        throw("mallocgc called with gcphase == _GCmarktermination")
    }

    if size == 0 {
        return unsafe.Pointer(&zerobase)
    }

    if debug.sbrk != 0 {
        align := uintptr(16)
        if typ != nil {
            align = uintptr(typ.align)
        }
        return persistentalloc(size, align, &memstats.other_sys)
    }

    // assistG is the G to charge for this allocation, or nil if
    // GC is not currently active.
    var assistG *g
    if gcBlackenEnabled != 0 {
        // Charge the current user G for this allocation.
        assistG = getg()
        if assistG.m.curg != nil {
            assistG = assistG.m.curg
        }
        // Charge the allocation against the G. We'll account
        // for internal fragmentation at the end of mallocgc.
        assistG.gcAssistBytes -= int64(size)

        if assistG.gcAssistBytes < 0 {
            // This G is in debt. Assist the GC to correct
            // this before allocating. This must happen
            // before disabling preemption.
            gcAssistAlloc(assistG)
        }
    }

    // Set mp.mallocing to keep from being preempted by GC.
    mp := acquirem()
    if mp.mallocing != 0 {
        throw("malloc deadlock")
    }
    if mp.gsignal == getg() {
        throw("malloc during signal")
    }
    mp.mallocing = 1

    shouldhelpgc := false
    dataSize := size
    c := gomcache()
    var x unsafe.Pointer
    noscan := typ == nil || typ.kind&kindNoPointers != 0
    if size <= maxSmallSize {
        if noscan && size < maxTinySize {
            // Tiny allocator.
            //
            // Tiny allocator combines several tiny allocation requests
            // into a single memory block. The resulting memory block
            // is freed when all subobjects are unreachable. The subobjects
            // must be noscan (don't have pointers), this ensures that
            // the amount of potentially wasted memory is bounded.
            //
            // Size of the memory block used for combining (maxTinySize) is tunable.
            // Current setting is 16 bytes, which relates to 2x worst case memory
            // wastage (when all but one subobjects are unreachable).
            // 8 bytes would result in no wastage at all, but provides less
            // opportunities for combining.
            // 32 bytes provides more opportunities for combining,
            // but can lead to 4x worst case wastage.
            // The best case winning is 8x regardless of block size.
            //
            // Objects obtained from tiny allocator must not be freed explicitly.
            // So when an object will be freed explicitly, we ensure that
            // its size >= maxTinySize.
            //
            // SetFinalizer has a special case for objects potentially coming
            // from tiny allocator, it such case it allows to set finalizers
            // for an inner byte of a memory block.
            //
            // The main targets of tiny allocator are small strings and
            // standalone escaping variables. On a json benchmark
            // the allocator reduces number of allocations by ~12% and
            // reduces heap size by ~20%.
            off := c.tinyoffset
            // Align tiny pointer for required (conservative) alignment.
            if size&7 == 0 {
                off = round(off, 8)
            } else if size&3 == 0 {
                off = round(off, 4)
            } else if size&1 == 0 {
                off = round(off, 2)
            }
            if off+size <= maxTinySize && c.tiny != 0 {
                // The object fits into existing tiny block.
                x = unsafe.Pointer(c.tiny + off)
                c.tinyoffset = off + size
                c.local_tinyallocs++
                mp.mallocing = 0
                releasem(mp)
                return x
            }
            // Allocate a new maxTinySize block.
            span := c.alloc[tinySpanClass]
            v := nextFreeFast(span)
            if v == 0 {
                v, _, shouldhelpgc = c.nextFree(tinySpanClass)
            }
            x = unsafe.Pointer(v)
            (*[2]uint64)(x)[0] = 0
            (*[2]uint64)(x)[1] = 0
            // See if we need to replace the existing tiny block with the new one
            // based on amount of remaining free space.
            if size < c.tinyoffset || c.tiny == 0 {
                c.tiny = uintptr(x)
                c.tinyoffset = size
            }
            size = maxTinySize
        } else {
            var sizeclass uint8
            if size <= smallSizeMax-8 {
                sizeclass = size_to_class8[(size+smallSizeDiv-1)/smallSizeDiv]
            } else {
                sizeclass = size_to_class128[(size-smallSizeMax+largeSizeDiv-1)/largeSizeDiv]
            }
            size = uintptr(class_to_size[sizeclass])
            spc := makeSpanClass(sizeclass, noscan)
            span := c.alloc[spc]
            v := nextFreeFast(span)
            if v == 0 {
                v, span, shouldhelpgc = c.nextFree(spc)
            }
            x = unsafe.Pointer(v)
            if needzero && span.needzero != 0 {
                memclrNoHeapPointers(unsafe.Pointer(v), size)
            }
        }
    } else {
        var s *mspan
        shouldhelpgc = true
        systemstack(func() {
            s = largeAlloc(size, needzero, noscan)
        })
        s.freeindex = 1
        s.allocCount = 1
        x = unsafe.Pointer(s.base())
        size = s.elemsize
    }

    var scanSize uintptr
    if !noscan {
        // If allocating a defer+arg block, now that we've picked a malloc size
        // large enough to hold everything, cut the "asked for" size down to
        // just the defer header, so that the GC bitmap will record the arg block
        // as containing nothing at all (as if it were unused space at the end of
        // a malloc block caused by size rounding).
        // The defer arg areas are scanned as part of scanstack.
        if typ == deferType {
            dataSize = unsafe.Sizeof(_defer{})
        }
        heapBitsSetType(uintptr(x), size, dataSize, typ)
        if dataSize > typ.size {
            // Array allocation. If there are any
            // pointers, GC has to scan to the last
            // element.
            if typ.ptrdata != 0 {
                scanSize = dataSize - typ.size + typ.ptrdata
            }
        } else {
            scanSize = typ.ptrdata
        }
        c.local_scan += scanSize
    }

    // Ensure that the stores above that initialize x to
    // type-safe memory and set the heap bits occur before
    // the caller can make x observable to the garbage
    // collector. Otherwise, on weakly ordered machines,
    // the garbage collector could follow a pointer to x,
    // but see uninitialized memory or stale heap bits.
    publicationBarrier()

    // Allocate black during GC.
    // All slots hold nil so no scanning is needed.
    // This may be racing with GC so do it atomically if there can be
    // a race marking the bit.
    if gcphase != _GCoff {
        gcmarknewobject(uintptr(x), size, scanSize)
    }

    if raceenabled {
        racemalloc(x, size)
    }

    if msanenabled {
        msanmalloc(x, size)
    }

    mp.mallocing = 0
    releasem(mp)

    if debug.allocfreetrace != 0 {
        tracealloc(x, size, typ)
    }

    if rate := MemProfileRate; rate > 0 {
        if rate != 1 && int32(size) < c.next_sample {
            c.next_sample -= int32(size)
        } else {
            mp := acquirem()
            profilealloc(mp, x, size)
            releasem(mp)
        }
    }

    if assistG != nil {
        // Account for internal fragmentation in the assist
        // debt now that we know it.
        assistG.gcAssistBytes -= int64(size - dataSize)
    }

    if shouldhelpgc {
        if t := (gcTrigger{kind: gcTriggerHeap}); t.test() {
            gcStart(t)
        }
    }

    return x
}

```
如果`hchanSize+ mem <= 32kb`,demo示例中传入mallocgc的size参数值为96+96=192bytes小于32kb，在P缓存的空闲列表中分配，因为本次调用传入参数的typ =nil，所以noscan值true，所以直接从缓存范围内直接返回一个可使用对象。


## 3.2 chansend
单步调试到putchan函数上后，反编译可以看到，在向channel中发送数据的时候，使用的是`runtime.chansend1`
```
func chansend1(c *hchan, elem unsafe.Pointer) {
    chansend(c, elem, true, getcallerpc())
}
/*
 * generic single channel send/recv
 * If block is not nil,
 * then the protocol will not
 * sleep but return if it could
 * not complete.
 *
 * sleep can wake up with g.param == nil
 * when a channel involved in the sleep has
 * been closed.  it is easiest to loop and re-run
 * the operation; we'll see that it's now closed.
 */
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
    if c == nil {
        if !block {
            return false
        }
        gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
        throw("unreachable")
    }

    if debugChan {
        print("chansend: chan=", c, "\n")
    }

    if raceenabled {
        racereadpc(c.raceaddr(), callerpc, funcPC(chansend))
    }

    // Fast path: check for failed non-blocking operation without acquiring the lock.
    //
    // After observing that the channel is not closed, we observe that the channel is
    // not ready for sending. Each of these observations is a single word-sized read
    // (first c.closed and second c.recvq.first or c.qcount depending on kind of channel).
    // Because a closed channel cannot transition from 'ready for sending' to
    // 'not ready for sending', even if the channel is closed between the two observations,
    // they imply a moment between the two when the channel was both not yet closed
    // and not ready for sending. We behave as if we observed the channel at that moment,
    // and report that the send cannot proceed.
    //
    // It is okay if the reads are reordered here: if we observe that the channel is not
    // ready for sending and then observe that it is not closed, that implies that the
    // channel wasn't closed during the first observation.
    if !block && c.closed == 0 && ((c.dataqsiz == 0 && c.recvq.first == nil) ||
        (c.dataqsiz > 0 && c.qcount == c.dataqsiz)) {
        return false
    }

    var t0 int64
    if blockprofilerate > 0 {
        t0 = cputicks()
    }

    lock(&c.lock)

    if c.closed != 0 {
        unlock(&c.lock)
        panic(plainError("send on closed channel"))
    }

    if sg := c.recvq.dequeue(); sg != nil {
        // Found a waiting receiver. We pass the value we want to send
        // directly to the receiver, bypassing the channel buffer (if any).
        send(c, sg, ep, func() { unlock(&c.lock) }, 3)
        return true
    }

    if c.qcount < c.dataqsiz {
        // Space is available in the channel buffer. Enqueue the element to send.
        qp := chanbuf(c, c.sendx)
        if raceenabled {
            raceacquire(qp)
            racerelease(qp)
        }
        typedmemmove(c.elemtype, qp, ep)
        c.sendx++
        if c.sendx == c.dataqsiz {
            c.sendx = 0
        }
        c.qcount++
        unlock(&c.lock)
        return true
    }

    if !block {
        unlock(&c.lock)
        return false
    }

    // Block on the channel. Some receiver will complete our operation for us.
    gp := getg()
    mysg := acquireSudog()
    mysg.releasetime = 0
    if t0 != 0 {
        mysg.releasetime = -1
    }
    // No stack splits between assigning elem and enqueuing mysg
    // on gp.waiting where copystack can find it.
    mysg.elem = ep
    mysg.waitlink = nil
    mysg.g = gp
    mysg.isSelect = false
    mysg.c = c
    gp.waiting = mysg
    gp.param = nil
    c.sendq.enqueue(mysg)
    goparkunlock(&c.lock, waitReasonChanSend, traceEvGoBlockSend, 3)
    // Ensure the value being sent is kept alive until the
    // receiver copies it out. The sudog has a pointer to the
    // stack object, but sudogs aren't considered as roots of the
    // stack tracer.
    KeepAlive(ep)

    // someone woke us up.
    if mysg != gp.waiting {
        throw("G waiting list is corrupted")
    }
    gp.waiting = nil
    if gp.param == nil {
        if c.closed == 0 {
            throw("chansend: spurious wakeup")
        }
        panic(plainError("send on closed channel"))
    }
    gp.param = nil
    if mysg.releasetime > 0 {
        blockevent(mysg.releasetime-t0, 2)
    }
    mysg.c = nil
    releaseSudog(mysg)
    return true
}
```
当i=4的时候由于执行了close操作，后续调度到相关写入channel的goroutine时，将会报错
```
Thread 6 hit Breakpoint 2, 0x0000000001093000 in main.putchan ()
(gdb) n
Single stepping until exit from function main.putchan,
which has no line number information.
0x00000000010043a0 in runtime.chansend1 ()
(gdb) n
Single stepping until exit from function runtime.chansend1,
which has no line number information.
0x00000000010043e0 in runtime.chansend ()
(gdb) n
Single stepping until exit from function runtime.chansend,
which has no line number information.
0x0000000001008d50 in runtime.lock ()
(gdb) n
Single stepping until exit from function runtime.lock,
which has no line number information.
0x0000000001004484 in runtime.chansend ()
(gdb) n
Single stepping until exit from function runtime.chansend,
which has no line number information.
[New Thread 0x1a53 of process 11969]
2
[Switching to Thread 0x1803 of process 11969]

Thread 7 hit Breakpoint 2, 0x0000000001093000 in main.putchan ()
(gdb) n
Single stepping until exit from function main.putchan,
which has no line number information.
panic: send on closed channel

goroutine 22 [running]:
main.putchan(0xc000084000)
    /Users/xxx/workspace/src/just.for.test/channeltest/demo.go:6 +0x37
created by main.main
    /Users/xxx/workspace/src/just.for.test/channeltest/demo.go:12 +0x7d
[Inferior 1 (process 11969) exited with code 02]
(gdb)
```
## 3.3 chanrecv 
从channel中读取数据使用的是`runtime.chanrecv2`,
```
// entry points for <- c from compiled code
//go:nosplit
func chanrecv1(c *hchan, elem unsafe.Pointer) {
    chanrecv(c, elem, true)
}

//go:nosplit
func chanrecv2(c *hchan, elem unsafe.Pointer) (received bool) {
    _, received = chanrecv(c, elem, true)
    return
}
// chanrev 从 channel中获取数据并写入到ep
// ep可能为nil,这种情况下，获得的数据将被忽略
// 如果block参数为false并且elements是可以使用的，返回false, false
// 否则，如果channel是关闭的，将ep置零值并返回true, false
// 否则，将一个元素填充到ep并返回true,true
// 一个非空的ep必须指向堆或者调用者栈
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
    // raceenabled: don't need to check ep, as it is always on the stack
    // or is new memory allocated by reflect.

    if debugChan {
        print("chanrecv: chan=", c, "\n")
    }

    if c == nil {
        if !block {
            return
        }
        gopark(nil, nil, waitReasonChanReceiveNilChan, traceEvGoStop, 2)
        throw("unreachable")
    }

    // Fast path: check for failed non-blocking operation without acquiring the lock.
    //
    // After observing that the channel is not ready for receiving, we observe that the
    // channel is not closed. Each of these observations is a single word-sized read
    // (first c.sendq.first or c.qcount, and second c.closed).
    // Because a channel cannot be reopened, the later observation of the channel
    // being not closed implies that it was also not closed at the moment of the
    // first observation. We behave as if we observed the channel at that moment
    // and report that the receive cannot proceed.
    //
    // The order of operations is important here: reversing the operations can lead to
    // incorrect behavior when racing with a close.
    if !block && (c.dataqsiz == 0 && c.sendq.first == nil ||
        c.dataqsiz > 0 && atomic.Loaduint(&c.qcount) == 0) &&
        atomic.Load(&c.closed) == 0 {
        return
    }

    var t0 int64
    if blockprofilerate > 0 {
        t0 = cputicks()
    }

    lock(&c.lock)

    if c.closed != 0 && c.qcount == 0 {
        if raceenabled {
            raceacquire(c.raceaddr())
        }
        unlock(&c.lock)
        if ep != nil {
            typedmemclr(c.elemtype, ep)
        }
        return true, false
    }

    if sg := c.sendq.dequeue(); sg != nil {
        // Found a waiting sender. If buffer is size 0, receive value
        // directly from sender. Otherwise, receive from head of queue
        // and add sender's value to the tail of the queue (both map to
        // the same buffer slot because the queue is full).
        recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
        return true, true
    }

    if c.qcount > 0 {
        // Receive directly from queue
        qp := chanbuf(c, c.recvx)
        if raceenabled {
            raceacquire(qp)
            racerelease(qp)
        }
        if ep != nil {
            typedmemmove(c.elemtype, ep, qp)
        }
        typedmemclr(c.elemtype, qp)
        c.recvx++
        if c.recvx == c.dataqsiz {
            c.recvx = 0
        }
        c.qcount--
        unlock(&c.lock)
        return true, true
    }

    if !block {
        unlock(&c.lock)
        return false, false
    }

    // no sender available: block on this channel.
    gp := getg()
    mysg := acquireSudog()
    mysg.releasetime = 0
    if t0 != 0 {
        mysg.releasetime = -1
    }
    // No stack splits between assigning elem and enqueuing mysg
    // on gp.waiting where copystack can find it.
    mysg.elem = ep
    mysg.waitlink = nil
    gp.waiting = mysg
    mysg.g = gp
    mysg.isSelect = false
    mysg.c = c
    gp.param = nil
    c.recvq.enqueue(mysg)
    goparkunlock(&c.lock, waitReasonChanReceive, traceEvGoBlockRecv, 3)

    // someone woke us up
    if mysg != gp.waiting {
        throw("G waiting list is corrupted")
    }
    gp.waiting = nil
    if mysg.releasetime > 0 {
        blockevent(mysg.releasetime-t0, 2)
    }
    closed := gp.param == nil
    gp.param = nil
    mysg.c = nil
    releaseSudog(mysg)
    return true, !closed
}
```

