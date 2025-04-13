package rclone

import (
	"fmt"
	"server/utils/cmd"
	"server/utils/logger"
	"strings"
	"time"

	"server/utils/config"
)

// CmdSync 执行RClone同步命令，将源目录同步到目标目录
// 同步操作会使目标目录与源目录的内容完全一致
// 参数:
//   - workdir: 工作目录，命令执行的基础目录
//   - logPath: 日志文件路径，用于记录执行过程的日志
//   - srcFs: 源文件系统路径，可以是本地路径或远程存储路径
//   - dstFs: 目标文件系统路径，可以是本地路径或远程存储路径
//   - createEmptySrcDirs: 是否在目标中创建源中的空目录
//   - showDebug: 是否显示详细调试信息
//   - includes: 包含的文件模式列表，格式为通配符
//   - excludes: 排除的文件模式列表，格式为通配符
//
// 返回值: 命令执行过程中的错误，成功则为nil
func CmdSync(workdir string, logPath string, srcFs string, dstFs string, createEmptySrcDirs bool, showDebug bool, includes []string, excludes []string) error {
	// 构建基础同步命令，设置源和目标路径
	shell := fmt.Sprintf("%s sync %s %s", GetRClonePath(), srcFs, dstFs)

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

	// 如果需要创建空源目录，添加相应参数
	if createEmptySrcDirs {
		shell += " --create-empty-src-dirs"
	}

	// 如果需要显示详细信息，添加详细输出参数
	if showDebug {
		shell += " -v" // 详细模式，输出更多日志
	}

	// 记录完整命令到日志
	logger.LOG.Infof("----------CmdSync------------:%s\n", shell)

	// 执行命令，设置超时时间
	err := cmd.ExecCronjobWithTimeOut(shell, workdir, logPath, time.Hour*time.Duration(config.CONF.RClone.CmdTimeOut))
	return err
}
