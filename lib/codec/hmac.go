package codec

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

// Hmac 返回指定密钥和正文的 sha256 类型的 hmac 哈希字节数组。
func Hmac(key []byte, body string) []byte {
	h := hmac.New(sha256.New, key)
	_, _ = io.WriteString(h, body)
	return h.Sum(nil)
}

// HmacBase64 返回指定密钥和正文经过 base64 编码后的 hmac 哈希值。
func HmacBase64(key []byte, body string) string {
	return base64.StdEncoding.EncodeToString(Hmac(key, body))
}
