package codec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"

	"git.zc0901.com/go/god/lib/logx"
)

var ErrPaddingSize = errors.New("padding size error")

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncryptor ecb

func NewECBEncryptor(b cipher.Block) cipher.BlockMode {
	return (*ecbEncryptor)(newECB(b))
}

func (x *ecbEncryptor) BlockSize() int { return x.blockSize }

// why we don't return error is because cipher.BlockMode doesn't allow this
func (x *ecbEncryptor) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		logx.Error("crypto/cipher: input not full blocks")
		return
	}
	if len(dst) < len(src) {
		logx.Error("crypto/cipher: output smaller than input")
		return
	}

	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecryptor ecb

func NewECBDecryptor(b cipher.Block) cipher.BlockMode {
	return (*ecbDecryptor)(newECB(b))
}

func (x *ecbDecryptor) BlockSize() int {
	return x.blockSize
}

// why we don't return error is because cipher.BlockMode doesn't allow this
func (x *ecbDecryptor) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		logx.Error("crypto/cipher: input not full blocks")
		return
	}
	if len(dst) < len(src) {
		logx.Error("crypto/cipher: output smaller than input")
		return
	}

	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

func EcbDecrypt(key, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		logx.Errorf("Decrypt key error: % x", key)
		return nil, err
	}

	decrypter := NewECBDecryptor(block)
	decrypted := make([]byte, len(src))
	decrypter.CryptBlocks(decrypted, src)

	return pkcs5Unpadding(decrypted, decrypter.BlockSize())
}

func EcbDecryptBase64(key, src string) (string, error) {
	keyBytes, err := getKeyBytes(key)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := EcbDecrypt(keyBytes, encryptedBytes)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(decryptedBytes), nil
}

func EcbEncrypt(key, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		logx.Errorf("Encrypt key error: % x", key)
		return nil, err
	}

	padded := pkcs5Padding(src, block.BlockSize())
	crypted := make([]byte, len(padded))
	encrypter := NewECBEncryptor(block)
	encrypter.CryptBlocks(crypted, padded)

	return crypted, nil
}

func EcbEncryptBase64(key, src string) (string, error) {
	keyBytes, err := getKeyBytes(key)
	if err != nil {
		return "", err
	}

	srcBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := EcbEncrypt(keyBytes, srcBytes)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func getKeyBytes(key string) ([]byte, error) {
	if len(key) > 32 {
		if keyBytes, err := base64.StdEncoding.DecodeString(key); err != nil {
			return nil, err
		} else {
			return keyBytes, nil
		}
	}

	return []byte(key), nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5Unpadding(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])
	if unpadding >= length || unpadding > blockSize {
		return nil, ErrPaddingSize
	}

	return src[:length-unpadding], nil
}
