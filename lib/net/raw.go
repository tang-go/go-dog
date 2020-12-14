// +build ttt

package net

import (
	"errors"
	sc "syscall"
	"net"
	"fmt"
)

const UDP = sc.SOCK_DGRAM
const TCP = sc.SOCK_STREAM


type Addr struct {
	ip   string
	port int
	addr sc.SockaddrInet4

}

func Address(ip string, port int) Addr {
	a := new(Addr)
	a.ip = ip
	a.port = port
	copy(a.addr.Addr[0:4], net.ParseIP(ip))
	a.addr.Port = port
	return *a
}

// 申请ip4 socket
func Socket4(t int) (int, error) {
	var clientsock int
	var err error
	if clientsock, err = sc.Socket(sc.AF_INET, t, sc.IPPROTO_IP); err != nil {
		return -1, err
	}
	//sc.SetsockoptInt(clientsock, sc.SOL_SOCKET, sc.SO_REUSEPORT, 1)
	return clientsock, nil
}

// 连接
func Connet(socket int, daddr Addr) error {
	if err := sc.Connect(socket, &daddr.addr); err != nil {
		return err
	}
	return nil
}

//绑定端口
func Bind(socket int, port int) error {
	return sc.Bind(socket, &sc.SockaddrInet4{
		Port: port,
	})
}

//接收消息
func Recvfrom(socket int , data []byte) (n int, from Addr, err error)  {
	var sa sc.Sockaddr
	n , sa , err = sc.Recvfrom(socket, data, 0)
	//copy(from.addr.Addr[0:4], f.Addr.Data[0:4])
	switch sa := sa.(type) {
	case *sc.SockaddrInet4:
		from = Addr{ip: fmt.Sprintf("%d.%d.%d.%d", sa.Addr[0], sa.Addr[1], sa.Addr[2], sa.Addr[3]), port: sa.Port}
	case *sc.SockaddrInet6:
		err =  errors.New("不支持ipv6")
	}
	return n, from, err

}

func Sendto(fd int, data []byte, to Addr) error {
	fmt.Println("-----", net.ParseIP(to.ip).To16())
	addr := sc.SockaddrInet4{
		Port: to.port,
		Addr: [4]byte{net.ParseIP(to.ip).To4()[0],
			net.ParseIP(to.ip).To4()[1],
			net.ParseIP(to.ip).To4()[2],
			net.ParseIP(to.ip).To4()[3]},
	}
	return sc.Sendto(fd,  data,0, &addr)
}
