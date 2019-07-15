# 1.字符串
一个字符串是一个不可改变的字节序列。和数组不同的是，字符串的元素不可修改，是一个只读的字节数组。每个字符串的长度虽然也是固定的，但是字符串的长度并不是字符串类型的一部分。
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
strings包下含有许多操作字符串的工具
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