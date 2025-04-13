package rclone

import (
	"fmt"
	"server/utils/cmd"
	"server/utils/config"
	"server/utils/logger"
	"strings"
	"time"
)

// CmdDelete 执行RClone删除命令，根据指定条件删除文件
// 可用于定期清理过期文件或备份
// 参数:
//   - workdir: 工作目录，命令执行的基础目录
//   - logPath: 日志文件路径，用于记录执行过程的日志
//   - workDir: 要清理的目录路径，支持远程存储路径
//   - includes: 包含的文件模式列表，格式为通配符
//   - excludes: 排除的文件模式列表，格式为通配符
//   - offset: 最小文件年龄，只删除比这个时间更早的文件
//   - isTest: 是否为测试模式(干运行)，设为true时只显示会删除哪些文件而不实际删除
//   - showDebug: 是否显示详细调试信息
//
// 返回值: 命令执行过程中的错误，成功则为nil
func CmdDelete(workdir string, logPath string, workDir string, includes []string, excludes []string, offset time.Duration, isTest bool, showDebug bool) error {

	// 构建基础删除命令，设置目标目录和最小文件年龄
	shell := fmt.Sprintf("%s delete %s --min-age %s", GetRClonePath(), workDir, offset)

	// 添加配置文件参数
	if GetRConfigPath() != "" {
		shell += " --config " + GetRConfigPath()
	}

	// 添加包含文件模式参数
	if len(includes) > 0 {
		for i := range includes {
			include := strings.TrimSpace(includes[i])
			if include != "" {
				shell += fmt.Sprintf(" --include \"%s\"", include)
			}
		}
	}

	// 添加排除文件模式参数
	if len(excludes) > 0 {
		for i := range excludes {
			exclude := strings.TrimSpace(excludes[i])
			if exclude != "" {
				shell += fmt.Sprintf(" --exclude \"%s\"", exclude)
			}
		}
	}

	// 如果是测试模式，添加干运行参数
	if isTest {
		shell += " --dry-run" // 干运行模式，不实际删除文件
	}

	// 如果需要显示详细信息，添加详细输出参数
	if showDebug {
		shell += " -v" // 详细模式，输出更多日志
	}

	// 示例时间格式注释: 26538h19m10.167625239s
	logger.LOG.Infof("----------CmdDelete------------:%s\n", shell)

	// 执行命令，设置超时时间
	err := cmd.ExecCronjobWithTimeOut(shell, workdir, logPath, time.Hour*time.Duration(config.CONF.RClone.CmdTimeOut))
	return err
}
