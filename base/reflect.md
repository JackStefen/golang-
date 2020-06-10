有关反射的内容，即多又重要，可能平时业务上用的较少，但是设计到抽象出来的功能时，比如基础性设施的开发时，会比较多的使用，本文我们将以一个示例开始我们的学习反射之旅，内容可能无法涵盖方方面面，咱们一点点的来。
```
package main

import (
    "fmt"
    "reflect"
)

type User struct {
    name string
    age  int
}

func (u User) Descrite() {
    fmt.Println("descrite method: ", u.name)
}

func main() {

    var user = User{
        name: "myname",
        age:  13,
    }
    t := reflect.TypeOf(user)
    fmt.Println(t.Method(0))
    v := reflect.ValueOf(user)
    mName:=v.Method(0)
    args := []reflect.Value{}
    mName.Call(args)
}

```

上面的程序很简单，就是定义了一个结构体，一个结构体方法，然后通过反射获取其类型，以及通过反射对其上的方法调用。

```
➜  reflectdemo go tool compile -S  demo3.go | grep -A 300 '.main STEXT size=676 args=0x0 locals=0x108'
"".main STEXT size=676 args=0x0 locals=0x108
	0x0000 00000 (demo3.go:17)	TEXT	"".main(SB), REFLECTMETHOD|ABIInternal, $264-0
	0x0000 00000 (demo3.go:17)	MOVQ	(TLS), CX
	0x0009 00009 (demo3.go:17)	LEAQ	-136(SP), AX
	0x0011 00017 (demo3.go:17)	CMPQ	AX, 16(CX)
	0x0015 00021 (demo3.go:17)	JLS	666
	0x001b 00027 (demo3.go:17)	SUBQ	$264, SP
	0x0022 00034 (demo3.go:17)	MOVQ	BP, 256(SP)
	0x002a 00042 (demo3.go:17)	LEAQ	256(SP), BP
	0x0032 00050 (demo3.go:17)	FUNCDATA	$0, gclocals·f14a5bc6d08bc46424827f54d2e3f8ed(SB)
	0x0032 00050 (demo3.go:17)	FUNCDATA	$1, gclocals·f0a6b772a0f7c24b4b4dd3ced41cc7ef(SB)
	0x0032 00050 (demo3.go:17)	FUNCDATA	$3, gclocals·8a0ea43e9c91b3513a59c3c1cfe6b709(SB)
	0x0032 00050 (demo3.go:17)	FUNCDATA	$4, "".main.stkobj(SB)
	0x0032 00050 (demo3.go:23)	PCDATA	$2, $1
	0x0032 00050 (demo3.go:23)	PCDATA	$0, $1
	0x0032 00050 (demo3.go:23)	LEAQ	go.string."myname"(SB), AX  
	0x0039 00057 (demo3.go:23)	PCDATA	$2, $0
	0x0039 00057 (demo3.go:23)	MOVQ	AX, ""..autotmp_26+152(SP)   // 结构体name字段赋值
	0x0041 00065 (demo3.go:23)	MOVQ	$6, ""..autotmp_26+160(SP)
	0x004d 00077 (demo3.go:23)	MOVQ	$13, ""..autotmp_26+168(SP)  // 结构体age字段赋值
	0x0059 00089 (demo3.go:23)	PCDATA	$2, $2
	0x0059 00089 (demo3.go:23)	LEAQ	type."".User(SB), CX         
	0x0060 00096 (demo3.go:23)	PCDATA	$2, $0
	0x0060 00096 (demo3.go:23)	MOVQ	CX, (SP)                     // User结构体类型作为第一个参数
	0x0064 00100 (demo3.go:23)	PCDATA	$2, $3
	0x0064 00100 (demo3.go:23)	PCDATA	$0, $0
	0x0064 00100 (demo3.go:23)	LEAQ	""..autotmp_26+152(SP), DX   // 结构体值作为第二个参数
	0x006c 00108 (demo3.go:23)	PCDATA	$2, $0
	0x006c 00108 (demo3.go:23)	MOVQ	DX, 8(SP)
	0x0071 00113 (demo3.go:23)	CALL	runtime.convT2E(SB)          // 调用runtime.convT2E,该函数将结构体类型转化为接口类型
	0x0076 00118 (demo3.go:23)	MOVQ	16(SP), AX
	0x007b 00123 (demo3.go:23)	PCDATA	$2, $2
	0x007b 00123 (demo3.go:23)	MOVQ	24(SP), CX
	0x0080 00128 (demo3.go:23)	PCDATA	$0, $2
	0x0080 00128 (demo3.go:23)	MOVQ	AX, reflect.i+104(SP)        // 将接口类型的参数进行赋值i变量，_type字段
	0x0085 00133 (demo3.go:23)	PCDATA	$2, $0
	0x0085 00133 (demo3.go:23)	MOVQ	CX, reflect.i+112(SP)       // 接口的data字段
	0x008a 00138 (demo3.go:23)	XCHGL	AX, AX
	0x008b 00139 ($GOROOT/src/reflect/type.go:1375)	PCDATA	$2, $1
	0x008b 00139 ($GOROOT/src/reflect/type.go:1375)	PCDATA	$0, $0
	0x008b 00139 ($GOROOT/src/reflect/type.go:1375)	MOVQ	reflect.i+104(SP), AX
	0x0090 00144 ($GOROOT/src/reflect/type.go:1376)	XCHGL	AX, AX
	0x0091 00145 ($GOROOT/src/reflect/type.go:3004)	TESTQ	AX, AX
	0x0094 00148 (:0)	JEQ	657
	0x009a 00154 (:0)	LEAQ	go.itab.*reflect.rtype,reflect.Type(SB), CX  // 设置go.itab.*reflect.rtype,reflect.Type
	0x00a1 00161 (demo3.go:24)	MOVQ	168(CX), CX                          // 168=160+8=0xa0+8 提取接口参数的Method方法
	0x00a8 00168 (demo3.go:24)	PCDATA	$2, $0
	0x00a8 00168 (demo3.go:24)	MOVQ	AX, (SP)                      // 参数i接口类型
	0x00ac 00172 (demo3.go:24)	MOVQ	$0, 8(SP)                     // Method方法调用参数：0
	0x00b5 00181 (demo3.go:24)	CALL	CX                           // 调用Type接口Method方法
	0x00b7 00183 (demo3.go:24)	PCDATA	$2, $4
	0x00b7 00183 (demo3.go:24)	PCDATA	$0, $3
	0x00b7 00183 (demo3.go:24)	LEAQ	""..autotmp_28+176(SP), DI
	0x00bf 00191 (demo3.go:24)	PCDATA	$2, $5
	0x00bf 00191 (demo3.go:24)	LEAQ	16(SP), SI
	0x00c4 00196 (demo3.go:24)	PCDATA	$2, $0
	0x00c4 00196 (demo3.go:24)	DUFFCOPY	$826
	0x00d7 00215 (demo3.go:24)	PCDATA	$2, $1
	0x00d7 00215 (demo3.go:24)	LEAQ	type.reflect.Method(SB), AX
	0x00de 00222 (demo3.go:24)	PCDATA	$2, $0
	0x00de 00222 (demo3.go:24)	MOVQ	AX, (SP)     // type.reflect.Method地址作为第一个参数
	0x00e2 00226 (demo3.go:24)	PCDATA	$2, $1
	0x00e2 00226 (demo3.go:24)	PCDATA	$0, $0
	0x00e2 00226 (demo3.go:24)	LEAQ	""..autotmp_28+176(SP), AX
	0x00ea 00234 (demo3.go:24)	PCDATA	$2, $0
	0x00ea 00234 (demo3.go:24)	MOVQ	AX, 8(SP)      // 
	0x00ef 00239 (demo3.go:24)	CALL	runtime.convT2E(SB)  // 将Method结构体进行接口转换以便调用输出
	0x00f4 00244 (demo3.go:24)	MOVQ	16(SP), AX
	0x00f9 00249 (demo3.go:24)	PCDATA	$2, $2
	0x00f9 00249 (demo3.go:24)	MOVQ	24(SP), CX
	0x00fe 00254 (demo3.go:24)	PCDATA	$0, $4
	0x00fe 00254 (demo3.go:24)	XORPS	X0, X0
	0x0101 00257 (demo3.go:24)	MOVUPS	X0, ""..autotmp_37+136(SP)
	0x0109 00265 (demo3.go:24)	MOVQ	AX, ""..autotmp_37+136(SP)
	0x0111 00273 (demo3.go:24)	PCDATA	$2, $0
	0x0111 00273 (demo3.go:24)	MOVQ	CX, ""..autotmp_37+144(SP)
	0x0119 00281 (demo3.go:24)	XCHGL	AX, AX
	0x011a 00282 ($GOROOT/src/fmt/print.go:275)	PCDATA	$2, $1
	0x011a 00282 ($GOROOT/src/fmt/print.go:275)	MOVQ	os.Stdout(SB), AX
	0x0121 00289 ($GOROOT/src/fmt/print.go:275)	PCDATA	$2, $6
	0x0121 00289 ($GOROOT/src/fmt/print.go:275)	LEAQ	go.itab.*os.File,io.Writer(SB), CX
	0x0128 00296 ($GOROOT/src/fmt/print.go:275)	PCDATA	$2, $1
	0x0128 00296 ($GOROOT/src/fmt/print.go:275)	MOVQ	CX, (SP)
	0x012c 00300 ($GOROOT/src/fmt/print.go:275)	PCDATA	$2, $0
	0x012c 00300 ($GOROOT/src/fmt/print.go:275)	MOVQ	AX, 8(SP)
	0x0131 00305 ($GOROOT/src/fmt/print.go:275)	PCDATA	$2, $1
	0x0131 00305 ($GOROOT/src/fmt/print.go:275)	PCDATA	$0, $0
	0x0131 00305 ($GOROOT/src/fmt/print.go:275)	LEAQ	""..autotmp_37+136(SP), AX
	0x0139 00313 ($GOROOT/src/fmt/print.go:275)	PCDATA	$2, $0
	0x0139 00313 ($GOROOT/src/fmt/print.go:275)	MOVQ	AX, 16(SP)
	0x013e 00318 ($GOROOT/src/fmt/print.go:275)	MOVQ	$1, 24(SP)
	0x0147 00327 ($GOROOT/src/fmt/print.go:275)	MOVQ	$1, 32(SP)
	0x0150 00336 ($GOROOT/src/fmt/print.go:275)	CALL	fmt.Fprintln(SB)
	0x0155 00341 (demo3.go:25)	PCDATA	$2, $1
	0x0155 00341 (demo3.go:25)	PCDATA	$0, $1
	0x0155 00341 (demo3.go:25)	LEAQ	go.string."myname"(SB), AX   // 此部分同上
	0x015c 00348 (demo3.go:25)	PCDATA	$2, $0
	0x015c 00348 (demo3.go:25)	MOVQ	AX, ""..autotmp_26+152(SP)
	0x0164 00356 (demo3.go:25)	MOVQ	$6, ""..autotmp_26+160(SP)
	0x0170 00368 (demo3.go:25)	MOVQ	$13, ""..autotmp_26+168(SP)
	0x017c 00380 (demo3.go:25)	PCDATA	$2, $1
	0x017c 00380 (demo3.go:25)	LEAQ	type."".User(SB), AX
	0x0183 00387 (demo3.go:25)	PCDATA	$2, $0
	0x0183 00387 (demo3.go:25)	MOVQ	AX, (SP)
	0x0187 00391 (demo3.go:25)	PCDATA	$2, $1
	0x0187 00391 (demo3.go:25)	PCDATA	$0, $0
	0x0187 00391 (demo3.go:25)	LEAQ	""..autotmp_26+152(SP), AX
	0x018f 00399 (demo3.go:25)	PCDATA	$2, $0
	0x018f 00399 (demo3.go:25)	MOVQ	AX, 8(SP)
	0x0194 00404 (demo3.go:25)	CALL	runtime.convT2E(SB)      // 参数接口化
	0x0199 00409 (demo3.go:25)	MOVQ	16(SP), AX
	0x019e 00414 (demo3.go:25)	PCDATA	$2, $2
	0x019e 00414 (demo3.go:25)	MOVQ	24(SP), CX
	0x01a3 00419 (demo3.go:25)	XCHGL	AX, AX
	0x01a4 00420 ($GOROOT/src/reflect/value.go:2253)	TESTQ	AX, AX
	0x01a7 00423 (:0)	JEQ	646
	0x01ad 00429 ($GOROOT/src/reflect/value.go:2261)	XCHGL	AX, AX
	0x01ae 00430 ($GOROOT/src/reflect/value.go:2708)	CMPB	reflect.dummy(SB), $0
	0x01b5 00437 ($GOROOT/src/reflect/value.go:2708)	JEQ	466
	0x01b7 00439 ($GOROOT/src/reflect/value.go:2709)	MOVQ	AX, reflect.dummy+8(SB)
	0x01be 00446 ($GOROOT/src/reflect/value.go:2709)	PCDATA	$2, $-2
	0x01be 00446 ($GOROOT/src/reflect/value.go:2709)	PCDATA	$0, $-2
	0x01be 00446 ($GOROOT/src/reflect/value.go:2709)	CMPL	runtime.writeBarrier(SB), $0
	0x01c5 00453 ($GOROOT/src/reflect/value.go:2709)	JNE	620
	0x01cb 00459 ($GOROOT/src/reflect/value.go:2709)	MOVQ	CX, reflect.dummy+16(SB)
	0x01d2 00466 ($GOROOT/src/reflect/value.go:2263)	PCDATA	$2, $2
	0x01d2 00466 ($GOROOT/src/reflect/value.go:2263)	PCDATA	$0, $5
	0x01d2 00466 ($GOROOT/src/reflect/value.go:2263)	MOVQ	AX, reflect.i+120(SP)
	0x01d7 00471 ($GOROOT/src/reflect/value.go:2263)	MOVQ	CX, reflect.i+128(SP)
	0x01df 00479 ($GOROOT/src/reflect/value.go:2263)	XCHGL	AX, AX
	0x01e0 00480 ($GOROOT/src/reflect/value.go:143)	PCDATA	$2, $6
	0x01e0 00480 ($GOROOT/src/reflect/value.go:143)	PCDATA	$0, $0
	0x01e0 00480 ($GOROOT/src/reflect/value.go:143)	MOVQ	reflect.i+120(SP), AX
	0x01e5 00485 ($GOROOT/src/reflect/value.go:144)	TESTQ	AX, AX
	0x01e8 00488 (:0)	JEQ	612
	0x01ea 00490 ($GOROOT/src/reflect/type.go:783)	MOVBLZX	23(AX), DX
	0x01ee 00494 ($GOROOT/src/reflect/type.go:783)	MOVL	DX, BX
	0x01f0 00496 ($GOROOT/src/reflect/type.go:783)	ANDL	$31, DX
	0x01f3 00499 ($GOROOT/src/reflect/value.go:149)	MOVQ	DX, SI
	0x01f6 00502 ($GOROOT/src/reflect/value.go:149)	BTSQ	$7, DX
	0x01fb 00507 ($GOROOT/src/reflect/type.go:3116)	TESTB	$32, BL
	0x01fe 00510 ($GOROOT/src/reflect/value.go:151)	CMOVQEQ	DX, SI
	0x0202 00514 ($GOROOT/src/reflect/value.go:147)	XCHGL	AX, AX
	0x0203 00515 ($GOROOT/src/reflect/value.go:148)	XCHGL	AX, AX
	0x0204 00516 (demo3.go:26)	PCDATA	$2, $2
	0x0204 00516 (demo3.go:26)	MOVQ	AX, (SP)
	0x0208 00520 (demo3.go:26)	PCDATA	$2, $0
	0x0208 00520 (demo3.go:26)	MOVQ	CX, 8(SP)
	0x020d 00525 (demo3.go:26)	MOVQ	SI, 16(SP)
	0x0212 00530 (demo3.go:26)	MOVQ	$0, 24(SP)
	0x021b 00539 (demo3.go:26)	CALL	reflect.Value.Method(SB)   // Method方法调用
	0x0220 00544 (demo3.go:26)	MOVQ	48(SP), AX                 // Method方法调用返回值Value
	0x0225 00549 (demo3.go:26)	PCDATA	$2, $2
	0x0225 00549 (demo3.go:26)	MOVQ	40(SP), CX
	0x022a 00554 (demo3.go:26)	PCDATA	$2, $7
	0x022a 00554 (demo3.go:26)	MOVQ	32(SP), DX
	0x022f 00559 (demo3.go:28)	PCDATA	$2, $2
	0x022f 00559 (demo3.go:28)	MOVQ	DX, (SP)
	0x0233 00563 (demo3.go:28)	PCDATA	$2, $0
	0x0233 00563 (demo3.go:28)	MOVQ	CX, 8(SP)
	0x0238 00568 (demo3.go:28)	MOVQ	AX, 16(SP)
	0x023d 00573 (demo3.go:28)	PCDATA	$2, $1
	0x023d 00573 (demo3.go:28)	LEAQ	""..autotmp_47+96(SP), AX
	0x0242 00578 (demo3.go:28)	PCDATA	$2, $0
	0x0242 00578 (demo3.go:28)	MOVQ	AX, 24(SP) // Call方法调用参数
	0x0247 00583 (demo3.go:28)	XORPS	X0, X0
	0x024a 00586 (demo3.go:28)	MOVUPS	X0, 32(SP)
	0x024f 00591 (demo3.go:28)	CALL	reflect.Value.Call(SB) // 进行反射调用
	0x0254 00596 (demo3.go:29)	MOVQ	256(SP), BP
	0x025c 00604 (demo3.go:29)	ADDQ	$264, SP
	0x0263 00611 (demo3.go:29)	RET
	0x0264 00612 (demo3.go:29)	XORL	SI, SI
	0x0266 00614 (demo3.go:29)	PCDATA	$2, $2
	0x0266 00614 (demo3.go:29)	XORL	CX, CX
	0x0268 00616 (demo3.go:29)	PCDATA	$2, $6
	0x0268 00616 (demo3.go:29)	XORL	AX, AX
	0x026a 00618 ($GOROOT/src/reflect/value.go:2263)	JMP	516
	0x026c 00620 ($GOROOT/src/reflect/value.go:2709)	PCDATA	$2, $-2
	0x026c 00620 ($GOROOT/src/reflect/value.go:2709)	PCDATA	$0, $-2
	0x026c 00620 ($GOROOT/src/reflect/value.go:2709)	LEAQ	reflect.dummy+16(SB), DI
	0x0273 00627 (demo3.go:25)	MOVQ	AX, DX
	0x0276 00630 ($GOROOT/src/reflect/value.go:2709)	MOVQ	CX, AX
	0x0279 00633 ($GOROOT/src/reflect/value.go:2709)	CALL	runtime.gcWriteBarrier(SB)
	0x027e 00638 ($GOROOT/src/reflect/value.go:2263)	MOVQ	DX, AX
	0x0281 00641 ($GOROOT/src/reflect/value.go:2709)	JMP	466
	0x0286 00646 ($GOROOT/src/reflect/value.go:2709)	PCDATA	$2, $1
	0x0286 00646 ($GOROOT/src/reflect/value.go:2709)	PCDATA	$0, $0
	0x0286 00646 ($GOROOT/src/reflect/value.go:2709)	XORL	AX, AX
	0x0288 00648 ($GOROOT/src/reflect/value.go:2709)	XORL	SI, SI
	0x028a 00650 ($GOROOT/src/reflect/value.go:2709)	PCDATA	$2, $6
	0x028a 00650 ($GOROOT/src/reflect/value.go:2709)	XORL	CX, CX
	0x028c 00652 (demo3.go:25)	JMP	516
	0x0291 00657 (demo3.go:25)	PCDATA	$2, $1
	0x0291 00657 (demo3.go:25)	XORL	AX, AX
	0x0293 00659 (demo3.go:25)	XORL	CX, CX
	0x0295 00661 ($GOROOT/src/reflect/type.go:1376)	JMP	161
	0x029a 00666 ($GOROOT/src/reflect/type.go:1376)	NOP
	0x029a 00666 (demo3.go:17)	PCDATA	$0, $-1
	0x029a 00666 (demo3.go:17)	PCDATA	$2, $-1
	0x029a 00666 (demo3.go:17)	CALL	runtime.morestack_noctxt(SB)
	0x029f 00671 (demo3.go:17)	JMP	0
	0x0000 65 48 8b 0c 25 00 00 00 00 48 8d 84 24 78 ff ff  eH..%....H..$x..
	0x0010 ff 48 3b 41 10 0f 86 7f 02 00 00 48 81 ec 08 01  .H;A.......H....
	0x0020 00 00 48 89 ac 24 00 01 00 00 48 8d ac 24 00 01  ..H..$....H..$..
	0x0030 00 00 48 8d 05 00 00 00 00 48 89 84 24 98 00 00  ..H......H..$...
	0x0040 00 48 c7 84 24 a0 00 00 00 06 00 00 00 48 c7 84  .H..$........H..
	0x0050 24 a8 00 00 00 0d 00 00 00 48 8d 0d 00 00 00 00  $........H......
	0x0060 48 89 0c 24 48 8d 94 24 98 00 00 00 48 89 54 24  H..$H..$....H.T$
	0x0070 08 e8 00 00 00 00 48 8b 44 24 10 48 8b 4c 24 18  ......H.D$.H.L$.
	0x0080 48 89 44 24 68 48 89 4c 24 70 90 48 8b 44 24 68  H.D$hH.L$p.H.D$h
	0x0090 90 48 85 c0 0f 84 f7 01 00 00 48 8d 0d 00 00 00  .H........H.....
	0x00a0 00 48 8b 89 a8 00 00 00 48 89 04 24 48 c7 44 24  .H......H..$H.D$
	0x00b0 08 00 00 00 00 ff d1 48 8d bc 24 b0 00 00 00 48  .......H..$....H
	0x00c0 8d 74 24 10 48 89 6c 24 f0 48 8d 6c 24 f0 e8 00  .t$.H.l$.H.l$...
	0x00d0 00 00 00 48 8b 6d 00 48 8d 05 00 00 00 00 48 89  ...H.m.H......H.
	0x00e0 04 24 48 8d 84 24 b0 00 00 00 48 89 44 24 08 e8  .$H..$....H.D$..
	0x00f0 00 00 00 00 48 8b 44 24 10 48 8b 4c 24 18 0f 57  ....H.D$.H.L$..W
	0x0100 c0 0f 11 84 24 88 00 00 00 48 89 84 24 88 00 00  ....$....H..$...
	0x0110 00 48 89 8c 24 90 00 00 00 90 48 8b 05 00 00 00  .H..$.....H.....
	0x0120 00 48 8d 0d 00 00 00 00 48 89 0c 24 48 89 44 24  .H......H..$H.D$
	0x0130 08 48 8d 84 24 88 00 00 00 48 89 44 24 10 48 c7  .H..$....H.D$.H.
	0x0140 44 24 18 01 00 00 00 48 c7 44 24 20 01 00 00 00  D$.....H.D$ ....
	0x0150 e8 00 00 00 00 48 8d 05 00 00 00 00 48 89 84 24  .....H......H..$
	0x0160 98 00 00 00 48 c7 84 24 a0 00 00 00 06 00 00 00  ....H..$........
	0x0170 48 c7 84 24 a8 00 00 00 0d 00 00 00 48 8d 05 00  H..$........H...
	0x0180 00 00 00 48 89 04 24 48 8d 84 24 98 00 00 00 48  ...H..$H..$....H
	0x0190 89 44 24 08 e8 00 00 00 00 48 8b 44 24 10 48 8b  .D$......H.D$.H.
	0x01a0 4c 24 18 90 48 85 c0 0f 84 d9 00 00 00 90 80 3d  L$..H..........=
	0x01b0 00 00 00 00 00 74 1b 48 89 05 00 00 00 00 83 3d  .....t.H.......=
	0x01c0 00 00 00 00 00 0f 85 a1 00 00 00 48 89 0d 00 00  ...........H....
	0x01d0 00 00 48 89 44 24 78 48 89 8c 24 80 00 00 00 90  ..H.D$xH..$.....
	0x01e0 48 8b 44 24 78 48 85 c0 74 7a 0f b6 50 17 89 d3  H.D$xH..tz..P...
	0x01f0 83 e2 1f 48 89 d6 48 0f ba ea 07 f6 c3 20 48 0f  ...H..H...... H.
	0x0200 44 f2 90 90 48 89 04 24 48 89 4c 24 08 48 89 74  D...H..$H.L$.H.t
	0x0210 24 10 48 c7 44 24 18 00 00 00 00 e8 00 00 00 00  $.H.D$..........
	0x0220 48 8b 44 24 30 48 8b 4c 24 28 48 8b 54 24 20 48  H.D$0H.L$(H.T$ H
	0x0230 89 14 24 48 89 4c 24 08 48 89 44 24 10 48 8d 44  ..$H.L$.H.D$.H.D
	0x0240 24 60 48 89 44 24 18 0f 57 c0 0f 11 44 24 20 e8  $`H.D$..W...D$ .
	0x0250 00 00 00 00 48 8b ac 24 00 01 00 00 48 81 c4 08  ....H..$....H...
	0x0260 01 00 00 c3 31 f6 31 c9 31 c0 eb 98 48 8d 3d 00  ....1.1.1...H.=.
	0x0270 00 00 00 48 89 c2 48 89 c8 e8 00 00 00 00 48 89  ...H..H.......H.
	0x0280 d0 e9 4c ff ff ff 31 c0 31 f6 31 c9 e9 73 ff ff  ..L...1.1.1..s..
	0x0290 ff 31 c0 31 c9 e9 07 fe ff ff e8 00 00 00 00 e9  .1.1............
	0x02a0 5c fd ff ff                                      \...
	rel 5+4 t=16 TLS+0
	rel 53+4 t=15 go.string."myname"+0
	rel 92+4 t=15 type."".User+0
	rel 114+4 t=8 runtime.convT2E+0
	rel 157+4 t=15 go.itab.*reflect.rtype,reflect.Type+0
	rel 181+0 t=11 +0
	rel 207+4 t=8 runtime.duffcopy+826
	rel 218+4 t=15 type.reflect.Method+0
	rel 240+4 t=8 runtime.convT2E+0
	rel 285+4 t=15 os.Stdout+0
	rel 292+4 t=15 go.itab.*os.File,io.Writer+0
	rel 337+4 t=8 fmt.Fprintln+0
	rel 344+4 t=15 go.string."myname"+0
	rel 383+4 t=15 type."".User+0
	rel 405+4 t=8 runtime.convT2E+0
	rel 432+4 t=15 reflect.dummy+-1
	rel 442+4 t=15 reflect.dummy+8
	rel 448+4 t=15 runtime.writeBarrier+-1
	rel 462+4 t=15 reflect.dummy+16
	rel 540+4 t=8 reflect.Value.Method+0
	rel 592+4 t=8 reflect.Value.Call+0
	rel 623+4 t=15 reflect.dummy+16
	rel 634+4 t=8 runtime.gcWriteBarrier+0
	rel 667+4 t=8 runtime.morestack_noctxt+0
```
看上面的汇编输出,

在进行Method方法调用时，如果你看不出来是否是真的Method方法调用，可以使用我们关于interface讲解中的关于
`readelf`和`objdump`的使用，读取其VMA以及数据大小

```
	0x009a 00154 (:0)	LEAQ	go.itab.*reflect.rtype,reflect.Type(SB), CX  // 设置go.itab.*reflect.rtype,reflect.Type
	0x00a1 00161 (demo3.go:24)	MOVQ	168(CX), CX                          // 168=160+8=0xa0+8
	0x00a8 00168 (demo3.go:24)	PCDATA	$2, $0
	0x00a8 00168 (demo3.go:24)	MOVQ	AX, (SP)                      // 参数i接口类型
	0x00ac 00172 (demo3.go:24)	MOVQ	$0, 8(SP)                     // Method方法调用参数：0
	0x00b5 00181 (demo3.go:24)	CALL	CX                           // 调用Type接口Method方法
```
```
➜  reflectdemo ./tool.sh iface.bin .rodata 'go.itab.*reflect.rtype,reflect.Type'
.rodata file-offset: 729088
.rodata VMA: 4923392
go.itab.*reflect.rtype,reflect.Type VMA: 5287296
go.itab.*reflect.rtype,reflect.Type SIZE: 272

0000000 40 24 4e 00 00 00 00 00 00 53 4e 00 00 00 00 00
0000010 d6 c9 33 e3 00 00 00 00 60 b3 46 00 00 00 00 00
0000020 40 ea 46 00 00 00 00 00 60 b2 46 00 00 00 00 00
0000030 e0 c0 46 00 00 00 00 00 d0 eb 46 00 00 00 00 00
0000040 20 eb 46 00 00 00 00 00 c0 c1 46 00 00 00 00 00
0000050 20 c3 46 00 00 00 00 00 70 b3 46 00 00 00 00 00
0000060 00 c4 46 00 00 00 00 00 10 c5 46 00 00 00 00 00
0000070 60 c6 46 00 00 00 00 00 50 e9 46 00 00 00 00 00
0000080 a0 c7 46 00 00 00 00 00 50 c1 46 00 00 00 00 00
0000090 80 c8 46 00 00 00 00 00 80 b3 46 00 00 00 00 00
00000a0 00 c9 46 00 00 00 00 00 d0 b4 46 00 00 00 00 00
#                               -----------------------
#                               0xa0+8 reflect.(*rtype).Method
00000b0 60 bc 46 00 00 00 00 00 30 c0 46 00 00 00 00 00
00000c0 70 c9 46 00 00 00 00 00 e0 c9 46 00 00 00 00 00
00000d0 60 b4 46 00 00 00 00 00 50 ca 46 00 00 00 00 00
00000e0 00 cb 46 00 00 00 00 00 60 bf 46 00 00 00 00 00
00000f0 50 b2 46 00 00 00 00 00 80 b1 46 00 00 00 00 00
0000100 a0 b3 46 00 00 00 00 00 c0 b0 46 00 00 00 00 00
0000110
➜  reflectdemo objdump -t -j .text iface.bin | grep 000000000046b0c0
000000000046b0c0 g     F .text	00000000000000b4 reflect.(*rtype).uncommon
➜  reflectdemo objdump -t -j .text iface.bin | grep 000000000046b3a0
000000000046b3a0 g     F .text	000000000000000b reflect.(*rtype).common
➜  reflectdemo objdump -t -j .text iface.bin | grep 000000000046b180
000000000046b180 g     F .text	00000000000000cb reflect.(*rtype).String
➜  reflectdemo objdump -t -j .text iface.bin | grep 000000000046b250
000000000046b250 g     F .text	000000000000000e reflect.(*rtype).Size
➜  reflectdemo objdump -t -j .text iface.bin | grep 000000000046bf60
000000000046bf60 g     F .text	00000000000000d0 reflect.(*rtype).PkgPath
➜  reflectdemo objdump -t -j .text iface.bin | grep 000000000046cb00
000000000046cb00 g     F .text	000000000000010a reflect.(*rtype).Out
➜  reflectdemo objdump -t -j .text iface.bin | grep 000000000046ca50
000000000046ca50 g     F .text	00000000000000a2 reflect.(*rtype).NumOut
➜  reflectdemo objdump -t -j .text iface.bin | grep 000000000046b460
000000000046b460 g     F .text	0000000000000068 reflect.(*rtype).NumMethod
➜  reflectdemo objdump -t -j .text iface.bin | grep 000000000046b4d0
000000000046b4d0 g     F .text	0000000000000788 reflect.(*rtype).Method
➜  reflectdemo
```

我们通过的上面的注释可以发现，在调用`TypeOf`的时候，有一个参数的类型转换过程,看了该函数的原型，大概就明白了为啥有这个过程：
```
// TypeOf返回表示i的动态类型的反射类型。
// 如果i是一个nil接口值，TypeOf返回nil.
func TypeOf(i interface{}) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}
```

将结构体进行`interface{}`转换的操作是由`runtime.convT2E`完成的。
```
// The conv and assert functions below do very similar things.
// The convXXX functions are guaranteed by the compiler to succeed.
// The assertXXX functions may fail (either panicking or returning false,
// depending on whether they are 1-result or 2-result).
// The convXXX functions succeed on a nil input, whereas the assertXXX
// functions fail on a nil input.

func convT2E(t *_type, elem unsafe.Pointer) (e eface) {
	if raceenabled {
		raceReadObjectPC(t, elem, getcallerpc(), funcPC(convT2E))
	}
	if msanenabled {
		msanread(elem, t.size)
	}
	x := mallocgc(t.size, t, true)
	// TODO: We allocate a zeroed object only to overwrite it with actual data.
	// Figure out how to avoid zeroing. Also below in convT2Eslice, convT2I, convT2Islice.
	typedmemmove(t, x, elem)
	e._type = t
	e.data = x
	return
}
```


在将参数接口转换后，赋值给了i变量，之后在`TypeOf`函数中做了什么，一看便知：
```
// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}
```

```
// toType converts from a *rtype to a Type that can be returned
// to the client of package reflect. In gc, the only concern is that
// a nil *rtype must be replaced by a nil Type, but in gccgo this
// function takes care of ensuring that multiple *rtype for the same
// type are coalesced into a single Type.
func toType(t *rtype) Type {
	if t == nil {
		return nil
	}
	return t
}
```

`rtype`是何物，是时候出场了
```
// rtype is the common implementation of most values.
// It is embedded in other struct types.
//
// rtype must be kept in sync with ../runtime/type.go:/^type._type.
type rtype struct {
	size       uintptr
	ptrdata    uintptr  // number of bytes in the type that can contain pointers
	hash       uint32   // hash of type; avoids computation in hash tables
	tflag      tflag    // extra type information flags
	align      uint8    // alignment of variable with this type
	fieldAlign uint8    // alignment of struct field with this type
	kind       uint8    // enumeration for C
	alg        *typeAlg // algorithm table
	gcdata     *byte    // garbage collection data
	str        nameOff  // string form
	ptrToThis  typeOff  // type for pointer to this type, may be zero
}

```
看过`interface`文章后，你大概知道了，他是干啥用的，这也是为什么`emptyInterface`和`eface`能够通过指针操作进行相互转换的原因操作。

继续看`fmt.Println(t.Method(0))`,`t`是`Type`类型，该`Type`是一个接口类型,完整形式如下：
```
// Type is the representation of a Go type.
//
// Not all methods apply to all kinds of types. Restrictions,
// if any, are noted in the documentation for each method.
// Use the Kind method to find out the kind of type before
// calling kind-specific methods. Calling a method
// inappropriate to the kind of type causes a run-time panic.
//
// Type values are comparable, such as with the == operator,
// so they can be used as map keys.
// Two Type values are equal if they represent identical types.
type Type interface {
	// Methods applicable to all types.

	// Align returns the alignment in bytes of a value of
	// this type when allocated in memory.
	Align() int

	// FieldAlign returns the alignment in bytes of a value of
	// this type when used as a field in a struct.
	FieldAlign() int

	// Method returns the i'th method in the type's method set.
	// It panics if i is not in the range [0, NumMethod()).
	//
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

	// NumMethod returns the number of exported methods in the type's method set.
	NumMethod() int

	// Name returns the type's name within its package for a defined type.
	// For other (non-defined) types it returns the empty string.
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

	// Kind returns the specific kind of this type.
	Kind() Kind

	// Implements reports whether the type implements the interface type u.
	Implements(u Type) bool

	// AssignableTo reports whether a value of the type is assignable to type u.
	AssignableTo(u Type) bool

	// ConvertibleTo reports whether a value of the type is convertible to type u.
	ConvertibleTo(u Type) bool

	// Comparable reports whether values of this type are comparable.
	Comparable() bool

	// Methods applicable only to some types, depending on Kind.
	// The methods allowed for each kind are:
	//
	//	Int*, Uint*, Float*, Complex*: Bits
	//	Array: Elem, Len
	//	Chan: ChanDir, Elem
	//	Func: In, NumIn, Out, NumOut, IsVariadic.
	//	Map: Key, Elem
	//	Ptr: Elem
	//	Slice: Elem
	//	Struct: Field, FieldByIndex, FieldByName, FieldByNameFunc, NumField

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
	//	t.NumIn() == 2
	//	t.In(0) is the reflect.Type for "int"
	//	t.In(1) is the reflect.Type for "[]float64"
	//	t.IsVariadic() == true
	//
	// IsVariadic panics if the type's Kind is not Func.
	IsVariadic() bool

	// Elem returns a type's element type.
	// It panics if the type's Kind is not Array, Chan, Map, Ptr, or Slice.
	Elem() Type

	// Field returns a struct type's i'th field.
	// It panics if the type's Kind is not Struct.
	// It panics if i is not in the range [0, NumField()).
	Field(i int) StructField

	// FieldByIndex returns the nested field corresponding
	// to the index sequence. It is equivalent to calling Field
	// successively for each index i.
	// It panics if the type's Kind is not Struct.
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

	// In returns the type of a function type's i'th input parameter.
	// It panics if the type's Kind is not Func.
	// It panics if i is not in the range [0, NumIn()).
	In(i int) Type

	// Key returns a map type's key type.
	// It panics if the type's Kind is not Map.
	Key() Type

	// Len returns an array type's length.
	// It panics if the type's Kind is not Array.
	Len() int

	// NumField returns a struct type's field count.
	// It panics if the type's Kind is not Struct.
	NumField() int

	// NumIn returns a function type's input parameter count.
	// It panics if the type's Kind is not Func.
	NumIn() int

	// NumOut returns a function type's output parameter count.
	// It panics if the type's Kind is not Func.
	NumOut() int

	// Out returns the type of a function type's i'th output parameter.
	// It panics if the type's Kind is not Func.
	// It panics if i is not in the range [0, NumOut()).
	Out(i int) Type

	common() *rtype
	uncommon() *uncommonType
}

```
`rtype`作为该接口的实现者，当我们在调用`t.Method(n int)`方法时，实际上执行的就是该类型上的方法
```
func (t *rtype) Method(i int) (m Method) {
	if t.Kind() == Interface {
		tt := (*interfaceType)(unsafe.Pointer(t))
		return tt.Method(i)
	}
	methods := t.exportedMethods()
	if i < 0 || i >= len(methods) {
		panic("reflect: Method index out of range")
	}
	p := methods[i]
	pname := t.nameOff(p.name)
	m.Name = pname.name()
	fl := flag(Func)
	mtyp := t.typeOff(p.mtyp)
	ft := (*funcType)(unsafe.Pointer(mtyp))
	in := make([]Type, 0, 1+len(ft.in()))
	in = append(in, t)
	for _, arg := range ft.in() {
		in = append(in, arg)
	}
	out := make([]Type, 0, len(ft.out()))
	for _, ret := range ft.out() {
		out = append(out, ret)
	}
	mt := FuncOf(in, out, ft.IsVariadic())
	m.Type = mt
	tfn := t.textOff(p.tfn)
	fn := unsafe.Pointer(&tfn)
	m.Func = Value{mt.(*rtype), fn, fl}

	m.Index = i
	return m
}
```
该方法返回值是一个结构体类型
```
/*
 * The compiler knows the exact layout of all the data structures above.
 * The compiler does not know about the data structures and methods below.
 */

// Method represents a single method.
type Method struct {
	// Name is the method name.
	// PkgPath is the package path that qualifies a lower case (unexported)
	// method name. It is empty for upper case (exported) method names.
	// The combination of PkgPath and Name uniquely identifies a method
	// in a method set.
	// See https://golang.org/ref/spec#Uniqueness_of_identifiers
	Name    string
	PkgPath string

	Type  Type  // 方法类型
	Func  Value // 接受者作为第一个参数的函数
	Index int   // 方法的索引
}

```
`mt := FuncOf(in, out, ft.IsVariadic())`这个函数的职能是返回一个函数类型带有给定的参数和返回值类型。
例如如果k代表int,e代表字符串，那么调用`FuncOf([]Type{k},[]Type{e},false)`代表函数`func(int) string`.
`variadic`控制该函数是否支持可变参数。如果`variadic`为`true`并且`in[len(in)-1]`不是一个切片,`FuncOf`将`panic`.
```
func FuncOf(in, out []Type, variadic bool) Type {
	if variadic && (len(in) == 0 || in[len(in)-1].Kind() != Slice) {
		panic("reflect.FuncOf: last arg of variadic func must be slice")
	}

	// Make a func type.
	var ifunc interface{} = (func())(nil)
	prototype := *(**funcType)(unsafe.Pointer(&ifunc))
	n := len(in) + len(out)

	var ft *funcType
	var args []*rtype
	switch {
	case n <= 4:
		fixed := new(funcTypeFixed4)
		args = fixed.args[:0:len(fixed.args)]
		ft = &fixed.funcType
	case n <= 8:
		fixed := new(funcTypeFixed8)
		args = fixed.args[:0:len(fixed.args)]
		ft = &fixed.funcType
	case n <= 16:
		fixed := new(funcTypeFixed16)
		args = fixed.args[:0:len(fixed.args)]
		ft = &fixed.funcType
	case n <= 32:
		fixed := new(funcTypeFixed32)
		args = fixed.args[:0:len(fixed.args)]
		ft = &fixed.funcType
	case n <= 64:
		fixed := new(funcTypeFixed64)
		args = fixed.args[:0:len(fixed.args)]
		ft = &fixed.funcType
	case n <= 128:
		fixed := new(funcTypeFixed128)
		args = fixed.args[:0:len(fixed.args)]
		ft = &fixed.funcType
	default:
		panic("reflect.FuncOf: too many arguments")
	}
	*ft = *prototype

	// Build a hash and minimally populate ft.
	var hash uint32
	for _, in := range in {
		t := in.(*rtype)
		args = append(args, t)
		hash = fnv1(hash, byte(t.hash>>24), byte(t.hash>>16), byte(t.hash>>8), byte(t.hash))
	}
	if variadic {
		hash = fnv1(hash, 'v')
	}
	hash = fnv1(hash, '.')
	for _, out := range out {
		t := out.(*rtype)
		args = append(args, t)
		hash = fnv1(hash, byte(t.hash>>24), byte(t.hash>>16), byte(t.hash>>8), byte(t.hash))
	}
	if len(args) > 50 {
		panic("reflect.FuncOf does not support more than 50 arguments")
	}
	ft.tflag = 0
	ft.hash = hash
	ft.inCount = uint16(len(in))
	ft.outCount = uint16(len(out))
	if variadic {
		ft.outCount |= 1 << 15
	}

	// Look in cache.
	if ts, ok := funcLookupCache.m.Load(hash); ok {
		for _, t := range ts.([]*rtype) {
			if haveIdenticalUnderlyingType(&ft.rtype, t, true) {
				return t
			}
		}
	}

	// Not in cache, lock and retry.
	funcLookupCache.Lock()
	defer funcLookupCache.Unlock()
	if ts, ok := funcLookupCache.m.Load(hash); ok {
		for _, t := range ts.([]*rtype) {
			if haveIdenticalUnderlyingType(&ft.rtype, t, true) {
				return t
			}
		}
	}

	addToCache := func(tt *rtype) Type {
		var rts []*rtype
		if rti, ok := funcLookupCache.m.Load(hash); ok {
			rts = rti.([]*rtype)
		}
		funcLookupCache.m.Store(hash, append(rts, tt))
		return tt
	}

	// Look in known types for the same string representation.
	str := funcStr(ft)
	for _, tt := range typesByString(str) {
		if haveIdenticalUnderlyingType(&ft.rtype, tt, true) {
			return addToCache(tt)
		}
	}

	// Populate the remaining fields of ft and store in cache.
	ft.str = resolveReflectName(newName(str, "", false))
	ft.ptrToThis = 0
	return addToCache(&ft.rtype)
}
```
本示例的返回值可以看出一下端倪：
```
{Descrite  func(main.User) <func(main.User) Value> 0}
```
上述是`fmt.Println(t.Method(0))`的输出，
```
reflect.Method{Name:"Descrite", PkgPath:"", Type:(*reflect.rtype)(0xc00008e0c0), Func:reflect.Value{typ:(*reflect.rtype)(0xc00008e0c0), ptr:(unsafe.Pointer)(0xc00008c028), flag:0x13}, Index:0}
```
上述是`fmt.Printf("%#v", t.Method(0))`的输出内容，可以看出，该内容即可`Method`结构体的内容，而该结构体的`Func`字段即是`<func(main.User) Value>`

到这里有关`Type`类型的`Method`方法调用就结束了，该调用可以查看具体的动态变量的方法类型。

继续下面的内容
```
    v := reflect.ValueOf(user)
    mName:=v.Method(0)
```

```
	0x0203 00515 (demo3.go:27)	MOVQ	"".user+240(SP), AX
	0x020b 00523 (demo3.go:27)	PCDATA	$2, $2
	0x020b 00523 (demo3.go:27)	MOVQ	"".user+224(SP), CX
	0x0213 00531 (demo3.go:27)	PCDATA	$0, $0
	0x0213 00531 (demo3.go:27)	MOVQ	"".user+232(SP), DX
	0x021b 00539 (demo3.go:27)	PCDATA	$2, $0
	0x021b 00539 (demo3.go:27)	PCDATA	$0, $9
	0x021b 00539 (demo3.go:27)	MOVQ	CX, ""..autotmp_5+320(SP)
	0x0223 00547 (demo3.go:27)	MOVQ	DX, ""..autotmp_5+328(SP)
	0x022b 00555 (demo3.go:27)	MOVQ	AX, ""..autotmp_5+336(SP)
	0x0233 00563 (demo3.go:27)	PCDATA	$2, $1
	0x0233 00563 (demo3.go:27)	LEAQ	type."".User(SB), AX
	0x023a 00570 (demo3.go:27)	PCDATA	$2, $0
	0x023a 00570 (demo3.go:27)	MOVQ	AX, (SP)
	0x023e 00574 (demo3.go:27)	PCDATA	$2, $1
	0x023e 00574 (demo3.go:27)	PCDATA	$0, $0
	0x023e 00574 (demo3.go:27)	LEAQ	""..autotmp_5+320(SP), AX
	0x0246 00582 (demo3.go:27)	PCDATA	$2, $0
	0x0246 00582 (demo3.go:27)	MOVQ	AX, 8(SP)
	0x024b 00587 (demo3.go:27)	CALL	runtime.convT2E(SB)
	0x0250 00592 (demo3.go:27)	MOVQ	16(SP), AX
	0x0255 00597 (demo3.go:27)	PCDATA	$2, $2
	0x0255 00597 (demo3.go:27)	MOVQ	24(SP), CX
	0x025a 00602 (demo3.go:27)	MOVQ	AX, ""..autotmp_12+168(SP)
	0x0262 00610 (demo3.go:27)	MOVQ	CX, ""..autotmp_12+176(SP)
	0x026a 00618 (demo3.go:27)	MOVQ	AX, (SP)
	0x026e 00622 (demo3.go:27)	PCDATA	$2, $0
	0x026e 00622 (demo3.go:27)	MOVQ	CX, 8(SP)
	0x0273 00627 (demo3.go:27)	CALL	reflect.ValueOf(SB)
	0x0278 00632 (demo3.go:27)	PCDATA	$2, $1
	0x0278 00632 (demo3.go:27)	MOVQ	16(SP), AX
	0x027d 00637 (demo3.go:27)	PCDATA	$2, $7
	0x027d 00637 (demo3.go:27)	MOVQ	24(SP), CX
	0x0282 00642 (demo3.go:27)	MOVQ	32(SP), DX
	0x0287 00647 (demo3.go:27)	PCDATA	$2, $2
	0x0287 00647 (demo3.go:27)	PCDATA	$0, $10
	0x0287 00647 (demo3.go:27)	MOVQ	AX, "".v+200(SP)
	0x028f 00655 (demo3.go:27)	PCDATA	$2, $0
	0x028f 00655 (demo3.go:27)	MOVQ	CX, "".v+208(SP)
	0x0297 00663 (demo3.go:27)	MOVQ	DX, "".v+216(SP)
	0x029f 00671 (demo3.go:28)	PCDATA	$2, $1
	0x029f 00671 (demo3.go:28)	MOVQ	"".v+200(SP), AX
	0x02a7 00679 (demo3.go:28)	PCDATA	$2, $7
	0x02a7 00679 (demo3.go:28)	MOVQ	"".v+208(SP), CX
	0x02af 00687 (demo3.go:28)	PCDATA	$0, $0
	0x02af 00687 (demo3.go:28)	MOVQ	"".v+216(SP), DX
	0x02b7 00695 (demo3.go:28)	PCDATA	$2, $2
	0x02b7 00695 (demo3.go:28)	MOVQ	AX, (SP)
	0x02bb 00699 (demo3.go:28)	PCDATA	$2, $0
	0x02bb 00699 (demo3.go:28)	MOVQ	CX, 8(SP)
	0x02c0 00704 (demo3.go:28)	MOVQ	DX, 16(SP)
	0x02c5 00709 (demo3.go:28)	MOVQ	$0, 24(SP)
	0x02ce 00718 (demo3.go:28)	CALL	reflect.Value.Method(SB)
	0x02d3 00723 (demo3.go:28)	PCDATA	$2, $1
	0x02d3 00723 (demo3.go:28)	MOVQ	32(SP), AX
	0x02d8 00728 (demo3.go:28)	PCDATA	$2, $7
	0x02d8 00728 (demo3.go:28)	MOVQ	40(SP), CX
	0x02dd 00733 (demo3.go:28)	MOVQ	48(SP), DX
```
第一部分都一样，都是将结构体接口化。
`ValueOf`返回一个新的`Value`使用存贮在i接口中的数据初始化，`ValueOf(nil)`返回零值
```
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
```
这里的`Value`是一个结构体,和上面的`TypeOf`返回的`Type`是接口类型不同。`Value`原型为：
```
// Value是Go value的反射接口。
//
// 并非所有方法都适用于所有类型的值。 在每种方法的文档中都注明了限制（如果有）。
// 请在调用特定于类型的方法之前，使用Kind方法查找值的类型。 
// 调用不适合该类型的方法将导致运行时恐慌
//
// 零值 Value 代表没有值.
// Its IsValid method returns false, its Kind method returns Invalid,
// its String method returns "<invalid Value>", and all other methods panic.
// 大多数函数和方法都不会返回一个非法的值。
// 如果有，则其文档中明确说明条件。
//
// 一个值可以由多个goroutine并发使用，前提是可以将基础Go值同时用于等效的直接操作。
//
// 要比较两个值，请比较Interface方法的返回值。 在两个Value上使用==不会比较它们表示的底层数据值。
type Value struct {
	// typ包含由Value表示的值的类型。
	typ *rtype

	// 指针值的数据，或者，如果设置了flagIndir，则指向数据的指针
	// 当设置了flagIndir或typ.pointers（）为true时有效
	ptr unsafe.Pointer

	// flag保存有关该值的元数据 值
	//  最低位是标志位:
	//	- flagStickyRO: 通过未导出的未嵌入字段获取，因此为只读
	//	- flagEmbedRO: 通过未导出的嵌入字段获取，因此为只读
	//	- flagIndir: val保存指向数据的指针
	//	- flagAddr: v.CanAddr 为 true (implies flagIndir)
	//	- flagMethod: v是一个方法值
	// 后五位给出值的种类。
	// This repeats typ.Kind() except for method values.
	// 
	// 其余的23+位给出方法值的方法号
	// 如果flag.kind() != Func,代码可以假定flagMethod未设置
	// 如果ifaceIndir(typ),代码可以假定flagIndir设置了
	flag

	// 方法值表示对于接受者r就像r.Read这样的调用。typ+val+flag位描述了接受者
	// 但标志的Kind位表示Func。方法是函数），并且标志的高位给出r的类型的方法表中的方法编号。
}

type flag uintptr

const (
	flagKindWidth        = 5 // there are 27 kinds
	flagKindMask    flag = 1<<flagKindWidth - 1
	flagStickyRO    flag = 1 << 5
	flagEmbedRO     flag = 1 << 6
	flagIndir       flag = 1 << 7
	flagAddr        flag = 1 << 8
	flagMethod      flag = 1 << 9
	flagMethodShift      = 10
	flagRO          flag = flagStickyRO | flagEmbedRO
)
```

看一下这个流程`unpackEface`,将`emptyInterface`i转化为一个`Value`

```
func unpackEface(i interface{}) Value {
	e := (*emptyInterface)(unsafe.Pointer(&i))
	// 注意：在我们明确是否是一个指针之前，不要读取e.word
	t := e.typ
	if t == nil {
		return Value{}
	}
	f := flag(t.Kind())
	// ifaceIndir报告t是否间接存储在接口值中。
	if ifaceIndir(t) {
		f |= flagIndir
	}
	return Value{t, e.word, f}
}

```
之后再该`Value`类型上调用`Method`方法，参数为0
```
// 方法返回与v的第i个方法相对应的函数值。 
// 返回函数上的Call的参数不应包含接收器； 
// 返回的函数将始终使用v作为接收者。 
// 如果i超出范围或v为nil接口值，则方法将panic。
func (v Value) Method(i int) Value {
	if v.typ == nil {
		panic(&ValueError{"reflect.Value.Method", Invalid})
	}
	if v.flag&flagMethod != 0 || uint(i) >= uint(v.typ.NumMethod()) {
		panic("reflect: Method index out of range")
	}
	if v.typ.Kind() == Interface && v.IsNil() {
		panic("reflect: Method on nil interface value")
	}
	fl := v.flag & (flagStickyRO | flagIndir) // Clear flagEmbedRO
	fl |= flag(Func)
	fl |= flag(i)<<flagMethodShift | flagMethod
	return Value{v.typ, v.ptr, fl}
}
```

之后的内容就是进行的反射调用了`CALL	reflect.Value.Call(SB)`
```
// Call调用函数v带有参数in
// 例如，如果len(in) == 3,v.Call(in)代表Go中的调用v(in[0], in[1], in[2])
// 如果v的Kind方法不是Func，Call将panic
// 将以Values返回输出结果
// 与Go中一样，每个输入参数必须匹配函数的相应输入参数的类型
// 如果v是可变参数函数，Call自己创建可变参数切片，并将相关值copy进去
func (v Value) Call(in []Value) []Value {
	v.mustBe(Func)
	v.mustBeExported()
	return v.call("Call", in)
}
```
首先进行类型检测，然后进行是否可导出的函数检测，通过后，即进入函数调用
```
func (v Value) call(op string, in []Value) []Value {
	// Get function pointer, type.
	t := (*funcType)(unsafe.Pointer(v.typ))
	var (
		fn       unsafe.Pointer
		rcvr     Value
		rcvrtype *rtype
	)
	if v.flag&flagMethod != 0 {
		rcvr = v
		rcvrtype, t, fn = methodReceiver(op, v, int(v.flag)>>flagMethodShift)
	} else if v.flag&flagIndir != 0 {
		fn = *(*unsafe.Pointer)(v.ptr)
	} else {
		fn = v.ptr
	}

	if fn == nil {
		panic("reflect.Value.Call: call of nil function")
	}

	isSlice := op == "CallSlice"
	n := t.NumIn()
	if isSlice {
		if !t.IsVariadic() {
			panic("reflect: CallSlice of non-variadic function")
		}
		if len(in) < n {
			panic("reflect: CallSlice with too few input arguments")
		}
		if len(in) > n {
			panic("reflect: CallSlice with too many input arguments")
		}
	} else {
		if t.IsVariadic() {
			n--
		}
		if len(in) < n {
			panic("reflect: Call with too few input arguments")
		}
		if !t.IsVariadic() && len(in) > n {
			panic("reflect: Call with too many input arguments")
		}
	}
	for _, x := range in {
		if x.Kind() == Invalid {
			panic("reflect: " + op + " using zero Value argument")
		}
	}
	for i := 0; i < n; i++ {
		if xt, targ := in[i].Type(), t.In(i); !xt.AssignableTo(targ) {
			panic("reflect: " + op + " using " + xt.String() + " as type " + targ.String())
		}
	}
	if !isSlice && t.IsVariadic() {
		// prepare slice for remaining values
		m := len(in) - n
		slice := MakeSlice(t.In(n), m, m)
		elem := t.In(n).Elem()
		for i := 0; i < m; i++ {
			x := in[n+i]
			if xt := x.Type(); !xt.AssignableTo(elem) {
				panic("reflect: cannot use " + xt.String() + " as type " + elem.String() + " in " + op)
			}
			slice.Index(i).Set(x)
		}
		origIn := in
		in = make([]Value, n+1)
		copy(in[:n], origIn)
		in[n] = slice
	}

	nin := len(in)
	if nin != t.NumIn() {
		panic("reflect.Value.Call: wrong argument count")
	}
	nout := t.NumOut()

	// Compute frame type.
	frametype, _, retOffset, _, framePool := funcLayout(t, rcvrtype)

	// Allocate a chunk of memory for frame.
	var args unsafe.Pointer
	if nout == 0 {
		args = framePool.Get().(unsafe.Pointer)
	} else {
		// Can't use pool if the function has return values.
		// We will leak pointer to args in ret, so its lifetime is not scoped.
		args = unsafe_New(frametype)
	}
	off := uintptr(0)

	// Copy inputs into args.
	if rcvrtype != nil {
		storeRcvr(rcvr, args)
		off = ptrSize
	}
	for i, v := range in {
		v.mustBeExported()
		targ := t.In(i).(*rtype)
		a := uintptr(targ.align)
		off = (off + a - 1) &^ (a - 1)
		n := targ.size
		if n == 0 {
			// Not safe to compute args+off pointing at 0 bytes,
			// because that might point beyond the end of the frame,
			// but we still need to call assignTo to check assignability.
			v.assignTo("reflect.Value.Call", targ, nil)
			continue
		}
		addr := add(args, off, "n > 0")
		v = v.assignTo("reflect.Value.Call", targ, addr)
		if v.flag&flagIndir != 0 {
			typedmemmove(targ, addr, v.ptr)
		} else {
			*(*unsafe.Pointer)(addr) = v.ptr
		}
		off += n
	}

	// Call.
	call(frametype, fn, args, uint32(frametype.size), uint32(retOffset))

	// For testing; see TestCallMethodJump.
	if callGC {
		runtime.GC()
	}

	var ret []Value
	if nout == 0 {
		typedmemclr(frametype, args)
		framePool.Put(args)
	} else {
		// Zero the now unused input area of args,
		// because the Values returned by this function contain pointers to the args object,
		// and will thus keep the args object alive indefinitely.
		typedmemclrpartial(frametype, args, 0, retOffset)

		// Wrap Values around return values in args.
		ret = make([]Value, nout)
		off = retOffset
		for i := 0; i < nout; i++ {
			tv := t.Out(i)
			a := uintptr(tv.Align())
			off = (off + a - 1) &^ (a - 1)
			if tv.Size() != 0 {
				fl := flagIndir | flag(tv.Kind())
				ret[i] = Value{tv.common(), add(args, off, "tv.Size() != 0"), fl}
				// Note: this does introduce false sharing between results -
				// if any result is live, they are all live.
				// (And the space for the args is live as well, but as we've
				// cleared that space it isn't as big a deal.)
			} else {
				// For zero-sized return value, args+off may point to the next object.
				// In this case, return the zero value instead.
				ret[i] = Zero(tv)
			}
			off += tv.Size()
		}
	}

	return ret
}

```


