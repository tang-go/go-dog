package io

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	gHeadLen = 4
)

//Write tcp黏包，写包
func Write(conn net.Conn, message []byte) (int, error) {
	buf := make([]byte, gHeadLen+len(message))
	binary.LittleEndian.PutUint32(buf[0:gHeadLen], uint32(len(message)))
	copy(buf[gHeadLen:], message)
	n, err := conn.Write(buf)
	if err != nil {
		return 0, err
	} else if n != len(buf) {
		return 0, fmt.Errorf("write %d less than %d", n, len(buf))
	}
	return n, err
}

//Read tcp黏包，读取包
func Read(conn net.Conn) (int, []byte, error) {
	l := make([]byte, gHeadLen)
	_, err := io.ReadFull(conn, l)
	if err != nil {
		return 0, nil, err
	}
	length := binary.LittleEndian.Uint32(l)
	data := make([]byte, length)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return 0, nil, err
	}
	return 0, data, nil
}

//WriteByTime tcp黏包，写包
func WriteByTime(conn net.Conn, message []byte, t time.Time) (int, error) {
	err := conn.SetWriteDeadline(t)
	if err != nil {
		return 0, err
	}
	defer conn.SetWriteDeadline(time.Time{})
	return Write(conn, message)
}

//ReadByTime tcp黏包，读取包  超时
func ReadByTime(conn net.Conn, t time.Time) (int, []byte, error) {
	err := conn.SetReadDeadline(t)
	if err != nil {
		return 0, nil, err
	}
	defer conn.SetReadDeadline(time.Time{})
	return Read(conn)
}

// BytesCombine 生成bytes数组
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

//IntToBytes 整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//BytesToInt 字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
