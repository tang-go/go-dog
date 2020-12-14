package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func AesDecrypt(echoStr string, EncodingAESKey string) string {
	wByte, err := base64.StdEncoding.DecodeString(echoStr)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	key, err := base64.StdEncoding.DecodeString(EncodingAESKey + "=")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	keyByte := []byte(key)
	x, err := _AesDecrypt(wByte, keyByte)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(x)
}

//AES解密
func _AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
