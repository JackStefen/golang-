这个包主要是处理`UTF-8`编码的字符串的

诸如`func Contains(s, substr string) bool`此类的功能我们经常使用就不过多的介绍了，大家可以看看下面这个简单的示例


```
package main

import (
	"strings"
	"fmt"
	"os"
	"io/ioutil"
)
func main() {
	var str = "China,您好"
        fmt.Println(strings.Repeat(str, 2)) //China,您好China,您好

        strPrt := fmt.Sprint(str)
        fmt.Println(strPrt)//China,您好

	fmt.Println(strings.Count("cheese", "e"))                           //3
	fmt.Println(strings.Count("five", ""))                              //5
	fmt.Println(strings.EqualFold("Go", "go"))                          // true

	fmt.Println(strings.Contains(str, "China"))                          // true
	fmt.Println(strings.ContainsRune(str, '您'))                         // true

        fmt.Println(strings.Replace(str, "China", "USA", -1))              //USA,您好

	fmt.Printf("Fields are: %q\n", strings.Fields("  foo bar  baz   ")) //Fields are: ["foo" "bar" "baz"]

	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	//Fields are: ["foo1" "bar2" "baz3"]
	fmt.Printf("Fields are: %q", strings.FieldsFunc("  foo1;bar2,baz3...", f))

	rot13 := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+13)%26																
                }
		return r
	}
	fmt.Println(strings.Map(rot13, "'Twas brillig and the slithy gopher...")) //'Gjnf oevyyvt naq gur fyvgul tbcure...

	reader := strings.NewReader("widuuv web")
	fmt.Printf("%#v\n",reader)
	fmt.Println(reader.Len())//10
	n, err := reader.Read(make([]byte, 10))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n)//10 该值依据定义的字节数组的长度，如果定义数组长度为6，该值也为6

	 reader1 := strings.NewReader("China,你好,你好，你好")
         b := make([]byte, 10) 
         if n1, err := reader1.ReadAt(b, 2); err != nil {
                fmt.Println(err) 
        } else {
                fmt.Println(string(b[:n1]))  //ina,你好
        }

	reader2 := strings.NewReader("hello shanghai China")
	b1 := make([]byte, 8)
	n2, _ := reader2.Read(b1)
	fmt.Println(string(b1[:n2])) //hello sh
	reader2.Seek(2, 1)
	n3,_ := reader2.Read(b1)
	fmt.Println(string(b1[:n3])) //ghai Chi



	reader3 := strings.NewReader("hello shanghai")
	b2 := make([]byte, 4)
	n4, _ := reader3.Read(b2)
	fmt.Println(string(b2[:n4])) //hell
	reader3.Seek(2, 1)
	reader3.UnreadByte()
	n5, _ := reader3.Read(b2)
	fmt.Println(string(b2[:n5]))// 空格sh

	reader4 := strings.NewReader("hello aniu")
	w, _ := os.Create("aniu.txt")
	defer w.Close()
	n6, err := reader4.WriteTo(w)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n6) //15

	// ---------------Replacer--------------------
	r := strings.NewReplacer("<", "&lt;", ">", "&gt;")
	fmt.Println(r.Replace("This is <b>HTML</b>")) //This is &lt;b&gt;HTML&lt;/b&gt;
	n7,err := r.WriteString(w, "This is <b>Html</b>!")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(n7)//32
	d, _ := ioutil.ReadFile("aniu.txt")
	fmt.Println(string(d))//hello aniuThis is &lt;b&gt;Html&lt;/b&gt;!
}
```
aniu.txt
````
hello aniuThis is &lt;b&gt;Html&lt;/b&gt;!
````

今天，我们来看一下`Builder`这个结构体,这个`Builder`用于有效的构建一个字符串，通过`Write`方法，**其最小化内存拷贝**，零值就能被使用。但是不要对零值的`Builder`进行拷贝

```
type Builder struct {
	addr *Builder // of receiver, to detect copies by value
	buf  []byte
}

```
为什么着重介绍这个东西，因为它还是比较有用的，你见到他的几率还是比较高的，当然了，最终我们还是要自己学会使用他的。



```
// String returns the accumulated string.
func (b *Builder) String() string {
	return *(*string)(unsafe.Pointer(&b.buf))
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *Builder) Len() int { return len(b.buf) }

// Cap returns the capacity of the builder's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (b *Builder) Cap() int { return cap(b.buf) }

// Reset resets the Builder to be empty.
func (b *Builder) Reset() {
	b.addr = nil
	b.buf = nil
}
```
上面几个方法都是比较简单的方式，实现了基本的操作和属性，


下面的特别注意了，也是核心功能
```
// Grow grows b's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to b
// without another allocation. If n is negative, Grow panics.
func (b *Builder) Grow(n int) {
	b.copyCheck()
	if n < 0 {
		panic("strings.Builder.Grow: negative count")
	}
	if cap(b.buf)-len(b.buf) < n {
		b.grow(n)
	}
}

func (b *Builder) Write(p []byte) (int, error) {
	b.copyCheck()
	b.buf = append(b.buf, p...)
	return len(p), nil
}

func (b *Builder) WriteByte(c byte) error {
	b.copyCheck()
	b.buf = append(b.buf, c)
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
// It returns the length of r and a nil error.
func (b *Builder) WriteRune(r rune) (int, error) {
	b.copyCheck()
	if r < utf8.RuneSelf {
		b.buf = append(b.buf, byte(r))
		return 1, nil
	}
	l := len(b.buf)
	if cap(b.buf)-l < utf8.UTFMax {
		b.grow(utf8.UTFMax)
	}
	n := utf8.EncodeRune(b.buf[l:l+utf8.UTFMax], r)
	b.buf = b.buf[:l+n]
	return n, nil
}

func (b *Builder) WriteString(s string) (int, error) {
	b.copyCheck()
	b.buf = append(b.buf, s...)
	return len(s), nil
}
```
先看最常用的吧，`WriteString`。它就是往`buffer`中追加数据，而追加的方式也比较直接，`append`操作，将字符串直接点点点。这样就说明了一个问题，字符串的特殊结构。
这个方法的返回值就更简单了，参数字符串的长度，还有`nil`.


再看看`Write`方法，干的一样的活，就是参数是个字节数组。`WriteByte`就更简单了，写入的是单个字节。

`WriteRune`方法写入`UTF-8`编码的`Unicode`码点`r`到`buffer`中。返回`r`的长度和`nil`.

到这里你可能觉得，卧槽，这没什么啊，不都很简单嘛，方法名都很说明问题，的确是。

不过相信你也发现了，在`WriteString`,`WriteRune`这些方法中，在追加数据到`buffer`中之前，都在做这个事情`b.copyCheck`。


```
func (b *Builder) copyCheck() {
	if b.addr == nil {
		// This hack works around a failing of Go's escape analysis
		// that was causing b to escape and be heap allocated.
		// See issue 23382.
		// TODO: once issue 7921 is fixed, this should be reverted to
		// just "b.addr = b".
		b.addr = (*Builder)(noescape(unsafe.Pointer(b)))
	} else if b.addr != b {
		panic("strings: illegal use of non-zero Builder copied by value")
	}
}
```

就是检测一下，这个`Builder`是否进行了复制操作。如果地址发生变，直接`panic`

最后，如果在使用之前明确知道自己需要多少内存，可以在使用之前，进行容量分配。如果当前所剩余空间小于参数`n`.将进行`grow`操作

```
func (b *Builder) grow(n int) {
	buf := make([]byte, len(b.buf), 2*cap(b.buf)+n)
	copy(buf, b.buf)
	b.buf = buf
}

```

重新申请空间，其容量大小为当前容量的两倍加上n，并进行数据拷贝.



以上内容就是关于`strings`包中关于`Builder`的全部内容了，发现了么，这个和`bytes`中的`Buffer`有比较类似的地方。


看看实际例子中关于这个结构体的使用吧
```
// Do executes the request and returns response or error.
//
func (r CatCountRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
	)

	method = "GET"

	path.Grow(1 + len("_cat") + 1 + len("count") + 1 + len(strings.Join(r.Index, ",")))
	path.WriteString("/")
	path.WriteString("_cat")
	path.WriteString("/")
	path.WriteString("count")
	if len(r.Index) > 0 {
		path.WriteString("/")
		path.WriteString(strings.Join(r.Index, ","))
	}

	params = make(map[string]string)

	if r.Format != "" {
		params["format"] = r.Format
	}

	if len(r.H) > 0 {
		params["h"] = strings.Join(r.H, ",")
	}

	if r.Help != nil {
		params["help"] = strconv.FormatBool(*r.Help)
	}

	if len(r.S) > 0 {
		params["s"] = strings.Join(r.S, ",")
	}

	if r.V != nil {
		params["v"] = strconv.FormatBool(*r.V)
	}

	if r.Pretty {
		params["pretty"] = "true"
	}

	if r.Human {
		params["human"] = "true"
	}

	if r.ErrorTrace {
		params["error_trace"] = "true"
	}

	if len(r.FilterPath) > 0 {
		params["filter_path"] = strings.Join(r.FilterPath, ",")
	}

	req, err := newRequest(method, path.String(), nil)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if len(r.Header) > 0 {
		if len(req.Header) == 0 {
			req.Header = r.Header
		} else {
			for k, vv := range r.Header {
				for _, v := range vv {
					req.Header.Add(k, v)
				}
			}
		}
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := transport.Perform(req)
	if err != nil {
		return nil, err
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	return &response, nil
}

```
还是
还是`go-elasticsearch`包，构建关于请求参数的数据。首先申请的空间，然后写入字符串，最后`path.String()`，`over`.一气呵成，
