## Buffer结构体
一个Buffer是一个具有读写方法，字节大小可变的缓冲区，零值是一个准备好使用的空的buffer。
```
// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
// The zero value for Buffer is an empty buffer ready to use.
type Buffer struct {
	buf       []byte   // contents are the bytes buf[off : len(buf)]
	off       int      // read at &buf[off], write at &buf[len(buf)]
	bootstrap [64]byte // memory to hold first slice; helps small buffers avoid allocation.
	lastRead  readOp   // last read operation, so that Unread* can work correctly.

	// FIXME: it would be advisable to align Buffer to cachelines to avoid false
	// sharing.
}
```
## NewBuffer
创建并初始化一个新的Buffer使用buf作为其初始化内容，新的缓冲区拥有buf的所有权，并且调用方不应在此调用后使用buf。大多数情况下，new(Buffer) 足以初始化一个Buffer
`func NewBuffer(buf []byte) *Buffer { return &Buffer{buf: buf} }`

## Bytes()
Bytes返回一个长度为b.Len()切片
```
func (b *Buffer) Bytes() []byte { return b.buf[b.off:] }
```

```
package main

import (
	"fmt"
	"bytes"
)

func main(){
	bys := []byte("this is first a byte slice")
	fmt.Println(bys)
	bys2 := new(bytes.Buffer)
	fmt.Println(bys2.Bytes())
	bys3 := bytes.NewBuffer([]byte("This is second a byte slice"))
	fmt.Println(bys3.Bytes()) //[84 104 105 115 32 105 115 32 97 32 98 121 116 101 32 115 108 105 99 101]
}
```

##Write
Write 将p的内容追加到buffer中，按需增加buffer的容量，返回的n值为p的长度，err永远是Nil,如果buffer变的太大，Write将会panic
`func (b *Buffer) Write(p []byte) (n int, err error) {`
## WriteString
WriteString 将参数s的内容追加到buffer中，同Write
`func (b *Buffer) WriteString(s string) (n int, err error) {`

```
   bys3 := bytes.NewBuffer([]byte("This is second a byte slice"))
	fmt.Println(bys3.Bytes())
	bys3.WriteString("this is append content.")
	fmt.Println(bys3.Bytes())
```


## WriteByte
`func (b *Buffer) WriteByte(c byte) error {`
示例:
```
//参数是字节类型，字节类型别名uint8，值为0-255,不进行byte()强制转化也是可行的
    var br bytes.Buffer
    br.WriteByte(byte(75))
    fmt.Printf("%v\n", string(br.Bytes()))
```
Output:
```
K
```
## WriteRune
`func (b *Buffer) WriteRune(r rune) (n int, err error) {`
```
    var br bytes.Buffer
    br.WriteRune('z') //Rune的字面量就是一个整型值，int32的别名br.WriteRune('赵')写入汉字会输出汉字
    fmt.Printf("%v\n", string(br.Bytes()))
```
Output:
```
z
```
## WriteTo
WriteTo向w中写数据，直到buffer耗尽或者遇到错误。返回值n表示写入的bytes数量，写时遇到的错误也将返回。
`func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {`

```
    var br bytes.Buffer
    br.WriteString("这是一个例子")
    br.WriteTo(os.Stdout)
```

## ReadFrom
`func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {`
ReadFrom从r中读取数据，直到遇到EOF.并将数据追加到buffer.返回读取的bytes数量，读取中遇到的错误除了io.EOF也将返回

## Read
`func (b *Buffer) Read(p []byte) (n int, err error) {`
Read从buffer中读取下一个len(p) bytes，直到buffer耗尽。返回值n代表读取的bytes数量，如果buffer没有数据返回， err==io.EOF.否则err== nil.

上面介绍了bytes.Buffer的一些主要方法，可以看到，该结构体实现了Read方法，Write方法，所以在Copy方法调用中，可以将此类的数据作为参数，示例：
```
func TestCopy(t *testing.T) {
	rb := new(Buffer)
	wb := new(Buffer)
	rb.WriteString("hello, world.")
	Copy(wb, rb)
	if wb.String() != "hello, world." {
		t.Errorf("Copy did not work properly")
	}
}

```