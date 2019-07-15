# 1.time.Time转string之Format方法
```
package main

import (
    "fmt"
    "time"
    "reflect"
)

const LayOut = "20060102"


func main(){
    now := time.Now()
    fmt.Println(reflect.TypeOf(now)) //time.Time
    nowStr := now.Format(LayOut)
    fmt.Println(nowStr) //20190508
    fmt.Println(reflect.TypeOf(nowStr)) //string

}
```
注意：使用time Format方法时，最好使用包内的常量类型。否则可能出现时间的变动。比如
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
如果上述常量无法满足需求，则可以自己定义常量

# 2. 时区问题
```
var ShanghaiTimeZone *time.Location

func InitTimeZone() {
	if location, err := time.LoadLocation("Asia/Shanghai"); err != nil {
		logrus.Panicln("Failed to load timezone Asia/Shanghai, err: ", err)
	} else {
		ShanghaiTimeZone = location
	}
}
now := time.Now().In(common.ShanghaiTimeZone)
nowStr := now.Format(TimeLayOut)
fmt.Println(nowStr)
```
# 3.string 转time.Time
```
package main

import (
    "fmt"
    "time"
    "reflect"
)

const LayOut = "20060102"


func main(){
    now := time.Now()
    fmt.Println(reflect.TypeOf(now))
    nowStr := now.Format(LayOut)
    fmt.Println(nowStr)
    fmt.Println(reflect.TypeOf(nowStr))
    nowTime, err := time.Parse(LayOut, nowStr)
    if err != nil {
    }
    fmt.Println(nowTime) //2019-05-08 00:00:00 +0000 UTC
    fmt.Println(reflect.TypeOf(nowTime))//time.Time
}
```
# 4.time.Time 转 int64
```
    now := time.Now()
    nowInt := now.Unix()
    fmt.Println(nowInt) // 1557280736
    fmt.Println(reflect.TypeOf(nowInt)) //int64
```
# 5.int64转time.Time
```
    nowTimeFromInt := time.Unix(nowInt, 0)
    fmt.Println(nowTimeFromInt) //2019-05-08 10:01:13 +0800 CST
    fmt.Println(reflect.TypeOf(nowTimeFromInt)) //time.Time
```

# 6.time.Time的时间加减计算
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
    fmt.Println(time.Since(t3)) // 从t3已经过去多长时间
}
```
输出结果：
```
2018-01-11 05:13:32 +0800 CST
2019-05-08 14:20:51.928066 +0800 CST m=+0.000256314 2019-05-08 14:21:51.928066 +0800 CST m=+60.000256314
t3 is before t1
3h0m0.000190914s
```
# 7.time.After超时应用
```
package main

import (
    "fmt"
    "time"
)

func main() {
    c := make(chan bool)
    select {
    case v:= <-c:
        fmt.Println(v)
    case <-time.After(3 * time.Second):
        fmt.Println("Timeout")
    }
}
```
# 8.简单的定时器NewTicker
每隔一秒钟，执行一次业务，直到定时器任务时间到期为止。
```
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    done := make(chan bool)
    go func() {
        time.Sleep(10 * time.Second)
        done <- true
    }()
    for {
        select {
        case <-done:
            fmt.Println("Done!")
            return
        case t := <-ticker.C:
            fmt.Println("Current time: ", t)
        }
    }
```
Output:
```
Current time:  2018-02-05 10:41:37.204211269 +0800 CST m=+1.001535613
Current time:  2018-02-05 10:41:38.204068714 +0800 CST m=+2.001320058
Current time:  2018-02-05 10:41:39.204252518 +0800 CST m=+3.001503862
Current time:  2018-02-05 10:41:40.204103403 +0800 CST m=+4.001281747
Current time:  2018-02-05 10:41:41.204360898 +0800 CST m=+5.001539242
Current time:  2018-02-05 10:41:42.204120805 +0800 CST m=+6.001227149
Current time:  2018-02-05 10:41:43.204406187 +0800 CST m=+7.001439531
Current time:  2018-02-05 10:41:44.203708482 +0800 CST m=+8.000741826
Current time:  2018-02-05 10:41:45.204431933 +0800 CST m=+9.001392277
Current time:  2018-02-05 10:41:46.204367246 +0800 CST m=+10.001327590
Done!
```
```
package main

import (
    "fmt"
    "time"
)

func main(){
    timer := time.NewTimer(12*time.Second)
    ticker := time.NewTicker(time.Second)
    go func () {
        for t := range ticker.C {
            fmt.Println("Tick at ", t)
        }
    }()
    <- timer.C
    ticker.Stop()
    fmt.Println("timer expired.")
}
```
# 9. Duration类型
```
// A Duration represents the elapsed time between two instants
// as an int64 nanosecond count. The representation limits the
// largest representable duration to approximately 290 years.
type Duration int64
```
Duration代表两个实例节点之间，经过的时间，以int64纳秒计数， 该表示将最大可表示的持续时间限制为大约290年。
```
package main

import (
        "fmt"
        "time"
        "reflect"
)

func expensiveCall() {}

func main() {
        t0 := time.Now()
        time.Sleep(1*time.Second)
        t1 := time.Now()
        fmt.Printf("The call took %v to run.\n", t1.Sub(t0)) //The call took 1.002671497s to run.
        fmt.Println(reflect.TypeOf(t1.Sub(t0))) //time.Duration
}
```
# 10. ParseDuration方法
ParseDuration解析持续时间字符串。 持续时间字符串是可能签名的十进制数字序列，每个都带有可选的分数和单位后缀，例如“300ms”，“ -  1.5h”或“2h45m”。 有效时间单位是“ns”，“us”（或“μs”），“ms”，“s”，“m”，“h”。
```
package main

import (
	"fmt"
	"time"
)

func main() {
	hours, _ := time.ParseDuration("10h")
	complex, _ := time.ParseDuration("1h10m10s")

	fmt.Println(hours)  //10h0m0s
	fmt.Println(complex) //1h10m10s
	fmt.Printf("there are %.0f seconds in %v\n", complex.Seconds(), complex) //there are 4210 seconds in 1h10m10s
}
```

# 11. Hours方法
```
func (d Duration) Hours() float64 {
    hour := d / Hour
    nsec := d % Hour
    return float64(hour) + float64(nsec)/(60*60*1e9)
}
```
```
t4,_ := time.ParseDuration("4h30m")
fmt.Printf("%v\n",t4.Hours())
```
Output:
```
4.5
```
# 12.Minutes方法
```
func (d Duration) Minutes() float64 {
    min := d / Minute
    nsec := d % Minute
    return float64(min) + float64(nsec)/(60*1e9)
}
```
```
t4,_ := time.ParseDuration("4h30m")
fmt.Printf("%v\n",t4.Minutes())
```
Output:
```
270
```
# 13.Nanoseconds方法（纳秒）
```
func (d Duration) Nanoseconds() int64 {
        return int64(d) 
}
```
```
ns, _ := time.ParseDuration("1000ns")
fmt.Printf("one microsecond has %d nanoseconds.", ns.Nanoseconds())
```
Output:
```
one microsecond has 1000 nanoseconds.
```
