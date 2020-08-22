net包为网络I/O提供一个便携式接口，包括TCP/IP, UDP, 域名解析，和Unix域套接字.
尽管该软件包提供了对低级网络的访问原语，大多数客户端将只需要提供的基本接口通过Dial，Listen和Accept函数以及相关的
Conn和Listener接口。crypto/tls包使用相同的接口和相似的Dial和Listen函数。
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
```
package main

import (
        "fmt"
        "log"
        "net"
)

func main() {
        conn, err := net.Dial("tcp", "127.0.0.1:8080")
        if err != nil {
                fmt.Println(err)
                return
        }
        defer conn.Close()
        log.Println("dial ok")
        conn.Write([]byte("hellodddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"))
}
```

Listen方法告知在本地网络地址，方法第一个参数必须是 "tcp", "tcp4", "tcp6", "unix" 或者 "unixpacket".
对于TCP网络，如果address参数中的主机为空或一个未指定的文字IP地址，Listen侦听所有可用的IP地址本地系统的单播和任意播IP地址。
要仅使用IPv4，请使用网络“ tcp4”。
该地址可以使用主机名，但是不建议使用，因为它将为主机的IP地址之一最多创建一个侦听器。
如果地址中的端口是空的，或者是0，就像"127.0.0.1:" or "[::1]:0"，端口号码就会自动选择，监听器的Addr方法可以用来发现选择的端口。


Dial连接到指定地址的网络上，预知的network参数的可选值为：
 "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only), "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4"(IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and "unixpacket".
对于TCP和UDP网络，地址的形式为"host:port".host必须是字段的IP地址，或者一个可以解析到IP地址的主机名。
port必须是一个字面量的端口号，或者一个服务的名称。如果host是一个IPv6的字面量地址，必须是如下形式的，例如："[2001:db8::1]:80" or "[fe80::1%zone]:80"
