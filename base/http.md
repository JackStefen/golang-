先看一下服务器是如何跑起来的，再慢慢分析
```
package main

import (
        "fmt"
        "net/http"
)

func main() {
        http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintf(w, "Hello, world...")
        }))
        http.ListenAndServe(":8080", nil)
}
```
看一下`ListenAndServe`到底是如何工作的. `ListenAndServe`侦听TCP网络地址addr，然后调用Serve with handler处理传入连接上的请求。接受的连接都配置为开启 TCP keep-alives.该处理程序通常为nil，在这种情况下，将使用`DefaultServeMux`。
```
// ListenAndServe always returns a non-nil error.
func ListenAndServe(addr string, handler Handler) error {
    server := &Server{Addr: addr, Handler: handler}
    return server.ListenAndServe()
}
```
关键还是`ListenAndServe`,而它的实质还是创建一个tcp Listener，然后`srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})`
```
func (srv *Server) ListenAndServe() error {
    if srv.shuttingDown() {
        return ErrServerClosed
    }
    addr := srv.Addr
    if addr == "" {
        addr = ":http"
    }
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}
```
首先看一下`tcpKeepAliveListener`,

```
type tcpKeepAliveListener struct {
    *net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
    tc, err := ln.AcceptTCP()
    if err != nil {
        return nil, err
    }
    tc.SetKeepAlive(true)
    tc.SetKeepAlivePeriod(3 * time.Minute)
    return tc, nil
}
```
其结构体实现了`Accept`方法。该方法设置了开启Tcp keep-alive属性。

再看一下Server的Serve()方法。其接受参数监听器上的连接，然后调用新的服务协程来进行处理请求。服务协程读取请求，然后调用`srv.Handler`返回给客户端。
```
func (srv *Server) Serve(l net.Listener) error {
    // 该方法是对其测试的钩子
    if fn := testHookServerServe; fn != nil {
        fn(srv, l) // call hook with unwrapped listener
    }

    l = &onceCloseListener{Listener: l}
    defer l.Close()

    if err := srv.setupHTTP2_Serve(); err != nil {
        return err
    }

    if !srv.trackListener(&l, true) {
        return ErrServerClosed
    }
    defer srv.trackListener(&l, false)

    var tempDelay time.Duration     // how long to sleep on accept failure
    baseCtx := context.Background() // base is always background, per Issue 16220
    // ServerContextKey = &contextKey{"http-server"},该变量是一个全局的上下文的key,
    // 可以在带有context.WithValue的HTTP处理程序中使用它来访问启动处理程序的服务器
    ctx := context.WithValue(baseCtx, ServerContextKey, srv)
    for {
        rw, e := l.Accept()
        if e != nil {
            select {
            case <-srv.getDoneChan():
                return ErrServerClosed
            default:
            }
            if ne, ok := e.(net.Error); ok && ne.Temporary() {
                if tempDelay == 0 {
                    tempDelay = 5 * time.Millisecond
                } else {
                    tempDelay *= 2
                }
                if max := 1 * time.Second; tempDelay > max {
                    tempDelay = max
                }
                srv.logf("http: Accept error: %v; retrying in %v", e, tempDelay)
                time.Sleep(tempDelay)
                continue
            }
            return e
        }
        tempDelay = 0
        c := srv.newConn(rw)
        c.setState(c.rwc, StateNew) // before Serve can return
        go c.serve(ctx)
    }
}
```
`onceCloseListener`包裹传进来的`Listener`.使用`sync.Once`原语对其进行多放`Close`调用。该原语支持仅一次方法调用。事实上，`sync.Once`是一个只执行一个操作的对象。多次调用其上的Do()方法，仅首次调用起效。
看一下Do方法的原型`func (o *Once) Do(f func())`
因为对Do的调用只有在对f的调用返回时才会返回，如果f导致Do被调用，那么它就会死锁。
如果f函数异常，Do将任务其已经结束返回，后续对Do的调用将不再调用f。
```
// onceCloseListener wraps a net.Listener, protecting it from
// multiple Close calls.
type onceCloseListener struct {
    net.Listener
    once     sync.Once
    closeErr error
}

func (oc *onceCloseListener) Close() error {
    oc.once.Do(oc.close)
    return oc.closeErr
}

func (oc *onceCloseListener) close() { oc.closeErr = oc.Listener.Close() }
```
通过查看Serve()方法，可以看出，其实质上还是一个加强版的服务器版本，剥离之后，和如下的服务器是不是大同小异啊
```
package main

import (
        "fmt"
        "io"
        "log"
        "net"
)

func main() {
        ln, err := net.Listen("tcp", ":8080")
        if err != nil {
                fmt.Println(err)
                return
        }
        conn, err := ln.Accept()
        if err != nil {
                fmt.Println(err)
                return
        }
        var buf = make([]byte, 10)
        var result string
        for {
                n, err := conn.Read(buf)
                if err == io.EOF {
                        fmt.Println(err.Error())
                        break
                }
                if err != nil {
                        fmt.Println(err)
                        return
                }
                log.Printf("client content is %s\n", string(buf[:n]))
                result += string(buf[:n])
        }
        log.Println(result)
}
```
具体是如何加强的，可以看看`c := srv.newConn(rw)`。其在原始的`net.Conn`上穿了一件衣服.
```
// Create new connection from rwc.
func (srv *Server) newConn(rwc net.Conn) *conn {
    c := &conn{
        server: srv,
        rwc:    rwc,
    }
    if debugServerConnections {
        c.rwc = newLoggingConn("server", c.rwc)
    }
    return c
}
```
穿完了衣服，干了一件事儿`c.setState(c.rwc, StateNew) // before Serve can return`
```
func (c *conn) setState(nc net.Conn, state ConnState) {
    srv := c.server
    switch state {
    case StateNew:
        srv.trackConn(c, true)
    case StateHijacked, StateClosed:
        srv.trackConn(c, false)
    }
    if state > 0xff || state < 0 {
        panic("internal error")
    }
    packedState := uint64(time.Now().Unix()<<8) | uint64(state)
    atomic.StoreUint64(&c.curState.atomic, packedState)
    // 当客户端连接发生状态变化时，进行一个回调函数的处理
    if hook := srv.ConnState; hook != nil {
        hook(nc, state)
    }
}
```
咱们传入的是StateNew,看看他干了啥
```
func (s *Server) trackConn(c *conn, add bool) {
    s.mu.Lock()
    defer s.mu.Unlock()
    if s.activeConn == nil {
        s.activeConn = make(map[*conn]struct{})
    }
    if add {
        s.activeConn[c] = struct{}{}
    } else {
        delete(s.activeConn, c)
    }
}
```
就是把当前新的连接加入到服务器的活跃的连接中。
然后起了一个goroutine，接管了请求处理。
```
// Serve a new connection.
func (c *conn) serve(ctx context.Context) {
    c.remoteAddr = c.rwc.RemoteAddr().String()
    ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
    defer func() {
        if err := recover(); err != nil && err != ErrAbortHandler {
            const size = 64 << 10
            buf := make([]byte, size)
            buf = buf[:runtime.Stack(buf, false)]
            c.server.logf("http: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
        }
        if !c.hijacked() {
            c.close()
            c.setState(c.rwc, StateClosed)
        }
    }()

    if tlsConn, ok := c.rwc.(*tls.Conn); ok {
        if d := c.server.ReadTimeout; d != 0 {
            c.rwc.SetReadDeadline(time.Now().Add(d))
        }
        if d := c.server.WriteTimeout; d != 0 {
            c.rwc.SetWriteDeadline(time.Now().Add(d))
        }
        if err := tlsConn.Handshake(); err != nil {
            // If the handshake failed due to the client not speaking
            // TLS, assume they're speaking plaintext HTTP and write a
            // 400 response on the TLS conn's underlying net.Conn.
            if re, ok := err.(tls.RecordHeaderError); ok && re.Conn != nil && tlsRecordHeaderLooksLikeHTTP(re.RecordHeader) {
                io.WriteString(re.Conn, "HTTP/1.0 400 Bad Request\r\n\r\nClient sent an HTTP request to an HTTPS server.\n")
                re.Conn.Close()
                return
            }
            c.server.logf("http: TLS handshake error from %s: %v", c.rwc.RemoteAddr(), err)
            return
        }
        c.tlsState = new(tls.ConnectionState)
        *c.tlsState = tlsConn.ConnectionState()
        if proto := c.tlsState.NegotiatedProtocol; validNPN(proto) {
            if fn := c.server.TLSNextProto[proto]; fn != nil {
                h := initNPNRequest{tlsConn, serverHandler{c.server}}
                fn(c.server, tlsConn, h)
            }
            return
        }
    }

    // HTTP/1.x from here on.

    ctx, cancelCtx := context.WithCancel(ctx)
    c.cancelCtx = cancelCtx
    defer cancelCtx()

    c.r = &connReader{conn: c}
    c.bufr = newBufioReader(c.r)
    c.bufw = newBufioWriterSize(checkConnErrorWriter{c}, 4<<10)

    for {
        w, err := c.readRequest(ctx)
        if c.r.remain != c.server.initialReadLimitSize() {
            // If we read any bytes off the wire, we're active.
            c.setState(c.rwc, StateActive)
        }
        if err != nil {
            const errorHeaders = "\r\nContent-Type: text/plain; charset=utf-8\r\nConnection: close\r\n\r\n"

            if err == errTooLarge {
                // Their HTTP client may or may not be
                // able to read this if we're
                // responding to them and hanging up
                // while they're still writing their
                // request. Undefined behavior.
                const publicErr = "431 Request Header Fields Too Large"
                fmt.Fprintf(c.rwc, "HTTP/1.1 "+publicErr+errorHeaders+publicErr)
                c.closeWriteAndWait()
                return
            }
            if isCommonNetReadError(err) {
                return // don't reply
            }

            publicErr := "400 Bad Request"
            if v, ok := err.(badRequestError); ok {
                publicErr = publicErr + ": " + string(v)
            }

            fmt.Fprintf(c.rwc, "HTTP/1.1 "+publicErr+errorHeaders+publicErr)
            return
        }

        // Expect 100 Continue support
        req := w.req
        if req.expectsContinue() {
            if req.ProtoAtLeast(1, 1) && req.ContentLength != 0 {
                // Wrap the Body reader with one that replies on the connection
                req.Body = &expectContinueReader{readCloser: req.Body, resp: w}
            }
        } else if req.Header.get("Expect") != "" {
            w.sendExpectationFailed()
            return
        }

        c.curReq.Store(w)

        if requestBodyRemains(req.Body) {
            registerOnHitEOF(req.Body, w.conn.r.startBackgroundRead)
        } else {
            w.conn.r.startBackgroundRead()
        }

        // HTTP cannot have multiple simultaneous active requests.[*]
        // Until the server replies to this request, it can't read another,
        // so we might as well run the handler in this goroutine.
        // [*] Not strictly true: HTTP pipelining. We could let them all process
        // in parallel even if their responses need to be serialized.
        // But we're not going to implement HTTP pipelining because it
        // was never deployed in the wild and the answer is HTTP/2.
        // HTTP不能同时有多个活动请求。
        // 直到服务器回复该请求，它才能读取另一个请求
        // 因此我们最好在此goroutine中运行处理程序。
        serverHandler{c.server}.ServeHTTP(w, w.req)
        w.cancelCtx()
        if c.hijacked() {
            return
        }
        w.finishRequest()
        if !w.shouldReuseConnection() {
            if w.requestBodyLimitHit || w.closedRequestBodyEarly() {
                c.closeWriteAndWait()
            }
            return
        }
        c.setState(c.rwc, StateIdle)
        c.curReq.Store((*response)(nil))

        if !w.conn.server.doKeepAlives() {
            // We're in shutdown mode. We might've replied
            // to the user without "Connection: close" and
            // they might think they can send another
            // request, but such is life with HTTP/1.1.
            return
        }

        if d := c.server.idleTimeout(); d != 0 {
            c.rwc.SetReadDeadline(time.Now().Add(d))
            if _, err := c.bufr.Peek(4); err != nil {
                return
            }
        }
        c.rwc.SetReadDeadline(time.Time{})
    }
}
```
先看看`context.WithValue`,`WithValue`返回`parent`的副本，其中与`key`关联的值为`val`。仅将上下文值用于传递过程和API的请求范围的数据，而不用于将可选参数传递给函数
```
func WithValue(parent Context, key, val interface{}) Context {
    if key == nil {
        panic("nil key")
    }
    if !reflect.TypeOf(key).Comparable() {
        panic("key is not comparable")
    }
    return &valueCtx{parent, key, val}
}
```
在每个goroutine中，如下我们处理请求：
```
// serverHandler delegates to either the server's Handler or
// DefaultServeMux and also handles "OPTIONS *" requests.
type serverHandler struct {
    srv *Server
}

func (sh serverHandler) ServeHTTP(rw ResponseWriter, req *Request) {
    handler := sh.srv.Handler
    if handler == nil {
        handler = DefaultServeMux
    }
    if req.RequestURI == "*" && req.Method == "OPTIONS" {
        handler = globalOptionsHandler{}
    }
    handler.ServeHTTP(rw, req)
}
```
首先其会从Server中检索`Handler`，如果handler为空，默认的会使用`DefaultServeMux`,如果请求的URI是`*`并且请求方式是`OPTIONS`,会使用`globalOptionsHandler`.
然后进行`handler`处理,以上是对整个流程的一个详细介绍，而对于像如下的服务端，也就大差不差了
```
package main

import (
        "fmt"
        "net/http"
        _ "net/http/pprof"
)

type UserRespone struct{}

func (*UserRespone) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome...")
}

func main() {
        mux := http.NewServeMux()
        mux.Handle("/", &UserRespone{})
        http.ListenAndServe(":8080", mux)
}
```
mux就是一个实现了`Handler`接口的结构体，自然可以在ListenAndServe中使用。
```
type ServeMux struct {
    mu    sync.RWMutex
    m     map[string]muxEntry
    es    []muxEntry // slice of entries sorted from longest to shortest.
    hosts bool       // whether any patterns contain hostnames
}
// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
    if r.RequestURI == "*" {
        if r.ProtoAtLeast(1, 1) {
            w.Header().Set("Connection", "close")
        }
        w.WriteHeader(StatusBadRequest)
        return
    }
    h, _ := mux.Handler(r)
    h.ServeHTTP(w, r)
}
```
