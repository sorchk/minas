package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"server/service/basic"
	"server/service/database"
	"server/utils/global"
	"server/utils/logger"
	"server/utils/mfa"
	"server/utils/tlscert"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var resetPwdCmd = &cobra.Command{
	Use:   "resetpwd",
	Short: "重置密码",
	Long:  `重置密码.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitConfig(cmd)
		// 初始化日志
		logger.Init()
		database.Init()
	},
	Run: func(cmd *cobra.Command, args []string) {
		resetPwd()
	},
}
var resetMfaCmd = &cobra.Command{
	Use:   "resetmfa",
	Short: "重置MFA",
	Long:  `重置MFA.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitConfig(cmd)
		// 初始化日志
		logger.Init()
		database.Init()
	},
	Run: func(cmd *cobra.Command, args []string) {
		resetMfa()
	},
}
var disableMfaCmd = &cobra.Command{
	Use:   "disablemfa",
	Short: "禁用MFA",
	Long:  `禁用MFA.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitConfig(cmd)
		// 初始化日志
		logger.Init()
		database.Init()
	},
	Run: func(cmd *cobra.Command, args []string) {
		disableMfa()
	},
}
var genCertCmd = &cobra.Command{
	Use:   "gencert",
	Short: "生成证书",
	Long:  `生成证书.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		InitConfig(cmd)
		// 初始化日志
		logger.Init()
	},
	Run: func(cmd *cobra.Command, args []string) {
		genCert()
	},
}

func init() {
	// 添加命令
	rootCmd.AddCommand(resetPwdCmd)
	rootCmd.AddCommand(resetMfaCmd)
	rootCmd.AddCommand(disableMfaCmd)
	rootCmd.AddCommand(genCertCmd)

}

func genCert() {
	parentDir := filepath.Dir(CERT_CRT_PATH)
	_, err := os.Stat(parentDir)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(parentDir, 0755)
		}
	}
	var domains string
	fmt.Print("请输入IP或域名，多个域名用逗号分隔：")
	fmt.Scanf("%s", &domains)
	log.Printf("%s\n", domains)
	tlscert.GenerateSelfSignedCert(CERT_CRT_PATH, CERT_KEY_PATH, "", "", strings.Split(domains, ",")...)
	log.Println("证书已生成")
}
func resetMfa() {
	user := basic.User{}
	// 查询数据库中的第一个用户
	err := global.DB.Model(user).First(&user).Order("ID asc").Limit(1).Error
	if err != nil {
		log.Println("用户不存在")
	}
	// 生成新的MFA密钥，长度为15
	user.MFACode = mfa.GenSecret(15)
	// 启用MFA认证
	user.MFAEnable = 1
	// 更新用户的MFA相关字段
	user.Update(&user, "mfa_enable", "mfa_code")
	// 输出新设置的MFA密钥
	log.Println("MFA已设置：" + user.MFACode)
}

func disableMfa() {
	user := basic.User{}
	// 查询数据库中的第一个用户
	err := global.DB.Model(user).First(&user).Order("ID asc").Limit(1).Error
	if err != nil {
		log.Println("用户不存在")
	}
	// 将MFA认证状态设置为禁用
	user.MFAEnable = 0
	// 更新用户的MFA启用状态
	user.Update(&user, "mfa_enable")
	// 输出操作结果
	log.Println("MFA已禁用")
}

func resetPwd() {
	var pwd string
	// 提示用户输入新密码
	fmt.Print("请输入新密码：")
	// 读取用户输入的密码
	fmt.Scanf("%s", &pwd)
	// 输出密码日志（注意：生产环境中应避免记录明文密码）
	log.Printf("%s\n", pwd)
	// 使用bcrypt算法对密码进行加密，使用默认复杂度
	pwdstr, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)

	if err != nil {
		// 密码加密失败时记录错误并返回
		log.Println("密码加密失败")
		return
	}
	log.Printf("%s\n", pwdstr)
	user := basic.User{}
	// 查询数据库中的第一个用户
	err = global.DB.Model(user).First(&user).Order("ID asc").Limit(1).Error
	if err != nil {
		log.Println("用户不存在")
	}
	// 更新用户的密码字段
	global.DB.Model(&user).Update("Password", pwdstr)
	// 输出操作结果
	log.Println("密码已修改")
}
