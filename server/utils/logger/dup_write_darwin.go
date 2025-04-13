package logger

import (
	"os"
	"runtime"

	"golang.org/x/sys/unix"
)

// stdErrFileHandler 标准错误文件处理器，持有重定向后的标准错误文件描述符
var stdErrFileHandler *os.File

// dupWrite 将标准错误输出重定向到指定文件
// 在macOS系统下，使用unix.Dup2系统调用实现文件描述符的复制
// 参数 file: 要重定向到的目标文件
// 返回值: 如果有错误发生返回error，否则返回nil
func dupWrite(file *os.File) error {
	// 保存文件句柄，防止被垃圾回收
	stdErrFileHandler = file

	// 使用unix.Dup2系统调用，将标准错误(stderr)的文件描述符复制为目标文件的文件描述符
	// 这样所有写入到stderr的内容都会被写入到指定文件
	if err := unix.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		return err
	}

	// 设置终结器，在垃圾回收时确保文件被正确关闭
	// 这可以防止资源泄露
	runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
		fd.Close()
	})
	return nil
}
