package util

import (
	crand "crypto/rand"
	mrand "math/rand"
	"time"
)

// 从 4 字节里解析出整数
func buildMsgLen(raw []byte) uint32 {
	//	var length uint32
	//	binary.Read(bytes.NewBuffer(raw), binary.BigEndian, &length)
	//	return length

	return uint32(raw[0])<<24 |
		uint32(raw[1])<<16 |
		uint32(raw[2])<<8 |
		uint32(raw[3])
}

func genarateMsgLen(raw []byte, n uint32) {
	raw[0] = byte(n >> 24)
	raw[1] = byte(n >> 16)
	raw[2] = byte(n >> 8)
	raw[3] = byte(n)
}

func PKCS7Unpadding(buf []byte) []byte {
	n := len(buf)
	last := buf[n-1]
	n -= int(last)

	return buf[:n]
}

// n == 0, return ""
func Rand(n int) []byte {
	buf := make([]byte, n)

	var tmp int
	for n > 0 {
		tmp, _ = crand.Read(buf)
		if tmp == n {
			break
		}
	}

	return buf
}

// n == 0, return ""
func RandNumber(l int) string {
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 0, l)
	for i := 0; i < l; i++ {
		b = append(b, 48+byte(r.Intn(10)))
	}

	return string(b)
}
