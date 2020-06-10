在Go中，可以使用time包提供的`time.Now()`获取当前时间。

Golang提供了用于测量和显示时间的时间包。 您可以根据所选时区获取当前时间，使用go time包添加持续时间与当前时区等。


## 1.基本案例

Go不使用`yyyy-mm-dd`布局来格式化时间。 而是，格式化一个特殊的布局参数

```
Mon Jan 2 15:04:05 -0700 MST 2006
```

与时间或日期的格式相同。

```
package main

import (
	"fmt"
	"time"
)

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

func main() {
	date := "1999-12-31"
	t, _ := time.Parse(layoutISO, date)
	fmt.Println(t)                  // 1999-12-31 00:00:00 +0000 UTC
	fmt.Println(t.Format(layoutUS)) // December 31, 1999
}

```

- `time.Parse`解析日期字符串
- `Format`格式化`time.Time`


它们具有以下签名：
```
func Parse(layout, value string) (Time, error) 
func (t Time) Format(layout string) string 
```

## 2.Date Time in Go

日期函数返回对应的时间

```
package main

import (
	"fmt"
	"time"
)

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

func main() {
	date := time.Date(2018, 01, 12, 22, 51, 48, 324359102, time.UTC)
	fmt.Printf("date is :%s \n", date)
	date = time.Date(2018, 01, 12, 22, 51, 48, 324359102, time.UTC)
	fmt.Printf("date is :%s \n", date)
	date = time.Now().UTC()
	fmt.Printf("current date is :%s", date) // run on local env
}
```
```
➜  timetest go run main.go
date is :2018-01-12 22:51:48.324359102 +0000 UTC
date is :2018-01-12 22:51:48.324359102 +0000 UTC
current date is :2020-05-28 06:25:07.991531 +0000 UTC
```

## 3. 对时间进行Add运算

```
package main

import (
	"fmt"
	"time"
)

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

func main() {
	date := time.Date(2018, 01, 12, 22, 51, 48, 324359102, time.UTC)
	next_date := date.AddDate(1, 2, 1)

	fmt.Printf("date is :%s\n", date)
	fmt.Printf("date after adding (1,2,1) is :%s \n", next_date)

	// use date.Add to add or substract time with (+ -) symbol
	next_date1 := date.Add(+time.Hour * 24)
	fmt.Printf("date after adding 24 hour is :%s \n", next_date1)

	next_date2 := date.Add(-time.Hour * 24)
	fmt.Printf("date after substracting 24 hour is :%s \n", next_date2)
}

```

```
➜  timetest go run main.go
date is :2018-01-12 22:51:48.324359102 +0000 UTC
date after adding (1,2,1) is :2019-03-13 22:51:48.324359102 +0000 UTC
date after adding 24 hour is :2018-01-13 22:51:48.324359102 +0000 UTC
date after substracting 24 hour is :2018-01-11 22:51:48.324359102 +0000 UTC
```

## 4.在Golang中获取当前的UNIX时间戳
您可以使用`now()`方法获取当前时间，它具有`Unix()`方法，可帮助将时间转换为golang中的UNIX时间戳。
```
package main

import (
	"fmt"
	"time"
)

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

func main() {
	cur_time := time.Now().Unix()
	fmt.Printf("current unix timestamp is :%v\n", cur_time)
}

```
```
➜  timetest go run main.go
current unix timestamp is :1590647369
```
逆运算是将Unix时间戳格式转化为Time格式

```
package main

import (
	"fmt"
	"time"
)

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

func main() {
	cur_time := time.Unix(time.Now().Unix(),0)
	fmt.Printf("current unix timestamp is :%v\n", cur_time)
}
```
```
➜  timetest go run main.go
current unix timestamp is :2020-05-28 14:54:26 +0800 CST
```
## 5.在Golang中解析日期字符串

您可以使用parse方法解析golang中的日期字符串。

```
package main

import (
	"fmt"
	"time"
)

func main() {
	date := "2018-10-24T18:50:23.541Z"
	parse_time, _ := time.Parse(time.RFC3339, date)

	fmt.Printf("current unix timestamp is :%s\n", date)
	fmt.Printf("parse_time is :%s", parse_time)
}

```

```
➜  timetest go run main.go
current unix timestamp is :2018-10-24T18:50:23.541Z
parse_time is :2018-10-24 18:50:23.541 +0000 UTC%
```

Go在time包中为常用格式提供了一些方便的常量
```
const (
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	// Handy time stamps.
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
)
```

## 6.如何在本地和其他时区使用时间戳获取当前日期和时间

`LoadLocation`返回具有给定名称的`Location`。
```
package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	fmt.Println("Location : ", t.Location(), " Time : ", t) // local time

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Location : ", location, " Time : ", t.In(location)) // America/New_York

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	fmt.Println("Location : ", loc, " Time : ", now) // Asia/Shanghai
}

```


```
➜  timetest go run main.go
Location :  Local  Time :  2020-05-28 14:34:54.683585 +0800 CST m=+0.000270177
Location :  America/New_York  Time :  2020-05-28 02:34:54.683585 -0400 EDT
Location :  Asia/Shanghai  Time :  2020-05-28 14:34:54.686005 +0800 CST
```

## 7.如何获取Weekday和YearDay

Weekday返回由t指定的星期几。YearDay返回由t指定的一年中的日期，非闰年为[1,365]，闰年为[1,366]。

```
package main

import (
	"fmt"
	"time"
)

func main() {
	t, _ := time.Parse("2006 01 02 15 04", "2015 11 11 16 50")
	fmt.Println(t.YearDay()) // 315
	fmt.Println(t.Weekday()) // Wednesday

	t, _ = time.Parse("2006 01 02 15 04", "2011 01 01 0 00")
	fmt.Println(t.YearDay())
	fmt.Println(t.Weekday())
}

```

```
➜  timetest go run main.go
315
Wednesday
1
Saturday
```

## 8.使用golang获取各种格式的当前日期和时间

在格式化时间格式的时候，可以自定义时间格式，而不使用time包中的定义的常量格式

```
package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	fmt.Println("Curret Time: ", t.Format("2006-01-02 15:04:05"))
}
```
```
➜  timetest go run main.go
Curret Time:  2020-05-28 14:41:19
```

## 9. 两个时间值的减法

```
package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	time.Sleep(1 * time.Second)
	stop := time.Now()
	fmt.Printf("The call took %v to run.\n", stop.Sub(start))
}

```

```
➜  timetest go run main.go
The call took 1.000261167s to run.
```

## 10.日期时间的前后关系判断
```
package main

import (
    "time"
    "fmt"
)
func main() {
    t1 := time.Now()
    t2 := t1.Add(time.Minute) //加一个分钟
    t3 := t1.Add(-3*time.Hour) // 减3个小时
    //分别指定年，月，日，时，分，秒，纳秒，时区
    ft := time.Date(2018, time.Month(1), 11, 5, 13, 32, 0, t1.Location())
    fmt.Println(ft)
    fmt.Println(t1, t2)
    if t3.After(t1) {
        fmt.Println("t3 is after t1")
    } else {
        fmt.Println("t3 is before t1")
    }
    fmt.Println(time.Since(t3))
}
```
```
➜  timetest go run demo1.go
2018-01-11 05:13:32 +0800 CST
2020-05-28 14:49:08.042055 +0800 CST m=+0.000311943 2020-05-28 14:50:08.042055 +0800 CST m=+60.000311943
t3 is before t1
3h0m0.000207336s
```

## 11.NewTimer
```
package main

import (
    "time"
    "fmt"
)

func main() {

	// Timers represent a single event in the future. You
	// tell the timer how long you want to wait, and it
	// provides a channel that will be notified at that
	// time. This timer will wait 2 seconds.
	timer1 := time.NewTimer(2 * time.Second)

	// The `<-timer1.C` blocks on the timer's channel `C`
	// until it sends a value indicating that the timer
	// fired.
	<-timer1.C
	fmt.Println("Timer 1 fired")
}

```

## 12.NewTicker

```
package main

import (
    "fmt"
    "time"
)

func main() {
    timer := time.NewTimer(10 * time.Second)

    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-timer.C:
            fmt.Println("Done")
            return
        case t := <-ticker.C:
            fmt.Println("ticker at: ", t)
        }
    }
}


```
