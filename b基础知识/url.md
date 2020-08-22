# 1.net/url简介
`import "net/url"`
url包解析URL并实现查询转义
- `URL`结构体
```
// Note that the Path field is stored in decoded form: /%47%6f%2f becomes /Go/.
// A consequence is that it is impossible to tell which slashes in the Path were
// slashes in the raw URL and which were %2f. This distinction is rarely important,
// but when it is, code must not use Path directly.
// The Parse function sets both Path and RawPath in the URL it returns,
// and URL's String method uses RawPath if it is a valid encoding of Path,
// by calling the EscapedPath method.
type URL struct {
	Scheme     string
	Opaque     string    // encoded opaque data
	User       *Userinfo // username and password information
	Host       string    // host or host:port
	Path       string    // path (relative paths may omit leading slash)
	RawPath    string    // encoded path hint (see EscapedPath method)
	ForceQuery bool      // append a query ('?') even if RawQuery is empty
	RawQuery   string    // encoded query values, without '?'
	Fragment   string    // fragment for references, without '#'
}
```


- `func Parse(rawurl string) (*URL, error)`
将原生的rawurl字符串解析成URL结构体
```
package main

import (
    "fmt"
    "log"
    "net/url"
)

func main() {
    u, err := url.Parse("http://www.baidu.com/search?q=dotnet")
    if err != nil {
        log.Fatal(err)
    }
    u.Scheme = "https"
    u.Host = "google.com"
    q := u.Query()
    q.Set("q", "golang")
    u.RawQuery = q.Encode()
    fmt.Println(u)
}
```
- Values 类型为map字典
通常用于查询参数，和表单值
```
// Values maps a string key to a list of values.
// It is typically used for query parameters and form values.
// Unlike in the http.Header map, the keys in a Values map
// are case-sensitive.
type Values map[string][]string
```

- `func (u *URL) Query() Values `返回的类型为Values
```
// Query parses RawQuery and returns the corresponding values.
// It silently discards malformed value pairs.
// To check errors use ParseQuery.
func (u *URL) Query() Values {
	v, _ := ParseQuery(u.RawQuery)
	return v
}
```
```
package main

import (
        "fmt"
        "log"
        "net/url"
)

func main() {
        u, err := url.Parse("https://example.org?q=golang&limit=10&offset=2")
        if err != nil {
                log.Fatal(err)
        }
        q := u.Query()
        fmt.Println(q["q"])   // [golang]
        fmt.Println(q.Get("limit"))  // 10
        fmt.Println(q.Get("offset"))  // 2 
}

```
- `func (v Values) Encode() string `
将参数列表转化为字符串拼接形式`q=golang&limit=10`
```
// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") sorted by key.
func (v Values) Encode() string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		keyEscaped := QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(QueryEscape(v))
		}
	}
	return buf.String()
}
```
- `func (u *URL) String() string`
String方法将URL重新组装成合法的string

一个完整的例子
```
func SendFromClientResp(method, uStr string, body io.Reader,
	headers, querys map[string]string) (io.Reader, error) {
	u, err := url.Parse(uStr)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	for k, v := range querys {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	var resp *http.Response
	for i := 0; i < reTryTimes; i++ {
		resp, err = http.DefaultClient.Do(req)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buff := bytes.NewBuffer(nil)
	io.Copy(buff, resp.Body)
	dr := decode2UTF8(httpRespCharset(resp), buff)
	buff.Reset()
	io.Copy(buff, dr)
	return buff, nil
}
```