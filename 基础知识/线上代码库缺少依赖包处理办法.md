## 1.本地如果为安装govender，需要安装一下，具体命令为：
`go get -u -v github.com/kardianos/govendor`

## 2.进到项目的vender目录下，如果新项目没有vendor目录，执行
`govendor init`
`cd vendor/`
vender目录下有个`vendor.json`文件，
这个 vendor.json 会类似 godep 工具中的描述文件版本的功能

## 3.执行命令`govendor add +external`来将依赖的第三方包加到vendor中

## 4.使用git 讲修改的verdor.json文件和vendor目录下的第三方包上传到线上代码库中