数据的json编码主要涉及两个方法调用：
# 1.Marshal方法
该方法返回参数v的json编码
```
func Marshal(v interface{}) ([]byte, error) {
	e := newEncodeState()

	err := e.marshal(v, encOpts{escapeHTML: true})
	if err != nil {
		return nil, err
	}
	buf := append([]byte(nil), e.Bytes()...)

	e.Reset()
	encodeStatePool.Put(e)

	return buf, nil
}
```
该方法可以接受任意类型的参数，并将其转化为字节序列，如果遇到的值实现了Marshaler接口，并且不是一个空指针，Marshal方法就会调用其上的MarshalJSON方法来产生JSON对象。如果没有MarshalJSON方法，但是它实现了encoding.TextMarshaler，Marshal调用marshalText方法，并将结果编码为JSON字符串。否则，Marshal方法使用以下的依赖于类型的默认编码；

- boolean值编码为JSON的boolean值
- 浮点数，整型，Number类型的值编码为JSON的number值
- 字符串类型编码为有效的UTF-8编码的JSON字符串值，使用rune类型替换无效的字节，尖括号<、>变为"\u003c","\u003e",以防止某些浏览器将JSON输出错误的解释为HTML."&"变成"\u0026"也是同样的原因。可以通过SetEscapeHTML(false)来禁止此转义。
- 数组和切片类型编码为JSON数组，除了[]byte字节数组编码为bash64的字符串。nil的切片会编码为null的JSON值。
```
    var a []int = []int{1, 2, 3, 4, 6}
	sliceJSON, err := json.Marshal(a)
	fmt.Println(string(sliceJSON)) //[1,2,3,4,6]
```
- 结构体编码为JSON对象。每一个导出的结构体字段变成JSON对象的成员，使用字段名作为JSON对象的一个key.除非字段以以下的方式被省略了：
每个字段的编码可以通过结构体的字段标签来定制。标签可以给出字段的名称，也可以为空， 来使用默认的字段名。omitempty可选项指定字符如果为空值，在编码的时候，需要忽略，以下值可以定义为空值：false,0,nil指针，nil接口，空数组，空字符串，空字典，空切片。特别的，如果字段标签为"-"",字段将在编时被忽略，如果标签为"-,"",编码后，出现在JSON对象中时，字段名将为-。
- "string"字符串标签标示字段将被编码为JSON字符串，该标签仅在字段类型为字符串，浮点数，整数，布尔类型时，使用。
- 指针类型的值编码为指针指向的值，空指针将被编码为null值。
- 接口类型的值将被编码为接口中包含的值，空的接口值将被编码为null值。
- channel,complex,函数类型的值不能被编码，试图编码将报错。
- 循环引用的数据Marshal将不处理它们，传入一个循环结构的数据到Marshal中将导致一个无线循环。




```
package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type User struct {
	Name  string `json:"user_name"`
	Age   int    `json:"user_age"`
	Addr  string `json:"address,omitempty"`
	Email string `json:",omitempty"`
	*IphoneInfo
}

type IphoneInfo struct {
	Id     int64       `json:"iphone_id"`
	Name   string      `json:"iphone_name"`
	IpAddr interface{} `json:"iphone_ip_addr"`
	IphoneType
}

type IphoneType struct {
	IphoneNumber string `json:"iphone_number"`
	Is5GIphone   bool   `json:"is_5g"`
}

func main() {
	var user User = User{
		Name: "yechongqiu",
		Age:  18,
		Addr: "Nanjiing",
		// Email: "yechongqiu@sunning.com",
	}
	fmt.Println(reflect.ValueOf(user) == reflect.ValueOf(&user))
	fmt.Println(reflect.ValueOf(&user))
	jsonUser, err := json.Marshal(&user)
	if err != nil {
		return
	}
	fmt.Println(string(jsonUser))
	var a []int = []int{1, 2, 3, 4, 6}
	sliceJSON, err := json.Marshal(a)
	fmt.Println(string(sliceJSON))

	iphoneInfo := &IphoneInfo{
		Id:     100001,
		Name:   "iphone 8",
		IpAddr: "192.168.10.101",
	}
	user.IphoneInfo = iphoneInfo
	jsonUser2, err := json.Marshal(user)
	if err != nil {
		return
	}
	fmt.Println(string(jsonUser2))
}

```
Output:
```
false
&{yechongqiu 18 Nanjiing  <nil>}
{"user_name":"yechongqiu","user_age":18,"address":"Nanjiing"}
[1,2,3,4,6]
{"user_name":"yechongqiu","user_age":18,"address":"Nanjiing","iphone_id":100001,"iphone_name":"iphone 8","iphone_ip_addr":"192.168.10.101","iphone_number":"","is_5g":false}

```

# Unmarshal

将JSON对象解码为结构体数据，方法的参数为：
- 需要解码的JSON对象的字节序列
- 将JOSN对象解码到的结构体，一般为指针。

```
func Unmarshal(data []byte, v interface{}) error {
	// Check for well-formedness.
	// Avoids filling out half a data structure
	// before discovering a JSON syntax error.
	var d decodeState
	err := checkValid(data, &d.scan)
	if err != nil {
		return err
	}

	d.init(data)
	return d.unmarshal(v)
}

```


如果第二个参数为空，或者不是一个指针，将返回一个错误。Unmarshal使用Marshal相反的编码。使用以下附加规则：

为了解码JSON到一个实现了Unmarshaler接口的值，Unmarshal方法调用值的UnmarshalJSON方法。包含输入为JSON null的情况。否则，如果值实现了encoding.TextUnmarshaler接口，Unmarshal方法将调用值的该方法。
为了解码一个JOSN到一个结构体，Unmarshal匹配输入类型的key到Marshal时使用的key.


解码JSON到一个接口值时，将如下的值保存到接口中去：

- json boolean -> bool
- json number -> float64
- json string -> string
- json arrays -> []interface{}
- json objects -> map[string]interface{}
- json null -> nil
- 

为了将一个JSON数组解码到一个切片中，Unmarshal将重置切片长度为0，然后再将每个元素追加到切片中。特别的，将一个空的JSON数组解码到一个切片中，Unmarshal使用一个空的切片替换切片。

为了将一个JSON对象解码到map中，Unmarshal 首先建立一个map，如果map为nil，将重新分配一个，否则直接使用它，保留之前的键值对，然后将JSON对象的key-value对存到map中去。map的键类型必须是一个字符串，或者整数，或者实现了excoding.TextUnmarshaler接口的类型。

如果一个JSON值不适合给的目标值的类型。或者JSON数字值范围超过了目标类型的范围，Unmarshal将跳过它，尽可能的完成其他解码。如果没有遇到更加严重的错误，Unmarshal将返回一个描述早前错误的异常。


```
    err = json.Unmarshal(jsonUser2, &userFromJSON)
	if err != nil {
		return
	}
	fmt.Println(userFromJSON)
```