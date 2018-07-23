package util

// 从 4 字节里解析出整数
func buildMsgLen(raw []byte) uint32 {
	return uint32(raw[0])<<24 |
		uint32(raw[1])<<16 |
		uint32(raw[2])<<8 |
		uint32(raw[3])
}

func Pkcs7Unpadding(buf []byte) []byte {
	n := len(buf)
	last := buf[n-1]
	n -= int(last)

	return buf[:n]
}
