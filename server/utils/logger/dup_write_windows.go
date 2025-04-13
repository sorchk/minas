package logger

import (
	"os"
	"runtime"

	"golang.org/x/sys/windows"
)

// stdErrFileHandler 标准错误文件处理器，持有重定向后的标准错误文件描述符
var stdErrFileHandler *os.File

// dupWrite 将标准错误输出重定向到指定文件
// 在Windows系统下，使用windows.SetStdHandle函数实现标准错误的重定向
// 参数 file: 要重定向到的目标文件
// 返回值: 如果有错误发生返回error，否则返回nil
func dupWrite(file *os.File) error {
	// 保存文件句柄，防止被垃圾回收
	stdErrFileHandler = file

	// 使用windows.SetStdHandle函数，将标准错误(STD_ERROR_HANDLE)设置为目标文件的句柄
	// 这样所有写入到stderr的内容都会被写入到指定文件
	if err := windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(file.Fd())); err != nil {
		return err
	}

	// 设置终结器，在垃圾回收时确保文件被正确关闭
	// 这可以防止资源泄露
	runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
		fd.Close()
	})
	return nil
}
