# 1.问题描述
model中一个表的字段为时间类型的字段time.Time.
在数据库设计中该字段可以为空。
可在实际写数据时，发生错误
```
INFO[0004]
[2018-02-11 10:04:39]  Error 1292: Incorrect datetime value: '0000-00-00' for column 'published_at' at row 1  
[2018-02-11 10:04:39]  [1.42ms]  INSERT INTO `xxxx` (`author_id`,`status`,`title`,`summary`,`content`,`enable_risk_levels`,`image`,`published_at`,`created_at`,`updated_at`,`deleted_at`) VALUES ('1','1','拉丝机弗兰克','方法孟老师','','','','0001-01-01 00:00:00','2018-02-11 10:04:39','2018-02-11 10:04:39',NULL) 
```
我们在表单中只提供了标题，摘要，内容，图片，级别的表单，其他都是默认值，create_at,update_at 为当前时间点。在golang 中，time.Time的默认时间值为'0001-01-01 00:00:00'
```
package main

import (
     "fmt"
     "time"
)

func main () {
    t := new(time.Time)
    fmt.Println(t)
}
```
Output:
```
0001-01-01 00:00:00 +0000 UTC
```
# 2,解决方案
根据别人的提议是因为数据库的sql_mode的配置
执行
`select @@sql_mode`
Output:
![image.png](http://upload-images.jianshu.io/upload_images/3004516-aafcbe2babc234cf.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

把模式中的NO_ZERO_IN_DATE, NO_ZERO_DATE除去即可。
```
mysql> set global sql_mode='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
    -> ;
Query OK, 0 rows affected, 1 warning (0.00 sec)

```
global 为全局设置，session为当前会话生效。
# 3.sql_mode介绍
[别人的介绍链接](http://blog.csdn.net/ccccalculator/article/details/70432123)