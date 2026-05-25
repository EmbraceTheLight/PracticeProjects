package tools

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Hash(str string) string {
	// 创建一个新的 sha256 哈希对象
	hash := sha256.New()

	// 将字符串写入哈希对象
	hash.Write([]byte(str))

	// 从哈希对象中获取哈希值
	hashBytes := hash.Sum(nil)

	// 将字节切片切换为十六进制字符串
	return hex.EncodeToString(hashBytes)
}
