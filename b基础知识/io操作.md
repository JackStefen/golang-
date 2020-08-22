包IO提供I/O原语的基本接口。它的主要任务是包装这些原语的现有实现
# 1.Reader
Reader是一个包裹基本读方法的接口,Read读取len(p)个bytes到p中，返回值n代表读取到的bytes数，值范围为大于等于0小于等于len(p).err为读取过程遇到的错误。即使read返回n<len（p），它也可以在调用期间使用p的所有部分作为临时空间。read通常返回可用的内容，而不是等待更多内容。
```
type Reader interface {
	Read(p []byte) (n int, err error)
}
```
# 2.Writer 
Writer是包裹了基本写方法的接口。
Write从p中写len(p)个字节到底层的数据流，返回值n代表成功写入的bytes数，值范围为`0<=n<=len(p)`,err代表写时遇到的错误。Write必须返回一个非空的错误，如果n<len(p),Write不能修改参数切片
```
type Writer interface {
	Write(p []byte) (n int, err error)
}
```
# 3.ReadWriter 

```
// ReadWriter is the interface that groups the basic Read and Write methods.
type ReadWriter interface {
	Reader
	Writer
}
```

# 4.ReaderFrom 
ReadFrom从r读取数据直到遇到EOF或者error.
```
type ReaderFrom interface {
	ReadFrom(r Reader) (n int64, err error)
}
```
# 5.WriterTo 
WriteTo向w中写数据，直到没有更多的数据可写或者遇到错误。
```
type WriterTo interface {
	WriteTo(w Writer) (n int64, err error)
}
```

# 6.ReaderAt 
ReadAt读取len(p)个字节到p从输入源的off偏移量开始。返回读取到的字节数和遇到的任何错误。如果ReadAt返回的n<len(p),同时将返回一个非空的错误，解释了为什么返回的字节数少，在这种层面上，ReadAt要严于Read.如果n=len(p), err == EOF或者err==nil.

```
type ReaderAt interface {
	ReadAt(p []byte, off int64) (n int, err error)
}
```
# 7.WriterAt 
WriterAt将len(p)个字节从p中写入到数据流中，写入时从偏移量off开始。返回写入的字节数和任何遇到错误。如果n<len(p),返回的err非空。
```
type WriterAt interface {
	WriteAt(p []byte, off int64) (n int, err error)
}
```
# 8.Copy
从src中复制数据到dst直到在src中遇到EOF,或者产生错误。返回复制的bytes数量，和复制中遇到的错误。
一个成功的复制，返回的err==nil,而不是err==EOF.因为Copy被定义为从src中读取直到EOF.它并不会把EOF当做一个错误来对待。如果src实现了WriterTo接口，copy通过调用src.WriteTo(dst)来实现。否则，如果dst实现了ReaderFrom接口，copy通过调用dst.ReadFrom(src)来实现。
```
func Copy(dst Writer, src Reader) (written int64, err error) {
	return copyBuffer(dst, src, nil)
}
```
示例：
```
package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "os"
    "io"
)
func main(){
    content, err := ioutil.ReadFile("./demo.go")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Printf("file content: %s", content)
    r := strings.NewReader("Hello World!")
    n, err := io.Copy(os.Stdout, r)
    fmt.Printf("\n%d, %v", n, err)   //12, <nil>
}
```
# 9.StringWriter
将字符串内容写到w中，接受一个字节切片，如果w实现了StringWriter接口，则直接调用，调用w.Write方法。
```
type StringWriter interface {
	WriteString(s string) (n int, err error)
}

func WriteString(w Writer, s string) (n int, err error) {
	if sw, ok := w.(StringWriter); ok {
		return sw.WriteString(s)
	}
	return w.Write([]byte(s))
}
```
示例：
```

    nstr, errstr := io.WriteString(os.Stdout, "Hello world")
    fmt.Printf("\n%d, %v", nstr, errstr)
```