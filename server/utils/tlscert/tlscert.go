package tlscert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	_ "embed"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

// 使用go:embed指令嵌入CA证书和私钥文件
//
//go:embed ca.crt
var caCertBytes []byte // CA证书的字节数据

//go:embed ca.key
var caKeyBytes []byte // CA私钥的字节数据

// GenerateSelfSignedCert 生成自签名证书
// 使用预嵌入的CA证书签发新的TLS证书
// 参数 certPath: 输出证书文件路径
// 参数 keyPath: 输出私钥文件路径
// 参数 commonName: 证书主体通用名称，通常是域名
// 参数 organization: 证书所属组织名称，多个组织用逗号分隔
// 参数 domains: 证书支持的其他域名或IP地址列表
// 返回值: 生成并加载的TLS证书
func GenerateSelfSignedCert(certPath string, keyPath string, commonName string, organization string, domains ...string) tls.Certificate {

	// 读取CA证书
	block, _ := pem.Decode(caCertBytes)
	// 解析证书
	caCert, _ := x509.ParseCertificate(block.Bytes)

	// 读取CA私钥
	blockKey, _ := pem.Decode(caKeyBytes)
	// 解析私钥
	caPrivKey, _ := x509.ParsePKCS1PrivateKey(blockKey.Bytes)

	// 生成用户私钥 (2048位RSA密钥)
	userPriv, _ := rsa.GenerateKey(rand.Reader, 2048)

	// 证书模板基础设置
	notBefore := time.Now()                                       // 证书生效时间为当前时间
	notAfter := notBefore.Add(10 * 365 * 24 * time.Hour)          // 证书有效期10年
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)     // 创建128位的随机序列号上限
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit) // 生成随机序列号
	if err != nil {
		panic(err)
	}

	// 设置默认值，如果参数为空
	if commonName == "" {
		commonName = "localhost"
	}
	if organization == "" {
		organization = "localhost"
	}

	// 创建用户证书模板
	user := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   commonName,                       // 通用名称可以是域名或者IP地址
			Organization: strings.Split(organization, ","), // 将组织名称字符串按逗号分割
		},
		NotBefore:             notBefore,                                                    // 证书生效时间
		NotAfter:              notAfter,                                                     // 证书过期时间
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature, // 密钥用途
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},               // 扩展用途(服务器认证)
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}, // 默认支持本地回环地址
	}

	// 获取本机所有网络接口地址，添加到证书中
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			// 检查IP地址类型
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				// 只添加IPv4地址且排除某些特定子网
				if ipnet.IP.To4() != nil && !strings.HasSuffix(ipnet.String(), ".1/16") {
					user.IPAddresses = append(user.IPAddresses, net.ParseIP(ipnet.IP.String()))
				}
			}
		}
	}

	// 添加自定义域名和IP地址
	if len(domains) > 0 {
		for _, domain := range domains {
			// 判断是IP地址还是域名
			if ip := net.ParseIP(domain); ip != nil {
				user.IPAddresses = append(user.IPAddresses, ip) // 添加IP地址
			} else {
				user.DNSNames = append(user.DNSNames, domain) // 添加域名
			}
		}
	}

	// 使用CA证书签发用户证书
	userCertBytes, err := x509.CreateCertificate(rand.Reader, &user, caCert, &userPriv.PublicKey, caPrivKey)
	if err != nil {
		log.Println("生成证书数据失败:" + err.Error())
		panic(err)
	}

	// 将证书写入文件
	certOut, err := os.Create(certPath)
	if err != nil {
		log.Println("创建cert失败:" + err.Error())
		panic(err)
	}
	// 写入用户证书和CA证书
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: userCertBytes})
	pem.Encode(certOut, block) // 同时写入CA证书
	certOut.Close()

	// 将私钥写入文件 (权限设为600，只有所有者可读写)
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Println("创建key失败:" + err.Error())
		panic(err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(userPriv)})
	keyOut.Close()

	// 测试加载证书和私钥是否成功
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Println("测试读取证书：" + err.Error())
		panic(err)
	}
	return cert
}
