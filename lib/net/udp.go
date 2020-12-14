package net

import (
	"fmt"
	"net"
)

func UdpSend(ip string, sport, dport int, data []byte) error {

	localip := net.ParseIP("127.0.0.1")
	remoteip := net.ParseIP(ip)
	lAddr := &net.UDPAddr{IP: localip, Port: sport}
	rAddr := &net.UDPAddr{IP: remoteip, Port: dport}
	conn, err := net.DialUDP("udp", lAddr, rAddr)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		return err
	}
	datalen := len(data)
	sndlen := 0

Send:
	if datalen > sndlen {
		if n, e := conn.Write(data); e != nil {
			return err
		} else {
			sndlen += n
			goto Send
		}
	}
	return nil
}

type UdpSvr struct {
	local  *net.UDPAddr
	conn   *net.UDPConn
	onRead func(data []byte, conn *net.UDPConn, remote *net.UDPAddr) error
	exit   bool
}

func (us *UdpSvr) Conn() *net.UDPConn { return us.conn }

func NewUdpSvr(ip string, port uint16, onRead func(data []byte, conn *net.UDPConn, remote *net.UDPAddr) error) (*UdpSvr, error) {
	svr := new(UdpSvr)
	svr.onRead = onRead
	var err error
	if svr.local, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port)); err != nil {
		return nil, err
	}
	svr.conn, err = net.ListenUDP("udp", svr.local)
	if err != nil {
		fmt.Println(err)
	}
	svr.exit = false

	go svr.recv()
	return svr, nil
}

func (this *UdpSvr) recv() {
	for {
		if this.exit {
			break
		}
		buf := make([]byte, 65535)
		n, raddr, err := this.conn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println("from ReadFromUDP:", err)
		} else {
			if this.onRead != nil {
				this.onRead(buf[0:n], this.conn, raddr)
			}
		}

	}

}

func (this *UdpSvr) Exit() {
	this.exit = true
}
