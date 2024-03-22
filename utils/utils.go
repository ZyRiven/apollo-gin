package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
)

func GetIp() net.IP {
	var ip net.IP
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Panicln("无法获取网络接口:", err)
		return nil
	}
	// 遍历每个网络接口
	for _, iface := range interfaces {
		if iface.Name != "lo" {
			// 获取该接口的地址列表
			addrs, err := iface.Addrs()
			if err != nil {
				log.Panicln("无法获取接口地址:", err)
				continue
			}
			if len(addrs) > 0 {
				ipnet, ok := addrs[0].(*net.IPNet)
				if !ok {
					log.Panicln("无效的 IP 地址:", err)
					continue
				}
				ip = ipnet.IP
			}
		}
	}

	return ip
}

// ValidatorMessage 验证报文
func ValidatorMessage(data string) (newVar string) {
	startIndex := strings.Index(data, "FA")
	endIndex := strings.Index(data, "FB")

	if startIndex == -1 || endIndex == -1 || startIndex >= endIndex {
		return ""
	}

	newVar = data[startIndex : endIndex+2]
	if len(newVar) < 28 || newVar[:2] != "FA" || newVar[len(newVar)-2:] != "FB" {
		return ""
	}

	return newVar
}

// convertSlice 切片类型转换
func convertSlice(data []interface{}, targetSlice interface{}) (interface{}, error) {
	targetType := reflect.TypeOf(targetSlice)
	if targetType.Kind() != reflect.Slice {
		return nil, fmt.Errorf("目标类型必须为切片类型")
	}
	sliceType := targetType.Elem()
	result := reflect.MakeSlice(targetType, len(data), len(data))

	for i, item := range data {
		value := reflect.New(sliceType).Elem()
		if reflect.ValueOf(item).Type().AssignableTo(sliceType) {
			value.Set(reflect.ValueOf(item))
		} else {
			return nil, fmt.Errorf("无法将接口转换为目标切片类型")
		}
		result.Index(i).Set(value)
	}

	return result.Interface(), nil
}

func AesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}

func AesDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

// PKCS7UnPadding 去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// PKCS7Padding 补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// crcModbus 报文crc校验
func crcModbus(data []byte) uint16 {
	var crc uint16 = 0xFFFF
	const poly uint16 = 0xA001
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if (crc & 0x0001) != 0 {
				crc = (crc >> 1) ^ poly
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}
