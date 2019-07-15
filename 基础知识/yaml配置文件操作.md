使用大神写的库`jinzhu/configor`
# 1.安装
`go get github.com/jinzhu/configor`
# 2.使用
```
package main

import (
	"fmt"
	"github.com/jinzhu/configor"
)

var Config = struct {
	APPName string `default:"app name"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"`
	}

	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}
}{}

func main() {
	configor.Load(&Config, "config.yml")
	fmt.Printf("config: %#v", Config)
}
```
yaml文件格式
```
#config.yml
appname: test

db:
    name:     test
    user:     test
    password: test
    port:     1234

contacts:
- name: i test
  email: test@test.com
```
输出:
```
config: struct { APPName string "default:\"app name\""; DB struct { Name string; User string "default:\"root\""; Password string "required:\"true\" env:\"DBPassword\""; Port uint "default:\"3306\"" }; Contacts []struct { Name string; Email string "required:\"true\"" } }{APPName:"test", DB:struct { Name string; User string "default:\"root\""; Password string "required:\"true\" env:\"DBPassword\""; Port uint "default:\"3306\"" }{Name:"test", User:"test", Password:"test", Port:0x4d2}, Contacts:[]struct { Name string; Email string "required:\"true\"" }{struct { Name string; Email string "required:\"true\"" }{Name:"i test", Email:"test@test.com"}}}%
```