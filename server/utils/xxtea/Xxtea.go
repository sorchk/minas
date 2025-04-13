package xxtea

import (
	"encoding/base64"
	"strings"
)

// delta XXTEA算法中使用的常量值
// 这是一个魔数，来源于黄金分割比例((sqrt(5) - 1) / 2) * 2^32
const delta = 0x9E3779B9

// toBytes 将uint32数组转换为字节数组
// 参数 v: 需要转换的uint32数组
// 参数 includeLength: 是否包含长度信息
// 返回值: 转换后的字节数组
func toBytes(v []uint32, includeLength bool) []byte {
	length := uint32(len(v))
	n := length << 2 // 计算字节长度 (每个uint32占4字节)
	if includeLength {
		// 如果包含长度信息，使用最后一个uint32作为长度值
		m := v[length-1]
		n -= 4
		// 验证长度合法性
		if (m < n-3) || (m > n) {
			return nil
		}
		n = m
	}
	bytes := make([]byte, n)
	// 将uint32数组转换为字节数组
	for i := uint32(0); i < n; i++ {
		bytes[i] = byte(v[i>>2] >> ((i & 3) << 3))
	}
	return bytes
}

// toUint32s 将字节数组转换为uint32数组
// 参数 bytes: 需要转换的字节数组
// 参数 includeLength: 是否在结果中包含长度信息
// 返回值: 转换后的uint32数组
func toUint32s(bytes []byte, includeLength bool) (v []uint32) {
	length := uint32(len(bytes))
	n := length >> 2 // 计算需要多少个uint32 (每4个字节一个uint32)
	// 如果字节长度不是4的倍数，需要额外一个uint32存储剩余字节
	if length&3 != 0 {
		n++
	}
	if includeLength {
		// 如果包含长度，额外分配一个uint32用于存储长度
		v = make([]uint32, n+1)
		v[n] = length
	} else {
		v = make([]uint32, n)
	}
	// 将字节数组转换为uint32数组
	for i := uint32(0); i < length; i++ {
		v[i>>2] |= uint32(bytes[i]) << ((i & 3) << 3)
	}
	return v
}

// mx XXTEA算法的核心混淆函数
// 用于混淆数据，增加加密强度
// 参数说明略 (内部算法实现细节)
func mx(sum uint32, y uint32, z uint32, p uint32, e uint32, k []uint32) uint32 {
	return ((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((sum ^ y) + (k[p&3^e] ^ z))
}

// fixk 确保密钥长度至少为4个uint32
// 如果密钥长度不足，通过复制已有密钥补全
// 参数 k: 原始密钥
// 返回值: 调整后的密钥
func fixk(k []uint32) []uint32 {
	if len(k) < 4 {
		key := make([]uint32, 4)
		copy(key, k)
		return key
	}
	return k
}

// encrypt 使用XXTEA算法加密数据
// 参数 v: 要加密的uint32数组
// 参数 k: 加密密钥
// 返回值: 加密后的uint32数组
func encrypt(v []uint32, k []uint32) []uint32 {
	length := uint32(len(v))
	n := length - 1
	k = fixk(k)
	var y, z, sum, e, p, q uint32
	z = v[n]
	sum = 0
	// 加密轮数由数据长度决定
	for q = 6 + 52/length; q > 0; q-- {
		sum += delta
		e = sum >> 2 & 3
		for p = 0; p < n; p++ {
			y = v[p+1]
			v[p] += mx(sum, y, z, p, e, k)
			z = v[p]
		}
		y = v[0]
		v[n] += mx(sum, y, z, p, e, k)
		z = v[n]
	}
	return v
}

// decrypt 使用XXTEA算法解密数据
// 参数 v: 要解密的uint32数组
// 参数 k: 解密密钥
// 返回值: 解密后的uint32数组
func decrypt(v []uint32, k []uint32) []uint32 {
	length := uint32(len(v))
	n := length - 1
	k = fixk(k)
	var y, z, sum, e, p, q uint32
	y = v[0]
	q = 6 + 52/length
	// 从最终的sum值开始，反向解密
	for sum = q * delta; sum != 0; sum -= delta {
		e = sum >> 2 & 3
		for p = n; p > 0; p-- {
			z = v[p-1]
			v[p] -= mx(sum, y, z, p, e, k)
			y = v[p]
		}
		z = v[n]
		v[0] -= mx(sum, y, z, p, e, k)
		y = v[0]
	}
	return v
}

// Encrypt 加密字节数组
// 参数 data: 要加密的字节数据
// 参数 key: 加密密钥
// 返回值: 加密后的字节数组
func Encrypt(data []byte, key []byte) []byte {
	if data == nil || len(data) == 0 {
		return data
	}
	return toBytes(encrypt(toUint32s(data, true), toUint32s(key, false)), false)
}

// Decrypt 解密字节数组
// 参数 data: 要解密的字节数据
// 参数 key: 解密密钥
// 返回值: 解密后的字节数组
func Decrypt(data []byte, key []byte) []byte {
	if data == nil || len(data) == 0 {
		return data
	}
	return toBytes(decrypt(toUint32s(data, false), toUint32s(key, false)), true)
}

// EncryptString 加密字符串
// 对字符串进行XXTEA加密，并进行Base64编码
// 参数 str: 要加密的字符串
// 参数 key: 加密密钥
// 返回值: Base64编码的加密字符串
func EncryptString(str, key string) string {
	s := []byte(str)
	k := []byte(key)
	b64 := base64.StdEncoding
	return b64.EncodeToString(Encrypt(s, k))
}

// DecryptString 解密字符串
// 对Base64编码的加密字符串进行解密
// 参数 str: 要解密的Base64编码字符串
// 参数 key: 解密密钥
// 返回值: 解密后的字符串和可能的错误
func DecryptString(str, key string) (string, error) {
	k := []byte(key)
	b64 := base64.StdEncoding
	decodeStr, err := b64.DecodeString(str)
	if err != nil {
		return "", err
	}
	result := Decrypt([]byte(decodeStr), k)
	return string(result), nil
}

// EncryptStdToURLString 加密字符串并转换为URL安全格式
// 将Base64标准编码转换为URL安全格式
// 参数 str: 要加密的字符串
// 参数 key: 加密密钥
// 返回值: URL安全的加密字符串
func EncryptStdToURLString(str, key string) string {
	return encryptBase64ToUrlFormat(EncryptString(str, key))
}

// DecryptURLToStdString 解密URL安全格式的字符串
// 将URL安全格式转换回标准Base64格式并解密
// 参数 str: 要解密的URL安全格式字符串
// 参数 key: 解密密钥
// 返回值: 解密后的字符串和可能的错误
func DecryptURLToStdString(str, key string) (string, error) {
	return DecryptString(decryptBase64ToStdFormat(str), key)
}

// encryptBase64ToUrlFormat 将标准Base64字符串转换为URL安全格式
// 替换 + 为 -, / 为 _, = 为 ~
// 参数 str: 标准Base64编码的字符串
// 返回值: URL安全格式的字符串
func encryptBase64ToUrlFormat(str string) string {
	str = strings.Replace(str, "+", "-", -1)
	str = strings.Replace(str, "/", "_", -1)
	str = strings.Replace(str, "=", "~", -1)
	return str
}

// decryptBase64ToStdFormat 将URL安全格式字符串转换为标准Base64格式
// 替换 - 为 +, _ 为 /, ~ 为 =
// 参数 str: URL安全格式的字符串
// 返回值: 标准Base64格式的字符串
func decryptBase64ToStdFormat(str string) string {
	str = strings.Replace(str, "-", "+", -1)
	str = strings.Replace(str, "_", "/", -1)
	str = strings.Replace(str, "~", "=", -1)
	return str
}

// ENC_PREFIX 加密字符串前缀，用于标识字符串是否已加密
const ENC_PREFIX = "enc:"

// EncryptAuto 自动加密字符串
// 如果字符串已经是加密的(以ENC_PREFIX开头)或为空，则不做处理
// 参数 str: 要加密的字符串
// 参数 key: 加密密钥
// 返回值: 加密后的字符串，带有前缀标记
func EncryptAuto(str string, key string) string {
	if str == "" || strings.HasPrefix(str, ENC_PREFIX) {
		return str
	}
	return ENC_PREFIX + EncryptString(str, key)
}

// DecryptAuto 自动解密字符串
// 检查字符串是否以ENC_PREFIX开头，是则解密，否则返回原字符串
// 参数 str: 可能加密的字符串
// 参数 key: 解密密钥
// 返回值: 解密后的字符串
func DecryptAuto(str string, key string) string {
	if str == "" || !strings.HasPrefix(str, ENC_PREFIX) {
		return str
	}
	str = str[len(ENC_PREFIX):]
	result, err := DecryptString(str, key)
	if err == nil && result != "" {
		return result
	} else {
		return str
	}
}
