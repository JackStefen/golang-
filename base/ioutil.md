# 1.ioutil实现一些I / O实用方法。
`import "io/ioutil"`
- `var Discard io.Writer = devNull(0)`
Discard是一个io.Writer，所有Write调用都可以在不执行任何操作的情况下成功完成。
- `func ReadAll(r io.Reader) ([]byte, error)`
从r读取数据直到遇到错误或者EOF,返回读到的数据，一个成功的调用返回的err==nil
```
package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
)

func main(){
    url := "http://www.baidu.com"
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(string(body))
}
```
- `func ReadFile(filename string) ([]byte, error) `
从以参数为文件名的文件中读取文件内容。成功调用返回err == nil.因为ReadFile读取整个文件，所以它不会将Read中的EOF视为要报告的错误。
```
package main

import (
    "fmt"
    "io/ioutil"
)
func main(){
    content, err := ioutil.ReadFile("./demo.go")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Printf("file content: %s", content)
}
```
- `func WriteFile(filename string, data []byte, perm os.FileMode) error`
WriteFile将数据写入由filename命名的文件。 如果该文件不存在，则WriteFile使用权限perm创建它; 否则WriteFile会在写入之前截断它。
```
package main

import (
    "io/ioutil"
    "fmt"
)

func main(){
    content:= "Hello world"

    if err := ioutil.WriteFile("./test.txt", []byte(content), 0644); err != nil {
        fmt.Println(err)
    }
    fmt.Println("write file success...")
}
```
- `func ReadDir(dirname string) ([]os.FileInfo, error)`
ReadDir读取由dirname命名的目录，并返回按filename排序的目录条目列表。
```
package main

import (
   "fmt"
   "io/ioutil"
)

func main() {
    files, err := ioutil.ReadDir(".")
    if err != nil {
        fmt.Println(err)
    }
    for _, file := range files {
        fmt.Println(file.Name())
    }
}
```
- `func TempDir(dir, prefix string) (name string, err error)`
TempDir在目录dir中创建一个新的临时目录，其名称以prefix开头，并返回新目录的路径。 如果dir是空字符串，TempDir将使用临时文件的默认目录（请参阅os.TempDir）。 同时调用TempDir的多个程序将不会选择相同的目录。 调用者有责任在不再需要时删除目录.
- `func TempFile(dir, pattern string) (f *os.File, err error)`

TempFile在目录dir中创建一个新的临时文件，打开文件进行读写，并返回生成的* os.File。 文件名是通过获取模式并在末尾添加随机字符串生成的。 如果pattern包含“*”，则随机字符串将替换最后一个“*”。 如果dir是空字符串，则TempFile使用临时文件的默认目录（请参阅os.TempDir）。 同时调用TempFile的多个程序不会选择相同的文件。 调用者可以使用f.Name（）来查找文件的路径名。 当不再需要时，调用者有责任删除该文件。

