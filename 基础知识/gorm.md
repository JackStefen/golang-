Gorm当前支持MySql, PostgreSql, Sqlite等主流数据库
# 1.安装驱动
`go get github.com/go-sql-driver/mysql`
# 2.小示例
```
package main

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql" // 包装
)

type User struct {
   Id int64
   UserId int64
   AddId int64
   Name string
   Address string
}

type Address struct {
    Id int64
    UserId int64
    AddId int64
    AddName string
    AddLocation string
}


func main() {
    db, err := gorm.Open("mysql", "root:123456@/guolianlc?charset=utf8&parseTime=True&loc=Local")
    if err != nil {
         fmt.Println("connect db error: ", err)
    }
    defer db.Close()
    if db.HasTable(&User{}) {
        db.AutoMigrate(&User{})
    } else {
        db.CreateTable(&User{})
    }
    //db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Address{})
    db.AutoMigrate(&Address{})
    db.Model(&User{}).AddForeignKey("add_id", "addresses(id)", "RESTRICT", "RESTRICT")
    db.Model(&User{}).AddForeignKey("add_id", "addresses(id)", "RESTRICT", "RESTRICT")
    db.Model(&User{}).AddIndex("idx_user_add_id", "add_id")
    db.Model(&User{}).AddUniqueIndex("idx_user_id", "user_id")
}
```
# 3.表级别操作
- `AutoMigrate() `
`db.AutoMigrate(&Address{})`
`AutoMigrate()`运行后，会自动migrate对应的model.仅仅新增新增的字段，不会进行修改已有的字段类型，删除字段的操作
- `HasTable()`
检查表是否存在
- `CreateTable()`
`db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Address{})`
创建表
默认情况下，表名为结构体名的复数形式，当然也可以禁用；
`db.SingularTable(true)`
- `DropTable()/ DropTableIfExists()`
删除表
- `ModifyColumn()`
 修改列
- `DropColumn()`
删除列
- `AddForeignKey()`
参数 : 1th:外键字段,2th:外键表(字段),3th:ONDELETE,4th:ONUPDATE
` db.Model(&User{}).AddForeignKey("add_id", "addresses(id)", "RESTRICT", "RESTRICT")`
两个表中的字段都必须存在，就像Users表中的add_id字段，如果不存在，无法自动新增字段，并自动创建外键
- `AddIndex() / AddUniqueIndex`
添加索引，添加唯一值索引
```
    db.Model(&User{}).AddForeignKey("add_id", "addresses(id)", "RESTRICT", "RESTRICT")
    db.Model(&User{}).AddIndex("idx_user_add_id", "add_id")
    db.Model(&User{}).AddUniqueIndex("idx_user_id", "user_id")
```
- `RemoveIndex()`
删除索引

示例中最终创建出的表结构：
```
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `add_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`),
  KEY `idx_user_add_id` (`add_id`),
  CONSTRAINT `users_add_id_addresses_id_foreign` FOREIGN KEY (`add_id`) REFERENCES `addresses` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```


# 4.表结构设计以及gorm标签的使用
go中使用结构体来作为表结构设计的载体，实例：
```
package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // 包装
)

type User struct {
	gorm.Model
	UserId    int64 `gorm:"index"`
	Birthday  time.Time
	Age       int           `gorm:"column:age"`                     //可定制列表名
	Name      string        `gorm:"size:255;index:idx_name_add_id"` // string默认长度为255, 使用这种tag重设。
	Num       int           `gorm:"AUTO_INCREMENT"`                 // 自增
	Email     string        `gorm:"type:varchar(100);unique_index"`
	AddressID sql.NullInt64 `gorm:"index:idx_name_add_id"`
	IgnoreMe  int           `gorm:"-"` // 忽略这个字段
        Desction  string        `gorm:"size:2049;comment:'用户描述字段'"`
Status       string `gorm:"type:enum('published','pending','deleted');default:'pending'"`
}

//设置表名，默认是结构体的名的复数形式
func (User) TableName() string {
	return "VIP_USER"
}

func main() {
	db, err := gorm.Open("mysql", "root:123456@/guolianlc?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("connect db error: ", err)
	}
	defer db.Close()
	if db.HasTable(&User{}) {
		db.AutoMigrate(&User{})
	} else {
		db.CreateTable(&User{})
	}
}
```
插入一条测试语句后，查询表结构如下：
![image.png](https://upload-images.jianshu.io/upload_images/3004516-c07920c8904e6031.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

gorm.Model为内建的结构体，结构如下：
```
// 基本模型的定义
type Model struct {
  ID        uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt *time.Time
}
```
创建出来的表结构为：
```
CREATE TABLE `VIP_USER` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `birthday` timestamp NULL DEFAULT NULL,
  `age` int(11) DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `num` int(11) DEFAULT NULL,
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address_id` bigint(20) DEFAULT NULL,
  `desction` varchar(2049) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户描述字段',
  `status` enum('published','pending','deleted') COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_VIP_USER_email` (`email`),
  KEY `idx_VIP_USER_deleted_at` (`deleted_at`),
  KEY `idx_VIP_USER_user_id` (`user_id`),
  KEY `idx_name_add_id` (`name`,`address_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

![组合索引.png](https://upload-images.jianshu.io/upload_images/3004516-ec53b38e5a3c010a.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

数据插入时，仅仅插入业务数据即可，`created_at` 和`updated_at`,`deleted_at`字段不用手动设置值，gorm会帮我们自动维护这些字段的值，当首次插入时，`created_at`和`updated_at`字段的值是相同的，都为当前数据记录插入的时间
```
  var user User = User{
             UserId: 1,
             Birthday: time.Now(),
             Age: 12,
             Name:"zhangsan",
             Num: 12,
             Email:"zhangsan@alibaba.com",
             AddressID:sql.NullInt64{Int64 : int64(1), Valid : err == nil},
             Desction:"first",
        }
        if err := db.Model(&User{}).Create(&user).Error; err != nil{

        }
```
![数据插入后的记录详情.png](https://upload-images.jianshu.io/upload_images/3004516-972b77f8c26b877d.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

当数据执行删除操作时，默认情况下执行的是软删除，仅仅设置`deleted_at`字段的值，为执行删除操作的时间
```
if err := db.Model(&User{}).Where("user_id=?", 1).Delete(&User{}).Error; err != nil {
        }
```

![执行删除操作.png](https://upload-images.jianshu.io/upload_images/3004516-a98b28b998e49fed.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

如果业务上需要，读取包含软删除的数据，可以在查询时加上
```
var usr = make([]*User,0)
        if err := db.Unscoped().Model(&User{}).Where("user_id=?", 1).Find(&usr).Error; err != nil {}
        for _, usser := range usr {
            fmt.Println(usser)
        }
```
Output:
```
&{{1 2019-05-15 13:34:23 +0800 CST 2019-05-15 13:34:23 +0800 CST 2019-05-15 13:42:13 +0800 CST} 1 2019-05-15 13:34:23 +0800 CST 12 zhangsan 12 zhangsan@alibaba.com {1 true} 0 first}
```

如果需要永久的删除数据，也就是物理删除，可以在`Unscoped()`的基础上，执行`Deleted()`
```
 if err := db.Unscoped().Model(&User{}).Where("user_id=?", 1).Delete(&User{}).Error; err != nil {}
```

![物理删除.png](https://upload-images.jianshu.io/upload_images/3004516-6eb0be721d35a888.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)


# 5.增删改查
##增
```
if err := tx.Model(&model.Teatures{}).Create(teatureRecord).Error; err != nil {
				ErrMsg := fmt.Sprintf("%s", err)
				if strings.HasPrefix(ErrMsg, "Error 1062: Duplicate entry") {
					continue
				}
				logrus.Errorln("updateUser, sava user teature err: ", err)
				tx.Rollback()
				return
			}
```
## 删
```
if err := tx.Where("created_at=?", int64Time).Delete(&model.User{}).Error; err != nil {
		tx.Rollback()
		logrus.Errorln("updateUser , delete user err: ", err)
		return
	}
```
## 改
```
if err := common.Db.Model(&model.User{}).
						Where("created_at=? and usr_id=? and usr_name=? and usr_code=?", int64Time, UserId, UserName, UserCode).
						Updates(
							map[string]interface{}{
								......
							}).Error; err != nil {
						logrus.Errorln("UpdateUsers, update user record err: ", err)
						return
					}
```
本示例中，给出的更改语法使用的是map字典，当然你也可以传入数据库字典结构体，但是需要注意的是：
在实际应用中，我们的文章编辑后台，在删除某个字段后，比如文章的摘要，进行更新提交时，发现更改并未生效，查看后台gorm转化成的sql语句，发现并没有更新摘要字段，这是因为，当结构体的某个字段为零值的时候，传入到updates方法中，并没有显示该字段，而udpates方法是根据该结构体有值的字段进行更新的，没有值的字段，并没有做任何操作，所以上述进行的更新也未起作用，这些细节需要格外注意
## 查
```
	dbUsers := make([]*model.User, 0)
	if err := common.Db.Model(&model.User{}).Where("created_at=?", int64Time).Find(&dbUsers).Error; err != nil {
		logrus.Errorln(err)
		return
	}
```
# 6.事务操作
```
        tx := common.Db.Begin()
	// 1.删除数据

	if err := tx.Where("created_at=?", int64Time).Delete(&model.User{}).Error; err != nil {
		tx.Rollback()
		logrus.Errorln("updateUser , delete user err: ", err)
		return
	}
	// 2. 插入新数据
	for useOrder, user := range users {
		for teatureOrder, teature := range user.Teatures {
			var teatureRecord = &model.Teatures{
				......
			}
			
			if err := tx.Model(&model.Teatures{}).Create(teatureRecord).Error; err != nil {
				ErrMsg := fmt.Sprintf("%s", err)
				if strings.HasPrefix(ErrMsg, "Error 1062: Duplicate entry") {
					continue
				}
				logrus.Errorln("updateUser, sava user teature err: ", err)
				tx.Rollback()
				return
			}
		}
	}
	// 提交事务操作
	tx.Commit()
```
