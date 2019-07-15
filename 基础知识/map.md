 map是一堆键值对的未排序集合，类似Python中字典的概念，它的格式为map[keyType]valueType，是一个key-value的hash结构。map的读取和设置也类似slice一样，通过key来操作，只是slice的index只能是int类型，而map多了很多类型，可以是int，可以是string及所有完全定义了==与!=操作的类型。
map线程不安全的数据结构，如果需要线程安全，需要加锁，或者直接只用sync包中的map
# 1.声明
`var 变量名 map[keytype] valuetype`
# 2.初始化
`var amay map[int] string = map[int] string {1:"zhao",2:"qian",3:"sun",4:"li"}`
`var bmap map[string] string = make(map[string] string)`
`var cmap map[string] string = make(map[string] string, 4)` //预先设置键值对数量，避免添加键值对时，重新分配内存
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
# 6.实例
```
package main

import "fmt"

func main() {
    //声明一个map类型的结果，key为int型，value为string类型
    var amap map[int] string
    amap = map[int] string{1:"zhao", 2:"qian"}
    fmt.Println(amap)
    //判断一个key是否在map中的方法。如果在true,不在ok为false
    if val, ok := amap[4]; ok {
        fmt.Println(val)
    }
    //range遍历map结构,ky为key,vl为对应的value
    for ky, vl := range amap {
        fmt.Println(ky, vl)
    }

    var bmap map[string] string
    bmap = make(map[string] string)
    fmt.Println(bmap) //bmap为一个空的map output: map[]
    bmap["name"] = "zhanglinpeng"
    bmap["address"] = "tianhetiyuchang"
    bmap["location"] = "ss"
    fmt.Println(bmap) //map[name:zhanglinpeng address:tianhetiyuchang location:ss]
    //delete(bmap, "location") // 删除一个key-value
    for k, v := range bmap {
        fmt.Println(k + " is " + v)
    }
    var cmap map[string][]map[string] string //声明一个嵌套的 map类型。value 是一个切片，元素类型还是切片
    cmap = make(map[string] []map[string]string)
    vmap := []map[string]string{map[string]string{"id": "12", "name":"zhao"}}
    cmap["local"] = vmap
    fmt.Println(cmap) //map[local:[map[id:12 name:zhao]]] 看出来，local的值为数组


    type Key struct {
        Path, Country string
    }

    // 创建一个key为结构体的map类型
    hits := make(map[Key]int)
    hits[Key{"/", "cn"}]++
    fmt.Println(hits) //map[{/ cn}:1] 看的出来key为一个结构体

    // 创建一个嵌套的map key为string类型， 值为map类型
    m:=make(map[string]map[string]int)
    c:=make(map[string]int)
    c["b"]=1
    m["a"]=c
    fmt.Println(m) //map[a:map[b:1]]
}

```
# 7.隔靴挠痒
```
package main

import (
    "fmt"
    "reflect"
    "unsafe"
)

func main(){
    mapa:= make(map[string]int, 1)
    mapa["zhao"] = 1
    mapa["qian"] = 2
    fmt.Println(mapa) //map[qian:2 zhao:1]
    fmt.Println(reflect.TypeOf(mapa)) //map[string]int
    fmt.Println(unsafe.Sizeof(mapa)) //8
}
```