package rclone

import (
	"fmt"
	"server/utils/cmd"
	"server/utils/config"
	"server/utils/logger"
	"strings"
	"time"
)

// CmdBisync 执行RClone双向同步(bisync)命令
// 在两个文件系统之间进行双向同步，并保留备份
// 参数:
//   - workdir: 工作目录，命令执行的基础目录
//   - logPath: 日志文件路径，用于记录执行过程的日志
//   - srcFs: 源文件系统路径
//   - dstFs: 目标文件系统路径
//   - createEmptySrcDirs: 是否在目标中创建源文件系统的空目录
//   - removeEmptyDirs: 是否删除空目录
//   - showDebug: 是否显示详细调试信息
//   - includes: 包含的文件模式列表，格式为通配符
//   - excludes: 排除的文件模式列表，格式为通配符
//
// 返回值: 命令执行过程中的错误，成功则为nil
func CmdBisync(workdir string, logPath string, srcFs string, dstFs string, createEmptySrcDirs bool, removeEmptyDirs bool, showDebug bool, includes []string, excludes []string) error {

	// 构建bisync基础命令，设置源、目标和备份目录
	shell := fmt.Sprintf("%s bisync %s %s --backup-dir1 %s --backup-dir2 %s", GetRClonePath(), srcFs, dstFs, srcFs+"_backupdir1", dstFs+"_backupdir2")

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

	// 根据选项添加额外参数
	if createEmptySrcDirs {
		shell += " --create-empty-src-dirs"
	}
	if removeEmptyDirs {
		shell += " --remove-empty-dirs"
	}
	if showDebug {
		shell += " -v" // 详细模式，输出更多日志
	}

	// 记录完整命令到日志
	logger.LOG.Infof("----------CmdSync------------:%s\n", shell)

	// 执行命令，设置超时时间
	err := cmd.ExecCronjobWithTimeOut(shell, workdir, logPath, time.Hour*time.Duration(config.CONF.RClone.CmdTimeOut))
	return err
}
