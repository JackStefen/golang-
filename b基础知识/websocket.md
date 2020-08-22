注：本文以`/gorilla/websocket`为基础进行说明
Server
```
package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))

```
Client

```
package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

```
# websocket.Upgrader{}
Upgrader指定升级http连接到websocket连接所需要的参数.

```
type Upgrader struct {
	// HandshakeTimeout指定握手完成的持续时间。
	HandshakeTimeout time.Duration

	// ReadBufferSize和WriteBufferSize指定I/O缓冲区大小（以字节为单位）。
    // 如果缓冲区大小为零，则使用HTTP服务器分配的缓冲区。
    // I/O缓冲区的大小不限制可以发送或接收的消息的大小。
	ReadBufferSize, WriteBufferSize int

	// WriteBufferPool是用于写操作的缓冲区池。
	// 如果未设置该值，则在连接的生存期内将写缓冲区分配给该连接。

	// 当应用程序在大量连接上的写入量适中时，池是最有用的。
	//
	// 应用程序应为WriteBufferSize的每个唯一值使用一个池。
	WriteBufferPool BufferPool

	// 子协议按优先顺序指定服务器支持的协议。 
    // 如果此字段不为nil，则Upgrade方法通过选择此列表中与客户端请求的协议的第一个匹配项来协商子协议。
    // 如果不匹配，则不会协商任何协议（握手响应中不包含Sec-Websocket-Protocol标头）。
	Subprotocols []string

	// 错误指定用于生成HTTP错误响应的函数。
    // 如果Error为nil，则使用http.Error生成HTTP响应。
	Error func(w http.ResponseWriter, r *http.Request, status int, reason error)

	// 如果请求Origin头是可接受的，则CheckOrigin返回true。
    // 如果CheckOrigin为nil，则使用安全的默认值：如果存在Origin请求标头且源主机与请求Host标头不相等，则返回false。
	CheckOrigin func(r *http.Request) bool

	// EnableCompression指定服务器是否应尝试对每个消息压缩进行协商（RFC 7692）。 
    // 将此值设置为true不能保证将支持压缩。 当前仅支持“无上下文接管”模式。
	EnableCompression bool
}
```
Upgrader的Upgrade方法是如何升级将HTTP服务器连接升级到WebSocket协议。
对客户端升级请求的响应中包含responseHeader。 
使用responseHeader指定cookie（Set-Cookie）和应用程序协商的子协议（Sec-WebSocket-Protocol）。
如果Upgrade失败，则Upgrade会通过HTTP错误响应来回复客户端。
```
func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*Conn, error) {
	const badHandshake = "websocket: the client is not using the websocket protocol: "

	if !tokenListContainsValue(r.Header, "Connection", "upgrade") {
		return u.returnError(w, r, http.StatusBadRequest, badHandshake+"'upgrade' token not found in 'Connection' header")
	}

	if !tokenListContainsValue(r.Header, "Upgrade", "websocket") {
		return u.returnError(w, r, http.StatusBadRequest, badHandshake+"'websocket' token not found in 'Upgrade' header")
	}

	if r.Method != "GET" {
		return u.returnError(w, r, http.StatusMethodNotAllowed, badHandshake+"request method is not GET")
	}

	if !tokenListContainsValue(r.Header, "Sec-Websocket-Version", "13") {
		return u.returnError(w, r, http.StatusBadRequest, "websocket: unsupported version: 13 not found in 'Sec-Websocket-Version' header")
	}

	if _, ok := responseHeader["Sec-Websocket-Extensions"]; ok {
		return u.returnError(w, r, http.StatusInternalServerError, "websocket: application specific 'Sec-WebSocket-Extensions' headers are unsupported")
	}

	checkOrigin := u.CheckOrigin
	if checkOrigin == nil {
		checkOrigin = checkSameOrigin
	}
	if !checkOrigin(r) {
		return u.returnError(w, r, http.StatusForbidden, "websocket: request origin not allowed by Upgrader.CheckOrigin")
	}

	challengeKey := r.Header.Get("Sec-Websocket-Key")
	if challengeKey == "" {
		return u.returnError(w, r, http.StatusBadRequest, "websocket: not a websocket handshake: `Sec-WebSocket-Key' header is missing or blank")
	}

	subprotocol := u.selectSubprotocol(r, responseHeader)

	// Negotiate PMCE
	var compress bool
	if u.EnableCompression {
		for _, ext := range parseExtensions(r.Header) {
			if ext[""] != "permessage-deflate" {
				continue
			}
			compress = true
			break
		}
	}

    // 这里是重点
	h, ok := w.(http.Hijacker)
	if !ok {
		return u.returnError(w, r, http.StatusInternalServerError, "websocket: response does not implement http.Hijacker")
	}
	var brw *bufio.ReadWriter
	netConn, brw, err := h.Hijack()
	if err != nil {
		return u.returnError(w, r, http.StatusInternalServerError, err.Error())
	}

	if brw.Reader.Buffered() > 0 {
		netConn.Close()
		return nil, errors.New("websocket: client sent data before handshake is complete")
	}

	var br *bufio.Reader
	if u.ReadBufferSize == 0 && bufioReaderSize(netConn, brw.Reader) > 256 {
		// Reuse hijacked buffered reader as connection reader.
		br = brw.Reader
	}

	buf := bufioWriterBuffer(netConn, brw.Writer)

	var writeBuf []byte
	if u.WriteBufferPool == nil && u.WriteBufferSize == 0 && len(buf) >= maxFrameHeaderSize+256 {
		// Reuse hijacked write buffer as connection buffer.
		writeBuf = buf
	}

	c := newConn(netConn, true, u.ReadBufferSize, u.WriteBufferSize, u.WriteBufferPool, br, writeBuf)
	c.subprotocol = subprotocol

	if compress {
		c.newCompressionWriter = compressNoContextTakeover
		c.newDecompressionReader = decompressNoContextTakeover
	}

	// Use larger of hijacked buffer and connection write buffer for header.
	p := buf
	if len(c.writeBuf) > len(p) {
		p = c.writeBuf
	}
	p = p[:0]

	p = append(p, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
	p = append(p, computeAcceptKey(challengeKey)...)
	p = append(p, "\r\n"...)
	if c.subprotocol != "" {
		p = append(p, "Sec-WebSocket-Protocol: "...)
		p = append(p, c.subprotocol...)
		p = append(p, "\r\n"...)
	}
	if compress {
		p = append(p, "Sec-WebSocket-Extensions: permessage-deflate; server_no_context_takeover; client_no_context_takeover\r\n"...)
	}
	for k, vs := range responseHeader {
		if k == "Sec-Websocket-Protocol" {
			continue
		}
		for _, v := range vs {
			p = append(p, k...)
			p = append(p, ": "...)
			for i := 0; i < len(v); i++ {
				b := v[i]
				if b <= 31 {
					// prevent response splitting.
					b = ' '
				}
				p = append(p, b)
			}
			p = append(p, "\r\n"...)
		}
	}
	p = append(p, "\r\n"...)

	// Clear deadlines set by HTTP server.
	netConn.SetDeadline(time.Time{})

	if u.HandshakeTimeout > 0 {
		netConn.SetWriteDeadline(time.Now().Add(u.HandshakeTimeout))
	}
	if _, err = netConn.Write(p); err != nil {
		netConn.Close()
		return nil, err
	}
	if u.HandshakeTimeout > 0 {
		netConn.SetWriteDeadline(time.Time{})
	}

	return c, nil
}

```
**该方法一大段都是用来进行错误检测，以及子协议选择和读写缓存，重点是劫持器。
请记住，您无法使用http.ResponseWriter编写响应，因为一旦您开始发送响应，它将关闭基础TCP连接。
因此，您需要使用HTTP劫持。 通过劫持，您可以接管基础的TCP连接处理程序和bufio.Writer。 
这使您可以在不关闭TCP连接的情况下读取和写入数据。**
最终
```
c := newConn(netConn, true, u.ReadBufferSize, u.WriteBufferSize, u.WriteBufferPool, br, writeBuf)
```

实现了连接之后，就是对连接的正常读写了，ReadMessage，WriteMessage。
WebSocket通信协议通过单个TCP连接提供全双工通信通道。 与HTTP相比，WebSocket不需要您发送请求即可获得响应。
 它们允许双向数据流，因此您只需等待服务器响应即可。 可用时，它将向您发送一条消息。


 对于需要连续数据交换的服务（例如即时通讯程序，在线游戏和实时交易系统），WebSockets是一个很好的解决方案。

 WebSocket连接由浏览器请求，并由服务器响应，然后建立连接。 此过程通常称为握手。 
 WebSockets中的特殊标头仅需要浏览器与服务器之间的一次握手即可建立连接，该连接将在其整个生命周期内保持活动状态。

 WebSockets解决了许多实时Web开发的难题，与传统的HTTP相比有很多好处：
 - 轻量级报头减少了数据传输开销。
 - 单个Web客户端仅需要一个TCP连接。
 - WebSocket服务器可以将数据推送到Web客户端。

 WebSocket协议实现起来相对简单。 它使用HTTP协议进行初始握手。
  握手成功后，将建立连接，并且WebSocket实质上使用原始TCP读取/写入数据。

下面是一个客户端请求
```
GET ws://localhost:8080/echo HTTP/1.1
Host: localhost:8080
Connection: Upgrade
Pragma: no-cache
Cache-Control: no-cache
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.122 Safari/537.36
Upgrade: websocket
Origin: http://localhost:8080
Sec-WebSocket-Version: 13
Accept-Encoding: gzip, deflate, br
Accept-Language: zh-CN,zh;q=0.9
Sec-WebSocket-Key: 4uEeKBqa1bL12W9Iazlu7w==
Sec-WebSocket-Extensions: permessage-deflate; client_max_window_bits
```
下面是一个服务端响应
```
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: XwZURdK+N7l+tVt7KVE5Pr4A+Ls=
```

“ Sec-WebSocket-key”是随机生成的，并且是Base64编码的。 接受请求后，服务器需要将此密钥附加到固定字符串。 假设您有x3JJHMbDL1EzLkh9GBhXDw ==键。 在这种情况下，可以使用SHA-1计算二进制值，并使用Base64对其进行编码。 您将获得HSmrc0sMlYUkAGmm5OPpG2HaGWk =。 使用它作为Sec-WebSocket-Accept响应标头的值。

握手成功完成后，您的应用程序可以从客户端读取数据或向客户端写入数据。 WebSocket规范定义了客户端和服务器之间使用的特定帧格式。 这是框架的位模式：

![模式](https://yalantis.com/uploads/ckeditor/pictures/4157/bit-pattern.png)