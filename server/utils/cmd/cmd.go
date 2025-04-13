// cmd包提供了执行外部命令的工具函数
package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// Exec 使用默认20秒超时执行命令
// 参数:
//   - cmdStr: 要执行的命令字符串
//
// 返回:
//   - string: 命令执行的标准输出
//   - error: 执行过程中的错误
func Exec(cmdStr string) (string, error) {
	return ExecWithTimeOut(cmdStr, 20*time.Second)
}

// ERR_CMD_TIMEOUT 命令超时错误的常量标识符
const ERR_CMD_TIMEOUT = "ErrCmdTimeout"

// handleErr 处理命令执行错误，格式化输出和错误信息
// 参数:
//   - stdout: 标准输出缓冲区
//   - stderr: 标准错误缓冲区
//   - err: 原始错误
//
// 返回:
//   - string: 格式化后的错误消息
//   - error: 原始错误
func handleErr(stdout, stderr bytes.Buffer, err error) (string, error) {
	errMsg := ""
	if len(stderr.String()) != 0 {
		errMsg = fmt.Sprintf("stderr: %s", stderr.String())
	}
	if len(stdout.String()) != 0 {
		if len(errMsg) != 0 {
			errMsg = fmt.Sprintf("%s; stdout: %s", errMsg, stdout.String())
		} else {
			errMsg = fmt.Sprintf("stdout: %s", stdout.String())
		}
	}
	return errMsg, err
}

// ExecWithTimeOut 执行命令并设置超时时间
// 参数:
//   - cmdStr: 要执行的命令字符串
//   - timeout: 超时时间
//
// 返回:
//   - string: 命令执行的标准输出
//   - error: 执行过程中的错误，如超时则返回ERR_CMD_TIMEOUT错误
func ExecWithTimeOut(cmdStr string, timeout time.Duration) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return "", err
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	after := time.After(timeout)
	select {
	case <-after:
		_ = cmd.Process.Kill()
		return "", errors.New(ERR_CMD_TIMEOUT)
	case err := <-done:
		if err != nil {
			return handleErr(stdout, stderr, err)
		}
	}

	return stdout.String(), nil
}

// ExecContainerScript 在Docker容器中执行命令
// 参数:
//   - containerName: 容器名称
//   - cmdStr: 要在容器中执行的命令
//   - timeout: 超时时间
//
// 返回:
//   - error: 执行过程中的错误
func ExecContainerScript(containerName, cmdStr string, timeout time.Duration) error {
	cmdStr = fmt.Sprintf("docker exec -i %s sh -c '%s'", containerName, cmdStr)
	out, err := ExecWithTimeOut(cmdStr, timeout)
	if err != nil {
		if out != "" {
			return fmt.Errorf("%s; err: %v", out, err)
		}
		return err
	}
	return nil
}

// ExecCronjobWithTimeOut 执行定时任务并将输出重定向到指定文件
// 参数:
//   - cmdStr: 要执行的命令
//   - workdir: 工作目录
//   - outPath: 输出文件路径
//   - timeout: 超时时间
//
// 返回:
//   - error: 执行过程中的错误
func ExecCronjobWithTimeOut(cmdStr, workdir, outPath string, timeout time.Duration) error {
	file, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}
	cmd.Dir = workdir
	cmd.Stdout = file
	cmd.Stderr = file
	if err := cmd.Start(); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	after := time.After(timeout)
	select {
	case <-after:
		_ = cmd.Process.Kill()
		return errors.New(ERR_CMD_TIMEOUT)
	case err := <-done:
		if err != nil {
			return err
		}
	}
	return nil
}

// Execf 使用格式化字符串执行命令
// 参数:
//   - cmdStr: 命令格式字符串
//   - a: 格式化参数
//
// 返回:
//   - string: 命令执行的标准输出
//   - error: 执行过程中的错误
func Execf(cmdStr string, a ...interface{}) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", fmt.Sprintf(cmdStr, a...))
	} else {
		cmd = exec.Command("sh", "-c", fmt.Sprintf(cmdStr, a...))
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return handleErr(stdout, stderr, err)
	}
	return stdout.String(), nil
}

// ExecWithCheck 以给定的命令名和参数执行命令
// 参数:
//   - name: 命令名称
//   - a: 命令参数
//
// 返回:
//   - string: 命令执行的标准输出
//   - error: 执行过程中的错误
func ExecWithCheck(name string, a ...string) (string, error) {
	cmd := exec.Command(name, a...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return handleErr(stdout, stderr, err)
	}
	return stdout.String(), nil
}

// ExecScript 执行脚本文件
// 参数:
//   - scriptPath: 脚本文件路径
//   - workDir: 工作目录
//
// 返回:
//   - string: 命令执行的标准输出
//   - error: 执行过程中的错误，默认10分钟超时
func ExecScript(scriptPath, workDir string) (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", scriptPath)
	} else {
		cmd = exec.Command("sh", scriptPath)
	}
	var stdout, stderr bytes.Buffer
	cmd.Dir = workDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return "", err
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	after := time.After(10 * time.Minute)
	select {
	case <-after:
		_ = cmd.Process.Kill()
		return "", errors.New(ERR_CMD_TIMEOUT)
	case err := <-done:
		if err != nil {
			return handleErr(stdout, stderr, err)
		}
	}

	return stdout.String(), nil
}

// ExecCmd 执行命令并返回错误
// 参数:
//   - cmdStr: 要执行的命令字符串
//
// 返回:
//   - error: 执行过程中的错误，包含命令输出
func ExecCmd(cmdStr string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error : %v, output: %s", err, output)
	}
	return nil
}

// ExecCmdWithDir 在指定目录执行命令
// 参数:
//   - cmdStr: 要执行的命令字符串
//   - workDir: 工作目录
//
// 返回:
//   - error: 执行过程中的错误，包含命令输出
func ExecCmdWithDir(cmdStr, workDir string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error : %v, output: %s", err, output)
	}
	return nil
}

// ExecShellWithTimeOut 使用日志记录器执行命令并设置超时
// 参数:
//   - cmdStr: 要执行的命令字符串
//   - workdir: 工作目录
//   - logger: 日志记录器
//   - timeout: 超时时间
//
// 返回:
//   - error: 执行过程中的错误
func ExecShellWithTimeOut(cmdStr, workdir string, logger *log.Logger, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/c", cmdStr)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", cmdStr)
	}
	cmd.Dir = workdir
	cmd.Stdout = logger.Writer()
	cmd.Stderr = logger.Writer()
	if err := cmd.Start(); err != nil {
		return err
	}
	err := cmd.Wait()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return errors.New(ERR_CMD_TIMEOUT)
	}
	return err
}
